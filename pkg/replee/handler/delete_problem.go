package handler

import (
	"github.com/perdasilva/replee/pkg/deppy"
	"github.com/perdasilva/replee/pkg/replee/action"
	"github.com/perdasilva/replee/pkg/replee/store"
)

var _ Handler = &DeleteProblemActionHandler{}

type DeleteProblemActionHandler struct{}

func (d DeleteProblemActionHandler) HandleAction(a action.Action, s store.ResolutionProblemStore) error {
	if err := validate(a, action.ActionTypeCreateProblem, []string{ParameterProblemID}); err != nil {
		return err
	}
	return s.Delete(deppy.Identifier(a.Parameters[ParameterProblemID].(string)))
}
