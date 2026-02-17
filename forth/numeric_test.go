// SPDX-License-Identifier: MIT

package forth

import (
	"testing"
)

func TestAdd(t *testing.T) {
	tstRunForth(t, "AddNum", `2 3 +  2 3.1 +`, 5, 5.1)
	tstRunForth(t, "AddStr", `" hi" "  there" +`, "hi there")
}

func TestMul(t *testing.T) {
	tstRunForth(t, "MulNum", "2 3 *   2 .25 *", 6, 0.5)
	tstRunForth(t, "MulStr", `" hi" 3 *  3 " yo" *`,
		"hihihi", "yoyoyo")
}

func TestPow(t *testing.T) {
	tstRunForth(t, "PowInt", `2 3 **`, 8.0)
	tstRunForth(t, "PowFloat", `2.0 3 **`, 8.0)
}

func TestSubtract(t *testing.T) {
	tstRunForth(t, "SubInt", `5 3 -`, 2)
	tstRunForth(t, "SubFloat", `5.0 3 -`, 2.0)
	tstRunForth(t, "SubNeg", `3 5 -`, -2)
}

func TestDivide(t *testing.T) {
	tstRunForth(t, "DivInt", `6 3 /`, 2)
	tstRunForth(t, "DivFloat", `5.0 2 /`, 2.5)
	tstRunForth(t, "DivTrunc", `3 2 /`, 1) // integer division
}

func TestSqrt(t *testing.T) {
	tstRunForth(t, "SqrtInt", `4 sqrt`, 2.0)
	tstRunForth(t, "SqrtFloat", `9.0 sqrt`, 3.0)
}

func TestLog(t *testing.T) {
	tstRunForth(t, "Log1", `1 log`, 0.0)
	// Note: log(e) ≈ 1, but due to floating point precision, we test with 1
}

func TestLog10(t *testing.T) {
	tstRunForth(t, "Log10_10", `10 log10`, 1.0)
	tstRunForth(t, "Log10_100", `100.0 log10`, 2.0)
}

func TestLog2(t *testing.T) {
	tstRunForth(t, "Log2_2", `2 log2`, 1.0)
	tstRunForth(t, "Log2_4", `4.0 log2`, 2.0)
}

func TestMax(t *testing.T) {
	tstRunForth(t, "MaxInt", `2 3 max`, 3)
	tstRunForth(t, "MaxFloat", `3.5 2 max`, 3.5)
	tstRunForth(t, "MaxEq", `4 4 max`, 4)
}

func TestMin(t *testing.T) {
	tstRunForth(t, "MinInt", `2 3 min`, 2)
	tstRunForth(t, "MinFloat", `3.5 2 min`, 2.0)
	tstRunForth(t, "MinEq", `4 4 min`, 4)
}

func TestSin(t *testing.T) {
	tstRunForth(t, "Sin0", `0 sin`, 0.0)
	// sin(π/2) ≈ 1, but let's use a small value
	tstRunForth(t, "SinPi2", `1.57079632679 sin`, 1.0) // approx π/2
}

func TestCos(t *testing.T) {
	tstRunForth(t, "Cos0", `0 cos`, 1.0)
}

func TestTan(t *testing.T) {
	tstRunForth(t, "Tan0", `0 tan`, 0.0)
}

func TestRound(t *testing.T) {
	tstRunForth(t, "RoundUp", `3.7 round`, 4.0)
	tstRunForth(t, "RoundDown", `3.2 round`, 3.0)
	tstRunForth(t, "RoundEq", `3 round`, 3.0)
}

func TestFloor(t *testing.T) {
	tstRunForth(t, "FloorUp", `3.7 floor`, 3.0)
	tstRunForth(t, "FloorDown", `3.2 floor`, 3.0)
	tstRunForth(t, "FloorEq", `3 floor`, 3.0)
}

func TestCeil(t *testing.T) {
	tstRunForth(t, "CeilUp", `3.7 ceil`, 4.0)
	tstRunForth(t, "CeilDown", `3.2 ceil`, 4.0)
	tstRunForth(t, "CeilEq", `3 ceil`, 3.0)
}
