// SPDX-License-Identifier: MIT

package forth

import (
	"io"
	"strings"
	"testing"
)

func TestImmediateWord(t *testing.T) {
	// This test verifies that `immediate` works when called after a definition
	// e.g. : foo ... ; immediate

	// We define a word `var-i` that uses postpone to hook into compilation
	// Then we use it in `test` to ensure it executed immediately during compilation of `test`

	code := `
	: var-i  postpone r> postpone i postpone swap postpone >r ; immediate
	: test 3 0 DO " %04d: test!\n" << var-i >> sprintf type LOOP ;
	`
	expected := "0000: test!\n0001: test!\n0002: test!\n"

	vm.ClearResetState()
	mark(vm)
	err := vm.RunFromSource(strings.NewReader(code), "test", io.Discard)
	if err != nil {
		t.Fatalf("Failed to run code: %v", err)
	}
	// check that var-i is immediate
	if !vm.words[vm.dict["var-i"]].IsImmediate() {
		t.Errorf("var-i is not immediate")
	}
	// this time capture the output into a buffer
	var buf strings.Builder
	err = vm.RunFromSource(strings.NewReader("test"), "test", &buf)
	if err != nil {
		t.Errorf("Failed to execute test word: %v", err)
	}
	if buf.String() != expected {
		t.Errorf("Expected output \"%s\", got \"%s\"", expected, buf.String())
	}
	forget(vm)

	// Test 2: Verify `immediate` works inside a definition (legacy/non-standard but supported)
	code2 := `
	: var-i-inner immediate postpone r> postpone i postpone swap postpone >r ;
	: test2 3 0 DO " %04d: test!\n" << var-i-inner >> sprintf type LOOP ;
    `
	vm.ResetState()
	mark(vm)
	err = vm.RunFromSource(strings.NewReader(code2), "test", io.Discard)
	if err != nil {
		t.Fatalf("Failed to run code: %v", err)
	}
	// check that var-i is immediate
	if !vm.words[vm.dict["var-i-inner"]].IsImmediate() {
		t.Errorf("var-i-inner is not immediate")
	}
	// this time capture the output into a buffer
	var buf2 strings.Builder
	err = vm.RunFromSource(strings.NewReader("test2"), "test", &buf2)
	if err != nil {
		t.Errorf("Failed to execute test2 word: %v", err)
	}
	if buf2.String() != expected {
		t.Errorf("Code2 Expected output \"%s\", got \"%s\"", expected, buf2.String())
	}
	forget(vm)
}
