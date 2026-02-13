// SPDX-License-Identifier: MIT

package forth

import (
	"math"
	"strings"
)

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

// : - ( a b -- a-b ) <code>
func subtract(vm *VM) (err error) {
	top := len(vm.Stack) - 1
	if top < 1 {
		return ErrUnderflow
	}
	switch op1 := vm.Stack[top].(type) {
	case int:
		switch op2 := vm.Stack[top-1].(type) {
		case int:
			vm.Stack[top-1] = op2 - op1
		case float64:
			vm.Stack[top-1] = op2 - float64(op1)
		default:
			err = ErrArgument
		}
	case float64:
		switch op2 := vm.Stack[top-1].(type) {
		case int:
			vm.Stack[top-1] = float64(op2) - op1
		case float64:
			vm.Stack[top-1] = op2 - op1
		default:
			err = ErrArgument
		}
	default:
		err = ErrArgument
	}
	vm.Stack = vm.Stack[:top]
	return
}

// : / ( a b -- a/b ) <code>
func divide(vm *VM) (err error) {
	top := len(vm.Stack) - 1
	if top < 1 {
		return ErrUnderflow
	}
	switch op1 := vm.Stack[top].(type) {
	case int:
		switch op2 := vm.Stack[top-1].(type) {
		case int:
			if op1 == 0 {
				err = ErrArgument
			} else {
				vm.Stack[top-1] = op2 / op1
			}
		case float64:
			if op1 == 0 {
				err = ErrArgument
			} else {
				vm.Stack[top-1] = op2 / float64(op1)
			}
		default:
			err = ErrArgument
		}
	case float64:
		switch op2 := vm.Stack[top-1].(type) {
		case int:
			if op1 == 0 {
				err = ErrArgument
			} else {
				vm.Stack[top-1] = float64(op2) / op1
			}
		case float64:
			if op1 == 0 {
				err = ErrArgument
			} else {
				vm.Stack[top-1] = op2 / op1
			}
		default:
			err = ErrArgument
		}
	default:
		err = ErrArgument
	}
	vm.Stack = vm.Stack[:top]
	return
}

// : sqrt ( a -- sqrt(a) ) <code>
func sqrt(vm *VM) (err error) {
	if len(vm.Stack) < 1 {
		return ErrUnderflow
	}
	switch op := vm.Stack[len(vm.Stack)-1].(type) {
	case int:
		if op < 0 {
			err = ErrArgument
		} else {
			vm.Stack[len(vm.Stack)-1] = math.Sqrt(float64(op))
		}
	case float64:
		if op < 0 {
			err = ErrArgument
		} else {
			vm.Stack[len(vm.Stack)-1] = math.Sqrt(op)
		}
	default:
		err = ErrArgument
	}
	return
}

// : log ( a -- log(a) ) <code>
func log(vm *VM) (err error) {
	if len(vm.Stack) < 1 {
		return ErrUnderflow
	}
	switch op := vm.Stack[len(vm.Stack)-1].(type) {
	case int:
		if op <= 0 {
			err = ErrArgument
		} else {
			vm.Stack[len(vm.Stack)-1] = math.Log(float64(op))
		}
	case float64:
		if op <= 0 {
			err = ErrArgument
		} else {
			vm.Stack[len(vm.Stack)-1] = math.Log(op)
		}
	default:
		err = ErrArgument
	}
	return
}

// : log10 ( a -- log10(a) ) <code>
func log10(vm *VM) (err error) {
	if len(vm.Stack) < 1 {
		return ErrUnderflow
	}
	switch op := vm.Stack[len(vm.Stack)-1].(type) {
	case int:
		if op <= 0 {
			err = ErrArgument
		} else {
			vm.Stack[len(vm.Stack)-1] = math.Log10(float64(op))
		}
	case float64:
		if op <= 0 {
			err = ErrArgument
		} else {
			vm.Stack[len(vm.Stack)-1] = math.Log10(op)
		}
	default:
		err = ErrArgument
	}
	return
}

// : log2 ( a -- log2(a) ) <code>
func log2(vm *VM) (err error) {
	if len(vm.Stack) < 1 {
		return ErrUnderflow
	}
	switch op := vm.Stack[len(vm.Stack)-1].(type) {
	case int:
		if op <= 0 {
			err = ErrArgument
		} else {
			vm.Stack[len(vm.Stack)-1] = math.Log2(float64(op))
		}
	case float64:
		if op <= 0 {
			err = ErrArgument
		} else {
			vm.Stack[len(vm.Stack)-1] = math.Log2(op)
		}
	default:
		err = ErrArgument
	}
	return
}

// : max ( a b -- max(a,b) ) <code>
func max(vm *VM) (err error) {
	top := len(vm.Stack) - 1
	if top < 1 {
		return ErrUnderflow
	}
	switch op1 := vm.Stack[top].(type) {
	case int:
		switch op2 := vm.Stack[top-1].(type) {
		case int:
			if op2 > op1 {
				vm.Stack[top-1] = op2
			} else {
				vm.Stack[top-1] = op1
			}
		case float64:
			vm.Stack[top-1] = math.Max(op2, float64(op1))
		default:
			err = ErrArgument
		}
	case float64:
		switch op2 := vm.Stack[top-1].(type) {
		case int:
			vm.Stack[top-1] = math.Max(float64(op2), op1)
		case float64:
			vm.Stack[top-1] = math.Max(op2, op1)
		default:
			err = ErrArgument
		}
	default:
		err = ErrArgument
	}
	vm.Stack = vm.Stack[:top]
	return
}

// : min ( a b -- min(a,b) ) <code>
func min(vm *VM) (err error) {
	top := len(vm.Stack) - 1
	if top < 1 {
		return ErrUnderflow
	}
	switch op1 := vm.Stack[top].(type) {
	case int:
		switch op2 := vm.Stack[top-1].(type) {
		case int:
			if op2 < op1 {
				vm.Stack[top-1] = op2
			} else {
				vm.Stack[top-1] = op1
			}
		case float64:
			vm.Stack[top-1] = math.Min(op2, float64(op1))
		default:
			err = ErrArgument
		}
	case float64:
		switch op2 := vm.Stack[top-1].(type) {
		case int:
			vm.Stack[top-1] = math.Min(float64(op2), op1)
		case float64:
			vm.Stack[top-1] = math.Min(op2, op1)
		default:
			err = ErrArgument
		}
	default:
		err = ErrArgument
	}
	vm.Stack = vm.Stack[:top]
	return
}

// : sin ( a -- sin(a) ) <code>
func sin(vm *VM) (err error) {
	if len(vm.Stack) < 1 {
		return ErrUnderflow
	}
	switch op := vm.Stack[len(vm.Stack)-1].(type) {
	case int:
		vm.Stack[len(vm.Stack)-1] = math.Sin(float64(op))
	case float64:
		vm.Stack[len(vm.Stack)-1] = math.Sin(op)
	default:
		err = ErrArgument
	}
	return
}

// : cos ( a -- cos(a) ) <code>
func cos(vm *VM) (err error) {
	if len(vm.Stack) < 1 {
		return ErrUnderflow
	}
	switch op := vm.Stack[len(vm.Stack)-1].(type) {
	case int:
		vm.Stack[len(vm.Stack)-1] = math.Cos(float64(op))
	case float64:
		vm.Stack[len(vm.Stack)-1] = math.Cos(op)
	default:
		err = ErrArgument
	}
	return
}

// : tan ( a -- tan(a) ) <code>
func tan(vm *VM) (err error) {
	if len(vm.Stack) < 1 {
		return ErrUnderflow
	}
	switch op := vm.Stack[len(vm.Stack)-1].(type) {
	case int:
		vm.Stack[len(vm.Stack)-1] = math.Tan(float64(op))
	case float64:
		vm.Stack[len(vm.Stack)-1] = math.Tan(op)
	default:
		err = ErrArgument
	}
	return
}

// : round ( a -- round(a) ) <code>
func round(vm *VM) (err error) {
	if len(vm.Stack) < 1 {
		return ErrUnderflow
	}
	switch op := vm.Stack[len(vm.Stack)-1].(type) {
	case int:
		vm.Stack[len(vm.Stack)-1] = float64(op)
	case float64:
		vm.Stack[len(vm.Stack)-1] = math.Round(op)
	default:
		err = ErrArgument
	}
	return
}

// : floor ( a -- floor(a) ) <code>
func floor(vm *VM) (err error) {
	if len(vm.Stack) < 1 {
		return ErrUnderflow
	}
	switch op := vm.Stack[len(vm.Stack)-1].(type) {
	case int:
		vm.Stack[len(vm.Stack)-1] = float64(op)
	case float64:
		vm.Stack[len(vm.Stack)-1] = math.Floor(op)
	default:
		err = ErrArgument
	}
	return
}

// : ceil ( a -- ceil(a) ) <code>
func ceil(vm *VM) (err error) {
	if len(vm.Stack) < 1 {
		return ErrUnderflow
	}
	switch op := vm.Stack[len(vm.Stack)-1].(type) {
	case int:
		vm.Stack[len(vm.Stack)-1] = float64(op)
	case float64:
		vm.Stack[len(vm.Stack)-1] = math.Ceil(op)
	default:
		err = ErrArgument
	}
	return
}

// numWordsInit adds numeric core words to the VM
func numWordsInit(vm *VM) {
	vm.Define("+", Word{add, false})
	vm.Define("-", Word{subtract, false})
	vm.Define("*", Word{multiply, false})
	vm.Define("/", Word{divide, false})
	vm.Define("sqrt", Word{sqrt, false})
	vm.Define("log", Word{log, false})
	vm.Define("log10", Word{log10, false})
	vm.Define("log2", Word{log2, false})
	vm.Define("max", Word{max, false})
	vm.Define("min", Word{min, false})
	vm.Define("sin", Word{sin, false})
	vm.Define("cos", Word{cos, false})
	vm.Define("tan", Word{tan, false})
	vm.Define("round", Word{round, false})
	vm.Define("floor", Word{floor, false})
	vm.Define("ceil", Word{ceil, false})
}
