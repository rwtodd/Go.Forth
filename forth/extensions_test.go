package forth

import (
	"errors"
	"io"
	"strings"
	"testing"
)

func TestExtensions(t *testing.T) {
	// Register a fake extension
	RegisterExtension("test-ext", func(vm *VM) error {
		vm.Define(&NativeWord{name: "test-word", run: func(vm *VM) error {
			vm.Push(123)
			return nil
		}, immediate: false})
		return nil
	})

	// Register a failing extension
	RegisterExtension("fail-ext", func(vm *VM) error {
		return errors.New("initialization failed")
	})

	t.Run("ExtensionList", func(t *testing.T) {
		vm := NewVM()
		vm.Run(strings.NewReader("extension-list"), io.Discard)
		val, err := vm.Pop()
		if err != nil {
			t.Fatalf("Expected value on stack, got error: %v", err)
		}
		list, ok := val.([]string)
		if !ok {
			t.Fatalf("Expected []string, got %T", val)
		}
		found := false
		for _, name := range list {
			if name == "test-ext" {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("Expected test-ext in extension list, got %v", list)
		}
	})

	t.Run("ActivateExtension", func(t *testing.T) {
		vm := NewVM()
		// "test-ext" 1 <activate-extensions>
		err := vm.Run(strings.NewReader(`<<" test-ext ">> <activate-extensions>`), io.Discard)
		if err != nil {
			t.Fatalf("Unexpected error: %v", err)
		}

		// Check if test-word is available
		_, ok := vm.dict["test-word"]
		if !ok {
			t.Errorf("Expected test-word to be defined after activation")
		}

		// Run test-word
		err = vm.Run(strings.NewReader("test-word"), io.Discard)
		if err != nil {
			t.Fatalf("Unexpected error running test-word: %v", err)
		}
		val, err := vm.Pop()
		if err != nil {
			t.Fatalf("Expected value from test-word")
		}
		if val != 123 {
			t.Errorf("Expected 123, got %v", val)
		}
	})

	t.Run("ActivateUnknownExtension", func(t *testing.T) {
		vm := NewVM()
		err := vm.Run(strings.NewReader(`<<" unknown-ext ">> <activate-extensions>`), io.Discard)
		if err == nil {
			t.Error("Expected error activating unknown extension")
		}
	})

	t.Run("ActivateFailingExtension", func(t *testing.T) {
		vm := NewVM()
		err := vm.Run(strings.NewReader(`<<" fail-ext ">> <activate-extensions>`), io.Discard)
		if err == nil {
			t.Error("Expected error activating failing extension")
		}
	})

	t.Run("ActivateWithOutVarlenSyntax", func(t *testing.T) {
		vm := NewVM()
		// <<" test-ext ">> <activate-extensions>
		err := vm.Run(strings.NewReader(`" test-ext" 1 <activate-extensions>`), io.Discard)
		if err != nil {
			t.Errorf("Unexpected error: %v", err)
		}
		_, ok := vm.dict["test-word"]
		if !ok {
			t.Errorf("Expected test-word to be defined after activation")
		}
	})

	t.Run("ActivateIdempotency", func(t *testing.T) {
		vm := NewVM()
		counter := 0
		RegisterExtension("counter-ext", func(vm *VM) error {
			counter++
			return nil
		})

		// First activation
		err := vm.Run(strings.NewReader(`" counter-ext" 1 <activate-extensions>`), io.Discard)
		if err != nil {
			t.Fatalf("Unexpected error: %v", err)
		}
		if counter != 1 {
			t.Errorf("Expected counter to be 1, got %d", counter)
		}

		// Second activation should be no-op
		err = vm.Run(strings.NewReader(`" counter-ext" 1 <activate-extensions>`), io.Discard)
		if err != nil {
			t.Fatalf("Unexpected error: %v", err)
		}
		if counter != 1 {
			t.Errorf("Expected counter to remain 1, got %d", counter)
		}
	})
}
