package forth

import (
	"math"
	"testing"
)

func TestFloat16Conversion(t *testing.T) {
	tests := []struct {
		val      float64
		expected uint16
		lossless bool
	}{
		{0.0, 0x0000, true},
		{math.Copysign(0, -1), 0x8000, true},
		{1.0, 0x3c00, true},
		{-1.0, 0xbc00, true},
		{1.5, 0x3e00, true},
		{0.5, 0x3800, true},
		{2.0, 0x4000, true},
		{3.0, 0x4200, true},
		{10.0, 0x4900, true},
		{-2.0, 0xc000, true},
		{65504.0, 0x7bff, true},          // Max normal
		{-65504.0, 0xfbff, true},         // Min normal
		{0.00006103515625, 0x0400, true}, // Smallest normal (2^-14)
		{math.Pow(2, -24), 0x0001, true}, // Smallest subnormal

		// Edge cases where it should NOT be lossless
		{0.1, 0, false},          // 0.1 cannot be exactly represented
		{100000.0, 0, false},     // Too large
		{-100000.0, 0, false},    // Too small
		{math.NaN(), 0, false},   // NaN
		{math.Inf(1), 0, false},  // Inf
		{math.Inf(-1), 0, false}, // -Inf
	}

	for _, tt := range tests {
		got, ok := float64ToFloat16(tt.val)
		if ok != tt.lossless {
			t.Errorf("float64ToFloat16(%g) lossless = %v, want %v", tt.val, ok, tt.lossless)
		}
		if ok {
			if got != tt.expected {
				t.Errorf("float64ToFloat16(%g) = 0x%04x, want 0x%04x", tt.val, got, tt.expected)
			}

			// Round trip check
			back := float16ToFloat64(got)
			if back != tt.val {
				// Handle -0.0 case specifically since direct comparison might be tricky with some test frameworks,
				// but Go == handles 0.0 == -0.0 as true. Wait, we want to distinguish signed zero?
				// Actually 0.0 == -0.0 is true in Go.
				// But we encoded sign bit.
				// Let's check bitwise for 0.0 and -0.0?
				// float16ToFloat64 should return the correct signed zero.
				// 0.0 == -0.0 in value.
				if tt.val == 0.0 && back == 0.0 {
					// Check sign bit
					if math.Signbit(tt.val) != math.Signbit(back) {
						t.Errorf("float16ToFloat64(0x%04x) sign mismatch. Got %g, want %g", got, back, tt.val)
					}
				} else {
					t.Errorf("Round trip failed: float16ToFloat64(0x%04x) = %g, want %g", got, back, tt.val)
				}
			}
		}
	}
}

func TestDecimalConversion(t *testing.T) {
	tests := []struct {
		val      float64
		expected uint16
		lossless bool
	}{
		// Scale 0 (div 10)
		{0.1, 0x0001, true},
		{1.2, 0x000c, true},
		{10.0, 0x0064, true},
		{819.1, 0x1fff, true},
		{-1.2, 0x800c, true},

		// Scale 1 (div 100)
		{0.01, 0x2001, true},
		{1.23, 0x207b, true},
		{81.91, 0x3fff, true},

		// Scale 2 (div 1000)
		{0.001, 0x4001, true},
		{3.141, 0x4c45, true},
		{8.191, 0x5fff, true},

		// Scale 3 (div 10000)
		{0.0001, 0x6001, true},
		{0.8191, 0x7fff, true},

		// Out of range / precision loss
		{820.0, 0, false},   // > 819.1
		{0.00001, 0, false}, // < 0.0001 step
		{1.2345, 0, false},  // Needs > 4 decimals

		// Zero is handled by float16, decimal shouldn't claim it necessarily, but valid 0 is valid.
		// Our logic might say 0.0 -> 0x0000 (0/10).
		{0.0, 0x0000, true},
	}

	for _, tt := range tests {
		got, ok := float64ToDecimal16(tt.val)
		if ok != tt.lossless {
			t.Errorf("float64ToDecimal16(%g) lossless = %v, want %v", tt.val, ok, tt.lossless)
		}
		if ok {
			if got != tt.expected {
				t.Errorf("float64ToDecimal16(%g) = 0x%04x, want 0x%04x", tt.val, got, tt.expected)
			}

			// Round trip
			back := decimal16ToFloat64(got)
			if math.Abs(back-tt.val) > 1e-9 {
				t.Errorf("Round trip failed: decimal16ToFloat64(0x%04x) = %g, want %g", got, back, tt.val)
			}
		}
	}
}
