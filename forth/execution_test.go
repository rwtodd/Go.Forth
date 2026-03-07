// SPDX-License-Identifier: MIT

package forth

import (
	"strings"
	"testing"
)

func TestTickAndExecute(t *testing.T) {
	// Defines a word, gets its xt, and executes it
	tstRunForth(t, "TickExecute", `: hello 42 ; " hello" lookup-xt execute`, 42)
	tstRunForth(t, "TickExecute2", `: hello 142 ; [ hello ] execute`, 142)
	tstRunForth(t, "QuotOptPlus", `3 3 [ + ] execute`, 6)

	t.Run("InspectCodesegTopLevel", func(t *testing.T) {
		vm.ClearResetState()
		startLen := len(vm.codeseg)
		tprog := strings.NewReader(`[ + ]`)
		_ = vm.Run(tprog, "test")
		if len(vm.codeseg) > startLen {
			t.Errorf("Codeseg grew by %d tokens! Expected 0 growth.", len(vm.codeseg)-startLen)
			for i := startLen; i < len(vm.codeseg); i++ {
				t.Logf("Token %d: %v", i, vm.codeseg[i])
			}
		}
	})

	tstRunForth(t, "QuotOptPlusComp", `: test 3 3 [ + ] execute ; test`, 6)

	t.Run("InspectCodesegNested", func(t *testing.T) {
		vm.ClearResetState()
		startLen := len(vm.codeseg)
		tprog := strings.NewReader(`: test 3 3 [ + ] execute ; test`)
		_ = vm.Run(tprog, "test")
		// We expect : test to generate SOME code, but let's see how much
		t.Logf("Codeseg grew by %d for nested compilation.", len(vm.codeseg)-startLen)
		for i := startLen; i < len(vm.codeseg); i++ {
			t.Logf("Token %d: %v", i, vm.codeseg[i])
		}
	})
}

func TestBracketTick(t *testing.T) {
	// Compiles a word that uses ['] to get an xt, then executes it
	tstRunForth(t, "BracketTick", `: hello 99 ; : runner [ hello ] execute ; runner`, 99)
}

func TestCompileComma(t *testing.T) {
	// Manually compiles a call to a word
	// : target 123 ;
	// : builder [[ ' target compile, ]] ;
	// builder -> should run target and push 123
	tstRunForth(t, "CompileComma", `: target 123 ; : builder [[ " target" lookup-xt ]] compile-xt ; builder`, 123)
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

func TestCompileClosure(t *testing.T) {
	// Tests compiling a closure directly
	// make-closure pushes a closure that computes 25
	// runner compiles that closure into itself
	// executing runner should execute the closure
	tstRunForth(t, "CompileClosure", `: make-closure [ 5 5 * ] ; : runner [[ make-closure compile-xt ]] ; runner`, 25)
}
