package repl

import (
	"context"
	"github.com/dop251/goja"
	"github.com/perdasilva/replee/pkg/deppy"
	"github.com/perdasilva/replee/pkg/deppy/resolution"
	"github.com/perdasilva/replee/pkg/deppy/resolver"
	"github.com/perdasilva/replee/pkg/deppy/variables"
	"reflect"
)

func BootstrapRepleeVM(ctx context.Context, vm *goja.Runtime) error {
	s := resolver.NewDeppyResolver()
	solveWrapper := func(p *resolution.MutableResolutionProblem, options ...resolver.Option) (*resolver.Solution, error) {
		return s.Solve(ctx, p, options...)
	}

	return vm.Set("deppy", map[string]interface{}{
		"newResolutionProblemBuilder": NewResolutionProblemBuilderWithCtx(ctx),
		"newProblem":                  resolution.NewMutableResolutionProblem,
		"newVariable":                 variables.NewMutableVariable,
		"solve":                       solveWrapper,
		"ctx":                         context.Background,
		"id":                          reflect.ValueOf(deppy.Identifierf),
		"newVariableSourceBuilder":    NewVariableSourceBuilder(ctx, vm),
		"opts": map[string]interface{}{
			"addAllVariablesToSolution": resolver.AddAllVariablesToSolution,
			"disableOrderPreference":    resolver.DisableOrderPreference,
		},
	})
}
