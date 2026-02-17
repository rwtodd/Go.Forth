// SPDX-License-Identifier: MIT

package forth

import (
	"testing"
)

func TestVariableDoesBasic(t *testing.T) {
	// 5 constant x -> x should be 5
	// internally: 5 ['] @ variable-does x
	// constant is now implemented via VariableWord with @ XT
	tstRunForth(t, "ConstantViaVarDoes", `5 constant x x`, 5)

	// Custom: 3.14 ['] @ variable-does pi
	tstRunForth(t, "ManualConstant", `3.14 ' @ variable-does pi pi`, 3.14)
}

func TestVariableDoesClosure(t *testing.T) {
	// : mk-counter 0 [ @ 1 + dup rot ! ] variable-does ;
	// mk-counter x
	// x -> 1
	// x -> 2
	code := `
	: mk-counter 0 [ dup @ 1 + dup rot ! ] variable-does ;
	mk-counter x
	x x
	`
	tstRunForth(t, "ClosureCounter", code, 1, 2)
}

func TestVariableDoesClosureNoLocals(t *testing.T) {
	// A simpler closure that just adds to the value without rotating (using dup)
	// : adder 10 [ @ + ] variable-does ;
	// adder add10
	// 5 add10 -> 15
	code := `
    : adder 10 [ @ + ] variable-does ;
    adder add10
    5 add10
    `
	tstRunForth(t, "ClosureAdder", code, 15)
}
