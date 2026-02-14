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
	// : builder [ ' target compile, ] ;
	// builder -> should run target and push 123
	tstRunForth(t, "CompileComma", `: target 123 ; : builder [ ' target compile, ] ; builder`, 123)
}
