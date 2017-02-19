package forth

import (
	"io/ioutil"
	"strings"
	"testing"
)

var vm = NewVM()

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
	vm.ResetState()
	var tprog = strings.NewReader("2 dup 3 dup")
	if err := vm.Run(tprog, ioutil.Discard); err != nil {
		t.Error(err)
	}
	if !stackEq(2, 2, 3, 3) {
		t.Fail()
	}
}

func TestSwap(t *testing.T) {
	vm.ResetState()
	var tprog = strings.NewReader("2 3 swap")
	if err := vm.Run(tprog, ioutil.Discard); err != nil {
		t.Error(err)
	}
	if !stackEq(3, 2) {
		t.Fail()
	}
}

func TestOver(t *testing.T) {
	vm.ResetState()
	var tprog = strings.NewReader("2 3 over")
	if err := vm.Run(tprog, ioutil.Discard); err != nil {
		t.Error(err)
	}
	if !stackEq(2, 3, 2) {
		t.Fail()
	}
	tprog = strings.NewReader("over")
	if err := vm.Run(tprog, ioutil.Discard); err != nil {
		t.Error(err)
	}
	if !stackEq(2, 3, 2, 3) {
		t.Fail()
	}
}

func TestDrop(t *testing.T) {
	vm.ResetState()
	var tprog = strings.NewReader("2 3 drop")
	if err := vm.Run(tprog, ioutil.Discard); err != nil {
		t.Error(err)
	}
	if !stackEq(2) {
		t.Fail()
	}
	tprog = strings.NewReader("drop drop")
	if err := vm.Run(tprog, ioutil.Discard); err != ErrUnderflow {
		t.Error("Should Underflow!")
	}
}

func TestRot(t *testing.T) {
	vm.ResetState()
	var tprog = strings.NewReader("2 3 rot")
	if err := vm.Run(tprog, ioutil.Discard); err != ErrUnderflow {
		t.Error("Should Underflow!")
	}
	tprog = strings.NewReader("4 rot")
	if err := vm.Run(tprog, ioutil.Discard); err != nil {
		t.Error(vm.Stack, err)
	}
	if !stackEq(3, 4, 2) {
		t.Fail()
	}
	tprog = strings.NewReader("-rot -rot")
	if err := vm.Run(tprog, ioutil.Discard); err != nil {
		t.Error(err)
	}
	if !stackEq(4, 2, 3) {
		t.Fail()
	}
}
