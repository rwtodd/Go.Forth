// SPDX-License-Identifier: MIT

package forth

import (
	"errors"
	"io"
	"strings"
	"testing"
)

// Helper to ignore stack check on error (since slices panic on comparison)
func tstRunForthErrIgnoreStack(t *testing.T, name string, code string, expectedErr error) {
	t.Helper()
	t.Run(name, func(t *testing.T) {
		vm.ResetState()
		tprog := strings.NewReader(code)
		if markerr := mark(vm); markerr != nil {
			t.Errorf("Unexpected error: %v", markerr)
		}
		if err := vm.Run(tprog, io.Discard); !errors.Is(err, expectedErr) {
			t.Errorf("Error mismatch. Expected: %v, Got: %v", expectedErr, err)
		}
		if forgeterr := forget(vm); forgeterr != nil {
			t.Errorf("Unexpected error: %v", forgeterr)
		}
	})
}

func TestVPush(t *testing.T) {
	// 1. Basic push to empty slice
	tstRunForth(t, "VPushInts", `0 ints << 1 2 3 >> <@push> dup @len swap 0 @`, 3, 1)

	// 2. Push to existing slice
	tstRunForth(t, "VPushAppend", `2 ints 42 over 0 ! 43 over 1 ! << 44 45 >> <@push> dup @len swap 2 @`, 4, 44)

	// 3. Float slice
	tstRunForth(t, "VPushFloats", `0 floats << 1.1 2.2 >> <@push> dup @len swap 0 @`, 2, 1.1)

	// 4. String slice
	tstRunForth(t, "VPushStrings", `0 strings << " a" " b" >> <@push> dup @len swap 1 @`, 2, "b")

	// 5. Byte slice
	// Note: byte array push typically requires int values in range? Or just int?
	// arrays.go: "byte array push expects integer"
	tstRunForth(t, "VPushBytes", `0 bytes << 65 66 >> <@push> dup @len swap 0 c@`, 2, 65)
}

func TestVMake(t *testing.T) {
	// 1. >> <ints>
	tstRunForth(t, "VMakeInts", `<< 10 20 30 >> <ints> dup @len swap 1 @`, 3, 20)

	// 2. >> <floats>
	tstRunForth(t, "VMakeFloats", `<< 1.5 2.5 >> <floats> dup @len swap 0 @`, 2, 1.5)

	// 3. >> <strings>
	tstRunForth(t, "VMakeStrings", `<< " hello" " world" >> <strings> dup @len swap 0 @`, 2, "hello")
}

func TestVariableVPush(t *testing.T) {
	// 1. Variable with ints
	tstRunForth(t, "VarVPush", `0 variable v 0 ints v ! v << 1 2 3 >> <@push> drop v @len`, 3)

	// 2. Variable check values
	tstRunForth(t, "VarVPushVal", `0 variable v 0 ints v ! v << 10 20 >> <@push> drop v 1 @`, 20)

	// 3. Stack effect: v << ... >> <@push> -> return var
	// In the test: "variable v" pushes *Variable. "0 ints v !" stores []int{} in it.
	// "v << 99 >> <@push>" pushes *Variable (same one).
	// "v @" gets the slice? No, "v" pushes var address.
	// Wait, "v @len". "arrayLen" handles *Variable by getting value.
	// So if stack has *Variable, "@len" works.
	// So "v << 99 >> <@push>" leaves *Variable on stack.
	// Then "@len" consumes it and returns length.
	tstRunForth(t, "VarStackLen", `0 variable v 0 ints v ! v << 99 >> <@push> @len`, 1)
}

func TestVPushErrors(t *testing.T) {
	// Test type mismatch
	tstRunForthErrIgnoreStack(t, "VPushIntError", `0 ints << 1.5 >> <@push>`, ErrArgument)
	tstRunForthErrIgnoreStack(t, "VPushFloatError", `0 floats << " s" >> <@push>`, ErrArgument)
	tstRunForthErrIgnoreStack(t, "VPushStringError", `0 strings << 1 >> <@push>`, ErrArgument)

	// Test vmake errors
	tstRunForthErrIgnoreStack(t, "VMakeIntError", `<< " s" >> <ints>`, ErrArgument)
}

func TestVMakeMixed(t *testing.T) {
	// 4. >> <things> (any)
	tstRunForth(t, "VMakeAny", `<< 1 " hi" 3.14 >> <things> dup @len swap dup 1 @ swap 0 @`, 3, "hi", 1)

	// 5. >> <bytes> (byte)
	tstRunForth(t, "VMakeBytes", `<< 65 66 67 >> <bytes> dup @len swap 0 c@`, 3, 65)

	// Test VMakeAny with variable push (append mixed types)
	tstRunForth(t, "VPushAnyVar", `0 variable v 0 things v ! v << " foo" 42 >> <@push> drop v 1 @ v 0 @`, 42, "foo")
}
