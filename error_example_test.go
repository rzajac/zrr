package zrr_test

import (
	"errors"
	"fmt"
	"time"

	"github.com/rzajac/zrr"
)

func ExampleError() {
	err := zrr.Wrap(errors.New("std error")).
		Code("ECode").
		Str("str", "string").
		Int("int", 5).
		Float64("float64", 1.23).
		Time("time", time.Date(2020, time.October, 7, 23, 47, 0, 0, time.UTC)).
		Bool("bool", true)

	fmt.Println(err.Error())

	// Output: std error :: bool=true code="ECode" float64=1.23 int=5 str="string" time=2020-10-07T23:47:00Z
}

func ExampleGetCode() {
	err := zrr.New("message", "ECode").Int("retry", 5)

	fmt.Println(zrr.GetCode(err))

	// Output: ECode
}

func ExampleGetInt() {
	err := zrr.New("message", "ECode").Int("retry", 5)

	fmt.Println(zrr.GetInt(err, "retry"))
	fmt.Println(zrr.GetInt(err, "not_set"))

	// Output:
	// 5 true
	// 0 false
}

func ExampleHasKey() {
	err := zrr.New("message", "ECode").Int("retry", 5)

	fmt.Println(zrr.GetInt(err, "retry"))
	fmt.Println(zrr.GetInt(err, "not_set"))

	// Output: 5 true
	// 0 false
}

func ExampleWrap() {
	err := errors.New("message")

	err = zrr.Wrap(err).Int("retry", 5)

	fmt.Println(err)

	// Output: message :: retry=5
}

func ExampleError_Cause() {
	err := zrr.New("message").Int("retry", 5)

	fmt.Println(err.Cause())

	// Output: message
}

func ExampleError_Fields() {
	err := zrr.Wrap(errors.New("std error")).
		Code("ECode").
		Str("str", "string").
		Int("int", 5).
		Float64("float64", 1.23).
		Time("time", time.Date(2020, time.October, 7, 23, 47, 0, 0, time.UTC)).
		Bool("bool", true)

	iter := err.Fields()
	for iter.Next() {
		key, val := iter.Get()
		fmt.Printf("%s = %v\n", key, val)
	}

	// Output:
	// bool = true
	// code = ECode
	// float64 = 1.23
	// int = 5
	// str = string
	// time = 2020-10-07 23:47:00 +0000 UTC
}
