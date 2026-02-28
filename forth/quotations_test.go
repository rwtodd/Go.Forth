// SPDX-License-Identifier: MIT

package forth

import (
	stdbytes "bytes"
	"strings"
	"testing"
)

func TestQuotationsBasic(t *testing.T) {
	// [ 1 2 + ] EXECUTE -> 3
	tstRunForth(t, "BasicQuot", "[ 1 2 + ] execute", 3)
}

func TestQuotationsCapture(t *testing.T) {
	// : test (| a |) [ a ] execute ; 10 test -> 10
	tstRunForth(t, "CaptureSimple", ": test (| a |) [ a ] execute ; 10 test", 10)
}

func TestQuotationsCaptureNested(t *testing.T) {
	// : test (| a |) [ (| b |) [ a b + ] execute ] execute ; 10 20 test -> 30
	tstRunForth(t, "CaptureNested", ": test (| a |) [ (| b |) [ a b + ] execute ] execute ; 10 20 test", 30)
}

func TestQuotationsModify(t *testing.T) {
	// : test (| a |) [ 20 a! ] execute a ; 10 test -> 20
	tstRunForth(t, "ModifyCapture", ": test (| a |) [ 20 a! ] execute a ; 10 test", 20)
}

func TestQuotationsEscaping(t *testing.T) {
	// : make-adder (| n |) [ n + ] ;
	// 5 make-adder constant add5
	// 10 add5 execute -> 15
	tstRunForth(t, "EscapingClosure", ": make-adder (| n |) [ n + ] ; 5 make-adder constant add5 10 add5 execute", 15)
	tstRunForth(t, "Counter", ": counter (| n |) [ n dup 1 + n! ] ; 1 counter dup execute swap execute", 1, 2)
	tstRunForth(t, "2 Counters", ": counter (| n |) [ n dup 1 + n! ] ; 1 counter 100 counter execute swap execute", 100, 1)
}

func TestQuotationsShadowing(t *testing.T) {
	// : test (| a |) [ (| a |) a ] execute a ;
	// 10 20 test -> 10 20 (Outer=20, Inner=10)
	tstRunForth(t, "Shadowing", ": test (| a |) [ (| a |) a ] execute a ; 10 20 test", 10, 20)
	tstRunForth(t, "Shadowing 2", ": test (| a |) a [ (| a |) 10 a + a! a ] execute a ; 10 test", 20, 10)
}

func TestRecursionInClosure(t *testing.T) {
	tstRunForth(t, "Recur in a Closure",
		`: quot+recur (| n | count -- count ) 0 count! " Starting n:" . n . cr [ (| n |) count 1 + count!  n 0 > IF n . cr n 1 - recur ELSE " done!" . cr THEN ] n swap execute " Bye!" type cr count ;
	   4 quot+recur`, 5)
	tstRunForth(t, "(tail-call) in a Closure",
		`: quot+tc (| n | count -- count ) 0 count! " Starting n:" . n . cr [ (| n |) count 1 + count!  n 0 > IF n . cr n 1 - (tail-call) ELSE " done!" . cr THEN ] n swap execute " Bye!" type cr count ;
	   4 quot+tc`, 5)
	tstRunForth(t, "Factorial w/RECUR", ": tst ( n -- fact ) 1 swap [ dup 1 = IF drop ELSE tuck * swap 1 - RECUR THEN ] execute ; 5 tst", 120)
	tstRunForth(t, "Factorial w/(tail-call)", ": tst ( n -- fact ) 1 swap [ dup 1 = IF drop ELSE tuck * swap 1 - (tail-call) THEN ] execute ; 5 tst", 120)
	tstRunForth(t, "Factorial w/RECUR and locals", ": tst2 ( n -- fact ) 1 swap [ (| acc n |) n 1 = IF acc ELSE acc n * n 1 - RECUR THEN ] execute ; 5 tst2", 120)
	tstRunForth(t, "Factorial w/(tail-call) and locals", ": tst2 ( n -- fact ) 1 swap [ (| acc n |) n 1 = IF acc ELSE acc n * n 1 - (tail-call) THEN ] execute ; 5 tst2", 120)
}

func TestGarbageWordDetected(t *testing.T) {
	tstRunForthErr(t, "Garbage Word Detected", ": tst3 ( n -- fact ) 1 swap [ (| acc n |) n 1 = IF acc ELSE acc n * n vbardg  RECUR THEN ] execute ;", ErrArgument)
}

func TestRecursionNoLocals(t *testing.T) {
	// Verify that returning from a recursive call doesn't break anything.
	// We need to capture stdout to verify the output.
	vm.ClearResetState()
	mark(vm)
	code := `[ dup 0= IF . " DONE!" type cr ELSE dup 1 - RECUR " WAS" . . cr THEN ]  4 swap execute`

	var buf stdbytes.Buffer
	if err := vm.RunFromSource(strings.NewReader(code), "test", &buf); err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	expected := "0 DONE!\nWAS 1 \nWAS 2 \nWAS 3 \nWAS 4 \n"
	if buf.String() != expected {
		t.Errorf("Output mismatch.\nExpected:\n%q\nGot:\n%q", expected, buf.String())
	}
	forget(vm)
}

func TestRecursionLocalsAccess(t *testing.T) {
	// Verify that we can access locals after returning from a recursive call
	// We need to capture stdout to verify the output.
	vm.ClearResetState()
	mark(vm)
	code := `[ (| n |) n 0= IF n . " DONE!" type cr ELSE n 1 - RECUR " WAS" . n . cr THEN ] 4 swap execute`

	var buf stdbytes.Buffer
	if err := vm.RunFromSource(strings.NewReader(code), "test", &buf); err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	expected := "0 DONE!\nWAS 1 \nWAS 2 \nWAS 3 \nWAS 4 \n"
	if buf.String() != expected {
		t.Errorf("Output mismatch.\nExpected:\n%q\nGot:\n%q", expected, buf.String())
	}
	forget(vm)
}

func TestTopLevelQuotation(t *testing.T) {
	// [ 10 20 + ] execute -> 30
	tstRunForth(t, "TopLevel", "[ 10 20 + ] execute", 30)
}

func TestLocalsAfterRecur(t *testing.T) {
	// Verify that recur works even if locals are defined AFTER the recur call
	// in the source text.
	vm.ClearResetState()
	mark(vm)
	// Factorial using recursion, but locals defined at the end.
	// [ dup 0 > IF dup 1 - recur * ELSE drop 1 THEN (| | temp ) ]
	code := `[ dup 0 > IF dup 1 - recur * ELSE drop 1 THEN (| | temp ) ] 5 swap execute`

	var buf stdbytes.Buffer
	// We expect 120 on stack, no output.
	if err := vm.RunFromSource(strings.NewReader(code), "test", &buf); err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if len(vm.Stack) != 1 {
		t.Fatalf("Stack length mismatch. Expected 1, got %d", len(vm.Stack))
	}
	if val, ok := vm.Stack[0].(int); !ok || val != 120 {
		t.Errorf("Value mismatch. Expected 120, got %v", vm.Stack[0])
	}
	forget(vm)
}

func TestTailCallLocalsAfter(t *testing.T) {
	// Verify that tail-call works even if locals are defined AFTER the call
	vm.ClearResetState()
	mark(vm)
	// [ dup 0 > IF 1 - (tail-call) ELSE THEN (| | n ) ]
	// This captures the tail-call scenario with late locals.
	code := `[ dup 0 > IF 1 - (tail-call) ELSE THEN (| | n ) ] 5 swap execute`

	var buf stdbytes.Buffer
	if err := vm.RunFromSource(strings.NewReader(code), "test", &buf); err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if len(vm.Stack) != 1 {
		t.Fatalf("Stack length mismatch. Expected 1, got %d", len(vm.Stack))
	}
	if val, ok := vm.Stack[0].(int); !ok || val != 0 {
		t.Errorf("Value mismatch. Expected 0, got %v", vm.Stack[0])
	}
	forget(vm)
}
