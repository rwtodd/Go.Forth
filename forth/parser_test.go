// SPDX-License-Identifier: MIT

package forth

import (
	"io"
	"os"
	"strings"
	"testing"
)

func TestLiterals(t *testing.T) {
	tstRunForth(t, "LitWithStr", `: tst [[ " 0" ORD ]] literal ; tst " 0" ORD =`, -1)
	tstRunForth(t, "LitRel", `: tst [[ -5 10 * ]] literal + ; 2 tst`, -48)
	tstRunForth(t, "LitLarge", `: tst [[ 55 1000 * ]] literal + ; 2 tst`, 55002)
}

func TestEval(t *testing.T) {
	tstRunForth(t, "EvalSimple", `" 5 10 + " eval`, 15)
	tstRunForth(t, "EvalNested", `" \" 3 4 * \" eval 2 + " eval 5 +`, 19)
	tstRunForth(t, "EvalDef", `" : evalword 42 ; " eval evalword`, 42)
	tstRunForth(t, "Eval During Def", `: testit 1 2 " +" eval ; testit`, 3)
}

func TestLoad(t *testing.T) {
	os.WriteFile("/tmp/testload.4th", []byte(": loadword 25 ;"), 0644)
	defer os.Remove("/tmp/testload.4th")
	tstRunForth(t, "LoadSimple", `" /tmp/testload.4th" load loadword`, 25)
}

func TestCompilationChecks(t *testing.T) {
	vm := NewVM(nil, io.Discard)
	err := vm.Run(strings.NewReader(`: outerword : innerword ;`), "test")
	if err == nil {
		t.Errorf("Expected error starting nested definition, got nil")
	}
}

func TestUpCase(t *testing.T) {
	tstRunForth(t, "CaseInsens", `: TST 3 4 sWaP ; tst`, 4, 3)
}
