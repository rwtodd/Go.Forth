// SPDX-License-Identifier: MIT

package forth

import (
	"testing"
)

func TestThings(t *testing.T) {
	tstRunForth(t, "ThingsLen", `5 things @len`, 5)
}

func TestAnyArrayGet(t *testing.T) {
	tstRunForth(t, "AnyGet", `3 things dup 42 swap 0 ! 0 @`, 42)
	tstRunForth(t, "AnyGetString", `2 things dup " hi" swap 1 ! 1 @`, "hi")
}

func TestAnyArraySet(t *testing.T) {
	tstRunForth(t, "AnySet", `3 things dup 42 swap 0 ! 0 @`, 42)
	tstRunForth(t, "AnySetMix", `3 things dup 42 swap 0 ! dup " hi" swap 1 ! dup 1 @ swap 0 @`, "hi", 42)
}

func TestAnyArrayPush(t *testing.T) {
	tstRunForth(t, "AnyPush", `0 things " first" @push 42 @push dup @len swap 0 @`, 2, "first")
}

func TestAnyArrayPop(t *testing.T) {
	tstRunForth(t, "AnyPop", `0 things 42 @push @pop drop @len`, 0)
	tstRunForth(t, "AnyPopVal", `0 things 42 @push @pop nip`, 42)
}
