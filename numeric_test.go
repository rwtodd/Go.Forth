package forth

import (
	"io/ioutil"
	"strings"
	"testing"
)

func TestAdd(t *testing.T) {
	vm.ResetState()
	var tprog = strings.NewReader("2 3 +  2 3.1 +")
	if err := vm.Run(tprog, ioutil.Discard); err != nil {
		t.Error(err)
	}
	if !stackEq(5, 5.1) {
		t.Fail()
	}
	vm.ResetState()
	tprog = strings.NewReader(`"hi" " there" +`)
	if err := vm.Run(tprog, ioutil.Discard); err != nil {
		t.Error(err)
	}
	if !stackEq("hi there") {
		t.Fail()
	}
}

func TestMul(t *testing.T) {
	vm.ResetState()
	var tprog = strings.NewReader("2 3 *   2 .25 *")
	if err := vm.Run(tprog, ioutil.Discard); err != nil {
		t.Error(err)
	}
	if !stackEq(6, 0.5) {
		t.Fail()
	}
	vm.ResetState()
	tprog = strings.NewReader(`"hi" 3 *  3 "yo" *`)
	if err := vm.Run(tprog, ioutil.Discard); err != nil {
		t.Error(err)
	}
	if !stackEq("hihihi", "yoyoyo") {
		t.Fail()
	}
}
