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

`ff4go` is a Go feature flag library. The implementation lives in four files:

- `ff4go.go` — the library itself
- `ff4go_test.go` — tests
- `helper.go` — `hashStringToFloat` (FNV-based deterministic hash) and `isExpired` (RFC 3339 date check)
- `helper_test.go` — helper tests

External dependencies: `github.com/spf13/viper` and `github.com/fsnotify/fsnotify` (used only by `NewManagerFromFileWithWatch`).

**Constructors:**
- `NewManagerFromFile() (*Manager, error)` — reads `ff4go.json` from the process's working directory
- `NewManagerFromBytes(data []byte) (*Manager, error)` — unmarshals JSON directly; useful for embedded configs or remote sources
- `NewManagerFromFileWithWatch() (*Manager, error)` — like `NewManagerFromFile` but uses viper to hot-reload on file changes via `fsnotify`

**Core types:**
- `Manager` — holds a slice of `FeatureFlag`; the receiver for all query methods
- `FeatureFlag` — has a name, enabled bool, description, and a `Rules` struct
- `Rules` — four targeting dimensions: `Users []string`, `Environments []string`, `Percentage float64`, `EndAt string` (RFC 3339)

**Public methods on `Manager`:**
- `IsEnabled(name)` — returns the flag's top-level `Enabled` bool; returns `false` if expired
- `HasFlag(name)` — reports whether a flag with the given name exists
- `IsEnabledForUser(name, user)` — checks the `Users` rule list (flag must also be enabled and not expired)
- `IsEnabledForEnvironment(name, env)` — checks the `Environments` rule list (flag must also be enabled and not expired)
- `IsEnabledForUserAndEnvironment(name, user, env)` — both user and environment checks must pass

`IsEnabledForUser` and `IsEnabledForEnvironment` route through `isEnabledForSomething`, which uses reflection to select the `Users` or `Environments` field from `Rules`. If `Percentage` is set (0 < p ≤ 100), it takes precedence over list-based checks and returns a **deterministic** result using an FNV-64a hash of `"flagName:id"` — the same user always gets the same result. `endAt` (RFC 3339) disables the flag for all methods once the timestamp is in the past.

**`ff4go.json` schema:**
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
        "percentage": 50.0,
        "endAt": "2026-12-31T23:59:59Z"
      }
    }
  ]
}
```
