package zrr

import (
	"errors"
	"fmt"
)

func ExampleError_Error() {
	err := Wrap(errors.New("std error")).Str("place", "here").Int("retry", 5)
	fmt.Println(err.Error())

	// Output: std error :: place="here" retry=5
}
