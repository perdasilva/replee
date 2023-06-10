package solver

import (
	"fmt"
	"github.com/perdasilva/replee/pkg/deppy"
	"strings"

	"github.com/go-air/gini/inter"
	"github.com/go-air/gini/logic"
	"github.com/go-air/gini/z"
)

type DuplicateIdentifier deppy.Identifier

func (e DuplicateIdentifier) Error() string {
	return fmt.Sprintf("duplicate identifier %q in input", deppy.Identifier(e))
}

type inconsistentLitMapping []error

func (inconsistentLitMapping) Error() string {
	return "internal solver failure"
}

var _ deppy.LitMapping = &litMapping{}

// litMapping performs translation between the input and output types of
// Solve (Constraints, Variables, etc.) and the variables that
// appear in the SAT formula.
type litMapping struct {
	inorder     []deppy.Variable
	variables   map[z.Lit]deppy.Variable
	lits        map[deppy.Identifier]z.Lit
	constraints map[z.Lit]deppy.AppliedConstraint
	c           *logic.C
	errs        inconsistentLitMapping
}

// newLitMapping returns a new litMapping with its state initialized based on
// the provided slice of Variables. This includes construction of
// the translation tables between Variables/Constraints and the
// inputs to the underlying solver.
func newLitMapping(variables []deppy.Variable) (*litMapping, error) {
	d := litMapping{
		inorder:     variables,
		variables:   make(map[z.Lit]deppy.Variable, len(variables)),
		lits:        make(map[deppy.Identifier]z.Lit, len(variables)),
		constraints: make(map[z.Lit]deppy.AppliedConstraint),
		c:           logic.NewCCap(len(variables)),
	}

	// First pass to assign lits:
	for _, variable := range variables {
		im := d.c.Lit()
		if _, ok := d.lits[variable.Identifier()]; ok {
			return nil, DuplicateIdentifier(variable.Identifier())
		}
		d.lits[variable.Identifier()] = im
		d.variables[im] = variable
	}

	for _, variable := range variables {
		for _, constraint := range variable.Constraints() {
			if constraint == nil || variable == nil {
				fmt.Println("nil constraint")
			}
			m := constraint.Apply(&d, variable.Identifier())
			if m == z.LitNull {
				// This constraint doesn't have a
				// useful representation in the SAT
				// inputs.
				continue
			}

			d.constraints[m] = deppy.AppliedConstraint{
				Variable:   variable,
				Constraint: constraint,
			}
		}
	}

	return &d, nil
}

// LogicCircuit returns the lit mappings internal logic circuit
// used by constraint for translation into boolean expressions processed by the solver
func (d *litMapping) LogicCircuit() *logic.C {
	return d.c
}

// LitOf returns the positive literal corresponding to the Variable
// with the given Identifier.
func (d *litMapping) LitOf(id deppy.Identifier) z.Lit {
	m, ok := d.lits[id]
	if ok {
		return m
	}
	d.errs = append(d.errs, fmt.Errorf("variable %q referenced but not provided", id))
	return z.LitNull
}

// VariableOf returns the Variable corresponding to the provided
// literal, or a zeroVariable if no such Variable exists.
func (d *litMapping) VariableOf(m z.Lit) deppy.Variable {
	i, ok := d.variables[m]
	if ok {
		return i
	}
	d.errs = append(d.errs, fmt.Errorf("no variable corresponding to %s", m))
	return zeroVariable{}
}

// ConstraintOf returns the constraint application corresponding to
// the provided literal, or a zeroConstraint if no such constraint
// exists.
func (d *litMapping) ConstraintOf(m z.Lit) deppy.AppliedConstraint {
	if a, ok := d.constraints[m]; ok {
		return a
	}
	d.errs = append(d.errs, fmt.Errorf("no constraint corresponding to %s", m))
	return deppy.AppliedConstraint{
		Variable:   zeroVariable{},
		Constraint: zeroConstraint{},
	}
}

// Error returns a single error value that is an aggregation of all
// errors encountered during a litMapping's lifetime, or nil if there have
// been no errors. A non-nil return value likely indicates a problem
// with the solver or constraint implementations.
func (d *litMapping) Error() error {
	if len(d.errs) == 0 {
		return nil
	}
	s := make([]string, len(d.errs))
	for i, err := range d.errs {
		s[i] = err.Error()
	}
	return fmt.Errorf("%d errors encountered: %s", len(s), strings.Join(s, ", "))
}

// AddConstraints adds the current constraints encoded in the embedded circuit to the
// solver g
func (d *litMapping) AddConstraints(g inter.S) {
	d.c.ToCnf(g)
}

func (d *litMapping) AssumeConstraints(s inter.S) {
	for m := range d.constraints {
		s.Assume(m)
	}
}

// CardinalityConstrainer constructs a sorting network to provide
// cardinality constraints over the provided slice of literals. Any
// new clauses and variables are translated to CNF and taught to the
// given inter.Adder, so this function will panic if it is in a test
// context.
func (d *litMapping) CardinalityConstrainer(g inter.Adder, ms []z.Lit) *logic.CardSort {
	clen := d.c.Len()
	cs := d.c.CardSort(ms)
	marks := make([]int8, clen, d.c.Len())
	for i := range marks {
		marks[i] = 1
	}
	for w := 0; w <= cs.N(); w++ {
		marks, _ = d.c.CnfSince(g, marks, cs.Leq(w))
	}
	return cs
}

// AnchorIdentifiers returns a slice containing the Identifiers of
// every Variable with at least one "Anchor" constraint, in the
// Order they appear in the input.
func (d *litMapping) AnchorIdentifiers() []deppy.Identifier {
	var ids []deppy.Identifier
	for _, variable := range d.inorder {
		for _, constraint := range variable.Constraints() {
			if constraint.Anchor() {
				ids = append(ids, variable.Identifier())
				break
			}
		}
	}
	return ids
}

func (d *litMapping) Variables(g inter.S) []deppy.Variable {
	var result []deppy.Variable
	for _, i := range d.inorder {
		if g.Value(d.LitOf(i.Identifier())) {
			result = append(result, i)
		}
	}
	return result
}

func (d *litMapping) Lits(dst []z.Lit) []z.Lit {
	if cap(dst) < len(d.inorder) {
		dst = make([]z.Lit, 0, len(d.inorder))
	}
	dst = dst[:0]
	for _, i := range d.inorder {
		m := d.LitOf(i.Identifier())
		dst = append(dst, m)
	}
	return dst
}

func (d *litMapping) Conflicts(g inter.Assumable) []deppy.AppliedConstraint {
	whys := g.Why(nil)
	as := make([]deppy.AppliedConstraint, 0, len(whys))
	for _, why := range whys {
		if a, ok := d.constraints[why]; ok {
			as = append(as, a)
		}
	}
	return as
}
