package handler

import (
	"github.com/perdasilva/replee/pkg/replee/action"
	"github.com/perdasilva/replee/pkg/replee/store"
)

type FilteredHandler struct {
	Handler Handler
	Filter  func(action.Action) bool
}

func (f FilteredHandler) HandleAction(a action.Action, s store.ResolutionProblemStore) error {
	if f.Filter(a) {
		return f.Handler.HandleAction(a, s)
	}
	return nil
}
