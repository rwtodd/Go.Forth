// SPDX-License-Identifier: MIT

package forth

import (
	"testing"
)

func TestBytes(t *testing.T) {
	tstRunForth(t, "BytesLen", `5 bytes @len`, 5)
}

func TestInts(t *testing.T) {
	tstRunForth(t, "IntsLen", `3 ints @len`, 3)
}

func TestFloats(t *testing.T) {
	tstRunForth(t, "FloatsLen", `4 floats @len`, 4)
}

func TestStrings(t *testing.T) {
	tstRunForth(t, "StringsLen", `2 strings @len`, 2)
}

func TestArrayGet(t *testing.T) {
	// Test int array
	tstRunForth(t, "IntGet", `3 ints dup 42 swap 0 ! 0 @`, 42)
	// Test float array
	tstRunForth(t, "FloatGet", `2 floats dup 3.14 swap 1 ! 1 @`, 3.14)
}

func TestArraySet(t *testing.T) {
	// Test int to int
	tstRunForth(t, "IntSet", `3 ints dup 42 swap 0 ! 0 @`, 42)
	// Test int to float (should work)
	tstRunForth(t, "FloatSet", `2 floats dup 5 swap 0 ! 0 @`, 5.0)
}

func TestByteGet(t *testing.T) {
	tstRunForth(t, "ByteGet", `5 bytes 255 over 0 c! 0 c@`, 255)
	// Test wrapping
	tstRunForth(t, "ByteWrap", `5 bytes 300 over 0 c! 0 c@`, 44) // 300 & 0xff = 44
}

func TestByteSet(t *testing.T) {
	tstRunForth(t, "ByteSet", `5 bytes 100 over 0 c! 0 c@`, 100)
	// Test wrapping
	tstRunForth(t, "ByteSetWrap", `5 bytes 500 over 0 c! 0 c@`, 244) // 500 & 0xff = 244
}

func TestArrayPush(t *testing.T) {
	// Push to int array
	tstRunForth(t, "PushInt", `3 ints 99 swap @push @len`, 4)
}

func TestArrayPop(t *testing.T) {
	// Pop from int array
	tstRunForth(t, "PopInt", `3 ints 42 over 0 ! @pop drop @len`, 2)
	// Pop value
	tstRunForth(t, "PopVal", `2 ints 99 over 1 ! @pop nip`, 99)
}

func TestArrayShift(t *testing.T) {
	// Shift from int array
	tstRunForth(t, "ShiftInt", `3 ints 42 over 0 ! @shift drop @len`, 2)
	// Shift value
	tstRunForth(t, "ShiftVal", `2 ints 77 over 0 ! @shift nip`, 77)
}

func TestArrayUnshift(t *testing.T) {
	// Unshift to int array
	tstRunForth(t, "UnshiftInt", `88 3 ints @unshift dup @len swap dup 0 @ swap @shift swap @len`, 4, 88, 88, 3)
	// Check first element
	tstRunForth(t, "UnshiftVal", `55 2 ints @unshift 0 @`, 55)
}

func TestArrayLen(t *testing.T) {
	tstRunForth(t, "LenByte", `10 bytes @len`, 10)
	tstRunForth(t, "LenInt", `7 ints @len`, 7)
	tstRunForth(t, "LenFloat", `3 floats @len`, 3)
	tstRunForth(t, "LenString", `4 strings @len`, 4)
}

func TestVariable(t *testing.T) {
	// Create variable and set value
	tstRunForth(t, "VarSet", `5 0 variable x 42 x !`, 5)
	// Get value
	tstRunForth(t, "VarGet", `0 variable y y @`, 0)
	// Update value
	tstRunForth(t, "VarUpdate", `0 variable z 99 z ! z @`, 99)
}

func TestVariableWithArrays(t *testing.T) {
	// Push to variable
	tstRunForth(t, "VarPush", `0 variable arr 3 ints arr ! 42 arr @push drop arr @len`, 4)
	// Pop from variable
	tstRunForth(t, "VarPop", `0 variable arr 3 ints arr ! arr @pop drop drop arr @len`, 2)
	// Test chaining @push on variable
	tstRunForth(t, "VarPushChain", `0 variable arr 3 ints arr ! 7 arr @push drop 6 arr @push drop 5 arr @push drop arr @len`, 6)
	tstRunForth(t, "VarPushChain2", `0 variable arr 3 ints arr ! 7 6 5 arr @push @push @push 5 @`, 7)
	tstRunForth(t, "VarPushChain3", `0 variable arr 0 ints arr ! 7 6 5 arr @push @push @push @pop swap @pop swap @pop nip`, 7, 6, 5)
}

func TestVariableArrayAccess(t *testing.T) {
	// Create variable with array and set element
	tstRunForth(t, "VarArrAccess", `0 variable myarr 3 ints myarr ! 42 myarr 0 ! myarr 0 @`, 42)
	// Set another element
	tstRunForth(t, "VarArrSet", `0 variable myarr 3 ints myarr ! 99 myarr 1 ! myarr 1 @`, 99)
	// Get default element
	tstRunForth(t, "VarArrDefault", `0 variable myarr 3 ints myarr ! myarr 2 @`, 0)
}
