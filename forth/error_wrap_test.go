package forth

import (
	"errors"
	"fmt"
	"strings"
	"testing"
)

func TestErrorBacktraceWrapping(t *testing.T) {
	// 1. Create a base error using your existing ErrUnderflowMsg helper
	baseErr := ErrUnderflowMsg("stack underflow")

	// 2. Wrap it just like we do in the Run() and interpret() methods
	wrappedOnce := fmt.Errorf("%w\n at `+`", baseErr)
	wrappedTwice := fmt.Errorf("%w\n at `dosums`", wrappedOnce)

	// 3. Verify that the final formatted string contains our keywords
	errStr := wrappedTwice.Error()
	if !strings.Contains(errStr, "stack underflow") || !strings.Contains(errStr, "at `+`") || !strings.Contains(errStr, "at `dosums`") {
		t.Errorf("Error string missing keywords: %s", errStr)
	}

	// 4. Verify that the "type" (sentinel error) is maintained via errors.Is
	if !errors.Is(wrappedTwice, ErrUnderflow) {
		t.Error("Should still be identified as ErrUnderflow via errors.Is")
	}

	// 5. Verify the ErrUnderflowMsg's specific Error struct is reachable via errors.As
	var forthErr *Error
	if !errors.As(wrappedTwice, &forthErr) {
		t.Error("Should be extractable as *forth.Error via errors.As")
	} else if forthErr.Err != ErrUnderflow {
		t.Errorf("Extracted *Error should point to ErrUnderflow, got %v", forthErr.Err)
	}
}
