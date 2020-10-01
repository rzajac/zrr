package zrr

import (
	"errors"
	"fmt"
	"testing"

	"github.com/stretchr/testify/suite"
)

func TestError(t *testing.T) { suite.Run(t, &ErrorSuite{}) }

type ErrorSuite struct{ suite.Suite }

func (ste *ErrorSuite) Test_Error_New() {
	// --- When ---
	E := New("em1")

	// --- Then ---
	ste.Exactly("em1", E.Error())
	ste.False(E.imm)
	ste.Exactly("em1", E.Error())
	ste.Exactly("em1", E.String())
}

func (ste *ErrorSuite) Test_Error_WithMeta() {
	// --- When ---
	E := New("em1").With("key", "val")

	// --- Then ---
	ste.Exactly("em1", E.Error())
	ste.False(E.imm)
	ste.Exactly("em1", E.Error())
	ste.Exactly("em1 --- key:val", E.String())
}

func (ste *ErrorSuite) TestName() {
	// --- Given ---
	e1 := errors.New("e1")
	e2 := fmt.Errorf("e2: %w", e1)
	e3 := fmt.Errorf("e3: %w", e2)
	e4 := fmt.Errorf("e4: %w", e3)

	// --- When ---
	// --- Then ---
	fmt.Println(e4)
	fmt.Println(errors.Unwrap(e4))
}

func (ste *ErrorSuite) TestName2() {
	// --- Given ---
	e1 := New("e1")
	e2 := Wrap(e1)
	e3 := Wrap(e2)
	e4 := Wrap(e3)

	// --- Then ---
	fmt.Println(e4)
	fmt.Println(errors.Unwrap(e4))
	fmt.Println(errors.Is(e4, e4))
	fmt.Println(errors.Is(e4, e3))
	fmt.Println(errors.Is(e4, e2))
	fmt.Println(errors.Is(e4, e1))

	fmt.Println("----")
	var tar *Error
	fmt.Println(errors.As(e4, &tar))
	fmt.Println(tar)
}
