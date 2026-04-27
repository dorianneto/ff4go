package ff4go

import (
	"encoding/json"
	"os"
	"reflect"
	"slices"
)

type Rules struct {
	Users        []string `json:"users"`
	Environments []string `json:"environments"`
	Percentage   float64  `json:"percentage"`
}

type FeatureFlag struct {
	Name        string `json:"name"`
	Enabled     bool   `json:"enabled"`
	Description string `json:"description"`
	Rules       Rules  `json:"rules"`
}

type Manager struct {
	Flags []FeatureFlag `json:"flags"`
}

const flagsFilePath = "ff4go.json"

func NewManagerFromBytes(data []byte) (*Manager, error) {
	var manager *Manager

	if err := json.Unmarshal(data, &manager); err != nil {
		return nil, err
	}

	return manager, nil
}

func NewManagerFromFile() (*Manager, error) {
	data, err := os.ReadFile(flagsFilePath)
	if err != nil {
		return nil, err
	}
	return NewManagerFromBytes(data)
}

func (m *Manager) IsEnabled(name string) bool {
	flag, found := m.getFlag(name)

	if !found {
		return false
	}

	return flag.Enabled
}

func (m *Manager) IsEnabledForUserAndEnvironment(name, user, environment string) bool {
	return m.IsEnabledForUser(name, user) && m.IsEnabledForEnvironment(name, environment)
}

func (m *Manager) IsEnabledForUser(name, user string) bool {
	return m.isEnabledForSomething(name, user, "Users")
}

func (m *Manager) IsEnabledForEnvironment(name, environment string) bool {
	return m.isEnabledForSomething(name, environment, "Environments")
}

func (m *Manager) isEnabledForSomething(name, something, field string) bool {
	flag, found := m.getFlag(name)
	if !found || !flag.Enabled {
		return false
	}

	if m.containsPercentage(flag) {
		return m.calculatePercentage(flag, something)
	}

	fieldValue, ok := reflect.ValueOf(flag.Rules).FieldByName(field).Interface().([]string)
	if !ok {
		return false
	}

	return slices.Contains(fieldValue, something)
}

func (m *Manager) getFlag(name string) (*FeatureFlag, bool) {
	for _, flag := range m.Flags {
		if flag.Name == name {
			return &flag, true
		}
	}

	return nil, false
}

func (m *Manager) containsPercentage(flag *FeatureFlag) bool {
	percentage := flag.Rules.Percentage

	if percentage <= 0 || percentage > 100 {
		return false
	}

	return true
}

func (m *Manager) calculatePercentage(flag *FeatureFlag, id string) bool {
	score := hashStringToFloat(flag.Name, id)
	return score < float64(flag.Rules.Percentage)/float64(100)
}
