package enum_test

import (
	"fmt"
	"testing"

	"github.com/0xcafe-io/enum"
)

type Status string
type Access int

var (
	StatusDraft  = enum.Def[Status]("draft")
	StatusOpen   = enum.Def[Status]("open")
	StatusMerged = enum.Def[Status]("merged")
	StatusClosed = enum.Def[Status]("closed")
)

var (
	AccessRead    = enum.Def[Access](1)
	AccessComment = enum.Def[Access](2)
	AccessWrite   = enum.Def[Access](4)
)

func Example_string() {
	var input Status = "postponed"
	if !enum.IsValid(input) {
		fmt.Println("bad input")
	}

	if err := enum.Validate(input); err != nil {
		fmt.Println(err)
	}

	input = "merged"
	if enum.IsValid(input) {
		fmt.Println("good input")
	}

	if input == StatusMerged {
		fmt.Println("nice job")
	}

	statuses := enum.ValuesOf[Status]()
	fmt.Println(statuses)

	// Output:
	// bad input
	// "postponed" is not a valid choice, allowed values are: "draft", "open", "merged", "closed"
	// good input
	// nice job
	// [draft open merged closed]
}

func Example_int() {
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

func Example_empty() {
	type Nothing int
	if err := enum.Validate[Nothing](1); err != nil {
		fmt.Println(err)
	}
	// Output:
	// Nothing doesn't have any definition
}

func Example_scope() {
	type OrderStatus string
	type UserStatus string
	var (
		OrderStatusPending    = enum.Def[OrderStatus]("pending")
		OrderStatusInProgress = enum.Def[OrderStatus]("in_progress")

		UserStatusActive = enum.Def[UserStatus]("pending")
		UserStatusBanned = enum.Def[UserStatus]("banned")
	)
	_, _, _, _ = OrderStatusPending, OrderStatusInProgress, UserStatusActive, UserStatusBanned

	if enum.IsValid[OrderStatus]("pending") {
		fmt.Println("pending is a valid order status")
	}
	if enum.IsValid[UserStatus]("pending") {
		fmt.Println("pending is a valid user status")
	}
	if !enum.IsValid[OrderStatus]("banned") {
		fmt.Println("banned is not a valid order status")
	}
	if !enum.IsValid[UserStatus]("in_progress") {
		fmt.Println("in_progress is not a valid user status")
	}
	// Output:
	// pending is a valid order status
	// pending is a valid user status
	// banned is not a valid order status
	// in_progress is not a valid user status
}

func BenchmarkIsValid(b *testing.B) {
	for i := 0; i < b.N; i++ {
		enum.IsValid(StatusDraft)
	}
}

func BenchmarkValidate(b *testing.B) {
	invalidStatus := Status("invalid")
	for i := 0; i < b.N; i++ {
		_ = enum.Validate(invalidStatus)
	}
}
