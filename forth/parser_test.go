package forth

import (
	"testing"
)

func TestLiterals(t *testing.T) {
    tstRunForth(t, `: tst [ -5 10 * ] literal + ; 2 tst`, -48)
	tstRunForth(t, `: tst [ 55 1000 * ] literal + ; 2 tst`, 55002) 
}

func TestUpCase(t *testing.T) {
	tstRunForth(t, `: TST 3 4 sWaP ; tst`, 4, 3) 
}
