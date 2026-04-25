# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Commands

```bash
# Build
go build -v ./...

# Run all tests
go test -v ./...

# Run a single test
go test -v -run TestWhenFeatureFlagIsEnabled ./...
```

## Architecture

`ff4go` is a zero-dependency Go feature flag library (single package, no external imports). The entire implementation lives in two files:

- `ff4go.go` — the library itself
- `ff4go_test.go` — tests

**Core types:**
- `Manager` — holds a slice of `FeatureFlag`; created via `NewManager([]byte)` which unmarshals a JSON payload
- `FeatureFlag` — has a name, enabled bool, description, and a `Rules` struct
- `Rules` — three targeting dimensions: `Users []string`, `Environments []string`, `Percentage float64`

**Check methods on `Manager`:**
- `IsEnabled(name)` — returns the flag's top-level `Enabled` bool
- `IsEnabledForUser(name, user)` — checks the `Users` rule list (flag must also be enabled)
- `IsEnabledForEnvironment(name, env)` — checks the `Environments` rule list (flag must also be enabled)

Both targeting methods route through `isEnabledForSomething`, which uses reflection to select the `Users` or `Environments` field from `Rules`. If `Percentage` is set (0 < p ≤ 100), it takes precedence over the list-based check and returns a probabilistic result via `rand.Float64()`.

**JSON schema expected by `NewManager`:**
```json
{
  "flags": [
    {
      "name": "flag-name",
      "enabled": true,
      "description": "optional",
      "rules": {
        "users": ["user1"],
        "environments": ["production"],
        "percentage": 50.0
      }
    }
  ]
}
```
