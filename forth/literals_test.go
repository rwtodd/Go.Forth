// SPDX-License-Identifier: MIT

package forth

import (
	"strings"
	"testing"
)

func TestLiteralOptimization(t *testing.T) {
	vm := NewVM()

	// Define a word that uses various literals
	src := `
	: test-lits 
		10
		40000
		100000
		" hello"
		" hello"
	;`

	// Run the definition
	if err := vm.Run(strings.NewReader(src), nil); err != nil {
		t.Fatalf("Failed to compile test-lits: %v", err)
	}

	// Verify vm.literals
	// We Expect: [100000, "hello"]
	// Note: 10 and 40000 are optimized into opcodes.
	if len(vm.literals) != 2 {
		t.Errorf("Expected 2 literals, got %d: %v", len(vm.literals), vm.literals)
	}

	if len(vm.strMap) != 1 {
		t.Errorf("Expected 1 interned string, got %d", len(vm.strMap))
	}

	if idx, ok := vm.strMap["hello"]; !ok {
		t.Errorf("Expected 'hello' to be interned. strMap keys:")
		for k, v := range vm.strMap {
			t.Logf("Key: %q, Val: %d", k, v)
		}
		t.Logf("Literals: %v", vm.literals)
	} else {
		// Verify the index points to "hello" in literals
		if idx >= len(vm.literals) {
			t.Errorf("Interned index %d out of bounds (len %d)", idx, len(vm.literals))
		} else if vm.literals[idx] != "hello" {
			t.Errorf("Interned index %d points to %v, expected 'hello'", idx, vm.literals[idx])
		}
	}

	// Verify 100000 is in literals
	foundLargeInt := false
	for _, lit := range vm.literals {
		if i, ok := lit.(int); ok && i == 100000 {
			foundLargeInt = true
			break
		}
	}
	if !foundLargeInt {
		t.Errorf("Expected 100000 in literals, but not found")
	}
}

func TestLiteralExecution(t *testing.T) {
	// Verify that the compiled code actually runs correctly
	tstRunForth(t, "LitOptimizationExec", `
	: test-exec 
		10 
		40000 
		100000 
		" hello" 
	; 
	test-exec
	`, 10, 40000, 100000, "hello")
}
