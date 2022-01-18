package zrr_test

import (
	"errors"
	"fmt"
	"time"

	"github.com/rzajac/zrr"
)

func ExampleError() {
	// Create an error and add a bunch of context fields to it.
	err := zrr.Wrap(errors.New("std error"), "ECode").
		Str("str", "string").
		Int("int", 5).
		Float64("float64", 1.23).
		Time("time", time.Date(2020, time.October, 7, 23, 47, 0, 0, time.UTC)).
		Bool("bool", true)

	fmt.Println(err.Error())
	fmt.Println(zrr.GetStr(err, "str"))
	fmt.Println(zrr.GetInt(err, "int"))
	fmt.Println(zrr.GetFloat64(err, "float64"))
	fmt.Println(zrr.GetTime(err, "time"))
	fmt.Println(zrr.GetBool(err, "bool"))
	fmt.Println(err.ErrCode())

	// Output:
	// std error
	// string true
	// 5 true
	// 1.23 true
	// 2020-10-07 23:47:00 +0000 UTC true
	// true true
	// ECode
}

func ExampleError_wrappingZrrError() {
	err := zrr.New("my error").Str("key", "value")

	e1 := fmt.Errorf("zrr wrapped: %w", err)

	fmt.Println(e1)
	fmt.Println(errors.Is(e1, err))

	// Output:
	// zrr wrapped: my error
	// true
}

func ExampleImm() {
	// Create immutable error.
	var ErrPackageLevel = zrr.Imm("package level error", "ECode")

	// Somewhere in the code use ErrPackageLevel and add context to it.
	err := zrr.Wrap(ErrPackageLevel, "ENewCode").Str("path", "/path/to/file")

	fmt.Println(ErrPackageLevel, zrr.GetCode(ErrPackageLevel)) // Notice the error code has not been changed.
	fmt.Println(err, zrr.GetCode(err))
	fmt.Println(errors.Is(err, ErrPackageLevel))

	// Output:
	// package level error ECode
	// package level error ENewCode
	// true
}

func ExampleWrap() {
	err := errors.New("some error")

	e1 := zrr.Wrap(err).Str("key", "value")

	fmt.Println(e1)
	fmt.Println(errors.Is(e1, err))

	// Output:
	// some error
	// true
}

func ExampleGetCode() {
	err := zrr.New("message", "ECode").Int("retry", 5)

	fmt.Println(zrr.GetCode(err))

	// Output: ECode
}

func ExampleGetInt() {
	var err error
	err = zrr.New("message", "ECode").Int("retry", 5)

	fmt.Println(zrr.GetInt(err, "retry"))   // 5 true
	fmt.Println(zrr.GetInt(err, "not_set")) // 0 false
	fmt.Println(zrr.HasKey(err, "retry"))   // true
	fmt.Println(zrr.HasKey(err, "not_set")) // false

	// Output:
	// 5 true
	// 0 false
	// true
	// false
}

func ExampleHasKey() {
	err := zrr.New("message", "ECode").Int("retry", 5)

	fmt.Println(zrr.GetInt(err, "retry"))
	fmt.Println(zrr.GetInt(err, "not_set"))

	// Output: 5 true
	// 0 false
}

func ExampleError_Error() {
	err := zrr.New("message").Int("retry", 5)

	fmt.Println(err.Error())

	// Output: message
}

func ExampleError_Fields() {
	// Create an error and add a bunch of context fields to it.
	err := zrr.Wrap(errors.New("std error"), "ECode").
		Str("str", "string").
		Int("int", 5).
		Float64("float64", 1.23).
		Time("time", time.Date(2020, time.October, 7, 23, 47, 0, 0, time.UTC)).
		Bool("bool", true)

	// Somewhere else (maybe during logging extract the context fields.
	iter := err.Fields()
	for iter.Next() {
		key, val := iter.Get()
		fmt.Printf("%s = %v\n", key, val)
	}

	// Output:
	// bool = true
	// float64 = 1.23
	// int = 5
	// str = string
	// time = 2020-10-07 23:47:00 +0000 UTC
}
