package forth

import "errors"

var (
	// ErrBadState reports bad VM states
	ErrBadState = errors.New("bad VM state")

	// ErrUnderflow reports stack underflow
	ErrUnderflow = errors.New("stack underflow")

	// ErrArgument reports a bad argument to an operation
	ErrArgument = errors.New("bad argument")
)
