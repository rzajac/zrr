package zrr

// Fielder represents an interface for setting metadata on Error instance.
type Fielder interface {
	// ZrrFields sets fields on Error instance.
	ZrrFields(e *Error) *Error
}
