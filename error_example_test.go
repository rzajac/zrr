package zrr_test

import (
	"errors"
	"fmt"
	"time"

	"github.com/rzajac/zrr"
)

func Example_error_Error() {
	err := zrr.Wrap(errors.New("std error")).
		Code("ECode").
		Str("str", "here").
		Int("int", 5).
		Float64("float64", 1.23).
		Time("time", time.Date(2020, time.October, 7, 23, 47, 0, 0, time.UTC)).
		Bool("bool", true)

	fmt.Println(err.Error())

	// Output: std error :: bool=true code="ECode" float64=1.23 int=5 str="here" time=2020-10-07T23:47:00Z
}

func Example_getCode() {
	err := zrr.New("message", "ECode").Int("retry", 5)
	fmt.Println(zrr.GetCode(err))

	// Output: ECode
}
func Example_getInt() {
	err := zrr.New("message", "ECode").Int("retry", 5)
	fmt.Println(zrr.GetInt(err, "retry"))
	fmt.Println(zrr.GetInt(err, "not_set"))
	fmt.Println(zrr.HasKey(err, "retry"))
	fmt.Println(zrr.HasKey(err, "not_set"))

	// Output: 5 true
	// 0 false
	// true
	// false
}
