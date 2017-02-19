package forth

import "strings"

func add(vm *VM) {
	top := len(vm.Stack) - 1
	if top < 1 {
		vm.Err = ErrUnderflow
		return
	}
	switch op1 := vm.Stack[top].(type) {
	case int:
		switch op2 := vm.Stack[top-1].(type) {
		case int:
			vm.Stack[top-1] = op1 + op2
		case float64:
			vm.Stack[top-1] = float64(op1) + op2
		default:
			vm.Err = ErrArgument
		}
	case float64:
		switch op2 := vm.Stack[top-1].(type) {
		case int:
			vm.Stack[top-1] = op1 + float64(op2)
		case float64:
			vm.Stack[top-1] = op1 + op2
		default:
			vm.Err = ErrArgument
		}
	case string:
		op2, ok := vm.Stack[top-1].(string)
		if ok {
			vm.Stack[top-1] = op2 + op1
		} else {
			vm.Err = ErrArgument
		}
	default:
		vm.Err = ErrArgument
	}
	vm.Stack = vm.Stack[:top]
}

func multiply(vm *VM) {
	top := len(vm.Stack) - 1
	if top < 1 {
		vm.Err = ErrUnderflow
		return
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
			vm.Err = ErrArgument
		}
	case float64:
		switch op2 := vm.Stack[top-1].(type) {
		case int:
			vm.Stack[top-1] = op1 * float64(op2)
		case float64:
			vm.Stack[top-1] = op1 * op2
		default:
			vm.Err = ErrArgument
		}
	case string:
		op2, ok := vm.Stack[top-1].(int)
		if ok {
			vm.Stack[top-1] = strings.Repeat(op1, op2)
		} else {
			vm.Err = ErrArgument
		}
	default:
		vm.Err = ErrArgument
	}
	vm.Stack = vm.Stack[:top]
}
