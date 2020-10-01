package zrr

import (
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
