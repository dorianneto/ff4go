# Changelog

## [0.3.0] — 2026-05-04

### Features

- **`NewManagerFromFileWithWatch`** — new constructor that reads `ff4go.json` and automatically reloads flags when the file changes, without restarting the process. The returned `*Manager` pointer is updated in place on each change.
- **`NewManagerFromBytes`** — new constructor that accepts raw JSON bytes, enabling initialization from embedded files (`//go:embed`), environment variables, or remote sources.
- **`IsEnabledForUserAndEnvironment`** — new method that returns `true` only when both the user and environment targeting checks pass simultaneously.
- **`HasFlag`** — new method that reports whether a flag with the given name exists in the manager.
- **`endAt` rule** — new `Rules` field (RFC 3339 string) that disables a flag after a given timestamp, regardless of other rules. Applies to `IsEnabled`, `IsEnabledForUser`, `IsEnabledForEnvironment`, and `IsEnabledForUserAndEnvironment`.
- **Deterministic percentage rollout** — `Percentage` targeting now uses an FNV-64a hash of `"flagName:id"` instead of a random value, ensuring the same user always gets the same result for a given flag.

### Bug Fixes

- **`IsEnabled` now respects `endAt`** — previously `IsEnabled` did not check the `endAt` rule and could return `true` for expired flags.
