// SPDX-License-Identifier: MIT

package forth

import (
	"testing"
)

func TestLocalsBasic(t *testing.T) {
	// : add3 (| a b c |) a b c + + ;
	// 1 2 3 add3 -> 6
	tstRunForth(t, "Basic Locals", ": add3 (| a b c |) a b c + + ; 1 2 3 add3", 6)
	tstRunForth(t, "Basic Locals Alternate", ": add3 (| a b c | ) a b c + + ; 1 2 3 add3", 6)
	tstRunForth(t, "Basic Locals Alternate 2", ": add3 (| a b | c ) c! a b c + + ; 1 2 3 add3", 6)
}

func TestLocalsModifying(t *testing.T) {
	// : inc (| a |) a 1 + a! a ;
	// 5 inc -> 6
	tstRunForth(t, "Modifying", ": inc (| a |) a 1 + a! a ; 5 inc", 6)
}

func TestLocalsUninitialized(t *testing.T) {
	// : test (| a | b ) a b ;
	// 1 test -> 1 nil
	tstRunForth(t, "Uninitialized", ": test (| a | b ) a b ; 1 test", 1, nil)

	// : test2 (| a | b ) 5 b! a b ;
	// 1 test2 -> 1 5
	tstRunForth(t, "InitializedLater", ": test2 (| a | b ) 5 b! a b ; 1 test2", 1, 5)
}

func TestLocalsShadowing(t *testing.T) {
	// Variable X
	// : test (| X |) X ;
	// 5 test -> 5.
	// X @ -> 0.
	tstRunForth(t, "Shadowing", "variable X : test (| X |) X ; 5 test X @", 5, 0)
}

func TestLocalsRecursion(t *testing.T) {
	// Factorial
	// : fact (| n |) n 1 <= IF 1 ELSE n n 1 - recur * THEN ;
	tstRunForth(t, "Recursion", ": fact (| n |) n 1 <= IF 1 ELSE n n 1 - recur * THEN ; 5 fact", 120)
}

func TestLocalsMixedControl(t *testing.T) {
	// Test locals inside IF (deferred definition)
	// : test 1 IF (| a |) a THEN ;
	tstRunForth(t, "MixedControl", ": test 1 IF (| a |) a THEN ; 10 test", 10)
}

func TestTailCall(t *testing.T) {
	// : tc-count (| n |) n n 1 - dup 0 > IF (tail-call) THEN drop ;
	// 5 tc-count -> 5 4 3 2 1
	tstRunForth(t, "TailCall", ": tc-count (| n |) n n 1 - dup 0 > IF (tail-call) THEN drop ; 5 tc-count", 5, 4, 3, 2, 1)
}
