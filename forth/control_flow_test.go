// SPDX-License-Identifier: MIT

package forth

import (
	"testing"
)

func TestExit(t *testing.T) {
	// Simple EXIT
	// 10 exit 20 -> should have 10 on stack
	tstRunForth(t, "Simple EXIT", ": test-exit 10 exit 20 ; test-exit", 10)

	// EXIT inside IF (taken)
	// 10 1 if exit then 20 -> should have 10
	tstRunForth(t, "EXIT inside IF taken", ": test-if-exit 10 1 if exit then 20 ; test-if-exit", 10)

	// EXIT inside IF (not taken)
	// 10 0 if exit then 20 -> should have 10 20
	tstRunForth(t, "EXIT inside IF not taken", ": test-if-no-exit 10 0 if exit then 20 ; test-if-no-exit", 10, 20)

	// EXIT inside DO (taken)
	// 5 0 do i 2 = if exit then i loop 99 -> 0 1
	tstRunForth(t, "EXIT inside DO", ": test-do-exit 5 0 do i 2 = if exit then i loop 99 ; test-do-exit", 0, 1)

	// Nested EXIT
	// inner: 1 exit 2
	// outer: inner 3
	// result: 1 3
	tstRunForth(t, "Nested EXIT", ": inner 1 exit 2 ; : outer inner 3 ; outer", 1, 3)
}

func TestLeave(t *testing.T) {
	// Simple LEAVE
	// 5 0 do i 2 = if leave then i loop -> 0 1
	tstRunForth(t, "Simple LEAVE", ": test-leave 5 0 do i 2 = if leave then i loop ; test-leave", 0, 1)

	// LEAVE at start
	// 5 0 do leave i loop 99 -> 99
	tstRunForth(t, "LEAVE at start", ": test-leave-start 5 0 do leave i loop 99 ; test-leave-start", 99)

	// Multiple LEAVEs
	// 10 0 do ...
	// i 2 = -> leave
	// i 4 = -> leave (never reached)
	// result: 0 1
	tstRunForth(t, "Multiple LEAVEs", ": test-multi-leave 10 0 do i 2 = if leave then i 4 = if leave then i loop ; test-multi-leave", 0, 1)

	// Nested Loop LEAVE
	// Outer 0..2 (3 times)
	// Inner 0..2 (3 times)
	// Condition: j i + 3 = if leave
	// j=0: i=0,1,2 (sum 0,1,2). No leave. Stack: 0 1 2
	// j=1: i=0,1,2 (sum 1,2,3). Leave at i=2 (sum=3). Stack: 0 1
	// j=2: i=0,1,2 (sum 2,3,4). Leave at i=1 (sum=3). Stack: 0
	// Total stack: 0 1 2 0 1 0
	tstRunForth(t, "Nested LEAVE", ": test-nested 3 0 do 3 0 do j i + 3 = if leave then i loop loop ; test-nested",
		0, 1, 2, // j=0
		0, 1, // j=1 (broke at i=2)
		0, // j=2 (broke at i=1)
	)

	// LEAVE with +LOOP
	// 10 0 do ... 2 +loop
	// i=0, 2, 4.
	// Condition: i 3 > if leave (taken at 4).
	// Stack: 0 2
	tstRunForth(t, "LEAVE with +LOOP", ": test-plus-loop 10 0 do i 3 > if leave then i 2 +loop ; test-plus-loop", 0, 2)
}

func TestControlFlowContext(t *testing.T) {
	// LEAVE outside loop
	// LEAVE is immediate, runs at compile time.
	// Error comes from opLeave.
	tstRunForthErr(t, "LEAVE outside loop", ": test-bad-leave leave ;", ErrBadState)

	// EXIT outside definition
	tstRunForthErr(t, "EXIT interpret mode", "exit", ErrBadState)
}
