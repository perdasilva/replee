package resolution

import (
	"fmt"
	"github.com/perdasilva/replee/pkg/deppy"
	"github.com/perdasilva/replee/pkg/deppy/utils"
	"github.com/perdasilva/replee/pkg/deppy/variables"
)

var _ deppy.MutableResolutionProblem = &MutableResolutionProblem{}

type MutableResolutionProblem struct {
	resolutionProblemID deppy.Identifier
	variables           *utils.ActivationMap[deppy.Identifier, deppy.MutableVariable]
}

func NewMutableResolutionProblem(resolutionProblemID deppy.Identifier) *MutableResolutionProblem {
	return &MutableResolutionProblem{
		resolutionProblemID: resolutionProblemID,
		variables:           utils.NewActivationMap[deppy.Identifier, deppy.MutableVariable](),
	}
}

func (m *MutableResolutionProblem) GetVariable(variableID deppy.Identifier, kind string) (deppy.Variable, error) {
	if v, ok := m.variables.GetValue(variableID); ok {
		if v.Kind() != kind {
			return nil, deppy.ConflictErrorf("variable %s is not of kind %s", variableID, kind)
		}
		return v, nil
	} else {
		v := variables.NewMutableVariable(variableID, kind, nil)
		m.variables.Put(variableID, v)
		return v, nil
	}
}

func (m *MutableResolutionProblem) GetMutableVariable(variableID deppy.Identifier, kind string) (deppy.MutableVariable, error) {
	if v, ok := m.variables.GetValue(variableID); ok {
		if v.Kind() != kind {
			return nil, deppy.ConflictErrorf("variable %s is not of kind %s", variableID, kind)
		}
		return v, nil
	} else {
		v := variables.NewMutableVariable(variableID, kind, nil)
		m.variables.Put(variableID, v)
		return v, nil
	}
}

func (m *MutableResolutionProblem) ActivateVariable(v deppy.MutableVariable) error {
	if vr, ok := m.variables.GetValue(v.Identifier()); ok {
		if vr.Kind() != v.Kind() {
			return fmt.Errorf("variable %s is not of kind %s", v.Identifier(), v.Kind())
		}
		if err := vr.Merge(v); err != nil {
			return err
		}
	} else {
		m.variables.Put(v.Identifier(), v)
	}
	m.variables.Activate(v.Identifier())
	return nil
}

func (m *MutableResolutionProblem) DeactivateVariable(variableID deppy.Identifier, kind string) error {
	if v, ok := m.variables.GetValue(variableID); ok {
		if v.Kind() != kind {
			return fmt.Errorf("variable %s is not of kind %s", variableID, kind)
		}
	} else {
		v := variables.NewMutableVariable(variableID, kind, nil)
		m.variables.Put(variableID, v)
	}
	m.variables.Deactivate(variableID)
	return nil
}

func (m *MutableResolutionProblem) ResolutionProblemID() deppy.Identifier {
	return m.resolutionProblemID
}

func (m *MutableResolutionProblem) GetMutableVariables() ([]deppy.MutableVariable, error) {
	return m.variables.Elements(), nil
}

func (m *MutableResolutionProblem) GetVariables() ([]deppy.Variable, error) {
	var vars []deppy.Variable
	for _, v := range m.variables.Elements() {
		vars = append(vars, v)
	}
	return vars, nil
}

func (m *MutableResolutionProblem) Options() []deppy.ResolutionOption {
	return nil
}