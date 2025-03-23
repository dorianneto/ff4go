package ff4go

import "testing"

func TestWhenFeatureFlagIsEnabled(t *testing.T) {
	m, err := NewManager([]byte(`{"flags":[{"name":"new-ui","enabled":true}]}`))
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
	m, err := NewManager([]byte(`{"flags":[{"name":"new-ui","enabled":true,"rules":{"users":["user1"]}}]}`))
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
	m, err := NewManager([]byte(`{"flags":[{"name":"new-ui","enabled":true,"rules":{"environments":["development"]}}]}`))
	if err != nil {
		t.Errorf("Error on initializing manager")
	}

	ff := m.IsEnabledForEnvironment("new-ui", "development")
	want := true

	if want != ff {
		t.Errorf("Expected %v but got %v", want, ff)
	}
}
