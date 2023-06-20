//go:build wasm

package terminal

import (
	"github.com/gdamore/tcell/v2"
	"syscall/js"
)

type regizableScreen struct {
	tcell.Screen
}

func (s *regizableScreen) Init() error {
	js.Global().Set("resizeTerminal", js.FuncOf(s.resizeTerminal))
	return s.Screen.Init()
}

func (s *regizableScreen) resizeTerminal(this js.Value, args []js.Value) interface{} {
	w := args[0].Int()
	h := args[1].Int()
	s.SetSize(w, h)
	return nil
}

func NewScreen() tcell.Screen {
	s, err := tcell.NewTerminfoScreen()
	if err != nil {
		panic(err)
	}
	return &regizableScreen{s}
}
