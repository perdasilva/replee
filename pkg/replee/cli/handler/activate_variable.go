package handler

import (
	"github.com/perdasilva/replee/pkg/deppy"
	"github.com/perdasilva/replee/pkg/deppy/variables"
	"github.com/perdasilva/replee/pkg/replee/cli/action"
	"github.com/perdasilva/replee/pkg/replee/cli/store"
)

var _ Handler = &ActivateVariableActionHandler{}

type ActivateVariableActionHandler struct {
}

func (h ActivateVariableActionHandler) HandleAction(a action.Action, s store.ResolutionProblemStore) error {
	if err := validate(a, action.ActionTypeActivateVariable, []string{ParameterProblemID, ParameterVariableID, ParameterVariableKind}); err != nil {
		return err
	}
	problemID := a.Parameters[ParameterProblemID].(string)
	variableID := deppy.Identifierf(a.Parameters[ParameterVariableID].(string))
	kind := a.Parameters[ParameterVariableKind].(string)
	var properties map[string]interface{}
	if _, ok := a.Parameters[ParameterProperties]; ok {
		properties = a.Parameters[ParameterProperties].(map[string]interface{})
	}
	if m, err := s.Get(deppy.Identifier(problemID)); err != nil {
		return err
	} else {
		if err := m.ActivateVariable(variables.NewMutableVariable(variableID, kind, properties)); err != nil {
			return err
		}
		return s.Save(m)
	}
}
