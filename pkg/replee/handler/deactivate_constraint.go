package handler

import (
	"github.com/perdasilva/replee/pkg/deppy"
	"github.com/perdasilva/replee/pkg/deppy/constraints"
	"github.com/perdasilva/replee/pkg/replee/action"
	"github.com/perdasilva/replee/pkg/replee/store"
)

var _ Handler = &DeactivateConstraintActionHandler{}

type DeactivateConstraintActionHandler struct{}

func (h DeactivateConstraintActionHandler) HandleAction(a action.Action, s store.ResolutionProblemStore) error {
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
		v, err := m.GetMutableVariable(variableID, variableKind)
		if err != nil {
			return err
		}
		switch constraintKind {
		case constraints.ConstraintKindMandatory:
			if err := v.RemoveMandatory(constraintID); err != nil {
				return err
			}
		case constraints.ConstraintKindProhibited:
			if err := v.RemoveProhibited(constraintID); err != nil {
				return err
			}
		case constraints.ConstraintKindConflict:
			if err := v.RemoveConflict(constraintID); err != nil {
				return err
			}
		case constraints.ConstraintKindDependency:
			var depIDs []deppy.Identifier
			deps, ok := a.Parameters[ParameterDependencies]
			if ok {
				for _, depID := range deps.([]interface{}) {
					depIDs = append(depIDs, deppy.Identifierf(depID.(string)))
				}
			}
			if err := v.RemoveDependency(constraintID, depIDs...); err != nil {
				return err
			}
		case constraints.ConstraintKindAtMost:
			var varIDs []deppy.Identifier
			vars, ok := a.Parameters[ParameterAtMostVariables]
			if ok {
				for _, varID := range vars.([]interface{}) {
					varIDs = append(varIDs, deppy.Identifierf(varID.(string)))
				}
			}
			if err := v.RemoveAtMost(constraintID, varIDs...); err != nil {
				return err
			}
		default:
			return deppy.Fatalf("unknown constraint kind: %s", constraintKind)
		}
		return s.Save(m)
	}
}
