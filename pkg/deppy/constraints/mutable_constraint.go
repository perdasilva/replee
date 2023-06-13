package constraints

import (
	"encoding/json"
	"github.com/perdasilva/replee/pkg/deppy"
	"sync"
)

type mutableConstraint struct {
	constraintID deppy.Identifier
	kind         string
	properties   map[string]interface{}
	lock         sync.RWMutex
}

func (m *mutableConstraint) ConstraintID() deppy.Identifier {
	return m.constraintID
}

func (m *mutableConstraint) Kind() string {
	return m.kind
}

func (m *mutableConstraint) Merge(other deppy.Constraint) error {
	for key, value := range other.GetProperties() {
		if err := m.setProperty(key, value); err != nil {
			return err
		}
	}
	return nil
}

func NewMutableConstraint(constraint deppy.Constraint) *mutableConstraint {
	return &mutableConstraint{
		properties: map[string]interface{}{},
		lock:       sync.RWMutex{},
	}
}

func (m *mutableConstraint) GetProperty(key string) (interface{}, bool) {
	m.lock.RLock()
	defer m.lock.RUnlock()
	prop, ok := m.properties[key]
	return prop, ok
}

func (m *mutableConstraint) GetProperties() map[string]interface{} {
	m.lock.RLock()
	defer m.lock.RUnlock()
	props := map[string]interface{}{}
	for k, v := range m.properties {
		props[k] = v
	}
	return props
}

func (m *mutableConstraint) MarshalJSON() ([]byte, error) {
	return json.Marshal(&struct {
		ConstraintID deppy.Identifier       `json:"constraintID"`
		Kind         string                 `json:"kind"`
		Properties   map[string]interface{} `json:"properties"`
	}{
		ConstraintID: m.ConstraintID(),
		Kind:         m.Kind(),
		Properties:   m.GetProperties(),
	})
}

func (m *mutableConstraint) UnmarshalJSON(jsonBytes []byte) error {
	data := &struct {
		ConstraintID deppy.Identifier       `json:"constraintID"`
		Kind         string                 `json:"kind"`
		Properties   map[string]interface{} `json:"properties"`
	}{}
	if err := json.Unmarshal(jsonBytes, data); err != nil {
		return err
	}
	m.constraintID = data.ConstraintID
	m.kind = data.Kind
	m.properties = data.Properties
	return nil
}

func (m *mutableConstraint) SetProperty(key string, value interface{}) error {
	m.lock.Lock()
	defer m.lock.Unlock()
	return m.setProperty(key, value)
}

func (m *mutableConstraint) setProperty(key string, value interface{}) error {
	if v, ok := m.properties[key]; ok {
		if v == value {
			return nil
		}
		return deppy.ConflictErrorf("merge conflict: property %s already set to %v", key, v)
	}
	m.properties[key] = value
	return nil
}
