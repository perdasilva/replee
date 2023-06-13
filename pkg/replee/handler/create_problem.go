package handler

import (
	"github.com/perdasilva/replee/pkg/deppy"
	"github.com/perdasilva/replee/pkg/replee/action"
	"github.com/perdasilva/replee/pkg/replee/store"
)

var _ Handler = &CreateProblemActionHandler{}

type CreateProblemActionHandler struct{}

func (c CreateProblemActionHandler) HandleAction(a action.Action, s store.ResolutionProblemStore) error {
	if err := validate(a, action.ActionTypeCreateProblem, []string{ParameterProblemID}); err != nil {
		return err
	}
	return s.New(deppy.Identifier(a.Parameters[ParameterProblemID].(string)))
}
