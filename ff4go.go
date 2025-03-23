package ff4go

import (
	"encoding/json"
	"slices"
)

type Rules struct {
	Users        []string `json:"users"`
	Environments []string `json:"environments"`
	// Percentage   float64  `json:"percentage"`
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

func NewManager(data []byte) (*Manager, error) {
	var Manager *Manager

	if err := json.Unmarshal(data, &Manager); err != nil {
		return nil, err
	}

	return Manager, nil
}

func (m *Manager) getFlag(name string) (*FeatureFlag, bool) {
	for _, flag := range m.Flags {
		if flag.Name == name {
			return &flag, true
		}
	}

	return nil, false
}

func (m *Manager) IsEnabled(name string) bool {
	flag, found := m.getFlag(name)

	if !found {
		return false
	}

	return flag.Enabled
}

func (m *Manager) IsEnabledForUser(name, user string) bool {
	flag, found := m.getFlag(name)

	if !found || !flag.Enabled {
		return false
	}

	return slices.Contains(flag.Rules.Users, user)
}

func (m *Manager) IsEnabledForEnvironment(name, environment string) bool {
	flag, found := m.getFlag(name)

	if !found || !flag.Enabled {
		return false
	}

	return slices.Contains(flag.Rules.Environments, environment)
}
