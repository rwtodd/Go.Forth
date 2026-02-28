package random

import (
	"io"
	"strings"
	"testing"

	"github.com/rwtodd/Go.Forth/forth"
)

func TestRandomExtension(t *testing.T) {
	vm := forth.NewVM()
	err := vm.RunFromSource(strings.NewReader(`<<" random ">> <activate-extensions>`), "test", io.Discard)
	if err != nil {
		t.Fatalf("Failed to activate random extension: %v", err)
	}

	t.Run("randint", func(t *testing.T) {
		err := vm.RunFromSource(strings.NewReader("10 4 randint"), "test", io.Discard)
		if err != nil {
			t.Fatalf("randint failed: %v", err)
		}
		val, err := vm.Pop()
		if err != nil {
			t.Fatalf("expected value on stack")
		}
		if n, ok := val.(int); !ok || n < 4 || n >= 10 {
			t.Errorf("expected int in [4, 10), got %v", val)
		}
	})

	t.Run("randfloat", func(t *testing.T) {
		err := vm.RunFromSource(strings.NewReader("randfloat"), "test", io.Discard)
		if err != nil {
			t.Fatalf("randfloat failed: %v", err)
		}
		val, err := vm.Pop()
		if err != nil {
			t.Fatalf("expected value on stack")
		}
		if f, ok := val.(float64); !ok || f < 0 || f >= 1 {
			t.Errorf("expected float64 in [0, 1), got %v", val)
		}
	})

	t.Run("@seed", func(t *testing.T) {
		err := vm.RunFromSource(strings.NewReader("42 randseed!"), "test", io.Discard)
		if err != nil {
			t.Fatalf("@seed failed: %v", err)
		}

		err = vm.RunFromSource(strings.NewReader("100 0 randint"), "test", io.Discard)
		if err != nil {
			t.Fatalf("randint after seed failed: %v", err)
		}
		val1, _ := vm.Pop()

		err = vm.RunFromSource(strings.NewReader("42 randseed! 100 0 randint"), "test", io.Discard)
		if err != nil {
			t.Fatalf("randint after re-seed failed: %v", err)
		}
		val2, _ := vm.Pop()

		if val1 != val2 {
			t.Errorf("expected deterministic output after seeding, got %v and %v", val1, val2)
		}
	})

	t.Run("@select", func(t *testing.T) {
		err := vm.RunFromSource(strings.NewReader(`<< 10 20 30 >> <ints> @select`), "test", io.Discard)
		if err != nil {
			t.Fatalf("@select failed: %v", err)
		}
		val, err := vm.Pop()
		if err != nil {
			t.Fatalf("expected value on stack")
		}
		n, ok := val.(int)
		if !ok || (n != 10 && n != 20 && n != 30) {
			t.Errorf("expected 10, 20, or 30 from array, got %v", val)
		}
	})

	t.Run("@shuffle", func(t *testing.T) {
		err := vm.RunFromSource(strings.NewReader(`<< 10 20 30 >> <ints> constant y  y @shuffle`), "test", io.Discard)
		if err != nil {
			t.Fatalf("@shuffle failed: %v", err)
		}
		// Since it's shuffled, elements should still be exactly 10, 20, 30.
		err = vm.RunFromSource(strings.NewReader(`y 0 @ y 1 @ y 2 @ + +`), "test", io.Discard)
		if err != nil {
			t.Fatalf("reading shuffled array failed: %v", err)
		}
		sum, _ := vm.Pop()
		if sum != 60 {
			t.Errorf("expected sum to be 60, got %v", sum)
		}
	})
}
