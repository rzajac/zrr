package zrr

import (
	"errors"
	"testing"
	"time"

	"github.com/ctx42/testing/pkg/assert"
)

func Test_IsImmutable(t *testing.T) {
	tt := []struct {
		testN string

		exp bool
		err error
	}{
		{"1", true, Imm("em0")},
		{"2", true, Imm("em0", "ECode")},
		{"3", false, New("em0", "ECode")},
		{"4", false, errors.New("message")},
		{"5", false, nil},
	}

	for _, tc := range tt {
		t.Run(tc.testN, func(t *testing.T) {
			assert.Equal(t, tc.exp, IsImmutable(tc.err))
		})
	}
}

func Test_HasKey(t *testing.T) {
	tt := []struct {
		testN string

		exp bool
		key string
		err error
	}{
		{"1", true, "key0", New("em0").Int("key0", 123)},
		{"2", false, "key0", New("em0")},
		{"3", false, "key0", errors.New("message")},
		{"4", false, "key0", nil},
	}

	for _, tc := range tt {
		t.Run(tc.testN, func(t *testing.T) {
			assert.Equal(t, tc.exp, HasKey(tc.err, tc.key))
		})
	}
}

func Test_HasCode(t *testing.T) {
	tt := []struct {
		testN string

		exp  bool
		code []string
		err  error
	}{
		{"1", true, []string{"ECode"}, New("em0", "ECode")},
		{"2", true, []string{"ECodeX", "ECode"}, New("em0", "ECode")},
		{"3", false, []string{"ECodeX"}, New("em0", "ECode")},
		{"4", false, []string{"ECode"}, errors.New("message")},
		{"5", false, []string{"ECode"}, nil},
	}

	for _, tc := range tt {
		t.Run(tc.testN, func(t *testing.T) {
			assert.Equal(t, tc.exp, HasCode(tc.err, tc.code...))
		})
	}
}

func Test_GetCode(t *testing.T) {
	tt := []struct {
		testN string

		exp string
		err error
	}{
		{"1", "", New("em0")},
		{"2", "ECode", New("em0", "ECode")},
		{"3", "", errors.New("message")},
		{"4", "", nil},
	}

	for _, tc := range tt {
		t.Run(tc.testN, func(t *testing.T) {
			assert.Equal(t, tc.exp, GetCode(tc.err))
		})
	}
}

func Test_GetStr(t *testing.T) {
	tt := []struct {
		testN string

		err   error
		key   string
		value string
		exist bool
	}{
		{"1", New("em0"), "key0", "", false},
		{"2", New("em0").Str("key0", ""), "key0", "", true},
		{"3", New("em0").Str("key0", "val0"), "key0", "val0", true},
		{"4", New("em0").Int("key0", 0), "key0", "", false},
		{"5", nil, "key0", "", false},
	}

	for _, tc := range tt {
		t.Run(tc.testN, func(t *testing.T) {
			// --- When ---
			value, exist := GetStr(tc.err, tc.key)

			// --- Then ---
			assert.Equal(t, tc.exist, exist)
			assert.Equal(t, tc.value, value)
		})
	}
}

func Test_GetInt(t *testing.T) {
	tt := []struct {
		testN string

		err   error
		key   string
		value int
		exist bool
	}{
		{"1", New("em0"), "key0", 0, false},
		{"2", New("em0").Int("key0", 0), "key0", 0, true},
		{"3", New("em0").Int("key0", 123), "key0", 123, true},
		{"4", New("em0").Str("key0", "val0"), "key0", 0, false},
		{"5", nil, "key0", 0, false},
	}

	for _, tc := range tt {
		t.Run(tc.testN, func(t *testing.T) {
			// --- When ---
			value, exist := GetInt(tc.err, tc.key)

			// --- Then ---
			assert.Equal(t, tc.exist, exist)
			assert.Equal(t, tc.value, value)
		})
	}
}

func Test_GetInt64(t *testing.T) {
	tt := []struct {
		testN string

		err   error
		key   string
		value int64
		exist bool
	}{
		{"1", New("em0"), "key0", 0, false},
		{"2", New("em0").Int64("key0", 0), "key0", 0, true},
		{"3", New("em0").Int64("key0", 123), "key0", 123, true},
		{"4", New("em0").Int("key0", 123), "key0", 0, false},
		{"5", New("em0").Str("key0", "val0"), "key0", 0, false},
		{"6", nil, "key0", 0, false},
	}

	for _, tc := range tt {
		t.Run(tc.testN, func(t *testing.T) {
			// --- When ---
			value, exist := GetInt64(tc.err, tc.key)

			// --- Then ---
			assert.Equal(t, tc.exist, exist)
			assert.Equal(t, tc.value, value)
		})
	}
}

func Test_GetFloat64(t *testing.T) {
	tt := []struct {
		testN string

		err   error
		key   string
		value float64
		exist bool
	}{
		{"1", New("em0"), "key0", 0.0, false},
		{"2", New("em0").Float64("key0", 0.0), "key0", 0.0, true},
		{"3", New("em0").Float64("key0", 0.123), "key0", 0.123, true},
		{"4", New("em0").Str("key0", "val0"), "key0", 0.0, false},
		{"5", nil, "key0", 0.0, false},
	}

	for _, tc := range tt {
		t.Run(tc.testN, func(t *testing.T) {
			// --- When ---
			value, exist := GetFloat64(tc.err, tc.key)

			// --- Then ---
			assert.Equal(t, tc.exist, exist)
			assert.Equal(t, tc.value, value)
		})
	}
}

func Test_GetTime(t *testing.T) {
	now := time.Now()

	tt := []struct {
		testN string

		err   error
		key   string
		value time.Time
		exist bool
	}{
		{"1", New("em0"), "key0", time.Time{}, false},
		{"2", New("em0").Time("key0", time.Time{}), "key0", time.Time{}, true},
		{"3", New("em0").Time("key0", now), "key0", now, true},
		{"4", New("em0").Str("key0", "val0"), "key0", time.Time{}, false},
		{"5", nil, "key0", time.Time{}, false},
	}

	for _, tc := range tt {
		t.Run(tc.testN, func(t *testing.T) {
			// --- When ---
			value, exist := GetTime(tc.err, tc.key)

			// --- Then ---
			assert.Equal(t, tc.exist, exist)
			assert.Equal(t, tc.value, value)
		})
	}
}

func Test_GetBool(t *testing.T) {
	tt := []struct {
		testN string

		err   error
		key   string
		value bool
		exist bool
	}{
		{"1", New("em0"), "key0", false, false},
		{"2", New("em0").Bool("key0", false), "key0", false, true},
		{"3", New("em0").Bool("key0", true), "key0", true, true},
		{"4", New("em0").Str("key0", "val0"), "key0", false, false},
		{"5", nil, "key0", false, false},
	}

	for _, tc := range tt {
		t.Run(tc.testN, func(t *testing.T) {
			// --- When ---
			value, exist := GetBool(tc.err, tc.key)

			// --- Then ---
			assert.Equal(t, tc.exist, exist)
			assert.Equal(t, tc.value, value)
		})
	}
}

func Test_NilPointerError(t *testing.T) {
	// --- Given ---
	var err *Error

	// --- Then ---
	assert.False(t, HasKey(err, "key0"))
	assert.False(t, HasCode(err, "ECode"))

	_, got := GetStr(err, "key0")
	assert.False(t, got)

	_, got = GetInt(err, "key0")
	assert.False(t, got)

	_, got = GetFloat64(err, "key0")
	assert.False(t, got)

	_, got = GetTime(err, "key0")
	assert.False(t, got)

	_, got = GetBool(err, "key0")
	assert.False(t, got)
}
