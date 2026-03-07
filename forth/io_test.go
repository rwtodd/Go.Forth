package forth

import (
	"strings"
	"testing"
)

func TestReadInput(t *testing.T) {
	code := `110 read read-line` // 110 is 'n'
	input := "interactive mode\nnext line of input!"

	vm := NewVM(strings.NewReader(input), nil)
	vm.ClearResetState()
	mark(vm)

	err := vm.Run(strings.NewReader(code), "test")
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	// stack should be: "interactive mode\nnext line of input!" parsed by read then read-line
	// "110 read" reads up to 'n' => "interactive mode\n" (reads until 'n' is found, omitting 'n')
	// Wait! the first 'n' is at "i[n]teractive mode\n"
	// So `110 read` gets "i"
	// Then `read-line` gets "teractive mode" (reads up to \n, strips \n)
	val1, _ := vm.Pop()
	val2, _ := vm.Pop()

	s2, _ := val1.(string) // top of stack, from read-line
	s1, _ := val2.(string) // bottom of stack, from read

	if s1 != "i" {
		t.Errorf("read expected 'i', got '%s'", s1)
	}

	if s2 != "teractive mode" {
		t.Errorf("read-line expected 'teractive mode', got '%s'", s2)
	}
}
