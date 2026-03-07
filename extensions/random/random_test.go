package random

import (
	"io"
	"strings"
	"testing"

	"github.com/rwtodd/Go.Forth/forth"
)

func TestRandomExtension(t *testing.T) {
	vm := forth.NewVM(nil, io.Discard)
	err := vm.Run(strings.NewReader(`<<" random ">> <activate-extensions>`), "test")
	if err != nil {
		t.Fatalf("Failed to activate random extension: %v", err)
	}

	t.Run("randint", func(t *testing.T) {
		err := vm.Run(strings.NewReader("10 4 randint"), "test")
		if err != nil {
			t.Fatalf("randint failed: %v", err)
		}
		val, err := vm.Pop()
		if err != nil {
			t.Fatalf("expected value on stack")
		}
		if n, ok := val.(int64); !ok || n < 4 || n >= 10 {
			t.Errorf("expected int in [4, 10), got %v", val)
		}
	})

	t.Run("randfloat", func(t *testing.T) {
		err := vm.Run(strings.NewReader("randfloat"), "test")
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
		err := vm.Run(strings.NewReader("42 randseed!"), "test")
		if err != nil {
			t.Fatalf("@seed failed: %v", err)
		}

		err = vm.Run(strings.NewReader("100 0 randint"), "test")
		if err != nil {
			t.Fatalf("randint after seed failed: %v", err)
		}
		val1, _ := vm.Pop()

		err = vm.Run(strings.NewReader("42 randseed! 100 0 randint"), "test")
		if err != nil {
			t.Fatalf("randint after re-seed failed: %v", err)
		}
		val2, _ := vm.Pop()

		if val1 != val2 {
			t.Errorf("expected deterministic output after seeding, got %v and %v", val1, val2)
		}
	})

	t.Run("@select", func(t *testing.T) {
		err := vm.Run(strings.NewReader(`<< 10 20 30 >> <ints> @select`), "test")
		if err != nil {
			t.Fatalf("@select failed: %v", err)
		}
		val, err := vm.Pop()
		if err != nil {
			t.Fatalf("expected value on stack")
		}
		n, ok := val.(int64)
		if !ok || (n != 10 && n != 20 && n != 30) {
			t.Errorf("expected 10, 20, or 30 from array, got %v", val)
		}
	})

	t.Run("@shuffle", func(t *testing.T) {
		err := vm.Run(strings.NewReader(`<< 10 20 30 >> <ints> " y" constant  y @shuffle`), "test")
		if err != nil {
			t.Fatalf("@shuffle failed: %v", err)
		}
		// Since it's shuffled, elements should still be exactly 10, 20, 30.
		err = vm.Run(strings.NewReader(`y 0 @ y 1 @ y 2 @ + +`), "test")
		if err != nil {
			t.Fatalf("reading shuffled array failed: %v", err)
		}
		sum, _ := vm.Pop()
		if sum != int64(60) {
			t.Errorf("expected sum to be 60, got %v", sum)
		}
	})
}
