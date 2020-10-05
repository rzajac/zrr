package zrr

import (
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func Test_Error_New(t *testing.T) {
	// --- When ---
	err := New("em0", "ECode")

	// --- Then ---
	assert.False(t, err.imm)
	assert.True(t, HasCode(err, "ECode"))
	assert.Exactly(t, `em0 :: code="ECode"`, err.Error())
}

func Test_Error_Newf(t *testing.T) {
	// --- When ---
	err := Newf("%s message", "error")

	// --- Then ---
	assert.False(t, err.imm)
	assert.Exactly(t, "error message", err.Error())
}

func Test_Error_Imm(t *testing.T) {
	// --- When ---
	err := Imm("em0")

	// --- Then ---
	assert.True(t, err.imm)
	assert.Exactly(t, "em0", err.Error())
}

func Test_Error_Imm_WithCode(t *testing.T) {
	// --- When ---
	err := Imm("em0", "ECode")

	// --- Then ---
	assert.True(t, err.imm)
	assert.Exactly(t, `em0 :: code="ECode"`, err.Error())
}

func Test_Error_Imm_WithCodes(t *testing.T) {
	// --- When ---
	err := Imm("em0", "ECode0", "ECode1")

	// --- Then ---
	assert.True(t, err.imm)
	assert.Exactly(t, `em0 :: code="ECode0"`, err.Error())
}

func Test_Error_Code(t *testing.T) {
	// --- When ---
	err := New("em0").Code("ECode")

	// --- Then ---
	assert.Exactly(t, `em0 :: code="ECode"`, err.Error())
}

func Test_Error_Str(t *testing.T) {
	// --- When ---
	err := Newf("em0").Str("key0", "val0")

	// --- Then ---
	assert.Exactly(t, `em0 :: key0="val0"`, err.Error())
}

func Test_Error_Int(t *testing.T) {
	// --- When ---
	err := Newf("em0").Int("key0", 0)

	// --- Then ---
	assert.Exactly(t, `em0 :: key0=0`, err.Error())
}

func Test_Error_Float64(t *testing.T) {
	// --- When ---
	err := Newf("em0").Float64("key0", 0.123)

	// --- Then ---
	assert.Exactly(t, `em0 :: key0=0.123`, err.Error())
}

func Test_Error_Time(t *testing.T) {
	// --- Given ---
	tim := time.Now()

	// --- When ---
	err := Newf("em0").Time("key0", tim)

	// --- Then ---
	exp := fmt.Sprintf(`em0 :: key0=%s`, tim.Format(time.RFC3339Nano))
	assert.Exactly(t, exp, err.Error())
}

func Test_Error_Bool(t *testing.T) {
	// --- When ---
	err0 := Newf("em0").Bool("key0", true)
	err1 := Newf("em0").Bool("key0", false)

	// --- Then ---
	assert.Exactly(t, `em0 :: key0=true`, err0.Error())
	assert.Exactly(t, `em0 :: key0=false`, err1.Error())
}

func Test_Error_Multi_Metadata(t *testing.T) {
	// --- When ---
	err := New("test msg", "ECode").Int("key0", 5).Str("key1", "I'm a string")

	// --- Then ---
	assert.Exactly(t, `test msg :: code="ECode" key0=5 key1="I'm a string"`, err.Error())
}

func Test_Error_Wrap(t *testing.T) {
	// --- Given ---
	e := errors.New("std error")

	// --- When ---
	err := Wrap(e)

	// --- Then ---
	assert.IsType(t, &Error{}, err)
	assert.False(t, err.imm)
	assert.Exactly(t, "std error", err.Error())
}

func Test_Error_Wrap_nil(t *testing.T) {
	// --- When ---
	err := Wrap(nil)

	// --- Then ---
	assert.Nil(t, err)
}

func Test_Error_Wrap_Error(t *testing.T) {
	// --- Given ---
	err0 := New("test msg")

	// --- When ---
	err1 := Wrap(err0)

	// --- Then ---
	assert.Same(t, err0, err1)
}

func Test_Error_Unwrap(t *testing.T) {
	// --- Given ---
	err0 := errors.New("std error")

	// --- When ---
	err1 := Wrap(err0).Unwrap()

	// --- Then ---
	assert.Same(t, err0, err1)
}

func Test_Error_with_immutable(t *testing.T) {
	// --- Given ---
	err0 := Imm("immutable error", "ECode")

	// --- When ---
	err1 := err0.Str("key0", "val0")

	// --- Then ---
	assert.NotSame(t, err0, err1)
	assert.False(t, err1.imm)
	assert.True(t, HasCode(err1, "ECode"))

	val, ok := GetStr(err1, "key0")
	assert.Exactly(t, "val0", val)
	assert.True(t, ok)
}

type implementor string

func (t implementor) ZrrFields(e *Error) *Error {
	return e.Str(KCode, string(t)).Int("key1", 123)
}

func Test_Error_FieldsFrom(t *testing.T) {
	// --- Given ---
	imp := implementor("test")

	// --- When ---
	err := New("test msg").FieldsFrom(imp)

	// --- Then ---
	assert.Exactly(t, `test msg :: code="test" key1=123`, err.Error())
}

func Test_Error_FieldsFrom_Immutable(t *testing.T) {
	// --- Given ---
	imm := Imm("test msg", "TCode")
	imp := implementor("test")

	// --- When ---
	err := imm.FieldsFrom(imp)

	// --- Then ---
	assert.True(t, errors.Is(err, imm))
	assert.Exactly(t, `test msg :: code="test" key1=123`, err.Error())
}

func Test_OnlyMessage(t *testing.T) {
	var nilErr *Error

	tt := []struct {
		testN string

		exp string
		err error
	}{
		{"1", "message", New("message", "ECode")},
		{"2", "message", New("message")},
		{"3", "message", errors.New("message")},
		{"4", "message", New("message").Str("key0", "val0")},
		{"5", "", nil},
		{"6", "", nilErr},
	}

	for _, tc := range tt {
		t.Run(tc.testN, func(t *testing.T) {
			assert.Exactly(t, tc.exp, OnlyMessage(tc.err), "test %s", tc.testN)
		})
	}
}
