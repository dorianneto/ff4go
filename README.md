![GitHub Actions Workflow Status](https://img.shields.io/github/actions/workflow/status/dorianneto/ff4go/go.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/dorianneto/ff4go)](https://goreportcard.com/report/github.com/dorianneto/ff4go)
![GitHub go.mod Go version](https://img.shields.io/github/go-mod/go-version/dorianneto/ff4go)
![GitHub License](https://img.shields.io/github/license/dorianneto/ff4go)

> Thank you for using ff4go! Feel free to report any [issue or improvement](https://github.com/dorianneto/ff4go/issues) 🙏


## Installation

```
go get github.com/dorianneto/ff4go
```


## Configuration

Create an `ff4go.json` file at your project root. Each flag supports the following fields:

```json
{
  "flags": [
    {
      "name": "new-ui",
      "enabled": true,
      "description": "Enables the redesigned UI",
      "rules": {
        "users": ["user1", "user2"],
        "environments": ["development", "staging"],
        "percentage": 25.0,
        "endAt": "2026-12-31T23:59:59Z"
      }
    }
  ]
}
```

### Rules

| Field          | Type       | Description |
|----------------|------------|-------------|
| `users`        | `[]string` | List of user identifiers that the flag is active for. |
| `environments` | `[]string` | List of environment names that the flag is active for. |
| `percentage`   | `float64`  | Percentage of users (0–100) to roll out to. When set, takes precedence over `users`/`environments` lists and uses a deterministic hash so the same user always gets the same result. |
| `endAt`        | `string`   | RFC 3339 timestamp after which the flag is treated as disabled, regardless of other rules. Omit to disable expiry. |

> `users` and `environments` are only evaluated when `percentage` is **not** set.


## Initialization

### From a file

Reads `ff4go.json` from the process working directory.

```go
m, err := ff4go.NewManagerFromFile()
if err != nil {
    panic(err)
}
```

### From bytes

Useful when you load the configuration from a remote source, an environment variable, or embed it with `//go:embed`.

```go
data := []byte(`{"flags":[{"name":"new-ui","enabled":true}]}`)

m, err := ff4go.NewManagerFromBytes(data)
if err != nil {
    panic(err)
}
```


### From a file with watch

Reads `ff4go.json` from the process working directory and automatically reloads flags when the file changes. Useful for long-running services that need to pick up flag updates without restarting.

```go
m, err := ff4go.NewManagerFromFileWithWatch()
if err != nil {
    panic(err)
}
```

The returned `*Manager` pointer is updated in place when the file changes — no need to re-initialize.


## API

```go
// IsEnabled returns the flag's top-level enabled value.
m.IsEnabled("new-ui")

// HasFlag reports whether a flag with the given name exists.
m.HasFlag("new-ui")

// IsEnabledForUser returns true when the flag is enabled and the user
// matches the users rule (or falls within the percentage rollout).
m.IsEnabledForUser("new-ui", "user1")

// IsEnabledForEnvironment returns true when the flag is enabled and the
// environment matches the environments rule (or falls within the percentage rollout).
m.IsEnabledForEnvironment("new-ui", "development")

// IsEnabledForUserAndEnvironment returns true only when both the user
// and environment checks pass simultaneously.
m.IsEnabledForUserAndEnvironment("new-ui", "user1", "development")
```

All targeting methods return `false` when the flag does not exist, is disabled, or has expired (`endAt` is in the past).


## Examples

### Basic flag check

```go
package main

import (
    "fmt"

    "github.com/dorianneto/ff4go"
)

func main() {
    m, err := ff4go.NewManagerFromFile()
    if err != nil {
        panic(err)
    }

    fmt.Println(m.IsEnabled("new-ui"))                                        // true
    fmt.Println(m.IsEnabledForUser("new-ui", "user1"))                        // true
    fmt.Println(m.IsEnabledForEnvironment("new-ui", "staging"))               // false
    fmt.Println(m.IsEnabledForUserAndEnvironment("new-ui", "user1", "development")) // true
}
```

### Loading from bytes (embedded config)

```go
package main

import (
    _ "embed"
    "fmt"

    "github.com/dorianneto/ff4go"
)

//go:embed ff4go.json
var flagsConfig []byte

func main() {
    m, err := ff4go.NewManagerFromBytes(flagsConfig)
    if err != nil {
        panic(err)
    }

    fmt.Println(m.IsEnabled("new-ui")) // true
}
```

### Percentage rollout

```json
{
  "flags": [
    {
      "name": "new-checkout",
      "enabled": true,
      "rules": { "percentage": 10.0 }
    }
  ]
}
```

```go
// The result is deterministic per (flag name, user id) pair.
fmt.Println(m.IsEnabledForUser("new-checkout", "user42")) // stable true/false
```

### Scheduled flag expiry

```json
{
  "flags": [
    {
      "name": "beta-banner",
      "enabled": true,
      "rules": {
        "users": ["tester1"],
        "endAt": "2026-06-01T00:00:00Z"
      }
    }
  ]
}
```

After `2026-06-01T00:00:00Z` all targeting methods return `false` for `beta-banner`.
