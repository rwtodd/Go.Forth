// SPDX-License-Identifier: MIT

package forth

import (
	"testing"
)

func TestTickAndExecute(t *testing.T) {
	// Defines a word, gets its xt with ', and executes it
	tstRunForth(t, "TickExecute", `: hello 42 ; ' hello execute`, 42)
}

func TestBracketTick(t *testing.T) {
	// Compiles a word that uses ['] to get an xt, then executes it
	tstRunForth(t, "BracketTick", `: hello 99 ; : runner ['] hello execute ; runner`, 99)
}

func TestCompileComma(t *testing.T) {
	// Manually compiles a call to a word
	// : target 123 ;
	// : builder [[ ' target compile, ]] ;
	// builder -> should run target and push 123
	tstRunForth(t, "CompileComma", `: target 123 ; : builder [[ ' target compile, ]] ; builder`, 123)
}

func TestRecurWithLiterals(t *testing.T) {
	tstRunForth(t, "RecurWithLiterals", `: test-recur (| n |) n 1 <= IF 1 ELSE n n 1 - recur " hi there!" . * THEN ; 3 test-recur`, 6)
}

func TestShadowing(t *testing.T) {
	tstRunForth(t, "BasicShadowing", `: foo 42 ; : foo foo 1 + ; foo`, 43)
	tstRunForth(t, "RevertShadow", `: foo 42 ; mark : foo 21 ; foo forget foo +`, 63)
}

func TestRecurMultipleWords(t *testing.T) {
	// Test that RECUR works correctly in multiple different words
	// Each word should recurse to itself, not to other words
	tstRunForth(t, "RecurMultipleWords", `: test-recur (| n |) n 1 <= IF 1 ELSE n n 1 - recur " hi there!" . * THEN ; : test-recur-2 (| n |) n 10 > IF 1 ELSE n n 1 + recur " hi there!" . * THEN ; 3 test-recur 8 test-recur-2`, 6, 720)
}
