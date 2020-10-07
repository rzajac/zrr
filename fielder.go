package zrr

// Fielder represents an interface for setting metadata on Error instance.
type Fielder interface {
	// ZrrFields sets fields on err instance and returns it.
	ZrrFields(err *Error) *Error
}
