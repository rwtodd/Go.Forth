// SPDX-License-Identifier: MIT

package forth

import "errors"

var (
	// ErrBadState reports bad VM states
	ErrBadState = errors.New("bad VM state")

	// ErrUnderflow reports stack underflow
	ErrUnderflow = errors.New("stack underflow")

	// ErrArgument reports a bad argument to an operation
	ErrArgument = errors.New("bad argument")

	// ErrRStackUnderflow reports when the Rstack is too low
	ErrRStackUnderflow = errors.New("r-stack underflow")

	// ErrIndexOutOfBounds reports array index out of bounds
	ErrIndexOutOfBounds = errors.New("index out of bounds")

	// ErrKeyNotFound reports when a dictionary key is not found
	ErrKeyNotFound = errors.New("key not found")
)
