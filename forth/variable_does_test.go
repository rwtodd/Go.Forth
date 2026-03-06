// SPDX-License-Identifier: MIT

package forth

import (
	"testing"
)

func TestVariableDoesBasic(t *testing.T) {
	// 5 " x" constant -> x should be 5
	// internally: 5 ['] @ " x" variable-does
	// constant is now implemented via VariableWord with @ XT
	tstRunForth(t, "ConstantViaVarDoes", `5 " x" constant x`, 5)

	// Custom: 3.14 " @" lookup-xt " pi" variable-does
	tstRunForth(t, "ManualConstant", `3.14 " @" lookup-xt " pi" variable-does pi`, 3.14)
}

func TestVariableDoesClosure(t *testing.T) {
	// : mk-counter 0 [ dup @ 1 + dup rot ! ] " counter" variable-does ;
	// mk-counter
	// counter -> 1
	// counter -> 2
	code := `
	: mk-counter 0 [ dup @ 1 + dup rot ! ] " counter" variable-does ;
	mk-counter
	counter counter
	`
	tstRunForth(t, "ClosureCounter", code, 1, 2)
}

func TestVariableDoesClosureNoLocals(t *testing.T) {
	// A simpler closure that just adds to the value without rotating (using dup)
	// : adder 10 swap [ @ + ] swap " add10" variable-does ;
	// adder
	// 5 add10 -> 15
	code := `
    : adder 10 [ @ + ] " add10" variable-does ;
    adder
    5 add10
    `
	tstRunForth(t, "ClosureAdder", code, 15)
}
