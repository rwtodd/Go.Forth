// SPDX-License-Identifier: MIT

package forth

import (
	"math"
	"strings"
)

// : + ( a b -- a+b ) <code>
func add(vm *VM) (err error) {
	top := len(vm.Stack) - 1
	if top < 1 {
		return ErrUnderflow
	}
	switch op1 := vm.Stack[top].(type) {
	case int64:
		switch op2 := vm.Stack[top-1].(type) {
		case int64:
			vm.Stack[top-1] = op1 + op2
		case float64:
			vm.Stack[top-1] = float64(op1) + op2
		default:
			err = ErrArgumentMsg("+ requires numeric arguments")
		}
	case float64:
		switch op2 := vm.Stack[top-1].(type) {
		case int64:
			vm.Stack[top-1] = op1 + float64(op2)
		case float64:
			vm.Stack[top-1] = op1 + op2
		default:
			err = ErrArgumentMsg("+ requires numeric arguments")
		}
	case string:
		op2, ok := vm.Stack[top-1].(string)
		if ok {
			vm.Stack[top-1] = op2 + op1
		} else {
			err = ErrArgumentMsg("+ requires two strings or two numbers")
		}
	default:
		err = ErrArgumentMsg("invalid types for +")
	}
	vm.Stack = vm.Stack[:top]
	return
}

// : ** ( a b -- a**b ) <code>
func raiseToPower(vm *VM) (err error) {
	top := len(vm.Stack) - 1
	if top < 1 {
		return ErrUnderflow
	}
	switch op1 := vm.Stack[top].(type) {
	case int64:
		switch op2 := vm.Stack[top-1].(type) {
		case int64:
			vm.Stack[top-1] = math.Pow(float64(op2), float64(op1))
		case float64:
			vm.Stack[top-1] = math.Pow(op2, float64(op1))
		default:
			err = ErrArgumentMsg("** requires numeric arguments")
		}
	case float64:
		switch op2 := vm.Stack[top-1].(type) {
		case int64:
			vm.Stack[top-1] = math.Pow(float64(op2), op1)
		case float64:
			vm.Stack[top-1] = math.Pow(op2, op1)
		default:
			err = ErrArgumentMsg("** requires numeric arguments")
		}
	default:
		err = ErrArgumentMsg("** requires numeric arguments")
	}
	vm.Stack = vm.Stack[:top]
	return
}

// : * ( a b -- a*b ) <code>
func multiply(vm *VM) (err error) {
	top := len(vm.Stack) - 1
	if top < 1 {
		return ErrUnderflow
	}
	switch op1 := vm.Stack[top].(type) {
	case int64:
		switch op2 := vm.Stack[top-1].(type) {
		case int64:
			vm.Stack[top-1] = op1 * op2
		case float64:
			vm.Stack[top-1] = float64(op1) * op2
		case string:
			vm.Stack[top-1] = strings.Repeat(op2, int(op1))
		default:
			err = ErrArgumentMsg("* requires numeric arguments or string+int")
		}
	case float64:
		switch op2 := vm.Stack[top-1].(type) {
		case int64:
			vm.Stack[top-1] = op1 * float64(op2)
		case float64:
			vm.Stack[top-1] = op1 * op2
		default:
			err = ErrArgumentMsg("* requires numeric arguments")
		}
	case string:
		op2, ok := vm.Stack[top-1].(int64)
		if ok {
			vm.Stack[top-1] = strings.Repeat(op1, int(op2))
		} else {
			err = ErrArgumentMsg("* string repetition requires integer count")
		}
	default:
		err = ErrArgumentMsg("invalid types for *")
	}
	vm.Stack = vm.Stack[:top]
	return
}

// : - ( a b -- a-b ) <code>
func subtract(vm *VM) (err error) {
	top := len(vm.Stack) - 1
	if top < 1 {
		return ErrUnderflow
	}
	switch op1 := vm.Stack[top].(type) {
	case int64:
		switch op2 := vm.Stack[top-1].(type) {
		case int64:
			vm.Stack[top-1] = op2 - op1
		case float64:
			vm.Stack[top-1] = op2 - float64(op1)
		default:
			err = ErrArgumentMsg("- requires numeric arguments")
		}
	case float64:
		switch op2 := vm.Stack[top-1].(type) {
		case int64:
			vm.Stack[top-1] = float64(op2) - op1
		case float64:
			vm.Stack[top-1] = op2 - op1
		default:
			err = ErrArgumentMsg("- requires numeric arguments")
		}
	default:
		err = ErrArgumentMsg("- requires numeric arguments")
	}
	vm.Stack = vm.Stack[:top]
	return
}

// : / ( a b -- a/b ) <code>
func divide(vm *VM) (err error) {
	top := len(vm.Stack) - 1
	if top < 1 {
		return ErrUnderflow
	}
	switch op1 := vm.Stack[top].(type) {
	case int64:
		switch op2 := vm.Stack[top-1].(type) {
		case int64:
			if op1 == 0 {
				err = ErrArgumentMsg("division by zero")
			} else {
				vm.Stack[top-1] = op2 / op1
			}
		case float64:
			if op1 == 0 {
				err = ErrArgumentMsg("division by zero")
			} else {
				vm.Stack[top-1] = op2 / float64(op1)
			}
		default:
			err = ErrArgumentMsg("/ requires numeric arguments")
		}
	case float64:
		switch op2 := vm.Stack[top-1].(type) {
		case int64:
			if op1 == 0 {
				err = ErrArgumentMsg("division by zero")
			} else {
				vm.Stack[top-1] = float64(op2) / op1
			}
		case float64:
			if op1 == 0 {
				err = ErrArgumentMsg("division by zero")
			} else {
				vm.Stack[top-1] = op2 / op1
			}
		default:
			err = ErrArgumentMsg("/ requires numeric arguments")
		}
	default:
		err = ErrArgumentMsg("/ requires numeric arguments")
	}
	vm.Stack = vm.Stack[:top]
	return
}

// : sqrt ( a -- sqrt(a) ) <code>
func sqrt(vm *VM) (err error) {
	if len(vm.Stack) < 1 {
		return ErrUnderflow
	}
	switch op := vm.Stack[len(vm.Stack)-1].(type) {
	case int64:
		if op < 0 {
			err = ErrArgumentMsg("sqrt requires non-negative argument")
		} else {
			vm.Stack[len(vm.Stack)-1] = math.Sqrt(float64(op))
		}
	case float64:
		if op < 0 {
			err = ErrArgumentMsg("sqrt requires non-negative argument")
		} else {
			vm.Stack[len(vm.Stack)-1] = math.Sqrt(op)
		}
	default:
		err = ErrArgumentMsg("sqrt requires numeric argument")
	}
	return
}

// : log ( a -- log(a) ) <code>
func log(vm *VM) (err error) {
	if len(vm.Stack) < 1 {
		return ErrUnderflow
	}
	switch op := vm.Stack[len(vm.Stack)-1].(type) {
	case int64:
		if op <= 0 {
			err = ErrArgumentMsg("log requires positive argument")
		} else {
			vm.Stack[len(vm.Stack)-1] = math.Log(float64(op))
		}
	case float64:
		if op <= 0 {
			err = ErrArgumentMsg("log requires positive argument")
		} else {
			vm.Stack[len(vm.Stack)-1] = math.Log(op)
		}
	default:
		err = ErrArgumentMsg("log requires numeric argument")
	}
	return
}

// : log10 ( a -- log10(a) ) <code>
func log10(vm *VM) (err error) {
	if len(vm.Stack) < 1 {
		return ErrUnderflow
	}
	switch op := vm.Stack[len(vm.Stack)-1].(type) {
	case int64:
		if op <= 0 {
			err = ErrArgumentMsg("log10 requires positive argument")
		} else {
			vm.Stack[len(vm.Stack)-1] = math.Log10(float64(op))
		}
	case float64:
		if op <= 0 {
			err = ErrArgumentMsg("log10 requires positive argument")
		} else {
			vm.Stack[len(vm.Stack)-1] = math.Log10(op)
		}
	default:
		err = ErrArgumentMsg("log10 requires numeric argument")
	}
	return
}

// : log2 ( a -- log2(a) ) <code>
func log2(vm *VM) (err error) {
	if len(vm.Stack) < 1 {
		return ErrUnderflow
	}
	switch op := vm.Stack[len(vm.Stack)-1].(type) {
	case int64:
		if op <= 0 {
			err = ErrArgumentMsg("log2 requires positive argument")
		} else {
			vm.Stack[len(vm.Stack)-1] = math.Log2(float64(op))
		}
	case float64:
		if op <= 0 {
			err = ErrArgumentMsg("log2 requires positive argument")
		} else {
			vm.Stack[len(vm.Stack)-1] = math.Log2(op)
		}
	default:
		err = ErrArgumentMsg("log2 requires numeric argument")
	}
	return
}

// : max ( a b -- max(a,b) ) <code>
func max(vm *VM) (err error) {
	top := len(vm.Stack) - 1
	if top < 1 {
		return ErrUnderflow
	}
	switch op1 := vm.Stack[top].(type) {
	case int64:
		switch op2 := vm.Stack[top-1].(type) {
		case int64:
			if op2 > op1 {
				vm.Stack[top-1] = op2
			} else {
				vm.Stack[top-1] = op1
			}
		case float64:
			vm.Stack[top-1] = math.Max(op2, float64(op1))
		default:
			err = ErrArgumentMsg("max requires numeric arguments")
		}
	case float64:
		switch op2 := vm.Stack[top-1].(type) {
		case int64:
			vm.Stack[top-1] = math.Max(float64(op2), op1)
		case float64:
			vm.Stack[top-1] = math.Max(op2, op1)
		default:
			err = ErrArgumentMsg("max requires numeric arguments")
		}
	default:
		err = ErrArgumentMsg("max requires numeric arguments")
	}
	vm.Stack = vm.Stack[:top]
	return
}

// : min ( a b -- min(a,b) ) <code>
func min(vm *VM) (err error) {
	top := len(vm.Stack) - 1
	if top < 1 {
		return ErrUnderflow
	}
	switch op1 := vm.Stack[top].(type) {
	case int64:
		switch op2 := vm.Stack[top-1].(type) {
		case int64:
			if op2 < op1 {
				vm.Stack[top-1] = op2
			} else {
				vm.Stack[top-1] = op1
			}
		case float64:
			vm.Stack[top-1] = math.Min(op2, float64(op1))
		default:
			err = ErrArgumentMsg("min requires numeric arguments")
		}
	case float64:
		switch op2 := vm.Stack[top-1].(type) {
		case int64:
			vm.Stack[top-1] = math.Min(float64(op2), op1)
		case float64:
			vm.Stack[top-1] = math.Min(op2, op1)
		default:
			err = ErrArgumentMsg("min requires numeric arguments")
		}
	default:
		err = ErrArgumentMsg("min requires numeric arguments")
	}
	vm.Stack = vm.Stack[:top]
	return
}

// : sin ( a -- sin(a) ) <code>
func sin(vm *VM) (err error) {
	if len(vm.Stack) < 1 {
		return ErrUnderflow
	}
	switch op := vm.Stack[len(vm.Stack)-1].(type) {
	case int64:
		vm.Stack[len(vm.Stack)-1] = math.Sin(float64(op))
	case float64:
		vm.Stack[len(vm.Stack)-1] = math.Sin(op)
	default:
		err = ErrArgumentMsg("sin requires numeric argument")
	}
	return
}

// : cos ( a -- cos(a) ) <code>
func cos(vm *VM) (err error) {
	if len(vm.Stack) < 1 {
		return ErrUnderflow
	}
	switch op := vm.Stack[len(vm.Stack)-1].(type) {
	case int64:
		vm.Stack[len(vm.Stack)-1] = math.Cos(float64(op))
	case float64:
		vm.Stack[len(vm.Stack)-1] = math.Cos(op)
	default:
		err = ErrArgumentMsg("cos requires numeric argument")
	}
	return
}

// : tan ( a -- tan(a) ) <code>
func tan(vm *VM) (err error) {
	if len(vm.Stack) < 1 {
		return ErrUnderflow
	}
	switch op := vm.Stack[len(vm.Stack)-1].(type) {
	case int64:
		vm.Stack[len(vm.Stack)-1] = math.Tan(float64(op))
	case float64:
		vm.Stack[len(vm.Stack)-1] = math.Tan(op)
	default:
		err = ErrArgumentMsg("tan requires numeric argument")
	}
	return
}

// : round ( a -- round(a) ) <code>
func round(vm *VM) (err error) {
	if len(vm.Stack) < 1 {
		return ErrUnderflow
	}
	switch op := vm.Stack[len(vm.Stack)-1].(type) {
	case int64:
		// already int64
	case float64:
		vm.Stack[len(vm.Stack)-1] = int64(math.Round(op))
	default:
		err = ErrArgumentMsg("round requires numeric argument")
	}
	return
}

// : floor ( a -- floor(a) ) <code>
func floor(vm *VM) (err error) {
	if len(vm.Stack) < 1 {
		return ErrUnderflow
	}
	switch op := vm.Stack[len(vm.Stack)-1].(type) {
	case int64:
		// already int64
	case float64:
		vm.Stack[len(vm.Stack)-1] = int64(math.Floor(op))
	default:
		err = ErrArgumentMsg("floor requires numeric argument")
	}
	return
}

// : ceil ( a -- ceil(a) ) <code>
func ceil(vm *VM) (err error) {
	if len(vm.Stack) < 1 {
		return ErrUnderflow
	}
	switch op := vm.Stack[len(vm.Stack)-1].(type) {
	case int64:
		// already int64
	case float64:
		vm.Stack[len(vm.Stack)-1] = int64(math.Ceil(op))
	default:
		err = ErrArgumentMsg("ceil requires numeric argument")
	}
	return
}

// numWordsInit adds numeric core words to the VM
func numWordsInit(vm *VM) {
	vm.Define(&NativeWord{name: "+", run: add, immediate: false})
	vm.Define(&NativeWord{name: "-", run: subtract, immediate: false})
	vm.Define(&NativeWord{name: "*", run: multiply, immediate: false})
	vm.Define(&NativeWord{name: "**", run: raiseToPower, immediate: false})
	vm.Define(&NativeWord{name: "/", run: divide, immediate: false})
	vm.Define(&NativeWord{name: "sqrt", run: sqrt, immediate: false})
	vm.Define(&NativeWord{name: "log", run: log, immediate: false})
	vm.Define(&NativeWord{name: "log10", run: log10, immediate: false})
	vm.Define(&NativeWord{name: "log2", run: log2, immediate: false})
	vm.Define(&NativeWord{name: "max", run: max, immediate: false})
	vm.Define(&NativeWord{name: "min", run: min, immediate: false})
	vm.Define(&NativeWord{name: "sin", run: sin, immediate: false})
	vm.Define(&NativeWord{name: "cos", run: cos, immediate: false})
	vm.Define(&NativeWord{name: "tan", run: tan, immediate: false})
	vm.Define(&NativeWord{name: "round", run: round, immediate: false})
	vm.Define(&NativeWord{name: "floor", run: floor, immediate: false})
	vm.Define(&NativeWord{name: "ceil", run: ceil, immediate: false})
}

// float16ToFloat64 converts a 16-bit floating point number to a float64
// Format: 1 bit sign, 5 bits exponent, 10 bits mantissa
func float16ToFloat64(u uint16) float64 {
	sign := uint64((u >> 15) & 0x1)
	exp := uint64((u >> 10) & 0x1f)
	mant := uint64(u & 0x3ff) // 0x3ff is 10 bits

	var fBits uint64

	if exp == 0 {
		if mant == 0 {
			// Zero
			fBits = sign << 63
		} else {
			// Subnormal
			// Value is (-1)^S * 2^-14 * (mant / 1024)
			// float64 subnormals are much smaller.
			// Just use float arithmetic for subnormals, it's easier and safe given float64 precision.
			val := float64(u&0x3ff) * math.Pow(2, -24) // 2^-14 * 2^-10 = 2^-24
			if sign == 1 {
				val = -val
			}
			return val
		}
	} else if exp == 31 {
		// Inf or NaN
		if mant == 0 {
			// Infinity
			fBits = (sign << 63) | (0x7ff << 52)
		} else {
			// NaN - preserve payload if possible, or just return generic NaN
			// generic NaN is fine
			fBits = (sign << 63) | (0x7ff << 52) | (1 << 51) | mant // Keep mantissa at bottom
		}
	} else {
		// Normalized
		// exp16 has bias 15. exp64 has bias 1023.
		// exp64 = exp16 - 15 + 1023 = exp16 + 1008
		newExp := exp + 1008
		newMant := mant << 42 // Shift 10 bits to top of 52-bit mantissa
		fBits = (sign << 63) | (newExp << 52) | newMant
	}

	if fBits != 0 || (sign != 0 && exp == 0 && mant == 0) { // optimization for non-calculated cases
		return math.Float64frombits(fBits)
	}
	return 0.0
}

// float64ToFloat16 tries to convert a float64 to a 16-bit float.
// Returns the uint16 representation and true if the conversion is lossless,
// false otherwise.
func float64ToFloat16(f float64) (uint16, bool) {
	bits := math.Float64bits(f)
	sign := uint16((bits >> 63) & 0x1)
	exp64 := (bits >> 52) & 0x7ff
	mant64 := bits & 0xfffffffffffff

	// 1. Check for Zero
	if exp64 == 0 && mant64 == 0 {
		return (sign << 15), true
	}

	// 2. Check for Inf/NaN
	if exp64 == 0x7ff {
		// We don't optimize Inf/NaN for now, relying on codeseg literal.
		// Although we COULD, let's treat them as non-optimizable to keep it simple/safe
		// unless specifically requested. User said "literals that can be accurately represented".
		// NaNs are tricky with payloads. Inf is fine.
		// Let's stick to Finite numbers for safety.
		return 0, false
	}

	// 3. Normalized numbers
	// target exp16 = exp64 - 1023 + 15 = exp64 - 1008
	exp16 := int(exp64) - 1008

	// Check exponent range for normalized float16 (1 to 30)
	if exp16 > 30 {
		return 0, false // Overflow
	}

	if exp16 > 0 {
		// Normalized in Float16
		// Check mantissa precision
		// We need the lower 42 bits of mant64 to be zero.
		if (mant64 & 0x3ffffffffff) != 0 {
			return 0, false // Precision loss
		}
		mant16 := uint16((mant64 >> 42) & 0x3ff)
		return (sign << 15) | (uint16(exp16) << 10) | mant16, true
	}

	// 4. Subnormals
	// exp16 is <= 0.
	// Float16 subnormals have conceptual exponent -14 (bias 15, encoded 0).
	// Current value is 1.mantissa * 2^(exp16 - 15)  (using raw exp16 value calculation isn't quite right here)
	// Let's rephrase:
	// Value = 1.mant64 * 2^(exp64 - 1023)
	// We want to represent as 0.mant16 * 2^-14
	// = 0.mant16 * 2^(1 - 15)

	// Shift needed:
	// 2^(exp64 - 1023) = M * 2^-14
	// M = 2^(exp64 - 1023 + 14) * (1 + mant64_fraction)
	// Exponent difference:
	// shift = -14 - (exp64 - 1023) = 1009 - exp64

	// Since exp16 = exp64 - 1008, and exp16 <= 0
	// exp64 <= 1008.

	// Example: exp64 = 1008 (exp16=0, 2^-15). Normal float16 range is 2^-14.
	// Smallest normal is 1.0 * 2^-14.
	// 1.0 * 2^-15 becomes 0.5 * 2^-14 -> mantissa 0.1000...

	// We need to shift the virtual "1" of the implicit mantissa right.
	// The implicit mantissa bits are '1' followed by the 52 mant64 bits.
	// full_mantissa = (1 << 52) | mant64.
	// We want to shift this right to fit into the 10 bits of float16 mantissa,
	// such that the effective exponent becomes -14.

	// exp_diff = -14 - (exp64 - 1023) = 1009 - exp64.
	// This is how many places we shift the binary point to the left (multiply by 2^-diff)
	// Or equivalent to right shifting the integer mantissa.

	// We have 53 bits of significant (1 implicit + 52 explicit).
	// We need to compress to 10 bits.
	// total shift = shift + (52 - 10) = shift + 42 ?
	// Let's verify.
	// Target: 10 bit integer M. Value = M * 2^-24.
	// Source: V = (1 + mant64/2^52) * 2^(exp64 - 1023).
	// V = ((1<<52 + mant64) / 2^52) * 2^(exp64 - 1023)
	//   = ((1<<52 + mant64) * 2^(exp64 - 1023 - 52))
	//   = ((1<<52 + mant64) * 2^(exp64 - 1075))

	// Target V = M * 2^-24.
	// M = V * 2^24
	//   = (1<<52 + mant64) * 2^(exp64 - 1075 + 24)
	//   = (1<<52 + mant64) * 2^(exp64 - 1051)

	// So we need to shift right by (1051 - exp64) bits.
	rshift := 1051 - int(exp64)

	if rshift < 0 {
		// Should be impossible if logic above matches normalized check?
		// exp64 > 1051 implies huge number. We already checked exp16 <= 0.
		// exp16 = exp64 - 1008. If exp16 <= 0, exp64 <= 1008.
		// So rshift >= 43.
		return 0, false
	}

	if rshift > 63 {
		// Too small to represent
		return 0, false
	}

	fullMant := (uint64(1) << 52) | mant64

	// Check if we lose precision
	// We are shifting right by rshift. We must ensure dropped bits are 0.
	mask := (uint64(1) << rshift) - 1
	if (fullMant & mask) != 0 {
		return 0, false // Precision loss
	}

	mant16 := fullMant >> rshift

	// mant16 must fit in 10 bits for subnormal
	if mant16 > 0x3ff {
		// Should prevent this via previous checks, but sanity check.
		// If it became normalized again?
		// e.g. if we rounded up? But we don't round.
		return 0, false
	}

	return (sign << 15) | uint16(mant16), true
}

// decimal16ToFloat64 converts a 16-bit decimal encoding to float64
// Format:
// Bit 15: Sign
// Bits 13-14: Scale (0=10, 1=100, 2=1000, 3=10000)
// Bits 0-12: Significand (0-8191)
func decimal16ToFloat64(u uint16) float64 {
	sign := (u >> 15) & 0x1
	scaleIdx := (u >> 13) & 0x3
	val := float64(u & 0x1fff)

	var divisor float64
	switch scaleIdx {
	case 0:
		divisor = 10.0
	case 1:
		divisor = 100.0
	case 2:
		divisor = 1000.0
	case 3:
		divisor = 10000.0
	}

	result := val / divisor
	if sign == 1 {
		result = -result
	}
	return result
}

// float64ToDecimal16 tries to convert a float64 to a 16-bit decimal encoding.
func float64ToDecimal16(f float64) (uint16, bool) {
	// 1. Handle Sign
	sign := uint16(0)
	if f < 0 {
		sign = 1
		f = -f
	} else if f == 0 {
		// Zero is always scale 0, val 0
		if math.Signbit(f) {
			return 0x8000, true // -0.0
		}
		return 0x0000, true
	}

	// 2. Iterate Scales to find fit
	// We want the lowest scale that works? Or highest?
	// Actually any scale that works is fine, but lower divisor means larger range...
	// Wait, divisor 10 range is 819.1. Divisor 10000 range is 0.8191.
	// So for a standard number like 1.2, scale 0 (12/10) works. Scale 1 (120/100) works.
	// We should probably prefer smaller divisor (larger range)?
	// Or maybe just iterate 0..3 and verify "is integer" and "fits in 13 bits".

	divisors := []float64{10.0, 100.0, 1000.0, 10000.0}

	for i, div := range divisors {
		scaled := f * div
		// Check if close to integer
		// 1.2 * 10 = 12.0.
		// 1.23 * 10 = 12.3 -> not integer.

		candidates := []float64{math.Floor(scaled), math.Ceil(scaled)}
		for _, cand := range candidates {
			if math.Abs(scaled-cand) < 1e-9 {
				// It is effectively an integer 'cand'
				// Check limits
				if cand <= 8191 {
					// Found it!
					scalebits := uint16(i) << 13
					valbits := uint16(cand)
					return (sign << 15) | scalebits | valbits, true
				}
			}
		}
	}

	return 0, false
}
