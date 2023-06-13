package action

const (
	ActionTypeActivateVariable     = "activate_variable"
	ActionTypeDeactivateVariable   = "deactivate_variable"
	ActionTypeMergeConstraint      = "merge_constraint"
	ActionTypeDeactivateConstraint = "deactivate_constraint"
	ActionTypeCreateProblem        = "create_problem"
	ActionTypeDeleteProblem        = "delete_problem"
	ActionTypeAppendVariableSource = "append_variable_source"
	ActionTypeRemoveVariableSource = "remove_variable_source"
	ActionTypeClean                = "clean"
	ActionTypeSolve                = "solve"
)

type Action struct {
	ActionType string                 `json:"actionType"`
	Parameters map[string]interface{} `json:"parameters"`
}
