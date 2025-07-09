package zrr

import (
	"encoding/json"
	"errors"
	"testing"
	"time"

	"github.com/ctx42/testing/pkg/assert"
)

func Test_Error_New(t *testing.T) {
	// --- When ---
	err := New("em0", "ECode")

	// --- Then ---
	assert.False(t, err.imm)
	assert.True(t, HasCode(err, "ECode"))
	assert.Equal(t, "em0", err.Error())
}

func Test_Error_Newf(t *testing.T) {
	// --- When ---
	err := Newf("%s message", "error")

	// --- Then ---
	assert.False(t, err.imm)
	assert.Equal(t, "error message", err.Error())
}

func Test_Error_Imm(t *testing.T) {
	// --- When ---
	err := Imm("em0")

	// --- Then ---
	assert.True(t, err.imm)
	assert.Equal(t, "em0", err.Error())
}

func Test_Error_Imm_WithCode(t *testing.T) {
	// --- When ---
	err := Imm("em0", "ECode")

	// --- Then ---
	assert.True(t, err.imm)
	assert.Equal(t, "em0", err.Error())
	assert.Equal(t, "ECode", GetCode(err))
}

func Test_Error_Imm_WithCodes(t *testing.T) {
	// --- When ---
	err := Imm("em0", "ECode0", "ECode1")

	// --- Then ---
	assert.True(t, err.imm)
	assert.Equal(t, "em0", err.Error())
	assert.Equal(t, "ECode0", GetCode(err))
}

func Test_Error_Imm_Wrap_withNewCode(t *testing.T) {
	// --- When ---
	im := Imm("em0", "ECode0")
	ne := Wrap(im, "ECode1")

	// --- Then ---
	assert.False(t, ne.imm)
	assert.NotSame(t, im, ne)
	assert.Same(t, im, ne.Unwrap())
	assert.Equal(t, "em0", ne.Error())
	assert.Equal(t, "ECode1", GetCode(ne))
}

func Test_Error_ErrCode(t *testing.T) {
	// --- When ---
	err := New("em0", "ECode")

	// --- Then ---
	assert.Equal(t, "em0", err.Error())
	assert.Equal(t, "ECode", GetCode(err))
}

func Test_Error_Str(t *testing.T) {
	// --- When ---
	err := Newf("em0").Str("key0", "val0")

	// --- Then ---
	assert.Equal(t, "em0", err.Error())
	assert.Equal(t, map[string]interface{}{"key0": "val0"}, err.GetMetadata())
}

func Test_Error_StrAppend_append(t *testing.T) {
	// --- When ---
	err := Newf("em0").Str("key0", "val0").StrAppend("key0", "1")

	// --- Then ---
	assert.Equal(t, "em0", err.Error())
	assert.Equal(t, map[string]interface{}{"key0": "val0;1"}, err.GetMetadata())
}

func Test_Error_StrAppend_create(t *testing.T) {
	// --- When ---
	err := Newf("em0").StrAppend("key0", "val1")

	// --- Then ---
	assert.Equal(t, "em0", err.Error())
	assert.Equal(t, map[string]interface{}{"key0": "val1"}, err.GetMetadata())
}

func Test_Error_StrAppend_overrideNonString(t *testing.T) {
	// --- When ---
	err := Newf("em0").Int("key0", 1).StrAppend("key0", "val1")

	// --- Then ---
	assert.Equal(t, "em0", err.Error())
	assert.Equal(t, map[string]interface{}{"key0": "val1"}, err.GetMetadata())
}

func Test_Error_Int(t *testing.T) {
	// --- When ---
	err := Newf("em0").Int("key0", 0)

	// --- Then ---
	assert.Equal(t, "em0", err.Error())
	assert.Equal(t, map[string]interface{}{"key0": 0}, err.GetMetadata())
}

func Test_Error_Int64(t *testing.T) {
	// --- When ---
	err := Newf("em0").Int64("key0", 1234)

	// --- Then ---
	assert.Equal(t, "em0", err.Error())
	assert.Equal(t, map[string]interface{}{"key0": int64(1234)}, err.GetMetadata())
}

func Test_Error_Float64(t *testing.T) {
	// --- When ---
	err := Newf("em0").Float64("key0", 0.123)

	// --- Then ---
	assert.Equal(t, "em0", err.Error())
	assert.Equal(t, map[string]interface{}{"key0": 0.123}, err.GetMetadata())
}

func Test_Error_Time(t *testing.T) {
	// --- Given ---
	tim := time.Now()

	// --- When ---
	err := Newf("em0").Time("key0", tim)

	// --- Then ---
	assert.Equal(t, "em0", err.Error())
	assert.Equal(t, map[string]interface{}{"key0": tim}, err.GetMetadata())
}

func Test_Error_Bool(t *testing.T) {
	// --- When ---
	err0 := Newf("em0").Bool("key0", true)
	err1 := Newf("em0").Bool("key0", false)

	// --- Then ---
	assert.Equal(t, "em0", err0.Error())
	assert.Equal(t, map[string]interface{}{"key0": true}, err0.GetMetadata())

	assert.Equal(t, "em0", err1.Error())
	assert.Equal(t, map[string]interface{}{"key0": false}, err1.GetMetadata())
}

func Test_Error_GetMetadata(t *testing.T) {
	// --- Given ---
	err0 := Imm("immutable error", "ECode").Int("key", 123)

	// --- When ---
	got := err0.GetMetadata()

	// --- Then ---
	exp := map[string]interface{}{
		"key": 123,
	}
	assert.Equal(t, exp, got)
}

func Test_Error_GetMetadata_multi(t *testing.T) {
	// --- When ---
	err := New("test msg", "ECode").Int("key0", 5).Str("key1", "I'm a string")

	// --- Then ---
	assert.Equal(t, "test msg", err.Error())
	assert.Equal(t, "ECode", GetCode(err))
	exp := map[string]interface{}{
		"key0": 5,
		"key1": "I'm a string",
	}
	assert.Equal(t, exp, err.GetMetadata())
}

func Test_Error_Wrap(t *testing.T) {
	// --- Given ---
	e := errors.New("std error")

	// --- When ---
	err := Wrap(e)

	// --- Then ---
	assert.Type(t, &Error{}, err)
	assert.False(t, err.imm)
	assert.Equal(t, "std error", err.Error())
	assert.NotNil(t, err.Unwrap())
}

func Test_Error_Wrap_withNewCode(t *testing.T) {
	// --- Given ---
	e := New("std error", "ECode0")

	// --- When ---
	err := Wrap(e, "ECode1")

	// --- Then ---
	assert.Type(t, &Error{}, err)
	assert.False(t, err.imm)
	assert.Equal(t, "std error", err.Error())
	assert.Equal(t, "ECode1", GetCode(err))
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
	assert.Equal(t, "val0", val)
	assert.True(t, ok)
}

func Test_Error_Error(t *testing.T) {
	tt := []struct {
		testN string

		exp string
		err *Error
	}{
		{"1", "message", New("message", "ECode")},
		{"2", "message", New("message")},
		{"3", "message", New("message").Str("key0", "val0")},
	}

	for _, tc := range tt {
		t.Run(tc.testN, func(t *testing.T) {
			assert.ErrorEqual(t, tc.exp, tc.err)
		})
	}
}

func Test_Error_ZrrFields(t *testing.T) {
	// --- Given ---
	err0 := New("test msg").Str("key0", "val0").Int("key1", 1)

	// --- When ---
	err1 := Wrap(err0.Unwrap()).SetErrMetadata(err0.GetMetadata())

	// --- Then ---
	assert.Equal(t, err0.Error(), err1.Error())
	assert.Same(t, err0.Unwrap(), err1.Unwrap())
}

func Test_Error_MarshalJSON(t *testing.T) {
	t.Run("without code", func(t *testing.T) {
		// --- Given ---
		e := New("test msg")

		// --- When ---
		data, err := json.Marshal(e)

		// --- Then ---
		assert.NoError(t, err)
		exp := `{"error":"test msg", "code":""}`
		assert.JSON(t, exp, string(data))
	})

	t.Run("with code no meta", func(t *testing.T) {
		// --- Given ---
		e := New("test msg", "ECTest")

		// --- When ---
		data, err := json.Marshal(e)

		// --- Then ---
		assert.NoError(t, err)
		exp := `{"error":"test msg", "code":"ECTest"}`
		assert.JSON(t, exp, string(data))
	})

	t.Run("with code and meta", func(t *testing.T) {
		// --- Given ---
		e := New("test msg", "ECTest").Str("key", "value")

		// --- When ---
		data, err := json.Marshal(e)

		// --- Then ---
		assert.NoError(t, err)
		exp := `{"error":"test msg", "code":"ECTest", "meta": {"key": "value"}}`
		assert.JSON(t, exp, string(data))
	})
}

func Test_Error_UnmarshalJSON(t *testing.T) {
	t.Run("without code", func(t *testing.T) {
		// --- Given ---
		data := []byte(`{"error":"test msg", "code":""}`)

		// --- When ---
		var e *Error
		err := json.Unmarshal(data, &e)

		// --- Then ---
		assert.NoError(t, err)
		assert.Equal(t, "test msg", e.error.Error())
		assert.Equal(t, "", e.ErrCode())
		assert.Len(t, 0, e.meta)
		assert.NotNil(t, e.meta)
	})

	t.Run("with code no meta", func(t *testing.T) {
		// --- Given ---
		data := []byte(`{"error":"test msg", "code":"ECode"}`)

		// --- When ---
		var e *Error
		err := json.Unmarshal(data, &e)

		// --- Then ---
		assert.NoError(t, err)
		assert.Equal(t, "test msg", e.error.Error())
		assert.Equal(t, "ECode", e.ErrCode())
		assert.Len(t, 0, e.meta)
		assert.NotNil(t, e.meta)
	})

	t.Run("with code and meta", func(t *testing.T) {
		// --- Given ---
		data := []byte(`{
			"error":"test msg", 
			"code":"ECode", 
			"meta": {
				"key": 123, 
				"tim": "2022-01-18T13:57:00Z"
			}
		}`)

		// --- When ---
		var e *Error
		err := json.Unmarshal(data, &e)

		// --- Then ---
		assert.NoError(t, err)
		assert.Equal(t, "test msg", e.error.Error())
		assert.Equal(t, "ECode", e.ErrCode())
		assert.Len(t, 2, e.meta)
		assert.HasKey(t, "key", e.meta)
		assert.Equal(t, float64(123), e.meta["key"])
		assert.Equal(t, "2022-01-18T13:57:00Z", e.meta["tim"])
	})

	t.Run("without error key", func(t *testing.T) {
		// --- Given ---
		data := []byte(`{"code":"code"}`)

		// --- When ---
		var e *Error
		err := json.Unmarshal(data, &e)

		// --- Then ---
		assert.ErrorIs(t, err, ErrInvJSON)
	})

	t.Run("unmarshal error", func(t *testing.T) {
		// --- Given ---
		data := []byte(`[1, 2, 3]`)

		// --- When ---
		var e *Error
		err := json.Unmarshal(data, &e)

		// --- Then ---
		assert.Type(t, err, &json.UnmarshalTypeError{})
	})
}

type implementor struct{ meta map[string]interface{} }

func (t implementor) GetMetadata() map[string]interface{} { return t.meta }

func Test_Error_SetMetadataFrom(t *testing.T) {
	// --- Given ---
	src := implementor{meta: map[string]interface{}{"k1": "v1", "k2": 2}}

	e := New("message", "code").Str("k2", "1")

	// --- When ---
	ne := e.SetMetadataFrom(src)

	// --- Then ---
	exp := map[string]interface{}{
		"k1": "v1",
		"k2": 2,
	}
	assert.Equal(t, exp, ne.GetMetadata())
}
