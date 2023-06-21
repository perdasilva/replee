package repl

import (
	"context"
	"github.com/perdasilva/replee/pkg/deppy"
	"github.com/perdasilva/replee/pkg/deppy/resolution"
)

type ResolutionProblemBuilder struct {
	ctx     context.Context
	builder resolution.ResolutionProblemBuilder
}

func NewResolutionProblemBuilderWithCtx(ctx context.Context) func(variableSourceID deppy.Identifier) *ResolutionProblemBuilder {
	return func(problemID deppy.Identifier) *ResolutionProblemBuilder {
		return &ResolutionProblemBuilder{
			ctx:     ctx,
			builder: resolution.NewResolutionProblemBuilder(problemID),
		}
	}
}
func (r *ResolutionProblemBuilder) WithVariableSources(variableSources ...deppy.VariableSource) *ResolutionProblemBuilder {
	r.builder.WithVariableSources(variableSources...)
	return r
}

func (r *ResolutionProblemBuilder) Build() (deppy.ResolutionProblem, error) {
	return r.builder.Build(r.ctx)
}
