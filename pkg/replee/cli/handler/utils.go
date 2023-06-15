package handler

import (
	"fmt"
	"github.com/perdasilva/replee/pkg/replee/cli/action"
)

func validate(a action.Action, expectedActionType string, expectedProperties []string) error {
	if a.ActionType != expectedActionType {
		return fmt.Errorf("action type %s is not supported by this handler", a.ActionType)
	}
	for _, property := range expectedProperties {
		if _, ok := a.Parameters[property]; !ok {
			return fmt.Errorf("action %s is missing property %s", a.ActionType, property)
		}
	}
	return nil
}
