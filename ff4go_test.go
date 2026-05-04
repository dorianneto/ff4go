package ff4go

import (
	"os"
	"testing"
	"time"
)

func TestNewManagerFromFileWithWatch(t *testing.T) {
	m, err := NewManagerFromFileWithWatch()
	if err != nil {
		t.Fatalf("Error initializing manager from file with watch: %v", err)
	}
	if !m.IsEnabled("new-ui") {
		t.Errorf("Expected new-ui to be enabled")
	}
}

func TestNewManagerFromFileWithWatch_ReloadsWhenFileChanges(t *testing.T) {
	originalContent, err := os.ReadFile(flagsFilePath)
	if err != nil {
		t.Fatalf("Failed to read original config: %v", err)
	}
	t.Cleanup(func() {
		os.WriteFile(flagsFilePath, originalContent, 0644)
	})

	initialConfig := `{"flags":[{"name":"watch-flag","enabled":true}]}`
	if err := os.WriteFile(flagsFilePath, []byte(initialConfig), 0644); err != nil {
		t.Fatalf("Failed to write initial config: %v", err)
	}

	m, err := NewManagerFromFileWithWatch()
	if err != nil {
		t.Fatalf("Error initializing manager: %v", err)
	}
	if !m.IsEnabled("watch-flag") {
		t.Errorf("Expected watch-flag to be enabled initially")
	}

	updatedConfig := `{"flags":[{"name":"watch-flag","enabled":false}]}`
	if err := os.WriteFile(flagsFilePath, []byte(updatedConfig), 0644); err != nil {
		t.Fatalf("Failed to write updated config: %v", err)
	}

	deadline := time.Now().Add(2 * time.Second)
	for time.Now().Before(deadline) {
		if !m.IsEnabled("watch-flag") {
			break
		}
		time.Sleep(50 * time.Millisecond)
	}

	if m.IsEnabled("watch-flag") {
		t.Errorf("Expected watch-flag to be disabled after file change")
	}
}

func TestNewManagerFromFile(t *testing.T) {
	m, err := NewManagerFromFile()
	if err != nil {
		t.Errorf("Error initializing manager from file: %v", err)
	}
	if !m.IsEnabled("new-ui") {
		t.Errorf("Expected new-ui to be enabled")
	}
}

func TestWhenFeatureFlagIsEnabled(t *testing.T) {
	m, err := NewManagerFromBytes([]byte(`{"flags":[{"name":"new-ui","enabled":true}]}`))
	if err != nil {
		t.Errorf("Error on initializing manager")
	}

	ff := m.IsEnabled("new-ui")
	want := true

	if want != ff {
		t.Errorf("Expected %v but got %v", want, ff)
	}
}

func TestWhenFeatureFlagIsEnabledForAnUser(t *testing.T) {
	m, err := NewManagerFromBytes([]byte(`{"flags":[{"name":"new-ui","enabled":true,"rules":{"users":["user1"]}}]}`))
	if err != nil {
		t.Errorf("Error on initializing manager")
	}

	ff := m.IsEnabledForUser("new-ui", "user1")
	want := true

	if want != ff {
		t.Errorf("Expected %v but got %v", want, ff)
	}
}

func TestWhenFeatureFlagIsEnabledForAnEnvironment(t *testing.T) {
	m, err := NewManagerFromBytes([]byte(`{"flags":[{"name":"new-ui","enabled":true,"rules":{"environments":["development"]}}]}`))
	if err != nil {
		t.Errorf("Error on initializing manager")
	}

	ff := m.IsEnabledForEnvironment("new-ui", "development")
	want := true

	if want != ff {
		t.Errorf("Expected %v but got %v", want, ff)
	}
}

func TestWhenFeatureFlagIsEnabledForAnUserWithPercentage(t *testing.T) {
	m, err := NewManagerFromBytes([]byte(`{"flags":[{"name":"new-ui","enabled":true,"rules":{"percentage": 50,"users":["user1"]}}]}`))
	if err != nil {
		t.Errorf("Error on initializing manager")
	}

	ff := m.IsEnabledForUser("new-ui", "user1")
	want := false

	if want != ff {
		t.Errorf("Expected %v but got %v", want, ff)
	}
}

func TestWhenFeatureFlagIsEnabledForAnUserAndEnvironment(t *testing.T) {
	m, err := NewManagerFromBytes([]byte(`{"flags":[{"name":"new-ui","enabled":true,"rules":{"users":["user1"],"environments":["development"]}}]}`))
	if err != nil {
		t.Errorf("Error on initializing manager")
	}

	ff := m.IsEnabledForUserAndEnvironment("new-ui", "user1", "development")
	want := true

	if want != ff {
		t.Errorf("Expected %v but got %v", want, ff)
	}
}

func TestWhenItHasFeatureFlag(t *testing.T) {
	m, err := NewManagerFromBytes([]byte(`{"flags":[{"name":"new-ui","enabled":true,"rules":{"users":["user1"],"environments":["development"]}}]}`))
	if err != nil {
		t.Errorf("Error on initializing manager")
	}

	ff := m.HasFlag("new-ui")
	want := true

	if want != ff {
		t.Errorf("Expected %v but got %v", want, ff)
	}
}

func TestWhenFeatureFlagHasAnExpireDate(t *testing.T) {
	tests := []struct {
		date     string
		expected bool
	}{
		{date: "2023-01-01T00:00:00Z", expected: false},
		{date: time.Now().Add(24 * time.Hour).Format(time.RFC3339), expected: true},
	}

	for _, tt := range tests {
		t.Run(tt.date, func(t *testing.T) {
			m, err := NewManagerFromBytes([]byte(`{"flags":[{"name":"new-ui","enabled":true,"rules":{"users":["user1"],"endAt":"` + tt.date + `","environments":["development"]}}]}`))
			if err != nil {
				t.Errorf("Error on initializing manager")
			}

			ff := m.IsEnabled("new-ui")

			if tt.expected != ff {
				t.Errorf("Expected %v but got %v", tt.expected, ff)
			}
		})
	}
}

func TestWhenFeatureFlagHasAnExpireDateInMethodsAddressingRules(t *testing.T) {
	tests := []struct {
		date     string
		expected bool
	}{
		{date: "2023-01-01T00:00:00Z", expected: false},
		{date: time.Now().Add(24 * time.Hour).Format(time.RFC3339), expected: true},
	}

	for _, tt := range tests {
		t.Run(tt.date, func(t *testing.T) {
			m, err := NewManagerFromBytes([]byte(`{"flags":[{"name":"new-ui","enabled":true,"rules":{"users":["user1"],"endAt":"` + tt.date + `","environments":["development"]}}]}`))
			if err != nil {
				t.Errorf("Error on initializing manager")
			}

			ff := m.IsEnabledForUser("new-ui", "user1")

			if tt.expected != ff {
				t.Errorf("Expected %v but got %v", tt.expected, ff)
			}

			ff = m.IsEnabledForEnvironment("new-ui", "development")

			if tt.expected != ff {
				t.Errorf("Expected %v but got %v", tt.expected, ff)
			}
		})
	}
}
