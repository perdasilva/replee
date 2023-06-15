package handler

import (
	"github.com/perdasilva/replee/pkg/replee/cli/action"
	"github.com/perdasilva/replee/pkg/replee/cli/store"
)

type Handler interface {
	HandleAction(action action.Action, store store.ResolutionProblemStore) error
}
