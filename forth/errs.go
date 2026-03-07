// SPDX-License-Identifier: MIT

package forth

import (
	"errors"
)

var (
	// ErrBadState reports bad VM states
	ErrBadState = errors.New("bad VM state")

	// ErrUnderflow reports stack underflow
	ErrUnderflow = errors.New("stack underflow")

	// ErrArgument reports a bad argument to an operation
	ErrArgument = errors.New("bad argument")

	// ErrIndexOutOfBounds reports array index out of bounds
	ErrIndexOutOfBounds = errors.New("index out of bounds")

	// ErrKeyNotFound reports when a dictionary key is not found
	ErrKeyNotFound = errors.New("key not found")

	// ErrUser reports a user-generated error
	ErrUser = errors.New("user exception")
)

// Error wraps a sentinel error with a specific message
type Error struct {
	Err error
	Msg string
}

func (e *Error) Error() string { return e.Msg }
func (e *Error) Unwrap() error { return e.Err }

func ErrBadStateMsg(msg string) error {
	return &Error{Err: ErrBadState, Msg: msg}
}

func ErrUnderflowMsg(msg string) error {
	return &Error{Err: ErrUnderflow, Msg: msg}
}

func ErrArgumentMsg(msg string) error {
	return &Error{Err: ErrArgument, Msg: msg}
}

func ErrIndexOutOfBoundsMsg(msg string) error {
	return &Error{Err: ErrIndexOutOfBounds, Msg: msg}
}

func ErrKeyNotFoundMsg(msg string) error {
	return &Error{Err: ErrKeyNotFound, Msg: msg}
}

func ErrUserMsg(msg string) error {
	return &Error{Err: ErrUser, Msg: msg}
}

// ErrInternalMsg wraps a generic error as an internal error, or just returns a new error
// if we don't have a specific sentinel for it.
// Actually, we said we'd remove bare fmt.Errorf.
// If it doesn't fit the others, maybe we need another category?
// For now, let's assume we can fit everything.
