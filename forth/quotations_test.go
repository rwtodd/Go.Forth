// SPDX-License-Identifier: MIT

package forth

import (
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
}

func TestQuotationsShadowing(t *testing.T) {
	// : test (| a |) [ (| a |) a ] execute a ;
	// 10 20 test -> 10 20 (Outer=20, Inner=10)
	tstRunForth(t, "Shadowing", ": test (| a |) [ (| a |) a ] execute a ; 10 20 test", 10, 20)
}

func TestRecursionInClosure(t *testing.T) {
	// Basic recursion inside a closure?
	// Since recur relies on name, anonymous recursion is tricky.
	// If we support :NONAME or just recursive closure via valid stack manipulation?
	// For now, let's test that we can call ourselves if we are assigned to a variable/constant?
	// : fact-closure [: dup 1 <= IF drop 1 ELSE dup 1 - fact-closure execute * THEN ;] ;
	// This requires forward reference or `defer`.
	// We don't have defer standard.
	// Let's skip recursion inside anonymous for now, consistent with plan.
}

func TestTopLevelQuotation(t *testing.T) {
	// [ 10 20 + ] execute -> 30
	tstRunForth(t, "TopLevel", "[ 10 20 + ] execute", 30)
}
