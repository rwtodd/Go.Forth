// SPDX-License-Identifier: MIT

package forth

import (
	"errors"
	"io"
	"strings"
	"testing"
)

var vm = NewVM()

// tstRunForth is the main test helper.. it lets you
// run a code string, and FAILS if there is an error or
// if the stack doesn't match the values passed in.
func tstRunForth(t *testing.T, name string, code string, vals ...any) {
	t.Helper()
	t.Run(name, func(t *testing.T) {
		vm.ClearResetState()
		if markerr := mark(vm); markerr != nil {
			t.Errorf("Unexpected error: %v", markerr)
		}
		tprog := strings.NewReader(code)
		if err := vm.RunFromSource(tprog, "test", io.Discard); err != nil {
			t.Errorf("Unexpected error: %v", err)
		}
		if !stackEq(vals...) {
			t.Errorf("Stack mismatch. Expected: %v, Got: %v", vals, vm.Stack)
		}
		if forgeterr := forget(vm); forgeterr != nil {
			t.Errorf("Unexpected error: %v", forgeterr)
		}
	})
}

// tstRunForthErr runs the code, expecting a specific error.
// It fails if the error doesn't match or the stack doesn't match.
func tstRunForthErr(t *testing.T, name string, code string, expectedErr error, vals ...any) {
	t.Helper()
	t.Run(name, func(t *testing.T) {
		vm.ClearResetState()
		tprog := strings.NewReader(code)
		if markerr := mark(vm); markerr != nil {
			t.Errorf("Unexpected error: %v", markerr)
		}
		if err := vm.RunFromSource(tprog, "test", io.Discard); !errors.Is(err, expectedErr) {
			t.Errorf("Error mismatch. Expected: %v, Got: %v", expectedErr, err)
		}
		if !stackEq(vals...) {
			t.Errorf("Stack mismatch. Expected: %v, Got: %v", vals, vm.Stack)
		}
		if forgeterr := forget(vm); forgeterr != nil {
			t.Errorf("Unexpected error: %v", forgeterr)
		}
	})
}

// stackEq is a helper function checking the stack contents
// against the arguments.
func stackEq(vals ...any) bool {
	if len(vals) != len(vm.Stack) {
		return false
	}

	for i := range vals {
		expected := vals[i]
		if e, ok := expected.(int); ok {
			expected = int64(e)
		}
		if expected != vm.Stack[i] {
			return false
		}
	}
	return true
}
