package forth

import "errors"

var (
	// ErrUnderflow reports stack underflow
	ErrUnderflow = errors.New("stack underflow")

	// ErrArgument reports a bad argument to an operation
	ErrArgument = errors.New("bad argument")
)
