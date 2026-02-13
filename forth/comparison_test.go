// SPDX-License-Identifier: MIT

package forth

import (
	"testing"
)

func TestComparisonOps(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected int // Expected top of stack (int boolean)
	}{
		// Equality
		{"EqIntTrue", "1 1 =", -1},
		{"EqIntFalse", "1 2 =", 0},
		{"EqFloatTrue", "1.5 1.5 =", -1},
		{"EqFloatFalse", "1.5 2.5 =", 0},
		{"EqIntFloatTrue", "1 1.0 =", -1},
		{"EqStringTrue", "\" abc\" \" abc\" =", -1},
		{"EqStringFalse", "\" abc\" \" def\" =", 0},

		// Less Than
		{"LtIntTrue", "1 2 <", -1},
		{"LtIntFalse", "2 1 <", 0},
		{"LtFloatTrue", "1.0 2.0 <", -1},
		{"LtIntFloatTrue", "1 2.0 <", -1},
		{"LtStringTrue", "\" abc\" \" def\" <", -1},
		{"LtStringFalse", "\" def\" \" abc\" <", 0},

		// Greater Than
		{"GtIntTrue", "2 1 >", -1},
		{"GtIntFalse", "1 2 >", 0},
		{"GtStringTrue", "\" def\" \" abc\" >", -1},

		// Not Equal
		{"NeIntTrue", "1 2 <>", -1},
		{"NeIntFalse", "1 1 <>", 0},

		// Zero Checks
		{"ZeroEqTrue", "0 0=", -1},
		{"ZeroEqFalse", "1 0=", 0},
		{"ZeroLtTrue", "-1 0<", -1},
		{"ZeroLtFalse", "0 0<", 0},
		{"ZeroGtTrue", "1 0>", -1},
		{"ZeroGtFalse", "0 0>", 0},
	}

	for _, tt := range tests {
		tstRunForth(t, tt.name, tt.input, tt.expected)
	}
}

func TestRefinedComparisons(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected int
	}{
		{"LteIntTrue", "1 1 <=", -1},
		{"LteIntTrue2", "1 2 <=", -1},
		{"LteIntFalse", "2 1 <=", 0},

		{"GteIntTrue", "1 1 >=", -1},
		{"GteIntTrue2", "2 1 >=", -1},
		{"GteIntFalse", "1 2 >=", 0},
	}

	for _, tt := range tests {
		tstRunForth(t, tt.name, tt.input, tt.expected)
	}
}

func TestBitwiseOps(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected int
	}{
		{"And", "-1 -1 and", -1},
		{"AndMix", "-1 0 and", 0},
		{"Or", "-1 0 or", -1},
		{"OrFalse", "0 0 or", 0},
		{"Xor", "-1 -1 xor", 0},
		{"XorTrue", "-1 0 xor", -1},
		{"Invert", "0 invert", -1},
		{"InvertFalse", "-1 invert", 0},
	}

	for _, tt := range tests {
		tstRunForth(t, tt.name, tt.input, tt.expected)
	}
}
