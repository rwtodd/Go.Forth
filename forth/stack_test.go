// SPDX-License-Identifier: MIT

package forth

import (
	"testing"
)

func TestDup(t *testing.T) {
	tstRunForth(t, "Dup2", `2 dup 3 dup`, 2, 2, 3, 3)
	tstRunForth(t, "DupFloat", `2.2 dup dup`, 2.2, 2.2, 2.2)
}

func TestSwap(t *testing.T) {
	tstRunForth(t, "Swap2", `2 3 swap`, 3, 2)
	tstRunForth(t, "SwapStrInt", `" hi"  43  swap`, 43, "hi")
}

func TestOver(t *testing.T) {
	tstRunForth(t, "OverBasic", `2 3 OVER `, 2, 3, 2)
	tstRunForthErr(t, "OverUnderflow", `over`, ErrUnderflow)
}

func TestNip(t *testing.T) {
	tstRunForth(t, "NipChain", `1 2 3 nip 4 nip 5 nip`, 1, 5)
}

func TestTuck(t *testing.T) {
	tstRunForth(t, "TuckChain", `1 2 tuck 3 tuck 4 tuck`,
		2, 1, 3, 2, 4, 3, 4)
}

func TestDrop(t *testing.T) {
	tstRunForth(t, "DropBasic", `2 3 drop`, 2)
	tstRunForthErr(t, "DropUnderflow", `drop drop`, ErrUnderflow)
}

func TestRot(t *testing.T) {
	tstRunForthErr(t, "RotUnderflow", `3 2  rot`, ErrUnderflow, 3, 2)
	tstRunForth(t, "RotBasic", ` 2 3 4 rot `, 3, 4, 2)
	tstRunForth(t, "RotMinus", ` 2 3.001 4 -rot `, 4, 2, 3.001)
	tstRunForth(t, "RotRotRot", ` 2 3.001 4 rot rot `, 4, 2, 3.001)
}

func TestPick(t *testing.T) {
	// Equivalencies
	// 0 pick == dup
	tstRunForth(t, "Pick0", `10 20 30 0 pick`, 10, 20, 30, 30)
	tstRunForth(t, "Pick0Dup", `10 20 30 dup`, 10, 20, 30, 30)

	// 1 pick == over
	tstRunForth(t, "Pick1", `10 20 30 1 pick`, 10, 20, 30, 20)
	tstRunForth(t, "Pick1Over", `10 20 30 over`, 10, 20, 30, 20)

	// Deeper
	tstRunForth(t, "Pick2", `10 20 30 2 pick`, 10, 20, 30, 10)

	// Errors
	tstRunForthErr(t, "PickUnderflow", `10 20 2 pick`, ErrUnderflow, 10, 20)
	tstRunForthErr(t, "PickType", `10 0.0 pick`, ErrArgument, 10)
}

func TestRoll(t *testing.T) {
	// Equivalencies
	// 0 roll == nop
	tstRunForth(t, "Roll0", `10 20 30 0 roll`, 10, 20, 30)

	// 1 roll == swap
	tstRunForth(t, "Roll1", `10 20 30 1 roll`, 10, 30, 20)
	tstRunForth(t, "Roll1Swap", `10 20 30 swap`, 10, 30, 20)

	// 2 roll == rot
	tstRunForth(t, "Roll2", `10 20 30 2 roll`, 20, 30, 10)
	tstRunForth(t, "Roll2Rot", `10 20 30 rot`, 20, 30, 10)

	// Deeper
	tstRunForth(t, "Roll3", `10 20 30 40 3 roll`, 20, 30, 40, 10)

	// Errors
	tstRunForthErr(t, "RollUnderflow", `10 20 2 roll`, ErrUnderflow, 10, 20)
}

func TestMinusRoll(t *testing.T) {
	// Equivalencies
	// 1 -roll == swap
	tstRunForth(t, "MinusRoll1", `10 20 30 1 -roll`, 10, 30, 20)

	// 2 -roll == -rot
	tstRunForth(t, "MinusRoll2", `10 20 30 2 -roll`, 30, 10, 20)
	tstRunForth(t, "MinusRoll2Rot", `10 20 30 -rot`, 30, 10, 20)

	// Deeper
	tstRunForth(t, "MinusRoll3", `10 20 30 40 3 -roll`, 40, 10, 20, 30)

	// Errors
	tstRunForthErr(t, "MinusRollUnderflow", `10 20 2 -roll`, ErrUnderflow, 10, 20)
}
