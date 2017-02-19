package forth

// stack words

func dup(vm *VM) (e error) {
	top := len(vm.Stack)
	if top >= 1 {
		vm.Stack = append(vm.Stack, vm.Stack[top-1])
	} else {
		e = ErrUnderflow
	}
	return
}

func over(vm *VM) (e error) {
	top := len(vm.Stack)
	if top >= 2 {
		vm.Stack = append(vm.Stack, vm.Stack[top-2])
	} else {
		e = ErrUnderflow
	}
	return
}

func drop(vm *VM) (e error) {
	top := len(vm.Stack)
	if top >= 1 {
		vm.Stack = vm.Stack[:top-1]
	} else {
		e = ErrUnderflow
	}
	return
}

func swap(vm *VM) (e error) {
	top := len(vm.Stack)
	if top >= 2 {
		vm.Stack[top-1], vm.Stack[top-2] = vm.Stack[top-2], vm.Stack[top-1]
	} else {
		e = ErrUnderflow
	}
	return
}

func rotate(vm *VM) (e error) {
	top := len(vm.Stack)
	if top >= 3 {
		vm.Stack[top-1], vm.Stack[top-2], vm.Stack[top-3] =
			vm.Stack[top-3], vm.Stack[top-1], vm.Stack[top-2]
	} else {
		e = ErrUnderflow
	}
	return
}

func minusRotate(vm *VM) (e error) {
	top := len(vm.Stack)
	if top >= 3 {
		vm.Stack[top-1], vm.Stack[top-2], vm.Stack[top-3] =
			vm.Stack[top-2], vm.Stack[top-3], vm.Stack[top-1]
	} else {
		e = ErrUnderflow
	}
	return
}
