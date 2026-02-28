// SPDX-License-Identifier: MIT

package forth

import (
	"testing"
)

func TestOrd(t *testing.T) {
	tstRunForth(t, "basic", `" a" ord`, 97)

	tstRunForth(t, "non-ascii", `" é" ord`, 233)

	tstRunForthErr(t, "empty", `" " ord`, ErrArgument)

	tstRunForthErr(t, "not-string", `123 ord`, ErrArgument)

	tstRunForth(t, "multi-char", `" abc" ord`, 97)
}

func TestChr(t *testing.T) {
	tstRunForth(t, "basic", `97 chr`, "a")

	tstRunForth(t, "non-ascii", `233 chr`, "é")

	tstRunForthErr(t, "not-int", `"a" chr`, ErrArgument)
}
