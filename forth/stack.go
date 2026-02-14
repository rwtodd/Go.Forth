// SPDX-License-Identifier: MIT

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

// >r push onto rstack
func toR(vm *VM) (e error) {
	var tos interface{}
	tos, e = vm.Pop()
	vm.RPush(tos)
	return
}

// r> pop from rstack
func fromR(vm *VM) (e error) {
	var tos interface{}
	tos, e = vm.RPop()
	vm.Push(tos)
	return
}

// r@ peek at rstack
func peekR(vm *VM) error {
	tos := len(vm.Rstack) - 1
	if tos < 0 {
		return ErrUnderflowMsg("return stack underflow")
	}
	vm.Push(vm.Rstack[tos])
	return nil
}

// rdrop removes the top item from the return stack
func rdrop(vm *VM) error {
	_, err := vm.RPop()
	return err
}

// stackWordsInit adds stack-related core words to the VM
func stackWordsInit(vm *VM) {
	vm.Define(Word{Name: "dup", Run: dup, Immediate: false})
	vm.Define(Word{Name: "drop", Run: drop, Immediate: false})
	vm.Define(Word{Name: "swap", Run: swap, Immediate: false})
	vm.Define(Word{Name: "over", Run: over, Immediate: false})
	vm.Define(Word{Name: "rot", Run: rotate, Immediate: false})
	vm.Define(Word{Name: "-rot", Run: minusRotate, Immediate: false})
	vm.Define(Word{Name: "nip", Run: nip, Immediate: false})
	vm.Define(Word{Name: "tuck", Run: tuck, Immediate: false})
	vm.Define(Word{Name: ">r", Run: toR, Immediate: false})
	vm.Define(Word{Name: "r>", Run: fromR, Immediate: false})
	vm.Define(Word{Name: "r@", Run: peekR, Immediate: false})
	vm.Define(Word{Name: "rdrop", Run: rdrop, Immediate: false})
}
