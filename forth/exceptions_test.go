package forth

import (
	"strings"
	"testing"
)

func TestExceptions(t *testing.T) {
	vm := NewVM()

	tests := []struct {
		name     string
		input    string
		expected []any
	}{
		{
			name:     "catch success",
			input:    "10 20 : non-throwing-word + ;  2 [ non-throwing-word ] catch?",
			expected: []any{int64(30), int64(0)}, // success flag is 0
		},
		{
			name:     "catch failure with string",
			input:    "10 20 : throwing-word + \" oh no!\" throw ;  2 [ throwing-word ] catch?",
			expected: []any{"oh no!", int64(-1)}, // failure flag is -1
		},
		{
			name:     "catch failure pops remaining working stack",
			input:    "5 5 10 20 : bad-add-and-throw + \" boom\" throw ; 2 [ bad-add-and-throw ] catch?",
			expected: []any{int64(5), int64(5), "boom", int64(-1)},
		},
		{
			name: "catch inside word",
			input: `
: turn-to-nil drop ;
: safe-add 
  2 [ + ] catch? 
  if 
    drop " caught error" turn-to-nil 0 
  then 
;
10 20 safe-add 0
			`,
			expected: []any{int64(30), int64(0)},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			vm.ClearResetState()
			err := vm.Eval(tc.input)
			if err != nil {
				t.Fatalf("unexpected eval error: %v", err)
			}

			if len(vm.Stack) != len(tc.expected) {
				t.Fatalf("expected stack len %d, got %d. Stack: %v", len(tc.expected), len(vm.Stack), vm.Stack)
			}

			for i, expected := range tc.expected {
				got := vm.Stack[i]

				// Handle Error struct comparison to string
				if expStr, ok := expected.(string); ok {
					if evalErr, isErr := got.(error); isErr {
						if !strings.HasPrefix(evalErr.Error(), expStr) {
							t.Errorf("stack[%d] error string expected to start with %q, got %q", i, expStr, evalErr.Error())
						}
						continue
					} else if gotStr, isStr := got.(string); isStr {
						if !strings.HasPrefix(gotStr, expStr) {
							t.Errorf("stack[%d] expected string to start with %q, got %q", i, expStr, gotStr)
						}
						continue
					}
				}

				if got != expected {
					t.Errorf("stack[%d] expected %v, got %v", i, expected, got)
				}
			}
		})
	}
}

func TestThrowInDefinition(t *testing.T) {
	tstRunForth(t, "throw in definition", `: try-it  0 [ 2 + ] catch? IF drop " yes" ELSE " no" THEN ; try-it`, "yes")
	tstRunForth(t, "throw in definition 2", `: try-it  0 [ " an error!" throw ] catch? IF drop " yes" ELSE " no" THEN ; try-it`, "yes")
}
