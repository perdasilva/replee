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

  v.builder.WithUpdateFn(func(ctx context.Context, problem deppy.MutableResolutionProblem, variable deppy.MutableVariable) error {
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

func (v *VariableSourceBuilder) WithFinalizeFn(call goja.FunctionCall) goja.Value {
  cb, ok := goja.AssertFunction(call.Argument(0))
  if !ok {
    return v.vm.ToValue(fmt.Errorf("first argument is not a function"))
  }

  v.builder.WithFinalizeFn(func(ctx context.Context, problem deppy.MutableResolutionProblem) error {
    // Call the passed function.
    p := v.vm.ToValue(problem)
    ret, err := cb(nil, p)
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

func (v *VariableSourceBuilder) WithVariableFilterFn(call goja.FunctionCall) goja.Value {
  cb, ok := goja.AssertFunction(call.Argument(0))
  if !ok {
    return v.vm.ToValue(fmt.Errorf("first argument is not a function"))
  }

  v.builder.WithVariableFilterFn(func(input deppy.Variable) bool {
    // Call the passed function.
    jsInput := v.vm.ToValue(input)
    ret, err := cb(nil, jsInput)
    if err != nil {
      return false
    }
    if goja.IsNull(ret) || goja.IsUndefined(ret) {
      return false
    }
    out, ok := ret.Export().(bool)
    if !ok {
      return false
    }
    return out
  })
  return v.vm.ToValue(v)
}

func (v *VariableSourceBuilder) Build() *VariableSourceWithContext {
  vs := v.builder.Build(v.ctx)
  return NewVariableSourceWithContext(v.ctx, vs)
}
