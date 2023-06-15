package repl

import (
	"context"
	"github.com/perdasilva/replee/pkg/deppy"
)

type VariableSourceWithContext struct {
	variableSource deppy.VariableSource
	ctx            context.Context
}

func NewVariableSourceWithContext(ctx context.Context, variableSource deppy.VariableSource) *VariableSourceWithContext {
	return &VariableSourceWithContext{
		variableSource: variableSource,
		ctx:            ctx,
	}
}

func (v *VariableSourceWithContext) VariableSourceID() deppy.Identifier {
	return v.variableSource.VariableSourceID()
}

func (v *VariableSourceWithContext) VariableFilterFunc() deppy.VarFilterFn {
	return v.variableSource.VariableFilterFunc()
}

func (v *VariableSourceWithContext) Update(resolution deppy.MutableResolutionProblem, nextVariable deppy.Variable) error {
	return v.variableSource.Update(v.ctx, resolution, nextVariable)
}

func (v *VariableSourceWithContext) Finalize(resolution deppy.MutableResolutionProblem) error {
	return v.variableSource.Finalize(v.ctx, resolution)
}
