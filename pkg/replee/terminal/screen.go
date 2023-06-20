//go:build !wasm

package terminal

import "github.com/gdamore/tcell/v2"

func NewScreen() tcell.Screen {
	s, err := tcell.NewScreen()
	if err != nil {
		panic(err)
	}
	return s
}
