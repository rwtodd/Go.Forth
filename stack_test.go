package forth

import (
	"io/ioutil"
	"strings"
	"testing"
)

var vm = NewVM()

// tstRunForth is the main test helper.. it lets you
// run a code string, and FAILS if there is an error or
// if the stack doesn't match the values passed in.
func tstRunForth(t *testing.T, code string, vals ...interface{}) {
	if err := tstRunForthErr(t,code, vals...); err != nil {
		t.Error(err)
	}
}

// tstRunForthErr only fails the test if the stack doesn't match,
// giving the error back to the caller for investigation.
func tstRunForthErr(t *testing.T, code string, vals ...interface{}) error {
	vm.ResetState()
	tprog := strings.NewReader(code)
	err := vm.Run(tprog, ioutil.Discard)
	if !stackEq(vals...)  {
		t.Fail()
	}
	return err
}

// stackEq is a helper function checking the stack contents
// against the arguments.
func stackEq(vals ...interface{}) bool {
	if len(vals) != len(vm.Stack) {
		return false
	}

	for i := range vals {
		if vals[i] != vm.Stack[i] {
			return false
		}
	}
	return true
}

func TestDup(t *testing.T) {
    tstRunForth(t, `2 dup 3 dup`, 2, 2, 3, 3)
    tstRunForth(t, `2.2 dup dup`, 2.2, 2.2, 2.2)
}

func TestSwap(t *testing.T) {
    tstRunForth(t, `2 3 swap`, 3, 2)
    tstRunForth(t, `" hi"  43  swap`, 43, "hi")
}

func TestOver(t *testing.T) {
    tstRunForth(t, `2 3 OVER `, 2, 3, 2)
    if e := tstRunForthErr(t, `over`) ; e != ErrUnderflow {
		t.Error(e)
	}
}

func TestNip(t *testing.T) {
    tstRunForth(t, `1 2 3 nip 4 nip 5 nip`, 1, 5)
}

func TestTuck(t *testing.T) {
    tstRunForth(t, `1 2 tuck 3 tuck 4 tuck`, 
	            2, 1, 3, 2, 4, 3, 4) 
}

func TestDrop(t *testing.T) {
    tstRunForth(t, `2 3 drop`, 2)
    if e := tstRunForthErr(t, `drop drop`) ; e != ErrUnderflow {
		t.Error(e)
	}
}

func TestRot(t *testing.T) {
    if e := tstRunForthErr(t, `3 2  rot`, 3, 2) ; e != ErrUnderflow {
		t.Error(e)
    }
    tstRunForth(t, ` 2 3 4 rot `, 3, 4, 2)
    tstRunForth(t, ` 2 3.001 4 -rot `,  4, 2, 3.001)
    tstRunForth(t, ` 2 3.001 4 rot rot `,  4, 2, 3.001)
}
