# Zero-maintenance enum solution for Go

Due to poor enum support in Go, a pattern has emerged where developers:

- create types with primitive [underlying types](https://go.dev/ref/spec#Underlying_types) (e.g `type Status string`,
  `type LogLevel int`)
- define set of values, e.g `const StatusDraft = Status("open")`
- implement `IsValid() bool / Validate() error` for validation (usually against user input).

The annoying part is maintaining the list of values while keeping validation logic in sync.  
Small thing, but the lack of single source of truth makes us feel dumb, detached, and sad.

## Features

- **Inline definitions**: register enum values as you declare them
- **Zero-maintenance**: automatic `IsValid() bool` and `Validate() error` functions without code generation
- **No interference**: no wrappers or no new type - definitions keep their original type and value
- **Scoped**: definitions are scoped to their type, e.g `"active"` can be valid value for `type UserStatus string`, but
  not
  necessarily for `type OrderStatus string`
- **User friendly error message**: validation error message is human-readable and helpful
- **Lightweight**: no dependencies, auditable (less than 200 lines of code)

## Installation

```bash
go get github.com/0xcafe-io/enum
```

## Usage

```go
package main

import (
	"fmt"

	"github.com/0xcafe-io/enum"
)

type Status string

var (
	StatusDraft  = enum.Def[Status]("draft")
	StatusOpen   = enum.Def[Status]("open")
	StatusMerged = enum.Def[Status]("merged")
	StatusClosed = enum.Def[Status]("closed")
)

func main() {
	userInput := "postponed"
	status := Status(userInput)

	if !enum.IsValid(status) {
		fmt.Println("bad status")
	}

	if err := enum.Validate(status); err != nil {
		fmt.Println(err)
	}

	status = "merged"
	if enum.IsValid(status) {
		fmt.Println("good status")
	}

	if status == StatusMerged {
		fmt.Println("nice job")
	}

	statuses := enum.ValuesOf[Status]()
	fmt.Println(statuses)

	// Output:
	// bad status
	// "postponed" is not a valid choice, allowed values are: "draft", "open", "merged", "closed"
	// good status
	// nice job
	// [draft open merged closed]
}
```

---

Integer types are also supported (`int`, `int8/byte`, `int16`, `int32/rune`, `int64`, and their unsigned siblings):

```go
package main

import (
	"fmt"

	"github.com/0xcafe-io/enum"
)

type Access int

var (
	AccessRead    = enum.Def[Access](1)
	AccessComment = enum.Def[Access](2)
	AccessWrite   = enum.Def[Access](4)
)

func main() {
	var access Access
	if !enum.IsValid(access) {
		fmt.Println("access denied")
	}

	access = 2
	if enum.IsValid(access) {
		fmt.Println("access granted")
	}

	if access == AccessComment {
		fmt.Println("can comment")
	}

	if err := enum.Validate[Access](99); err != nil {
		fmt.Println(err)
	}

	// Output:
	// access denied
	// access granted
	// can comment
	// 99 is not a valid choice, allowed values are: 1, 2, 4
}
```

## Limitations

- `var` instead of `const` for definitions due to function call
