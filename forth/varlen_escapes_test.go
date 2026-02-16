package forth

import (
	"testing"
)

func TestVarLen(t *testing.T) {
	// Note: The parser consumes the delimiter space after string words.
	// So " Hello" -> "Hello" on stack.

	tstRunForth(t, "Simple Sprintf", `" Hello %s!" << " World" >>sprintf`, "Hello World!")

	tstRunForth(t, "Integers", `" %d + %d = %d" << 1 2 1 2 + >>sprintf`, "1 + 2 = 3")

	tstRunForth(t, "Mixed Types", `" Float: %.1f, Int: %d, Str: %s" << 1.5 42 " Foo" >>sprintf`, "Float: 1.5, Int: 42, Str: Foo")

	// Nested
	tstRunForth(t, "Nested",
		`" Outer: %s" << " Inner: %s" << " Val" >>sprintf >>sprintf`,
		"Outer: Inner: Val")

	// Test depth separately
	tstRunForth(t, "Depth", `1 2 3 depth`, 1, 2, 3, 3)
}

func TestEscapes(t *testing.T) {
	// Leading space in code is consumed.
	tstRunForth(t, "Newline", `" Line1\nLine2"`, "Line1\nLine2")
	tstRunForth(t, "Tab", `" Col1\tCol2"`, "Col1\tCol2")
	tstRunForth(t, "Quote", `" \"Quoted\""`, `"Quoted"`)
	tstRunForth(t, "Backslash", `" Back\\slash"`, `Back\slash`)
	tstRunForth(t, "Mixed", `" A\nB\tC"`, "A\nB\tC")

	// Unknown escape \a -> a
	tstRunForth(t, "Unknown Escape Becomes Char", `" \a"`, `a`)
}
