package forth

import "errors"

// stack words

func dup(vm *VM) {
	top := len(vm.Stack)
	if top >= 1 {
		vm.Stack = append(vm.Stack, vm.Stack[top-1])
	} else {
		vm.Err = errors.New("dup: stack underflow")
	}
}

func over(vm *VM) {
	top := len(vm.Stack)
	if top >= 2 {
		vm.Stack = append(vm.Stack, vm.Stack[top-2])
	} else {
		vm.Err = errors.New("over: stack underflow")
	}
}

func drop(vm *VM) {
	top := len(vm.Stack)
	if top >= 1 {
		vm.Stack = vm.Stack[:top-1]
	} else {
		vm.Err = errors.New("drop: stack underflow")
	}
}

func swap(vm *VM) {
	top := len(vm.Stack)
	if top >= 2 {
		vm.Stack[top-1], vm.Stack[top-2] = vm.Stack[top-2], vm.Stack[top-1]
	} else {
		vm.Err = errors.New("swap: stack underflow")
	}
}
