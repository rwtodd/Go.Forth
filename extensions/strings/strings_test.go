package strings

import (
	"io"
	"strings"
	"testing"

	"github.com/rwtodd/Go.Forth/forth"
)

func TestStringsExtension(t *testing.T) {
	vm := forth.NewVM()
	err := vm.RunFromSource(strings.NewReader(`<<" strings ">> <activate-extensions>`), "test", io.Discard)
	if err != nil {
		t.Fatalf("Failed to activate strings extension: %v", err)
	}

	runTest := func(name, code string, expected any) {
		t.Run(name, func(t *testing.T) {
			err := vm.RunFromSource(strings.NewReader(code), name, io.Discard)
			if err != nil {
				t.Fatalf("run failed: %v", err)
			}
			val, err := vm.Pop()
			if err != nil {
				t.Fatalf("expected value on stack")
			}
			// the stack shouldn't need clearing as we popped the expected result

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
					t.Errorf("expected %v, got %v", expected, val)
				}
			}
		})
	}

	runTest("blank-t", `"   " blank""?`, true)
	runTest("blank-f", `" a " blank""?`, false)

	runTest("len", `" hello" len""`, 5)

	runTest("trim", `"  foo  " trim""`, "foo")
	runTest("triml", `"  foo  " triml""`, "foo  ")
	runTest("trimr", `"  foo  " trimr""`, " foo")

	runTest("startsWith-t", `" hello" " he" starts""?`, true)
	runTest("startsWith-f", `" hello" " x" starts""?`, false)

	runTest("endsWith-t", `" hello" " lo" ends""?`, true)
	runTest("endsWith-f", `" hello" " l" ends""?`, false)

	runTest("upper", `" hello" upper""`, "HELLO")
	runTest("lower", `" HELLO" lower""`, "hello")

	runTest("split", `" a,b,c" " ," split""`, []string{"a", "b", "c"})

	runTest("@join", `<< " a" " b" " c" >> <strings> " ," @join""`, "a,b,c")
	runTest("<join>", `" a" " b" " c" 3 " ," <join"">`, "a,b,c")

	runTest("sub-1", `" hello" -1 0 sub""`, "o")
	runTest("sub-2", `" hello" 0 0 sub""`, "")
	runTest("sub-3", `" hello" 3 0 sub""`, "lo")
	runTest("sub-4", `" hello" 3 -1 sub""`, "l")
	runTest("sub-5", `" hello" -3 -1 sub""`, "ll")

	runTest("replace", `" foo bar foo" " foo" " baz" replace""`, "baz bar baz")

	runTest("index", `" hello" " e" index""`, 1)
	runTest("index-fail", `" hello" " z" index""`, -1)

	runTest("contains-t", `" hello" " ll" contains""?`, true)
	runTest("contains-f", `" hello" " z" contains""?`, false)
}
