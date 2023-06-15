package handler

import (
	"github.com/perdasilva/replee/pkg/deppy"
	"github.com/perdasilva/replee/pkg/deppy/constraints"
	"github.com/perdasilva/replee/pkg/deppy/variables"
	"github.com/perdasilva/replee/pkg/replee/cli/action"
	"github.com/perdasilva/replee/pkg/replee/cli/store"
)

var _ Handler = &MergeConstraintActionHandler{}

type MergeConstraintActionHandler struct {
}

func (h MergeConstraintActionHandler) HandleAction(a action.Action, s store.ResolutionProblemStore) error {
	if err := validate(a, action.ActionTypeMergeConstraint, []string{
		ParameterProblemID,
		ParameterVariableID,
		ParameterVariableKind,
		ParameterConstraintID,
		ParameterConstraintKind,
	}); err != nil {
		return err
	}
	problemID := a.Parameters[ParameterProblemID].(string)
	variableID := deppy.Identifierf(a.Parameters[ParameterVariableID].(string))
	variableKind := a.Parameters[ParameterVariableKind].(string)
	constraintID := deppy.Identifierf(a.Parameters[ParameterConstraintID].(string))
	constraintKind := a.Parameters[ParameterConstraintKind].(string)

	if m, err := s.Get(deppy.Identifier(problemID)); err != nil {
		return err
	} else {
		params, ok := a.Parameters[ParameterProperties]
		if !ok {
			params = map[string]interface{}{}
		}
		v := variables.NewMutableVariable(variableID, variableKind, params.(map[string]interface{}))
		switch constraintKind {
		case constraints.ConstraintKindMandatory:
			if err := v.AddMandatory(constraintID); err != nil {
				return err
			}
		case constraints.ConstraintKindProhibited:
			if err := v.AddProhibited(constraintID); err != nil {
				return err
			}
		case constraints.ConstraintKindConflict:
			conflictingVariableID, ok := a.Parameters[ParameterConflictingVariableID]
			if !ok {
				return deppy.Fatalf("missing parameter %s", ParameterConflictingVariableID)
			}
			if err := v.AddConflict(constraintID, deppy.Identifierf(conflictingVariableID.(string))); err != nil {
				return err
			}
		case constraints.ConstraintKindDependency:
			dependencies, ok := a.Parameters[ParameterDependencies]
			if !ok {
				return deppy.Fatalf("missing parameter %s", ParameterDependencies)
			}
			var depIDs []deppy.Identifier
			for _, dep := range dependencies.([]interface{}) {
				depIDs = append(depIDs, deppy.Identifierf(dep.(string)))
			}
			if err := v.AddDependency(constraintID, depIDs...); err != nil {
				return err
			}
		case constraints.ConstraintKindAtMost:
			vars, ok := a.Parameters[ParameterAtMostVariables]
			if !ok {
				vars = []interface{}{}
			}
			n, ok := a.Parameters[ParameterAtMostN]
			if !ok {
				n = -1
			}
			var varIDs []deppy.Identifier
			for _, dep := range vars.([]interface{}) {
				varIDs = append(varIDs, deppy.Identifierf(dep.(string)))
			}
			if err := v.AddAtMost(constraintID, n.(int), varIDs...); err != nil {
				return err
			}
		}
		if err := m.ActivateVariable(v); err != nil {
			return err
		}
		return s.Save(m)
	}
}
