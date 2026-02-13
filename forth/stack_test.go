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
