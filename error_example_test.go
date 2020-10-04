package zrr

import (
	"errors"
	"fmt"
)

func ExampleWrap() {
	err := Wrap(errors.New("std error")).Str("place", "here").Int("retry", 5)
	fmt.Println(err.Error())

	// Output: std error
}

func ExampleError_String() {
	err := Wrap(errors.New("std error")).Str("place", "here").Int("retry", 5)
	fmt.Println(err.String())

	// Output: std error :: place="here" retry=5
}
