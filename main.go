package main

import (
	"context"
	"github.com/dop251/goja"
	"github.com/perdasilva/replee/pkg/replee/repl"
	"github.com/perdasilva/replee/pkg/replee/terminal"
	"github.com/rivo/tview"
	"strings"
)

type ReplUI struct {
	app      *tview.Application
	terminal *terminal.RepleeTerminal
	vm       *goja.Runtime
}

func NewReplUI(app *tview.Application, vm *goja.Runtime) *ReplUI {
	ui := &ReplUI{
		app: app,
		vm:  vm,
	}
	ui.terminal = terminal.NewRepleeTerminal(app, ui.execute)
	return ui
}

func (ui *ReplUI) execute(command string) *terminal.Output {
	response := &terminal.Output{
		IsErr:       false,
		IsSyntaxErr: false,
		Output:      "",
	}

	value, err := ui.vm.RunString(command)
	if err != nil {
		response.IsErr = true
		response.Output = err.Error()
		if exception, ok := err.(*goja.Exception); ok && strings.Index(exception.Value().String(), "Unexpected end of input") != -1 {
			response.IsSyntaxErr = true
		}
	} else {
		if !goja.IsNull(value) && !goja.IsUndefined(value) {
			response.Output = value.String()
		}
	}
	return response
}

func main() {
	ctx := context.Background()
	vm := goja.New()
	vm.SetFieldNameMapper(goja.TagFieldNameMapper("json", true))

	if err := repl.BootstrapRepleeVM(ctx, vm); err != nil {
		panic(err)
	}
	app := tview.NewApplication().SetScreen(terminal.NewScreen())
	ui := NewReplUI(app, vm)

	if err := app.SetRoot(ui.terminal, true).EnableMouse(true).SetFocus(ui.terminal).Run(); err != nil {
		panic(err)
	}
}
