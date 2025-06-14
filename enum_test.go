package enum_test

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

func Example() {
	input1 := Status("postponed")
	if !enum.IsValid(input1) {
		fmt.Println("bad input")
	}

	input2 := Status("merged")
	if enum.IsValid(input2) {
		fmt.Println("good input")
	}

	if input2 == StatusMerged {
		fmt.Println("all good")
	}

	values := enum.ValuesOf[Status]()
	fmt.Println(values)

	// Output:
	// bad input
	// good input
	// all good
	// [draft open merged closed]
}
