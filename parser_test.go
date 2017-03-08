package forth

import (
	"io/ioutil"
	"strings"
	"testing"
)

func TestLiterals(t *testing.T) {
	vm.ResetState()
	var tprog = strings.NewReader(`: tst [ -5 10 * ] literal + ; 2 tst`)
	if err := vm.Run(tprog, ioutil.Discard); err != nil {
		t.Error(err)
	}
	if !stackEq(-48) {
		t.Fail()
	}

	vm.ResetState()
	tprog = strings.NewReader(`: tst [ 55 1000 * ] literal + ; 2 tst`)
	if err := vm.Run(tprog, ioutil.Discard); err != nil {
		t.Error(err)
	}
	if !stackEq(55002) {
		t.Fail()
	}
}

func TestUpCase(t *testing.T) {
	vm.ResetState()
	var tprog = strings.NewReader(`: TST 3 4 sWaP ; tst`)
	if err := vm.Run(tprog, ioutil.Discard); err != nil {
		t.Error(err)
	}
	if !stackEq(4, 3) {
		t.Fail()
	}
}
