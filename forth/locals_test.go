// SPDX-License-Identifier: MIT

package forth

import (
	"strings"
	"testing"
)

func TestLocalsBasic(t *testing.T) {
	vm := NewVM()
	vm.Define("test", Word{func(vm *VM) error {
		vm.Push(1)
		vm.Push(2)
		return nil
	}, false})

	// : add3 {: a b c :} a b c + + . ;
	input := ": add3 {: a b c :} a b c + + . ; 1 2 3 add3"
	r := strings.NewReader(input)
	var out strings.Builder

	if err := vm.Run(r, &out); err != nil {
		t.Fatalf("Run failed: %v", err)
	}

	expected := "6 "
	if out.String() != expected {
		t.Errorf("Expected output %q, got %q", expected, out.String())
	}
}

func TestLocalsModifying(t *testing.T) {
	vm := NewVM()

	// : inc {: a :} a 1 + a! a . ;
	input := ": inc {: a :} a 1 + a! a . ; 5 inc"
	r := strings.NewReader(input)
	var out strings.Builder

	if err := vm.Run(r, &out); err != nil {
		t.Fatalf("Run failed: %v", err)
	}

	expected := "6 "
	if out.String() != expected {
		t.Errorf("Expected output %q, got %q", expected, out.String())
	}
}

func TestLocalsUninitialized(t *testing.T) {
	vm := NewVM()

	// : test {: a | b :} a . b . ;
	// a gets 1. b is local but uninitialized (0 or nil, usually 0/nil in Go slice).
	// But Pop() returns arbitrary types.
	// In Go.Forth VM stack is []any.
	// locals = make([]any, count). So nil.
	// We need to check what printing nil does.
	// ' . ' uses fmt.Sprint.

	input := ": test {: a | b :} a . b . ; 1 test"
	r := strings.NewReader(input)
	var out strings.Builder

	if err := vm.Run(r, &out); err != nil {
		t.Fatalf("Run failed: %v", err)
	}

	// output depends on how nil is printed. Assuming "<nil>"?
	// Let's verify what happens.
	// Actually, let's set b explicitely to confirm access.
	// : test {: a | b :} 5 b! a . b . ; -> 1 5

	input2 := ": test2 {: a | b :} 5 b! a . b . ; 1 test2"
	r2 := strings.NewReader(input2)
	out.Reset()
	if err := vm.Run(r2, &out); err != nil {
		t.Fatalf("Run failed: %v", err)
	}

	expected := "1 5 "
	if out.String() != expected {
		t.Errorf("Expected output %q, got %q", expected, out.String())
	}
}

func TestLocalsShadowing(t *testing.T) {
	vm := NewVM()

	// Variable X
	// : test {: X :} X . ;
	// 5 test -> 5.
	// X @ -> 0.

	input := "variable X : test {: X :} X . ; 5 test X @"
	r := strings.NewReader(input)
	var out strings.Builder // ' . ' prints with space. 'X @' puts on stack.
	// wait, I need to print X @.

	input = "variable X : test {: X :} X . ; 5 test X @ ."
	r = strings.NewReader(input)
	out.Reset()

	if err := vm.Run(r, &out); err != nil {
		t.Fatalf("Run failed: %v", err)
	}

	expected := "5 0 "
	if out.String() != expected {
		t.Errorf("Expected output %q, got %q", expected, out.String())
	}
}

func TestLocalsRecursion(t *testing.T) {
	vm := NewVM()

	// Factorial
	// : fact {: n :} n 1 <= IF 1 ELSE n n 1 - recur * THEN ;

	input := ": fact {: n :} n 1 <= IF 1 ELSE n n 1 - recur * THEN ; 5 fact ."
	r := strings.NewReader(input)
	var out strings.Builder

	if err := vm.Run(r, &out); err != nil {
		t.Fatalf("Run failed: %v", err)
	}

	expected := "120 "
	if out.String() != expected {
		t.Errorf("Expected output %q, got %q", expected, out.String())
	}
}

func TestLocalsMixedControl(t *testing.T) {
	vm := NewVM()

	// Test locals inside IF (deferred definition)
	// : test 1 IF {: a :} a . THEN ;

	input := ": test 1 IF {: a :} a . THEN ; 10 test"
	r := strings.NewReader(input)
	var out strings.Builder

	if err := vm.Run(r, &out); err != nil {
		t.Fatalf("Run failed: %v", err)
	}

	expected := "10 "
	if out.String() != expected {
		t.Errorf("Expected output %q, got %q", expected, out.String())
	}
}
