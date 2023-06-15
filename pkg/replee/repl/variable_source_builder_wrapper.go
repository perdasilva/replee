package repl

import (
	"context"
	"fmt"
	"github.com/dop251/goja"
	"github.com/perdasilva/replee/pkg/deppy"
	"github.com/perdasilva/replee/pkg/deppy/variable_sources"
)

type VariableSourceBuilder struct {
	builder variable_sources.VariableSourceBuilder
	ctx     context.Context
	vm      *goja.Runtime
}

func NewVariableSourceBuilder(ctx context.Context, vm *goja.Runtime) func(variableSourceID deppy.Identifier) *VariableSourceBuilder {
	return func(variableSourceID deppy.Identifier) *VariableSourceBuilder {
		return &VariableSourceBuilder{
			ctx:     ctx,
			vm:      vm,
			builder: variable_sources.NewVariableSourceBuilder(variableSourceID),
		}
	}
}

func (v *VariableSourceBuilder) WithUpdateFn(call goja.FunctionCall) goja.Value {
	cb, ok := goja.AssertFunction(call.Argument(0))
	if !ok {
		return v.vm.ToValue(fmt.Errorf("first argument is not a function"))
	}

	v.builder.WithUpdateFn(func(ctx context.Context, problem deppy.MutableResolutionProblem, variable deppy.Variable) error {
		// Call the passed function.
		p := v.vm.ToValue(problem)
		v := v.vm.ToValue(variable)
		ret, err := cb(nil, p, v)
		if err != nil {
			return err
		}
		if goja.IsNull(ret) || goja.IsUndefined(ret) {
			return nil
		}
		errString, ok := ret.Export().(string)
		if !ok {
			return fmt.Errorf("expected string return value")
		}
		return fmt.Errorf(errString)
	})
	return v.vm.ToValue(v)
}

func (v *VariableSourceBuilder) WithFinalizeFn(finalizeFunc func(problem deppy.MutableResolutionProblem) error) *VariableSourceBuilder {
	v.builder.WithFinalizeFn(func(ctx context.Context, problem deppy.MutableResolutionProblem) error {
		return finalizeFunc(problem)
	})
	return v
}

func (v *VariableSourceBuilder) WithVariableFilterFn(variableSourceFilterFn deppy.VarFilterFn) *VariableSourceBuilder {
	v.builder.WithVariableFilterFn(variableSourceFilterFn)
	return v
}

func (v *VariableSourceBuilder) Build() *VariableSourceWithContext {
	vs := v.builder.Build(v.ctx)
	return NewVariableSourceWithContext(v.ctx, vs)
}
