package forth

// stack words

// : dup ( a -- a a ) <code>
func dup(vm *VM) (e error) {
	top := len(vm.Stack)
	if top >= 1 {
		vm.Stack = append(vm.Stack, vm.Stack[top-1])
	} else {
		e = ErrUnderflow
	}
	return
}

// : over swap dup -rot ;
func over(vm *VM) (e error) {
	top := len(vm.Stack)
	if top >= 2 {
		vm.Stack = append(vm.Stack, vm.Stack[top-2])
	} else {
		e = ErrUnderflow
	}
	return
}

// : drop ( a -- ) <code>
func drop(vm *VM) (e error) {
	top := len(vm.Stack)
	if top >= 1 {
		vm.Stack = vm.Stack[:top-1]
	} else {
		e = ErrUnderflow
	}
	return
}

// : swap ( a b -- b a )  <code>
func swap(vm *VM) (e error) {
	top := len(vm.Stack)
	if top >= 2 {
		vm.Stack[top-1], vm.Stack[top-2] = vm.Stack[top-2], vm.Stack[top-1]
	} else {
		e = ErrUnderflow
	}
	return
}

// : rot  ( a b c -- b c a ) <code>
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

// : -rot  rot rot ;
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

// : nip swap drop ;
func nip(vm *VM) (e error) {
	top := len(vm.Stack)
	if top >= 2 {
		vm.Stack[top-2] = vm.Stack[top-1]
		vm.Stack = vm.Stack[:top-1]
	} else {
		e = ErrUnderflow
	}
	return
}

// : tuck swap over ;
func tuck(vm *VM) (e error) {
	top := len(vm.Stack)
	if top >= 2 {
		vm.Stack = append(vm.Stack, vm.Stack[top-1])
		vm.Stack[top-1], vm.Stack[top-2] = vm.Stack[top-2], vm.Stack[top-1]
	} else {
		e = ErrUnderflow
	}
	return
}

// stackWordsInit adds stack-related core words to the VM
func stackWordsInit(vm *VM) {
	vm.Define("dup", Word{dup, false})
	vm.Define("drop", Word{drop, false})
	vm.Define("swap", Word{swap, false})
	vm.Define("over", Word{over, false})
	vm.Define("rot", Word{rotate, false})
	vm.Define("-rot", Word{minusRotate, false})
	vm.Define("nip", Word{nip, false})
	vm.Define("tuck", Word{tuck, false})
}
