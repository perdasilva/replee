package handler

import (
	"github.com/perdasilva/replee/pkg/deppy"
	"github.com/perdasilva/replee/pkg/replee/cli/action"
	"github.com/perdasilva/replee/pkg/replee/cli/store"
)

var _ Handler = &DeactivateVariableActionHandler{}

type DeactivateVariableActionHandler struct {
}

func (h DeactivateVariableActionHandler) HandleAction(a action.Action, s store.ResolutionProblemStore) error {
	if err := validate(a, action.ActionTypeDeactivateVariable, []string{ParameterProblemID, ParameterVariableID, ParameterVariableKind}); err != nil {
		return err
	}
	problemID := a.Parameters[ParameterProblemID].(string)
	variableID := deppy.Identifierf(a.Parameters[ParameterVariableID].(string))
	kind := a.Parameters[ParameterVariableKind].(string)
	if m, err := s.Get(deppy.Identifier(problemID)); err != nil {
		return err
	} else {
		if err := m.DeactivateVariable(variableID, kind); err != nil {
			return err
		}
		return s.Save(m)
	}
}
