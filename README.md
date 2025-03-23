[![Go Report Card](https://goreportcard.com/badge/github.com/dorianneto/ff4go)](https://goreportcard.com/report/github.com/dorianneto/ff4go)
![GitHub go.mod Go version](https://img.shields.io/github/go-mod/go-version/dorianneto/ff4go)
![GitHub License](https://img.shields.io/github/license/dorianneto/ff4go)

> Thank you for using ff4go! Feel free to report any [issue or improvement](https://github.com/dorianneto/ff4go/issues) 🙏

## Documentation

https://pkg.go.dev/github.com/dorianneto/ff4go#section-documentation


## Demo

First, you need to install the dependency
```
go get github.com/dorianneto/ff4go
```

Once the dependency is installed, you can use it like this:
```
package main

import (
	"fmt"

	"github.com/dorianneto/ff4go"
)

func main() {
	m, err := ff4go.NewManager([]byte(`{"flags":[{"name":"new-ui","enabled":true,"rules":{"users":["user1"],"environments":["development"]}}]}`))
	if err != nil {
		panic(err)
	}

	fmt.Println(m.IsEnabled("new-ui"))                          // true
	fmt.Println(m.IsEnabledForUser("new-ui", "user1"))          // true
	fmt.Println(m.IsEnabledForEnvironment("new-ui", "staging")) // false
}
```

## SDK

| Language | SDK |
| -------- | ------- |
| Go | [ff4go](https://github.com/dorianneto/ff4go) |
| Typescript | coming soon |


## License

[MIT](https://github.com/dorianneto/ff4go/blob/main/LICENSE)
