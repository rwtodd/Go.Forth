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
	if err := vm.RunFromSource(strings.NewReader(src), "test", nil); err != nil {
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
		if i, ok := lit.(int64); ok && i == 100000 {
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

func TestFloat16Optimization(t *testing.T) {
	vm := NewVM()

	// 1.0 -> optimizable
	// 1.5 -> optimizable
	// 0.1 -> not optimizable by Float16, BUT optimizable by Decimal.
	// 70000.0 -> not optimizable (overflow both)

	src := `
	: test-float-lits
		1.0
		1.5
		0.1
		70000.0
	;`

	if err := vm.RunFromSource(strings.NewReader(src), "test", nil); err != nil {
		t.Fatalf("Failed to compile test-float-lits: %v", err)
	}

	// Verify vm.literals
	// Expected: [70000.0]
	// 1.0, 1.5 (float16) and 0.1 (decimal) should be opcodes.
	if len(vm.literals) != 1 {
		t.Errorf("Expected 1 float literal, got %d: %v", len(vm.literals), vm.literals)
	}

	foundHuge := false

	for _, lit := range vm.literals {
		if f, ok := lit.(float64); ok {
			if f == 70000.0 {
				foundHuge = true
			}
		}
	}

	if !foundHuge {
		t.Errorf("Expected 70000.0 in literals")
	}
}

func TestFloat16Execution(t *testing.T) {
	tstRunForth(t, "Float16Exec", `
	: test-f16 
		1.0 
		1.5 
		-2.0
		0.0
		0.1 
		70000.0
	;
	test-f16
	`, 1.0, 1.5, -2.0, 0.0, 0.1, 70000.0)
}

func TestDecimalOptimization(t *testing.T) {
	vm := NewVM()

	// 0.1 -> Decimal (Scale 0)
	// 1.23 -> Decimal (Scale 1)
	// 0.005 -> Decimal (Scale 2)
	// 12.3 -> Decimal (Scale 0)
	// 70000.0 -> Float16? No (overflow). Decimal? 70000*10 > 8191. No. Not optimized.
	// 123.4567 -> Returns false (too many digits for scale 3)

	src := `
	: test-dec-lits
		0.1
		1.23
		0.005
		12.3
		70000.0
		123.4567
	;`

	if err := vm.RunFromSource(strings.NewReader(src), "test", nil); err != nil {
		t.Fatalf("Failed to compile test-dec-lits: %v", err)
	}

	// Verify vm.literals
	// Expected: [70000.0, 123.4567]
	// 0.1, 1.23, 0.005, 12.3 should be opcodes.
	if len(vm.literals) != 2 {
		t.Errorf("Expected 2 decimal literals, got %d: %v", len(vm.literals), vm.literals)
	}

	foundHuge := false
	foundComplex := false
	for _, lit := range vm.literals {
		if f, ok := lit.(float64); ok {
			if f == 70000.0 {
				foundHuge = true
			}
			if f == 123.4567 {
				foundComplex = true
			}
		}
	}

	if !foundHuge {
		t.Errorf("Expected 70000.0 in literals")
	}
	if !foundComplex {
		t.Errorf("Expected 123.4567 in literals")
	}
}

func TestDecimalExecution(t *testing.T) {
	tstRunForth(t, "DecimalExec", `
	: test-dec 
		0.1 
		1.23 
		-0.005
		70000.0
	;
	test-dec
	`, 0.1, 1.23, -0.005, 70000.0)
}
