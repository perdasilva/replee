package main

import (
	"context"
	"fmt"
	"github.com/dop251/goja"
	"github.com/perdasilva/replee/pkg/replee/repl"
	"log"
	"syscall/js"
)

func main() {
	ctx := context.Background()
	vm := goja.New()
	vm.SetFieldNameMapper(goja.TagFieldNameMapper("json", true))

	if err := repl.BootstrapRepleeVM(ctx, vm); err != nil {
		log.Fatal(err)
	}

	js.Global().Set("replee", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		input := args[0].String()
		value, err := vm.RunString(input)
		if err != nil {
			return err.Error() // Convert error to string
		}
		switch v := value.Export().(type) { // Convert Goja value to a basic Go type
		case string:
			return v
		case int:
			return v
		case bool:
			return v
		default:
			// This will catch other types (like Goja's native objects, arrays, etc.)
			// You may need to add more cases here depending on your needs
			return fmt.Sprintf("%v", v)
		}
	}))

	<-(make(chan struct{}))
}
