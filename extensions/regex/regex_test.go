package regex

import (
	"io"
	"regexp"
	"strings"
	"testing"

	"github.com/rwtodd/Go.Forth/forth"
)

func TestRegexExtension(t *testing.T) {
	vm := forth.NewVM(nil, io.Discard)
	err := vm.Run(strings.NewReader(`<<" regex ">> <activate-extensions>`), "test")
	if err != nil {
		t.Fatalf("Failed to activate regex extension: %v", err)
	}

	runTest := func(name, code string, expected any) {
		t.Run(name, func(t *testing.T) {
			err := vm.Run(strings.NewReader(code), name)
			if err != nil {
				t.Fatalf("run failed: %v", err)
			}
			val, err := vm.Pop()
			if err != nil {
				t.Fatalf("expected value on stack")
			}

			if expectedBool, ok := expected.(bool); ok {
				expectedInt := 0
				if expectedBool {
					expectedInt = -1
				}
				if valInt, ok := val.(int64); !ok || valInt != int64(expectedInt) {
					t.Errorf("expected bool %v (forth %d), got %v", expectedBool, expectedInt, val)
				}
			} else if expSlice, ok := expected.([]string); ok {
				valSlice, ok := val.([]string)
				if !ok || strings.Join(valSlice, "|") != strings.Join(expSlice, "|") {
					t.Errorf("expected slice %v, got %v", expSlice, val)
				}
			} else {
				expectedAny := expected
				if e, ok := expected.(int); ok {
					expectedAny = int64(e)
				}
				if val != expectedAny {
					// handle regex pointers
					if expRe, ok := expectedAny.(*regexp.Regexp); ok {
						valRe, vok := val.(*regexp.Regexp)
						if !vok || valRe.String() != expRe.String() {
							t.Errorf("expected regex %v, got %v", expRe, val)
						}
					} else {
						t.Errorf("expected %v, got %v", expected, val)
					}
				}
			}
		})
	}

	runTest("rx:", `rx: /hello/`, "hello")
	runTest("rx:-custom", `rx: {world}`, "world")
	runTest("[rx:]", `: make-pat [rx:] /test/ ; make-pat`, "test")
	runTest("rx-compile", `rx: /^foo/ rx-compile`, regexp.MustCompile("^foo"))

	runTest("match?-t", `" foo" rx: /^f/ rx-match?`, true)
	runTest("match?-f", `" bar" rx: /^f/ rx-match?`, false)
	runTest("match?-compiled", `" foo" rx: /^f/ rx-compile rx-match?`, true)

	runTest("gsub", `" hello lolo" rx: /lo/ " x" rx-gsub`, "helx xx")
	runTest("gsub-compiled", `" hello lolo" rx: /lo/ rx-compile " x" rx-gsub`, "helx xx")
	runTest("gsub-backref", `" hello" rx: /(e)/ " x${1}x" rx-gsub`, "hxexllo")

	runTest("sub", `" hello lolo" rx: /lo/ " x" rx-sub`, "helx lolo")
	runTest("sub-compiled", `" hello lolo" rx: /lo/ rx-compile " x" rx-sub`, "helx lolo")
	runTest("sub-backref", `" hello" rx: /(e)/ " x${1}x" rx-sub`, "hxexllo")

	runTest("split", `" a,b,c" rx: /,/ rx-split`, []string{"a", "b", "c"})
	runTest("split-empty", `" a,,c" rx: /,/ rx-split`, []string{"a", "", "c"})

	runTest("find", `"   hello " rx: /[a-z]+/ rx-find`, "hello")
	runTest("find-fail", `"   123 " rx: /[a-z]+/ rx-find`, "")

	runTest("gfind", `"  a b 1 c " rx: /[a-z]/ rx-gfind`, []string{"a", "b", "c"})
	runTest("gfind-fail", `"  1 2 3 " rx: /[a-z]/ rx-gfind`, []string{})

	runTest("match-xt", `" hello" rx: /(e)(l)/ [ 1 @ nip ] rx-match`, "e")
}
