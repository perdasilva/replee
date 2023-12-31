package solver

import (
	"github.com/perdasilva/replee/pkg/deppy"
	"github.com/perdasilva/replee/pkg/deppy/constraints"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestOrder(t *testing.T) {
	type tc struct {
		Name       string
		Constraint deppy.Constraint
		Expected   []deppy.Identifier
	}

	for _, tt := range []tc{
		{
			Name:       "mandatory",
			Constraint: constraints.Mandatory("a"),
		},
		{
			Name:       "prohibited",
			Constraint: constraints.Prohibited("a"),
		},
		{
			Name:       "dependency",
			Constraint: constraints.Dependency("dcid", "a", "b", "c"),
			Expected:   []deppy.Identifier{"a", "b", "c"},
		},
		{
			Name:       "conflict",
			Constraint: constraints.Conflict("cid", "a"),
		},
	} {
		t.Run(tt.Name, func(t *testing.T) {
			assert.Equal(t, tt.Expected, tt.Constraint.Order())
		})
	}
}
