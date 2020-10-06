package zrrtest

import (
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/rzajac/zrr"
)

// AssertCause asserts err is instance of zrr.Error and has error message
// equal to cause.
func AssertCause(t *testing.T, err error, cause string, args ...interface{}) {
	t.Helper()
	args = mArgs(args...)

	E, ok := err.(*zrr.Error)
	if !ok {
		t.Error("expected err to ne instance of zrr.Error")
		return
	}

	got := E.Unwrap().Error()
	if got != cause {
		t.Errorf("expected cause '%s' but got '%s'", cause, got)
		return
	}
}

// AssertCode asserts err is instance of zrr.Error and has error code exp.
func AssertCode(t *testing.T, err error, exp string, args ...interface{}) {
	t.Helper()
	args = mArgs(args...)
	AssertStr(t, err, zrr.KCode, exp)
}

// AssertStr asserts err is instance of zrr.Error and has key with value exp.
func AssertStr(t *testing.T, err error, key, exp string, args ...interface{}) {
	t.Helper()
	args = mArgs(args...)

	got, ok := zrr.GetStr(err, key)
	require.True(t, ok, "expected key '%s' is present", key)
	assert.Exactly(t, exp, got, "expected %s='%s' got %s='%s'", key, exp, key, got)
}

// AssertInt asserts err is instance of zrr.Error and has key with value exp.
func AssertInt(t *testing.T, err error, key string, exp int, args ...interface{}) {
	t.Helper()
	args = mArgs(args...)

	got, ok := zrr.GetInt(err, key)
	require.True(t, ok, "expected key '%s' is present", key)
	assert.Exactly(t, exp, got, "expected %s=%d got %s=%d", key, exp, key, got)
}

// AssertFloat64 asserts err is instance of zrr.Error and has key with value exp.
func AssertFloat64(t *testing.T, err error, key string, exp float64, args ...interface{}) {
	t.Helper()
	args = mArgs(args...)

	got, ok := zrr.GetFloat64(err, key)
	require.True(t, ok, "expected key '%s' is present", key)
	assert.Exactly(t, exp, got, "expected %s=%f got %s=%f", key, exp, key, got)
}

// AssertTime asserts err is instance of zrr.Error and has key with value exp.
func AssertTime(t *testing.T, err error, key string, exp time.Time, args ...interface{}) {
	t.Helper()
	args = mArgs(args...)

	got, ok := zrr.GetTime(err, key)
	require.True(t, ok, "expected key '%s' is present", key)
	assert.Exactly(t, exp, got, "expected %s='%s' got %s='%s'", key, exp, key, got)
}

// AssertBool asserts err is instance of zrr.Error and has key with value exp.
func AssertBool(t *testing.T, err error, key string, exp bool, args ...interface{}) {
	t.Helper()
	args = mArgs(args...)

	got, ok := zrr.GetBool(err, key)
	require.True(t, ok, "expected key '%s' is present", key)
	assert.Exactly(t, exp, got, "expected %s=%v got %s=%v", key, exp, key, got)
}

// msg builds and returns assertion message.
func msg(args ...interface{}) string {
	if len(args) == 0 {
		return ""
	}

	var msg string
	var uses int

	format := args[0]
	if smsg, ok := format.(string); ok {
		if uses = strings.Count(smsg, "%"); uses > 0 {
			nargs := args[1 : 1+uses]
			msg = fmt.Sprintf(smsg, nargs...)
		} else {
			msg += smsg
		}
	} else {
		msg += fmt.Sprintf("%+v", args[0])
	}

	for i := uses + 1; i < len(args); i++ {
		msg += fmt.Sprintf(".%+v", args[i])
	}
	return msg
}

// mArgs returns message as a first item in empty interfaces slice.
func mArgs(args ...interface{}) []interface{} {
	return []interface{}{msg(args...)}
}
