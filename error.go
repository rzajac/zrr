// Package ero provides errors with metadata.
package zrr

import (
	"errors"
	"fmt"

	"github.com/rs/zerolog"
)

// FCode error code field name.
const FCode = "code"

// FSevere severe error field name.
const FSevere = "severe"

// Error represents an error with metadata fields.
type Error struct {
	// Wrapped error.
	error

	// Is error immutable.
	imm bool

	// Metadata associated with an error.
	meta map[string]interface{}
}

// New is a constructor returning new Error instance.
func New(msg string, args ...interface{}) *Error {
	return base(fmt.Errorf(msg, args...))
}

// Imm is a constructor returning new immutable Error.
// Usually used for package level errors. Error code is optional, if more
// then one code is set only the first one will be used.
func Imm(msg string, code ...string) *Error {
	e := base(fmt.Errorf(msg))
	e.imm = true
	if len(code) > 0 {
		e.meta[FCode] = code[0]
	}
	return e
}

func base(err error) *Error {
	e := &Error{
		error: err,
		imm:   false,
		meta:  make(map[string]interface{}, 1),
	}
	return e
}

// Wrap wraps error in Error instance. It returns nil if err is nil.
func Wrap(err error) *Error {
	if err == nil {
		return nil
	}
	if e, ok := err.(*Error); ok {
		return e
	}
	if E := findError(err); E != nil {
		E2 := base(err)
		// Copy fields from immutable error.
		for k, v := range E.meta {
			E2.meta[k] = v
		}
		return E2
	}
	return base(err)
}

// Error implements error interface.
func (e *Error) Error() string { return e.error.Error() }

// String implements fmt.Stringer interface. The returned string will have
// message followed by space and key value pairs in alphabetical order:
// "error message -- place:123 key:value severe:true".
func (e *Error) String() string {
	msg := e.error.Error()
	fields := ""
	div := ""
	for k, v := range e.meta {
		fields += fmt.Sprintf(" %s:%v", k, v)
	}
	if fields != "" {
		div = " ---"
	}
	return msg + div + fields
}

// With adds context field to the error.
func (e *Error) With(key string, val interface{}) *Error {
	if e.imm {
		// Error code and severity on immutable cannot be changed.
		if key == FCode || key == FSevere {
			return e
		}
		if _, ok := e.error.(*Error); ok {
			// We already wrapped original error.
			e.meta[key] = val
		} else {
			// We wrap original immutable error so errors.Is works properly.
			ne := base(e)
			// Copy fields from immutable error.
			for k, v := range e.meta {
				ne.meta[k] = v
			}
			ne.meta[key] = val
			return ne
		}
	}
	e.meta[key] = val
	return e
}

// Unwrap unwraps original error.
func (e *Error) Unwrap() error { return e.error }

// Cause returns underlying error.
func (e *Error) Cause() error { return e.error }

// Clone returns deep clone of Error instance.
// WARNING: Cloned error is no longer immutable.
func (e *Error) Clone() *Error {
	ne := &Error{
		error: e.error,
		meta:  make(map[string]interface{}, len(e.meta)),
	}
	for k, v := range e.meta {
		ne.meta[k] = v
	}
	return ne
}

func (e *Error) MarshalZerologObject(evt *zerolog.Event) {
	evt.Str(zerolog.MessageFieldName, e.Error())
	evt.Fields(e.meta)
}

// findError finds first instance of Error in the chain.
func findError(err error) *Error {
	if e, ok := err.(*Error); ok {
		return e
	}

	var ret *Error
	for err != nil {
		if err = errors.Unwrap(err); err == nil {
			return nil
		}

		if e, ok := err.(*Error); ok {
			return e
		}
	}
	return ret
}
