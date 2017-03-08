package forth

import "strings"

// : + ( a b -- a+b ) <code>
func add(vm *VM) (err error) {
	top := len(vm.Stack) - 1
	if top < 1 {
		return ErrUnderflow
	}
	switch op1 := vm.Stack[top].(type) {
	case int:
		switch op2 := vm.Stack[top-1].(type) {
		case int:
			vm.Stack[top-1] = op1 + op2
		case float64:
			vm.Stack[top-1] = float64(op1) + op2
		default:
			err = ErrArgument
		}
	case float64:
		switch op2 := vm.Stack[top-1].(type) {
		case int:
			vm.Stack[top-1] = op1 + float64(op2)
		case float64:
			vm.Stack[top-1] = op1 + op2
		default:
			err = ErrArgument
		}
	case string:
		op2, ok := vm.Stack[top-1].(string)
		if ok {
			vm.Stack[top-1] = op2 + op1
		} else {
			err = ErrArgument
		}
	default:
		err = ErrArgument
	}
	vm.Stack = vm.Stack[:top]
	return
}

// : * ( a b -- a*b ) <code>
func multiply(vm *VM) (err error) {
	top := len(vm.Stack) - 1
	if top < 1 {
		return ErrUnderflow
	}
	switch op1 := vm.Stack[top].(type) {
	case int:
		switch op2 := vm.Stack[top-1].(type) {
		case int:
			vm.Stack[top-1] = op1 * op2
		case float64:
			vm.Stack[top-1] = float64(op1) * op2
		case string:
			vm.Stack[top-1] = strings.Repeat(op2, op1)
		default:
			err = ErrArgument
		}
	case float64:
		switch op2 := vm.Stack[top-1].(type) {
		case int:
			vm.Stack[top-1] = op1 * float64(op2)
		case float64:
			vm.Stack[top-1] = op1 * op2
		default:
			err = ErrArgument
		}
	case string:
		op2, ok := vm.Stack[top-1].(int)
		if ok {
			vm.Stack[top-1] = strings.Repeat(op1, op2)
		} else {
			err = ErrArgument
		}
	default:
		err = ErrArgument
	}
	vm.Stack = vm.Stack[:top]
	return
}

// numWordsInit adds numeric core words to the VM
func numWordsInit(vm *VM) {
	vm.Define("+", Word{add, false})
	vm.Define("*", Word{multiply, false})
}
