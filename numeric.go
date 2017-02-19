package forth

import "errors"

func add(vm *VM) {
	top := len(vm.Stack) - 1
	if top < 0 {
		vm.Err = errors.New("+: stack underflow")
		return
	}
	i1, ok1 := vm.Stack[top].(int)
	i2, ok2 := vm.Stack[top-1].(int)
	if ok1 && ok2 {
		vm.Stack[top-1] = i1 + i2
		vm.Stack = vm.Stack[:top]
	} else {
		vm.Err = errors.New("+: bad argument")
	}
}
