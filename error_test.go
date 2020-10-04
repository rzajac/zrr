package zrr

import (
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

func TestError(t *testing.T) { suite.Run(t, &ErrorSuite{}) }

type ErrorSuite struct{ suite.Suite }

func (ste *ErrorSuite) Test_Error_New() {
	// --- When ---
	err := New("em0", "ECode")

	// --- Then ---
	ste.False(err.Immutable())
	ste.Exactly("em0", err.Error())
	ste.Exactly(`em0 :: code="ECode"`, err.String())
	ste.True(err.HasCode("ECode"))
}

func (ste *ErrorSuite) Test_Error_HasKey() {
	// --- When ---
	err := New("em0", "ECode")

	// --- Then ---
	ste.True(err.HasKey(KCode))
}

func (ste *ErrorSuite) Test_Error_Newf() {
	// --- When ---
	err := Newf("%s message", "error")

	// --- Then ---
	ste.False(err.Immutable())
	ste.Exactly("error message", err.Error())
	ste.Exactly("error message", err.String())
}

func (ste *ErrorSuite) Test_Error_Imm() {
	// --- When ---
	err := Imm("em0")

	// --- Then ---
	ste.True(err.Immutable())
	ste.Exactly("em0", err.Error())
	ste.Exactly("em0", err.String())
}

func (ste *ErrorSuite) Test_Error_Imm_WithCode() {
	// --- When ---
	err := Imm("em0", "ECode")

	// --- Then ---
	ste.True(err.Immutable())
	ste.Exactly("em0", err.Error())
	ste.Exactly(`em0 :: code="ECode"`, err.String())
}

func (ste *ErrorSuite) Test_Error_Imm_WithCodes() {
	// --- When ---
	err := Imm("em0", "ECode0", "ECode1")

	// --- Then ---
	ste.True(err.Immutable())
	ste.Exactly("em0", err.Error())
	ste.Exactly(`em0 :: code="ECode0"`, err.String())
}

func (ste *ErrorSuite) Test_Error_Code() {
	// --- When ---
	err := New("em0").Code("ECode")

	// --- Then ---
	ste.Exactly(`em0 :: code="ECode"`, err.String())
}

func (ste *ErrorSuite) Test_Error_HasCode() {
	// --- Given ---
	err := New("em0")
	ste.False(err.HasCode("ECode"))

	// --- When ---
	err = err.Code("ECode")

	// --- Then ---
	ste.True(err.HasCode("ECode"))
	ste.Exactly(`em0 :: code="ECode"`, err.String())
}

func (ste *ErrorSuite) Test_Error_Str() {
	// --- When ---
	err := Newf("em0").Str("key0", "val0")

	// --- Then ---
	ste.Exactly(`em0 :: key0="val0"`, err.String())
}

func (ste *ErrorSuite) Test_Error_GetStr() {
	tt := []struct {
		testN string

		err   *Error
		key   string
		value string
		exist bool
	}{
		{"1", New("em0"), "key0", "", false},
		{"2", New("em0").Str("key0", ""), "key0", "", true},
		{"3", New("em0").Str("key0", "val0"), "key0", "val0", true},
		{"4", New("em0").Int("key0", 0), "key0", "", false},
	}

	for _, tc := range tt {
		ste.T().Run(tc.testN, func(t *testing.T) {
			// --- When ---
			value, exist := tc.err.GetStr(tc.key)

			// --- Then ---
			assert.Exactly(t, tc.exist, exist, "test %s", tc.testN)
			assert.Exactly(t, tc.value, value, "test %s", tc.testN)
		})
	}
}

func (ste *ErrorSuite) Test_Error_Int() {
	// --- When ---
	err := Newf("em0").Int("key0", 0)

	// --- Then ---
	ste.Exactly(`em0 :: key0=0`, err.String())
}

func (ste *ErrorSuite) Test_Error_GetInt() {
	tt := []struct {
		testN string

		err   *Error
		key   string
		value int
		exist bool
	}{
		{"1", New("em0"), "key0", 0, false},
		{"2", New("em0").Int("key0", 0), "key0", 0, true},
		{"3", New("em0").Int("key0", 123), "key0", 123, true},
		{"4", New("em0").Str("key0", "val0"), "key0", 0, false},
	}

	for _, tc := range tt {
		ste.T().Run(tc.testN, func(t *testing.T) {
			// --- When ---
			value, exist := tc.err.GetInt(tc.key)

			// --- Then ---
			assert.Exactly(t, tc.exist, exist, "test %s", tc.testN)
			assert.Exactly(t, tc.value, value, "test %s", tc.testN)
		})
	}
}

func (ste *ErrorSuite) Test_Error_Float64() {
	// --- When ---
	err := Newf("em0").Float64("key0", 0.123)

	// --- Then ---
	ste.Exactly(`em0 :: key0=0.123`, err.String())
}

func (ste *ErrorSuite) Test_Error_GetFloat64() {
	tt := []struct {
		testN string

		err   *Error
		key   string
		value float64
		exist bool
	}{
		{"1", New("em0"), "key0", 0.0, false},
		{"2", New("em0").Float64("key0", 0.0), "key0", 0.0, true},
		{"3", New("em0").Float64("key0", 0.123), "key0", 0.123, true},
		{"4", New("em0").Str("key0", "val0"), "key0", 0.0, false},
	}

	for _, tc := range tt {
		ste.T().Run(tc.testN, func(t *testing.T) {
			// --- When ---
			value, exist := tc.err.GetFloat64(tc.key)

			// --- Then ---
			assert.Exactly(t, tc.exist, exist, "test %s", tc.testN)
			assert.Exactly(t, tc.value, value, "test %s", tc.testN)
		})
	}
}

func (ste *ErrorSuite) Test_Error_Time() {
	// --- Given ---
	tim := time.Now()

	// --- When ---
	err := Newf("em0").Time("key0", tim)

	// --- Then ---
	exp := fmt.Sprintf(`em0 :: key0=%s`, tim.Format(time.RFC3339Nano))
	ste.Exactly(exp, err.String())
}

func (ste *ErrorSuite) Test_Error_GetTime() {
	now := time.Now()

	tt := []struct {
		testN string

		err   *Error
		key   string
		value time.Time
		exist bool
	}{
		{"1", New("em0"), "key0", time.Time{}, false},
		{"2", New("em0").Time("key0", time.Time{}), "key0", time.Time{}, true},
		{"3", New("em0").Time("key0", now), "key0", now, true},
		{"4", New("em0").Str("key0", "val0"), "key0", time.Time{}, false},
	}

	for _, tc := range tt {
		ste.T().Run(tc.testN, func(t *testing.T) {
			// --- When ---
			value, exist := tc.err.GetTime(tc.key)

			// --- Then ---
			assert.Exactly(t, tc.exist, exist, "test %s", tc.testN)
			assert.Exactly(t, tc.value, value, "test %s", tc.testN)
		})
	}
}

func (ste *ErrorSuite) Test_Error_Bool() {
	// --- When ---
	err0 := Newf("em0").Bool("key0", true)
	err1 := Newf("em0").Bool("key0", false)

	// --- Then ---
	ste.Exactly(`em0 :: key0=true`, err0.String())
	ste.Exactly(`em0 :: key0=false`, err1.String())
}

func (ste *ErrorSuite) Test_Error_GetBool() {
	tt := []struct {
		testN string

		err   *Error
		key   string
		value bool
		exist bool
	}{
		{"1", New("em0"), "key0", false, false},
		{"2", New("em0").Bool("key0", false), "key0", false, true},
		{"3", New("em0").Bool("key0", true), "key0", true, true},
		{"4", New("em0").Str("key0", "val0"), "key0", false, false},
	}

	for _, tc := range tt {
		ste.T().Run(tc.testN, func(t *testing.T) {
			// --- When ---
			value, exist := tc.err.GetBool(tc.key)

			// --- Then ---
			assert.Exactly(t, tc.exist, exist, "test %s", tc.testN)
			assert.Exactly(t, tc.value, value, "test %s", tc.testN)
		})
	}
}

func (ste *ErrorSuite) Test_Error_Multi_Metadata() {
	// --- When ---
	err := New("test msg", "ECode").Int("key0", 5).Str("key1", "I'm a string")

	// --- Then ---
	ste.Exactly("test msg", err.Error())
	ste.Exactly(`test msg :: code="ECode" key0=5 key1="I'm a string"`, err.String())
}

func (ste *ErrorSuite) Test_Error_Wrap() {
	// --- Given ---
	e := errors.New("std error")

	// --- When ---
	err := Wrap(e)

	// --- Then ---
	ste.IsType(&Error{}, err)
	ste.False(err.Immutable())
	ste.Exactly("std error", err.Error())
	ste.Exactly("std error", err.String())
}

func (ste *ErrorSuite) Test_Error_Wrap_nil() {
	// --- When ---
	err := Wrap(nil)

	// --- Then ---
	ste.Nil(err)
}

func (ste *ErrorSuite) Test_Error_Wrap_Error() {
	// --- Given ---
	err0 := New("test msg")

	// --- When ---
	err1 := Wrap(err0)

	// --- Then ---
	ste.Same(err0, err1)
}

func (ste *ErrorSuite) Test_Error_Unwrap() {
	// --- Given ---
	err0 := errors.New("std error")

	// --- When ---
	err1 := Wrap(err0).Unwrap()

	// --- Then ---
	ste.Same(err0, err1)
}

func (ste *ErrorSuite) Test_Error_with_immutable() {
	// --- Given ---
	err0 := Imm("immutable error", "ECode")

	// --- When ---
	err1 := err0.Str("key0", "val0")

	// --- Then ---
	ste.NotSame(err0, err1)
	ste.False(err1.Immutable())
	ste.True(err1.HasCode("ECode"))

	val, ok := err1.GetStr("key0")
	ste.Exactly("val0", val)
	ste.True(ok)
}

type TestT string

func (t TestT) ZrrFields(e *Error) *Error {
	return e.Str(KCode, string(t)).Int("key1", 123)
}

func (ste *ErrorSuite) Test_Error_FieldsFrom() {
	// --- Given ---
	t := TestT("test")

	// --- When ---
	err := New("test msg").FieldsFrom(t)

	// --- Then ---
	ste.Exactly(`test msg :: code="test" key1=123`, err.String())
}

func (ste *ErrorSuite) Test_Error_FieldsFrom_Immutable() {
	// --- Given ---
	imm := Imm("test msg", "TCode")
	t := TestT("test")

	// --- When ---
	err := imm.FieldsFrom(t)

	// --- Then ---
	ste.True(errors.Is(err, imm))
	ste.Exactly(`test msg :: code="test" key1=123`, err.String())
}
