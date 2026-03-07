// SPDX-License-Identifier: MIT

package forth

import (
	"strings"
	"testing"
)

func TestImmediateWord(t *testing.T) {
	// This test verifies that `immediate` works when called after a definition
	// e.g. : foo ... ; immediate

	// We define a word `var-i` that uses postpone to hook into compilation
	// Then we use it in `test` to ensure it executed immediately during compilation of `test`

	code := `
	: var-i  [[ <<" r> i swap >r ">> ]] <postpone> ; immediate
	: test 3 0 DO " %04d: test!\n" << var-i >> sprintf type LOOP ;
	`
	expected := "0000: test!\n0001: test!\n0002: test!\n"

	vm.ClearResetState()
	mark(vm)
	err := vm.Run(strings.NewReader(code), "test")
	if err != nil {
		t.Fatalf("Failed to run code: %v", err)
	}
	// check that var-i is immediate
	if !vm.words[vm.dict["var-i"]].IsImmediate() {
		t.Errorf("var-i is not immediate")
	}
	// this time capture the output into a buffer
	var buf strings.Builder
	vm.SetOutput(&buf)
	err = vm.Run(strings.NewReader("test"), "test")
	if err != nil {
		t.Errorf("Failed to execute test word: %v", err)
	}
	if buf.String() != expected {
		t.Errorf("Expected output \"%s\", got \"%s\"", expected, buf.String())
	}
	forget(vm)

	// Test 2: Verify `immediate` works inside a definition (legacy/non-standard but supported)
	code2 := `
	: var-i-inner immediate [[ <<" R> I SWAP >R ">> ]] <postpone> ;
	: test2 3 0 DO " %04d: test!\n" << var-i-inner >> sprintf type LOOP ;
    `
	vm.ResetState()
	mark(vm)
	err = vm.Run(strings.NewReader(code2), "test")
	if err != nil {
		t.Fatalf("Failed to run code: %v", err)
	}
	// check that var-i is immediate
	if !vm.words[vm.dict["var-i-inner"]].IsImmediate() {
		t.Errorf("var-i-inner is not immediate")
	}
	// this time capture the output into a buffer
	var buf2 strings.Builder
	vm.SetOutput(&buf2)
	err = vm.Run(strings.NewReader("test2"), "test")
	if err != nil {
		t.Errorf("Failed to execute test2 word: %v", err)
	}
	if buf2.String() != expected {
		t.Errorf("Code2 Expected output \"%s\", got \"%s\"", expected, buf2.String())
	}
	forget(vm)
}
