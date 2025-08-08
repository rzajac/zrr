// Package zrr provides a way to add and inspect type safe error context.
//
// The error context might be useful for example when logging errors which were
// created in some deeper parts of your code.
package zrr

import (
	"encoding/json"
	"errors"
	"fmt"
	"time"
)

// MetadataGetter is an interface wrapping GetMetadata method.
type MetadataGetter interface {
	// GetMetadata returns error metadata. The returned metadata map should be
	// considered read-only.
	GetMetadata() map[string]interface{}
}

// MetadataErrSetter is an interface wrapping SetErrMetadata method.
type MetadataErrSetter interface {
	// SetErrMetadata sets error metadata. The returned instance might be
	// different from the one this method is called if the error is immutable.
	SetErrMetadata(map[string]interface{}) *Error
}

// ECInvJSON represents invalid JSON error code.
const ECInvJSON = "ECInvJSON"

// ErrInvJSON represents package level error indicating JSON structure or
// format error.
var ErrInvJSON = Imm("invalid JSON", ECInvJSON)

// Wrap wraps err in Error instance. It returns nil if err is nil.
func Wrap(err error, code ...string) *Error {
	if err == nil {
		return nil
	}
	if e, ok := err.(*Error); ok {
		if len(code) > 0 {
			return e.setCode(code[0])
		}
		return e
	}
	return base(err, false, code...)
}

// Error represents an error with metadata key value pairs.
type Error struct {
	// Wrapped error.
	error

	// Error code.
	code string

	// Is error immutable.
	// The immutable error instance is never being changed.
	imm bool

	// Key value metadata associated with the error.
	meta map[string]interface{}
}

// New is a constructor returning new Error instance.
func New(msg string, code ...string) *Error {
	return base(errors.New(msg), false, code...)
}

// Newf is a constructor returning new Error instance.
// Arguments are handled in the same manner as in fmt.Errorf.
func Newf(msg string, args ...interface{}) *Error {
	return base(fmt.Errorf(msg, args...), false)
}

// Imm is a constructor returning new immutable Error instance.
//
// Immutable error instances are never changed when adding / changing fields.
// They are good choice for package level errors.
//
// Error code is optional, if more than one code is provided the first
// one will be used.
func Imm(msg string, code ...string) *Error {
	return base(errors.New(msg), true, code...)
}

// base is a base constructor for Error.
func base(err error, imm bool, code ...string) *Error {
	return &Error{
		error: err,
		code:  fistCode(code...),
		imm:   imm,
		meta:  make(map[string]interface{}),
	}
}

// Error implements error interface and returns error message and key value
// pairs associated with it separated by MsgSep.
func (e *Error) Error() string { return e.error.Error() }

// ErrCode returns error code.
func (e *Error) ErrCode() string { return e.code }

// setCode sets error code to the error.
func (e *Error) setCode(c string) *Error {
	if e.imm {
		return base(e, false, c).SetErrMetadata(e.GetMetadata())
	}
	e.code = c
	return e
}

// Str adds the key with string val to the error.
func (e *Error) Str(key string, s string) *Error { return e.with(key, s) }

// StrAppend appends the string s (prefixed with semicolon) to the string
// represented by key k. The key will be added if it doesn't exist. If the
// key already exists and is not a string the old key will be overwritten.
func (e *Error) StrAppend(key string, s string) *Error {
	if si, ok := e.meta[key]; ok {
		if ss, ok := si.(string); ok {
			s = ss + ";" + s
		}
	}
	return e.with(key, s)
}

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

// SetErrMetadata sets error metadata. The returned instance might be
// different from the one this method is called if the error is immutable.
func (e *Error) SetErrMetadata(src map[string]interface{}) *Error {
	ee := e
	for k, v := range src {
		ee = ee.with(k, v)
	}
	return ee
}

// SetMetadataFrom is a convenience method setting metadata from an instance
// which implements MetadataGetter.
func (e *Error) SetMetadataFrom(src MetadataGetter) *Error {
	return e.SetErrMetadata(src.GetMetadata())
}

// GetMetadata returns error metadata. The returned metadata map should be
// considered read-only.
func (e *Error) GetMetadata() map[string]interface{} { return e.meta }

// MetaAll returns error metadata. The returned metadata map should be
// considered read-only.
func (e *Error) MetaAll() map[string]any { return e.meta }

// with adds context to the error.
func (e *Error) with(key string, v interface{}) *Error {
	// Handle immutable error.
	if e.imm {
		ne := base(e, false, e.code)
		ne.meta[key] = v
		return ne
	}

	e.meta[key] = v
	return e
}

// Unwrap unwraps original error.
func (e *Error) Unwrap() error { return e.error }

func (e *Error) MarshalJSON() ([]byte, error) {
	m := map[string]interface{}{
		"error": e.Error(),
		"code":  e.code,
	}
	if len(e.meta) > 0 {
		m["meta"] = e.meta
	}
	return json.Marshal(m)
}

// UnmarshalJSON unmarshal error's JSON representation.
// Notes:
//   - all metadata numeric values will be unmarshalled as float64
func (e *Error) UnmarshalJSON(data []byte) error {
	m := make(map[string]interface{}, 3)

	if err := json.Unmarshal(data, &m); err != nil {
		return err
	}

	msgI, _ := m["error"]
	msg, _ := msgI.(string)
	if msg == "" {
		return ErrInvJSON
	}

	codeI, _ := m["code"]
	code, _ := codeI.(string)

	metaI, _ := m["meta"]
	var meta map[string]interface{}
	if metaI != nil {
		meta, _ = metaI.(map[string]interface{})
	}
	if meta == nil {
		meta = make(map[string]interface{})
	}

	e.error = errors.New(msg)
	e.code = code
	e.meta = meta
	return nil
}

// isNil returns true if v is nil or v is nil interface.
// func isNil(v interface{}) bool {
//	defer func() { recover() }()
//	return v == nil || reflect.ValueOf(v).IsNil()
// }

// firstCode returns first code from the slice.
func fistCode(code ...string) string {
	if len(code) > 0 {
		return code[0]
	}
	return ""
}
