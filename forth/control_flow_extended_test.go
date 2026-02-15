package forth

import (
	"testing"
)

func TestBeginAgain(t *testing.T) {
	// Infinite loop broken by LEAVE
	// Stack trace:
	// 0 -> 1+ -> 1
	// 1 -> 1+ -> 2
	// 2 -> 1+ -> 3 (break)
	// Result: 3
	tstRunForth(t, "BEGIN AGAIN LEAVE", ": test-again 0 begin 1 + dup 3 = if leave then again ; test-again", 3)
}

func TestBeginUntil(t *testing.T) {
	// Loop until condition true
	// 0 -> 1 -> 2 -> 3 (until true)
	// Result: 3
	tstRunForth(t, "BEGIN UNTIL", ": test-until 0 begin 1 + dup 3 = until ; test-until", 3)

	// Check LEAVE works in UNTIL
	// 0 -> 1 -> 2 (leave)
	// Result: 2
	tstRunForth(t, "BEGIN UNTIL LEAVE", ": test-until-leave 0 begin 1 + dup 2 = if leave then dup 5 = until ; test-until-leave", 2)
}

func TestBeginWhileRepeat(t *testing.T) {
	// Standard WHILE loop
	// 0 -> 1 -> 2 -> 3 (while false)
	// Result: 3
	tstRunForth(t, "BEGIN WHILE REPEAT", ": test-while 0 begin dup 3 < while 1 + repeat ; test-while", 3)

	// LEAVE in WHILE loop
	// 0 -> 1 -> 2 (leave)
	// Result: 2
	tstRunForth(t, "BEGIN WHILE REPEAT LEAVE", ": test-while-leave 0 begin dup 5 < while dup 2 = if leave then 1 + repeat ; test-while-leave", 2)
}

func TestNestedControl(t *testing.T) {
	// Nested BEGIN loops
	// Outer 0
	//   Inner 0 -> 1 -> 2 (break) -> drop 2 -> Outer 1
	//   Inner 0 -> 1 -> 2 (break) -> drop 2 -> Outer 2
	// Outer loop breaks
	// Result: 2
	tstRunForth(t, "Nested BEGIN Loops", ": test-nested 0 begin dup 2 < while 0 begin dup 2 < while 1 + repeat drop 1 + repeat ; test-nested", 2)
	tstRunForth(t, "BEGIN inside DO", `: tst ( n -- n' ) 0 swap 0 do i BEGIN dup 10 > IF drop leave THEN {: v :} v + v 1 + AGAIN " DID:" . i . cr LOOP ; 100 tst`, 440)
}

func TestLoopWithLocals(t *testing.T) {
	// Ensure locals work inside BEGIN loops
	// : test-locals {: a :} 0 begin dup a < while 1 + repeat ;
	// 3 test-locals -> 3
	tstRunForth(t, "Locals in BEGIN", ": test-locals {: a :} 0 begin dup a < while 1 + repeat ; 3 test-locals", 3)
}
