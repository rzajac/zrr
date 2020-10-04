// Package zrr provides errors with metadata as key value pairs.
//
// The errors.Wrap function returns a new error to which we can add metadata.
// For example
//
//     _, err := ioutil.ReadAll(r)
//     if err != nil {
//             return zrr.Wrap(err).Str("origin", "67b75223ad8c2183")
//     }
//
package zrr

import (
	"errors"
	"fmt"
	"sort"
	"strings"
	"time"
)

// KCode represents key name for storing error code.
const KCode = "code"

// Fielder represents an interface
type Fielder interface {
	// ZrrFields sets fields on Error instance.
	ZrrFields(e *Error) *Error
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

// Immutable returns true if the error is immutable.
func (e *Error) Immutable() bool { return e.imm }

// Error implements error interface.
func (e *Error) Error() string { return e.error.Error() }

// String implements fmt.Stringer interface.
//
// The returned key value pairs will be in alphabetical order:
// `error message :: aaa=123 bbb="string value" ccc=true`.
func (e *Error) String() string {
	msg := e.error.Error()
	if len(e.meta) == 0 {
		return msg
	}

	keys := make([]string, 0, len(e.meta))
	for fn := range e.meta {
		keys = append(keys, fn)
	}
	sort.Strings(keys)

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

	div := " "
	if len(parts) > 0 {
		div = " :: "
	}

	return msg + div + strings.Join(parts, " ")
}

// Wrap wraps error in Error instance. It returns nil if err is nil.
func Wrap(err error) *Error {
	if err == nil {
		return nil
	}
	if e, ok := err.(*Error); ok {
		return e
	}
	return base(err, false)
}

// Code adds the key KCode with string val to the *Error metadata.
func (e *Error) Code(val string) *Error { return e.with(KCode, val) }

// HasCode returns true if error has KCode key and it equals code.
func (e *Error) HasCode(code string) bool {
	if val, ok := e.meta[KCode]; ok {
		return val == code
	}
	return false
}

// Str adds the key with string val to the *Error metadata.
func (e *Error) Str(key string, val string) *Error {
	return e.with(key, val)
}

// GetStr returns the key as a string. If key does not exist
// or it is not a string it will return false as the second return value.
func (e *Error) GetStr(key string) (string, bool) {
	if val, ok := e.meta[key]; ok {
		if ret, ok := val.(string); ok {
			return ret, true
		}
	}
	return "", false
}

// Int adds the key with integer val to the *Error metadata.
func (e *Error) Int(key string, val int) *Error {
	return e.with(key, val)
}

// GetInt returns the key as an integer. If key does not exist
// or it is not an integer it will return false as the second return value.
func (e *Error) GetInt(key string) (int, bool) {
	if val, ok := e.meta[key]; ok {
		if ret, ok := val.(int); ok {
			return ret, true
		}
	}
	return 0, false
}

// Float64 adds the key with float64 val to the *Error metadata.
func (e *Error) Float64(key string, val float64) *Error {
	return e.with(key, val)
}

// GetFloat64 returns the key as a float64. If key does not exist
// or it is not a float64 it will return false as the second return value.
func (e *Error) GetFloat64(key string) (float64, bool) {
	if val, ok := e.meta[key]; ok {
		if ret, ok := val.(float64); ok {
			return ret, true
		}
	}
	return 0, false
}

// Time adds the key with val as a time to the *Error metadata.
func (e *Error) Time(key string, val time.Time) *Error {
	return e.with(key, val)
}

// GetTime returns the key as a time.Time. If key does not exist
// or it is not a time.Time it will return false as the second return value.
func (e *Error) GetTime(key string) (time.Time, bool) {
	if val, ok := e.meta[key]; ok {
		if ret, ok := val.(time.Time); ok {
			return ret, true
		}
	}
	return time.Time{}, false
}

// Bool adds the key with val as a boolean to the *Error metadata.
func (e *Error) Bool(key string, val bool) *Error {
	return e.with(key, val)
}

// GetBool returns the key as a boolean. If key does not exist
// or it is not a boolean it will return false as the second return value.
func (e *Error) GetBool(key string) (bool, bool) {
	if val, ok := e.meta[key]; ok {
		if ret, ok := val.(bool); ok {
			return ret, true
		}
	}
	return false, false
}

// FieldsFrom set fields from src.
func (e *Error) FieldsFrom(src Fielder) *Error {
	// Handle immutable error.
	if e.imm {
		ne := base(e, false)
		if e.HasKey(KCode) {
			ne.meta[KCode] = e.meta[KCode]
		}
		return src.ZrrFields(ne)
	}

	return src.ZrrFields(e)
}

// with adds context to the error.
func (e *Error) with(key string, val interface{}) *Error {
	// Handle immutable error.
	if e.imm {
		ne := base(e, false)
		if e.HasKey(KCode) {
			ne.meta[KCode] = e.meta[KCode]
		}
		ne.meta[key] = val
		return ne
	}

	e.meta[key] = val
	return e
}

// Unwrap unwraps original error.
func (e *Error) Unwrap() error {
	return e.error
}

// HasKey returns true if error has the key set.
func (e *Error) HasKey(key string) bool {
	_, ok := e.meta[key]
	return ok
}
