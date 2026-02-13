// SPDX-License-Identifier: MIT

package forth

import (
	"testing"
)

func TestEmptyDict(t *testing.T) {
	// Just check that empty-dict works without error
	if e := tstRunForthErr(t, `empty-dict drop`); e != nil {
		t.Errorf("empty-dict failed: %v", e)
	}
}

func TestDictSetGet(t *testing.T) {
	// Set a value and get it back
	tstRunForth(t, `" value" empty-dict tuck " key" d! " key" d@`, "value")
}

func TestDictDelete(t *testing.T) {
	// Set a value, delete it, verify it's gone with d@|
	tstRunForth(t, `empty-dict dup " key" 42 -rot d! dup " key" ddel " key" 99 d@|`, 99)
}

func TestDictKeys(t *testing.T) {
	// Create dict with some keys and get keys
	tstRunForth(t, `empty-dict dup " a" 1 -rot d! dup " b" 2 -rot d! dup " c" 3 -rot d! dkeys @len`, 3)
}

func TestDictGetOr(t *testing.T) {
	// Test with existing key
	tstRunForth(t, `empty-dict dup " key" 42 -rot d! " key" 99 d@|`, 42)
	// Test with missing key
	tstRunForth(t, `empty-dict " missing" 99 d@|`, 99)
}

func TestDictGetQuery(t *testing.T) {
	// Test with existing key
	tstRunForth(t, `empty-dict dup " key" 42 -rot d! " key" d@?`, 42, -1)
	// Test with missing key
	tstRunForth(t, `empty-dict " missing" d@?`, 0)
}

func TestDictErrors(t *testing.T) {
	// Test d@ with missing key
	if e := tstRunForthErr(t, `empty-dict " missing" d@`); e != ErrKeyNotFound {
		t.Errorf("Expected ErrKeyNotFound, got %v", e)
	}
	// Test with wrong types
	if e := tstRunForthErr(t, `42 " key" d@`); e != ErrArgument {
		t.Errorf("Expected ErrArgument, got %v", e)
	}
	if e := tstRunForthErr(t, `empty-dict 42 d@`); e != ErrArgument {
		t.Errorf("Expected ErrArgument, got %v", e)
	}
}
