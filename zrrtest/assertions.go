package zrrtest

import (
	"testing"
	"time"

	"github.com/ctx42/testing/pkg/assert"

	"github.com/rzajac/zrr"
)

// AssertCause asserts err is instance of zrr.Error and has error
// message (without key value pairs) equal to cause.
func AssertCause(t *testing.T, err error, cause string, _ ...any) {
	t.Helper()

	assert.NotNil(t, err)
	var E *zrr.Error
	assert.ErrorAs(t, &E, err)

	got := E.Error()
	if got != cause {
		t.Errorf("expected cause '%s' but got '%s'", cause, got)
		return
	}
}

// AssertContains asserts err is instance of zrr.Error and has error
// message (without key value pairs) which contains cause string.
func AssertContains(t *testing.T, err error, cause string, _ ...any) {
	t.Helper()

	assert.NotNil(t, err)
	var E *zrr.Error
	assert.ErrorAs(t, &E, err)

	got := E.Error()
	assert.Contain(t, cause, got)
}

// AssertCode asserts err is instance of zrr.Error and has error code exp.
func AssertCode(t *testing.T, err error, exp string, _ ...any) {
	t.Helper()

	assert.NotNil(t, err)
	var E *zrr.Error
	assert.ErrorAs(t, &E, err)
	assert.Equal(t, exp, E.ErrCode())
}

// AssertEqual asserts err and got are instance of zrr.Error and their error
// messages (with key value pairs) are equal.
func AssertEqual(t *testing.T, exp, got error, _ ...any) {
	t.Helper()

	assert.NotNil(t, exp)
	assert.NotNil(t, got)
	var EE, EG *zrr.Error
	assert.ErrorAs(t, &EE, exp)
	assert.ErrorAs(t, &EG, got)
	assert.Equal(t, EE.Error(), EG.Error())
}

// AssertNoKey asserts err is instance of zrr.Error and has no key set.
func AssertNoKey(t *testing.T, err error, key string, _ ...any) {
	t.Helper()

	assert.NotNil(t, err)
	assert.False(t, zrr.HasKey(err, key))
}

// AssertStr asserts err is instance of zrr.Error and has key with value exp.
func AssertStr(t *testing.T, err error, key, exp string, _ ...any) {
	t.Helper()

	assert.NotNil(t, err)
	got, ok := zrr.GetStr(err, key)
	assert.True(t, ok)
	assert.Equal(t, exp, got)
}

// AssertInt asserts err is instance of zrr.Error and has key with value exp.
func AssertInt(t *testing.T, err error, key string, exp int, _ ...any) {
	t.Helper()

	assert.NotNil(t, err)
	got, ok := zrr.GetInt(err, key)
	assert.True(t, ok)
	assert.Equal(t, exp, got)
}

// AssertInt64 asserts err is instance of zrr.Error and has key with value exp.
func AssertInt64(t *testing.T, err error, key string, exp int64, _ ...any) {
	t.Helper()

	assert.NotNil(t, err)
	got, ok := zrr.GetInt64(err, key)
	assert.True(t, ok)
	assert.Equal(t, exp, got)
}

// AssertFloat64 asserts err is instance of zrr.Error and has key with value exp.
func AssertFloat64(t *testing.T, err error, key string, exp float64, _ ...any) {
	t.Helper()

	assert.NotNil(t, err)
	got, ok := zrr.GetFloat64(err, key)
	assert.True(t, ok)
	assert.Equal(t, exp, got)
}

// AssertTime asserts err is instance of zrr.Error and has key with value exp.
func AssertTime(t *testing.T, err error, key string, exp time.Time, _ ...any) {
	t.Helper()

	assert.NotNil(t, err)
	got, ok := zrr.GetTime(err, key)
	assert.True(t, ok)
	assert.Equal(t, exp, got)
}

// AssertBool asserts err is instance of zrr.Error and has key with value exp.
func AssertBool(t *testing.T, err error, key string, exp bool, _ ...any) {
	t.Helper()

	assert.NotNil(t, err)
	got, ok := zrr.GetBool(err, key)
	assert.True(t, ok)
	assert.Equal(t, exp, got)
}
