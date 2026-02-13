// SPDX-License-Identifier: MIT

package forth

import (
	"testing"
)

func TestEmptyDict(t *testing.T) {
	// Just check that empty-dict works without error
	tstRunForth(t, "EmptyDict", `empty-dict drop`)
}

func TestDictSetGet(t *testing.T) {
	// Set a value and get it back
	tstRunForth(t, "SetGet", `" value" empty-dict tuck " key" d! " key" d@`, "value")
}

func TestDictDelete(t *testing.T) {
	// Set a value, delete it, verify it's gone with d@|
	tstRunForth(t, "Delete", `empty-dict dup " key" 42 -rot d! dup " key" ddel " key" 99 d@|`, 99)
}

func TestDictKeys(t *testing.T) {
	// Create dict with some keys and get keys
	tstRunForth(t, "Keys", `empty-dict dup " a" 1 -rot d! dup " b" 2 -rot d! dup " c" 3 -rot d! dkeys @len`, 3)
}

func TestDictGetOr(t *testing.T) {
	// Test with existing key
	tstRunForth(t, "GetOrExist", `empty-dict dup " key" 42 -rot d! " key" 99 d@|`, 42)
	// Test with missing key
	tstRunForth(t, "GetOrMissing", `empty-dict " missing" 99 d@|`, 99)
}

func TestDictGetQuery(t *testing.T) {
	// Test with existing key
	tstRunForth(t, "QueryExist", `empty-dict dup " key" 42 -rot d! " key" d@?`, 42, -1)
	// Test with missing key
	tstRunForth(t, "QueryMissing", `empty-dict " missing" d@?`, 0)
}

func TestDictErrors(t *testing.T) {
	// Test d@ with missing key
	tstRunForthErr(t, "MissingKey", `empty-dict " missing" d@`, ErrKeyNotFound)
	// Test with wrong types
	tstRunForthErr(t, "WrongArgs1", `42 " key" d@`, ErrArgument)
	tstRunForthErr(t, "WrongArgs2", `empty-dict 42 d@`, ErrArgument)
}
