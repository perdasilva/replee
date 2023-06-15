package handler

import (
	"github.com/perdasilva/replee/pkg/replee/cli/action"
	"github.com/perdasilva/replee/pkg/replee/cli/store"
)

var _ Handler = &RepleeHandler{}

type RepleeHandler struct {
	handlers []Handler
}

func NewRepleeHandler() *RepleeHandler {
	return &RepleeHandler{
		handlers: []Handler{
			&FilteredHandler{
				Handler: &DeleteProblemActionHandler{},
				Filter: func(a action.Action) bool {
					return a.ActionType == action.ActionTypeDeleteProblem
				},
			},
			&FilteredHandler{
				Handler: &CreateProblemActionHandler{},
				Filter: func(a action.Action) bool {
					return a.ActionType == action.ActionTypeCreateProblem
				},
			},
			&FilteredHandler{
				Handler: &ActivateVariableActionHandler{},
				Filter: func(a action.Action) bool {
					return a.ActionType == action.ActionTypeActivateVariable
				},
			},
			&FilteredHandler{
				Handler: &DeactivateVariableActionHandler{},
				Filter: func(a action.Action) bool {
					return a.ActionType == action.ActionTypeDeactivateVariable
				},
			},
			&FilteredHandler{
				Handler: &MergeConstraintActionHandler{},
				Filter: func(a action.Action) bool {
					return a.ActionType == action.ActionTypeMergeConstraint
				},
			},
			&FilteredHandler{
				Handler: &SolveActionHandler{},
				Filter: func(a action.Action) bool {
					return a.ActionType == action.ActionTypeSolve
				},
			},
		},
	}
}

func (r RepleeHandler) HandleAction(a action.Action, s store.ResolutionProblemStore) error {
	for _, h := range r.handlers {
		if err := h.HandleAction(a, s); err != nil {
			return err
		}
	}
	return nil
}
