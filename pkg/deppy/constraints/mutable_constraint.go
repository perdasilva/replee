package constraints

import (
	"encoding/json"
	"github.com/perdasilva/replee/pkg/deppy"
	"sync"
)

var _ deppy.MutableConstraint = &MutableConstraint{}

type MutableConstraint struct {
	deppy.Constraint
	constraintID deppy.Identifier
	kind         string
	properties   map[string]interface{}
	lock         sync.RWMutex
}

func (m *MutableConstraint) ConstraintID() deppy.Identifier {
	return m.constraintID
}

func (m *MutableConstraint) Merge(other deppy.Constraint) error {
	for key, value := range other.GetProperties() {
		if err := m.setProperty(key, value); err != nil {
			return err
		}
	}
	return nil
}

func NewMutableConstraint(constraint deppy.Constraint) *MutableConstraint {
	return &MutableConstraint{
		Constraint: constraint,
		properties: map[string]interface{}{},
		lock:       sync.RWMutex{},
	}
}

func (m *MutableConstraint) GetProperty(key string) (interface{}, bool) {
	m.lock.RLock()
	defer m.lock.RUnlock()
	prop, ok := m.properties[key]
	return prop, ok
}

func (m *MutableConstraint) GetProperties() map[string]interface{} {
	m.lock.RLock()
	defer m.lock.RUnlock()
	props := map[string]interface{}{}
	for k, v := range m.properties {
		props[k] = v
	}
	return props
}

func (m *MutableConstraint) MarshalJSON() ([]byte, error) {
	return json.Marshal(&struct {
		Kind       string                 `json:"kind"`
		Properties map[string]interface{} `json:"properties"`
	}{
		Kind:       m.Kind(),
		Properties: m.GetProperties(),
	})
}

func (m *MutableConstraint) SetProperty(key string, value interface{}) error {
	m.lock.Lock()
	defer m.lock.Unlock()
	return m.setProperty(key, value)
}

func (m *MutableConstraint) setProperty(key string, value interface{}) error {
	if v, ok := m.properties[key]; ok {
		if v == value {
			return nil
		}
		return deppy.ConflictErrorf("merge conflict: property %s already set to %v", key, v)
	}
	m.properties[key] = value
	return nil
}
