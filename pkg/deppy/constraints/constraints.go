package constraints

import (
	"fmt"
	"github.com/go-air/gini/z"
	"github.com/perdasilva/replee/pkg/deppy"
	"github.com/perdasilva/replee/pkg/deppy/utils"
	"strings"
	"sync"
)

const (
	ConstraintKindMandatory  = "deppy.constraint.mandatory"
	ConstraintKindProhibited = "deppy.constraint.prohibited"
	ConstraintKindConflict   = "deppy.constraint.conflict"
	ConstraintKindDependency = "deppy.constraint.dependency"
	ConstraintKindAtMost     = "deppy.constraint.atmost"
)

var _ deppy.Constraint = &MutableConstraint{}

type MandatoryConstraint struct {
	MutableConstraint
}

func Mandatory() *MandatoryConstraint {
	return &MandatoryConstraint{
		MutableConstraint: MutableConstraint{
			kind:       ConstraintKindMandatory,
			properties: map[string]interface{}{},
			lock:       sync.RWMutex{},
		},
	}
}

func (constraint *MandatoryConstraint) String(subject deppy.Identifier) string {
	return fmt.Sprintf("%s is mandatory", subject)
}

func (constraint *MandatoryConstraint) Apply(lm deppy.LitMapping, subject deppy.Identifier) z.Lit {
	return lm.LitOf(subject)
}

func (constraint *MandatoryConstraint) Order() []deppy.Identifier {
	return nil
}

func (constraint *MandatoryConstraint) Anchor() bool {
	return true
}

func (constraint *MandatoryConstraint) Merge(other deppy.Constraint) error {
	if err := constraint.MutableConstraint.Merge(other); err != nil {
		return err
	}
	if _, ok := other.(*MandatoryConstraint); !ok {
		return deppy.ConflictErrorf("cannot merge constraints of different kind [%T != %T]", constraint, other)
	}
	return nil
}

type ProhibitedConstraint struct {
	MutableConstraint
}

func Prohibited() *ProhibitedConstraint {
	return &ProhibitedConstraint{
		MutableConstraint: MutableConstraint{
			kind:       ConstraintKindProhibited,
			properties: map[string]interface{}{},
			lock:       sync.RWMutex{},
		},
	}
}

func (constraint *ProhibitedConstraint) Merge(other deppy.Constraint) error {
	if err := constraint.MutableConstraint.Merge(other); err != nil {
		return err
	}
	if _, ok := other.(*ProhibitedConstraint); !ok {
		return deppy.ConflictErrorf("cannot merge constraints of different kind [%T != %T]", constraint, other)
	}
	return nil
}

func (constraint *ProhibitedConstraint) String(subject deppy.Identifier) string {
	return fmt.Sprintf("%s is ProhibitedConstraint", subject)
}

func (constraint *ProhibitedConstraint) Apply(lm deppy.LitMapping, subject deppy.Identifier) z.Lit {
	return lm.LitOf(subject).Not()
}

func (constraint *ProhibitedConstraint) Order() []deppy.Identifier {
	return nil
}

func (constraint *ProhibitedConstraint) Anchor() bool {
	return false
}

type ConflictConstraint struct {
	MutableConstraint
	conflictingVariableID deppy.Identifier
	lock                  sync.RWMutex
}

func Conflict(conflict deppy.Identifier) *ConflictConstraint {
	return &ConflictConstraint{
		MutableConstraint: MutableConstraint{
			kind:       ConstraintKindConflict,
			properties: map[string]interface{}{},
			lock:       sync.RWMutex{},
		},
		conflictingVariableID: conflict,
	}
}

func (constraint *ConflictConstraint) Merge(other deppy.Constraint) error {
	if err := constraint.MutableConstraint.Merge(other); err != nil {
		return err
	}
	if cc, ok := other.(*ConflictConstraint); ok {
		if cc.conflictingVariableID != constraint.conflictingVariableID {
			return deppy.ConflictErrorf("cannot merge constraints with different conflicting variable [%s != %s]", constraint.conflictingVariableID, cc.conflictingVariableID)
		}
		return nil
	}
	return deppy.ConflictErrorf("cannot merge constraints of different kind [%T != %T]", constraint, other)
}

func (constraint *ConflictConstraint) String(subject deppy.Identifier) string {
	return fmt.Sprintf("%s conflicts with %s", subject, constraint.conflictingVariableID)
}

func (constraint *ConflictConstraint) Apply(lm deppy.LitMapping, subject deppy.Identifier) z.Lit {
	return lm.LogicCircuit().Or(lm.LitOf(subject).Not(), lm.LitOf(constraint.conflictingVariableID).Not())
}

func (constraint *ConflictConstraint) Order() []deppy.Identifier {
	return nil
}

func (constraint *ConflictConstraint) Anchor() bool {
	return false
}

func (constraint *ConflictConstraint) SetConflictingVariableID(id deppy.Identifier) error {
	constraint.lock.Lock()
	defer constraint.lock.Unlock()
	if constraint.conflictingVariableID == "" || constraint.conflictingVariableID == id {
		constraint.conflictingVariableID = id
		return nil
	}
	return deppy.FatalError("conflicting variable id already set")
}

type DependencyConstraint struct {
	MutableConstraint
	utils.ActivationSet[deppy.Identifier]
}

func Dependency(dependencies ...deppy.Identifier) *DependencyConstraint {
	c := &DependencyConstraint{
		MutableConstraint: MutableConstraint{
			kind:       ConstraintKindDependency,
			properties: map[string]interface{}{},
			lock:       sync.RWMutex{},
		},
		ActivationSet: *utils.NewActivationSet[deppy.Identifier](),
	}
	for _, dependency := range dependencies {
		c.Activate(dependency)
	}
	return c
}

func (d *DependencyConstraint) Merge(other deppy.Constraint) error {
	if err := d.MutableConstraint.Merge(other); err != nil {
		return err
	}
	if cc, ok := other.(*DependencyConstraint); ok {
		for _, element := range cc.Elements() {
			if active, _ := cc.IsActivated(element); active {
				d.Activate(element)
			} else {
				d.Deactivate(element)
			}
		}
		return nil
	}
	return deppy.ConflictErrorf("cannot merge constraints of different kind [%T != %T]", d, other)
}

func (d *DependencyConstraint) Apply(lm deppy.LitMapping, subject deppy.Identifier) z.Lit {
	dependencyIDs := d.Elements()
	m := lm.LitOf(subject).Not()
	for _, each := range dependencyIDs {
		m = lm.LogicCircuit().Or(m, lm.LitOf(each))
	}
	return m
}

func (d *DependencyConstraint) Order() []deppy.Identifier {
	return d.Elements()
}

func (d *DependencyConstraint) Anchor() bool {
	return false
}

func (d *DependencyConstraint) String(subject deppy.Identifier) string {
	dependencyIDs := d.Elements()
	if len(dependencyIDs) == 0 {
		return fmt.Sprintf("%s has a DependencyConstraint without any candidates to satisfy it", subject)
	}
	s := make([]string, len(dependencyIDs))
	for i, each := range dependencyIDs {
		s[i] = string(each)
	}
	return fmt.Sprintf("%s requires at least one of %s", subject, strings.Join(s, ", "))
}

type AtMostConstraint struct {
	MutableConstraint
	utils.ActivationSet[deppy.Identifier]
	n    int
	lock sync.RWMutex
}

func AtMost(n int, variables ...deppy.Identifier) *AtMostConstraint {
	c := &AtMostConstraint{
		MutableConstraint: MutableConstraint{
			kind:       ConstraintKindAtMost,
			properties: map[string]interface{}{},
			lock:       sync.RWMutex{},
		},
		ActivationSet: *utils.NewActivationSet[deppy.Identifier](),
		n:             n,
	}
	for _, variable := range variables {
		c.Activate(variable)
	}
	return c
}

func (a *AtMostConstraint) Merge(other deppy.Constraint) error {
	if err := a.MutableConstraint.Merge(other); err != nil {
		return err
	}
	if cc, ok := other.(*AtMostConstraint); ok {
		if a.n != -1 && cc.n != -1 && a.n != cc.n {
			return deppy.ConflictErrorf("cannot merge constraints with different n [%d != %d]", a.n, cc.n)
		}
		if a.n == -1 {
			a.n = cc.n
		}
		for _, element := range cc.Elements() {
			if active, _ := cc.IsActivated(element); active {
				a.Activate(element)
			} else {
				a.Deactivate(element)
			}
		}
		return nil
	}
	return deppy.ConflictErrorf("cannot merge constraints of different kind [%T != %T]", a, other)
}

func (a *AtMostConstraint) Apply(lm deppy.LitMapping, _ deppy.Identifier) z.Lit {
	ids := a.Elements()
	ms := make([]z.Lit, len(ids))
	for i, each := range ids {
		ms[i] = lm.LitOf(each)
	}
	return lm.LogicCircuit().CardSort(ms).Leq(a.n)
}

func (a *AtMostConstraint) Order() []deppy.Identifier {
	return a.Elements()
}

func (a *AtMostConstraint) Anchor() bool {
	return false
}

func (a *AtMostConstraint) String(subject deppy.Identifier) string {
	ids := a.Elements()
	s := make([]string, len(ids))
	for i, each := range ids {
		s[i] = string(each)
	}
	return fmt.Sprintf("%s permits at most %d of %s", subject, a.n, strings.Join(s, ", "))
}

func (a *AtMostConstraint) SetN(n int) error {
	a.lock.Lock()
	defer a.lock.Unlock()
	if n < 0 {
		return deppy.FatalError("n must be greater than or equal to 0")
	}
	if a.n > 0 {
		return deppy.FatalError(fmt.Sprintf("n is already set to %d", a.n))
	}
	a.n = n
	return nil
}
