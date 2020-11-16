// Package zrr provides a way to add and inspect type safe error context.
//
// The error context might be useful for example when logging errors which were
// created in some deeper parts of your code.
//
package zrr

import (
	"errors"
	"fmt"
	"reflect"
	"sort"
	"strings"
	"time"
)

// KCode represents the key name used for error code.
const KCode = "code"

// MsgSep is the separator between error message and key value pairs.
const MsgSep = " :: "

// Wrap wraps err in Error instance. It returns nil if err is nil.
func Wrap(err error) *Error {
	if err == nil {
		return nil
	}
	if e, ok := err.(*Error); ok {
		return e
	}
	return base(err, false)
}

// Error represents an error with metadata key value pairs.
type Error struct {
	// Wrapped error.
	error

	// Is error immutable.
	// The immutable error instance is never being changed.
	imm bool

	// Key value metadata associated with the error.
	meta map[string]interface{}
}

// New is a constructor returning new Error instance.
func New(msg string, code ...string) *Error {
	e := base(errors.New(msg), false)
	if len(code) > 0 {
		e = e.Code(code[0])
	}
	return e
}

// Newf is a constructor returning new Error instance.
// Arguments are handled in the same manner as in fmt.Errorf.
func Newf(msg string, args ...interface{}) *Error {
	return base(fmt.Errorf(msg, args...), false)
}

// Imm is a constructor returning new immutable Error instance.
// Usually used for package level errors. Error code is optional, if more
// then one code is provided the first one will be used.
func Imm(msg string, code ...string) *Error {
	e := base(errors.New(msg), true)
	e.imm = true
	if len(code) > 0 {
		e.meta[KCode] = code[0]
	}
	return e
}

// base is a base constructor for Error.
func base(err error, imm bool) *Error {
	return &Error{
		error: err,
		imm:   imm,
		meta:  make(map[string]interface{}, 1),
	}
}

// Error implements error interface and returns error message and key value
// pairs associated with it separated by MsgSep.
func (e *Error) Error() string { return e.msg(true) }

// Cause returns error message without key value pairs.
func (e *Error) Cause() string { return e.msg(false) }

// Code adds error code to the error.
func (e *Error) Code(c string) *Error { return e.with(KCode, c) }

// Str adds the key with string val to the error.
func (e *Error) Str(key string, s string) *Error { return e.with(key, s) }

// Int adds the key with integer val to the error.
func (e *Error) Int(key string, i int) *Error { return e.with(key, i) }

// Int64 adds the key with int64 val to the error.
func (e *Error) Int64(key string, i int64) *Error { return e.with(key, i) }

// Float64 adds the key with float64 val to the error.
func (e *Error) Float64(key string, f float64) *Error { return e.with(key, f) }

// Time adds the key with val as a time to the error.
func (e *Error) Time(key string, t time.Time) *Error { return e.with(key, t) }

// Bool adds the key with val as a boolean to the error.
func (e *Error) Bool(key string, b bool) *Error { return e.with(key, b) }

// FieldsFrom set fields from src.
func (e *Error) FieldsFrom(src Fielder) *Error {
	// Handle immutable error.
	if e.imm {
		ne := base(e, false)
		if HasKey(e, KCode) {
			ne.meta[KCode] = e.meta[KCode]
		}
		return src.ZrrFields(ne)
	}

	return src.ZrrFields(e)
}

// ZrrFields implements Fielder interface.
func (e *Error) ZrrFields(err *Error) *Error {
	for k, v := range e.meta {
		_ = err.with(k, v)
	}
	return err
}

// with adds context to the error.
func (e *Error) with(key string, v interface{}) *Error {
	// Handle immutable error.
	if e.imm {
		ne := base(e, false)
		if HasKey(e, KCode) {
			ne.meta[KCode] = e.meta[KCode]
		}
		ne.meta[key] = v
		return ne
	}

	e.meta[key] = v
	return e
}

// Unwrap unwraps original error.
func (e *Error) Unwrap() error { return e.error }

// msg constructs error message. If meta is set to false it will return only
// error message without metadata.
func (e *Error) msg(meta bool) string {
	var msg string
	var w *Error
	if errors.As(e.error, &w) {
		msg = w.error.Error()
	} else {
		msg = e.error.Error()
	}

	// Return only error message.
	if len(e.meta) == 0 || !meta {
		return msg
	}

	// Sort metadata keys.
	keys := make([]string, 0, len(e.meta))
	for fn := range e.meta {
		keys = append(keys, fn)
	}
	sort.Strings(keys)

	// Construct metadata based on sorted keys.
	parts := make([]string, 0, len(e.meta))
	for _, fn := range keys {
		val := e.meta[fn]
		switch v := val.(type) {
		case string:
			// If value is a string escape quotes.
			val = `"` + strings.ReplaceAll(v, `"`, `\"`) + `"`
		case time.Time:
			val = v.Format(time.RFC3339Nano)
		}
		parts = append(parts, fn+"="+fmt.Sprintf("%v", val))
	}

	if len(parts) > 0 {
		return msg + MsgSep + strings.Join(parts, " ")
	}

	return msg
}

// Fields returns metadata iterator. Caller must not hold to the iterator
// longer then it is necessary to loop over metadata key-value pairs.
func (e *Error) Fields() *iter { return newIter(e) }

// isNil returns true if a is nil or a is nil interface.
func isNil(a interface{}) bool {
	defer func() { recover() }()
	return a == nil || reflect.ValueOf(a).IsNil()
}
