// SPDX-License-Identifier: MIT

package forth

import "testing"

func TestStackVariable(t *testing.T) {
	// 10 "x" (variable)
	// x @ -> 10
	tstRunForth(t, "StackVariable", `10 " x" (variable) x @`, 10)
	tstRunForth(t, "StackVariableUC", `10 " X" (variable) x @ X @`, 10, 10)
}

func TestStackConstant(t *testing.T) {
	// 20 "y" (constant)
	// y -> 20
	tstRunForth(t, "StackConstant", `20 " y" (constant) y`, 20)
	tstRunForth(t, "StackConstantUC", `20 " Y" (constant) y Y`, 20, 20)
}

func TestStackVariableDoes(t *testing.T) {
	// : adder 10 [ @ + ] "add10" (variable-does) ;
	// adder
	// 5 add10 -> 15
	code := `
    : adder 10 [ @ + ] " Add10" (variable-does) ;
    adder
    5 add10
    `
	tstRunForth(t, "StackVariableDoes", code, 15)
}
