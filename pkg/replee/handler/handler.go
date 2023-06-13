package handler

import (
	"github.com/perdasilva/replee/pkg/replee/action"
	"github.com/perdasilva/replee/pkg/replee/store"
)

type Handler interface {
	HandleAction(action action.Action, store store.ResolutionProblemStore) error
}
