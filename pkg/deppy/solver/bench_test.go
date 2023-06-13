package solver

import (
	"context"
	"github.com/perdasilva/replee/pkg/deppy"
	"github.com/perdasilva/replee/pkg/deppy/constraints"
	"math/rand"
	"strconv"
	"testing"
)

var BenchmarkInput = func() []deppy.Variable {
	const (
		length      = 256
		seed        = 9
		pMandatory  = .1
		pDependency = .15
		nDependency = 6
		pConflict   = .05
		nConflict   = 3
	)

	id := func(i int) deppy.Identifier {
		return deppy.Identifier(strconv.Itoa(i))
	}

	variable := func(i int) TestVariable {
		var c []deppy.Constraint
		if rand.Float64() < pMandatory {
			c = append(c, constraints.Mandatory())
		}
		if rand.Float64() < pDependency {
			n := rand.Intn(nDependency-1) + 1
			var d []deppy.Identifier
			for x := 0; x < n; x++ {
				y := i
				for y == i {
					y = rand.Intn(length)
				}
				d = append(d, id(y))
			}
			c = append(c, constraints.Dependency(d...))
		}
		if rand.Float64() < pConflict {
			n := rand.Intn(nConflict-1) + 1
			for x := 0; x < n; x++ {
				y := i
				for y == i {
					y = rand.Intn(length)
				}
				c = append(c, constraints.Conflict(id(y)))
			}
		}
		return TestVariable{
			identifier:  id(i),
			constraints: c,
		}
	}

	rand.Seed(seed)
	result := make([]deppy.Variable, length)
	for i := range result {
		result[i] = variable(i)
	}
	return result
}()

func BenchmarkSolve(b *testing.B) {
	for i := 0; i < b.N; i++ {
		s, err := NewSolver(WithInput(BenchmarkInput))
		if err != nil {
			b.Fatalf("failed to initialize solver: %s", err)
		}
		_, err = s.Solve(context.Background())
		if err != nil {
			b.Fatalf("failed to initialize solver: %s", err)
		}
	}
}

func BenchmarkNewInput(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, err := NewSolver(WithInput(BenchmarkInput))
		if err != nil {
			b.Fatalf("failed to initialize solver: %s", err)
		}
	}
}
