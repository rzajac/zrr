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

// AssertCause asserts err is instance of zrr.Error and has error
// message (without key value pairs) equal to cause.
func AssertCause(t *testing.T, err error, cause string, args ...interface{}) {
	t.Helper()
	args = mArgs(args...)

	require.NotNil(t, err, m("err=nil", args...))
	E, ok := err.(*zrr.Error)
	if !ok {
		t.Error("expected err to be instance of zrr.Error")
		return
	}

	got := E.Cause()
	if got != cause {
		t.Errorf("expected cause '%s' but got '%s'", cause, got)
		return
	}
}

// AssertContains asserts err is instance of zrr.Error and has error
// message (without key value pairs) which contains cause string.
func AssertContains(t *testing.T, err error, cause string, args ...interface{}) {
	t.Helper()
	args = mArgs(args...)

	require.NotNil(t, err, m("err=nil", args...))
	E, ok := err.(*zrr.Error)
	if !ok {
		t.Error("expected err to be instance of zrr.Error")
		return
	}

	got := E.Cause()
	assert.Contains(t, got, cause)
}

// AssertCode asserts err is instance of zrr.Error and has error code exp.
func AssertCode(t *testing.T, err error, exp string, args ...interface{}) {
	t.Helper()
	args = mArgs(args...)
	require.NotNil(t, err, m("err=nil", args...))
	AssertStr(t, err, zrr.KCode, exp)
}

// AssertEqual asserts err and got are instance of zrr.Error and their error
// messages (with key value pairs) are equal.
func AssertEqual(t *testing.T, exp, got error, args ...interface{}) {
	t.Helper()
	args = mArgs(args...)

	require.NotNil(t, exp, m("exp=nil", args...))
	require.NotNil(t, got, m("got=nil", args...))
	expE, ok := exp.(*zrr.Error)
	if !ok {
		t.Error("expected exp to be instance of zrr.Error")
		return
	}

	gotE, ok := got.(*zrr.Error)
	if !ok {
		t.Error("expected got to be instance of zrr.Error")
		return
	}

	assert.Exactly(t, expE.Error(), gotE.Error())
}

// AssertStr asserts err is instance of zrr.Error and has key with value exp.
func AssertStr(t *testing.T, err error, key, exp string, args ...interface{}) {
	t.Helper()
	args = mArgs(args...)

	require.NotNil(t, err, m("err=nil", args...))
	got, ok := zrr.GetStr(err, key)
	require.True(t, ok, "expected key '%s' is present", key)
	assert.Exactly(t, exp, got, "expected %s='%s' got %s='%s'", key, exp, key, got)
}

// AssertInt asserts err is instance of zrr.Error and has key with value exp.
func AssertInt(t *testing.T, err error, key string, exp int, args ...interface{}) {
	t.Helper()
	args = mArgs(args...)

	require.NotNil(t, err, m("err=nil", args...))
	got, ok := zrr.GetInt(err, key)
	require.True(t, ok, "expected key '%s' is present", key)
	assert.Exactly(t, exp, got, "expected %s=%d got %s=%d", key, exp, key, got)
}

// AssertInt64 asserts err is instance of zrr.Error and has key with value exp.
func AssertInt64(t *testing.T, err error, key string, exp int64, args ...interface{}) {
	t.Helper()
	args = mArgs(args...)

	require.NotNil(t, err, m("err=nil", args...))
	got, ok := zrr.GetInt64(err, key)
	require.True(t, ok, "expected key '%s' is present", key)
	assert.Exactly(t, exp, got, "expected %s=%d got %s=%d", key, exp, key, got)
}

// AssertFloat64 asserts err is instance of zrr.Error and has key with value exp.
func AssertFloat64(t *testing.T, err error, key string, exp float64, args ...interface{}) {
	t.Helper()
	args = mArgs(args...)

	require.NotNil(t, err, m("err=nil", args...))
	got, ok := zrr.GetFloat64(err, key)
	require.True(t, ok, "expected key '%s' is present", key)
	assert.Exactly(t, exp, got, "expected %s=%f got %s=%f", key, exp, key, got)
}

// AssertTime asserts err is instance of zrr.Error and has key with value exp.
func AssertTime(t *testing.T, err error, key string, exp time.Time, args ...interface{}) {
	t.Helper()
	args = mArgs(args...)

	require.NotNil(t, err, m("err=nil", args...))
	got, ok := zrr.GetTime(err, key)
	require.True(t, ok, "expected key '%s' is present", key)
	assert.Exactly(t, exp, got, "expected %s='%s' got %s='%s'", key, exp, key, got)
}

// AssertBool asserts err is instance of zrr.Error and has key with value exp.
func AssertBool(t *testing.T, err error, key string, exp bool, args ...interface{}) {
	t.Helper()
	args = mArgs(args...)

	require.NotNil(t, err, m("err=nil", args...))
	got, ok := zrr.GetBool(err, key)
	require.True(t, ok, "expected key '%s' is present", key)
	assert.Exactly(t, exp, got, "expected %s=%v got %s=%v", key, exp, key, got)
}

// m builds and returns assertion message with prefix name.
func m(name string, args ...interface{}) string {
	msg := msg(args...)
	if msg != "" {
		name = "." + name
	}
	return msg + name
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
