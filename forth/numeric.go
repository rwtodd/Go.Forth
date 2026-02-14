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
			err = ErrArgumentMsg("+ requires numeric arguments")
		}
	case float64:
		switch op2 := vm.Stack[top-1].(type) {
		case int:
			vm.Stack[top-1] = op1 + float64(op2)
		case float64:
			vm.Stack[top-1] = op1 + op2
		default:
			err = ErrArgumentMsg("+ requires numeric arguments")
		}
	case string:
		op2, ok := vm.Stack[top-1].(string)
		if ok {
			vm.Stack[top-1] = op2 + op1
		} else {
			err = ErrArgumentMsg("+ requires two strings or two numbers")
		}
	default:
		err = ErrArgumentMsg("invalid types for +")
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
			err = ErrArgumentMsg("* requires numeric arguments or string+int")
		}
	case float64:
		switch op2 := vm.Stack[top-1].(type) {
		case int:
			vm.Stack[top-1] = op1 * float64(op2)
		case float64:
			vm.Stack[top-1] = op1 * op2
		default:
			err = ErrArgumentMsg("* requires numeric arguments")
		}
	case string:
		op2, ok := vm.Stack[top-1].(int)
		if ok {
			vm.Stack[top-1] = strings.Repeat(op1, op2)
		} else {
			err = ErrArgumentMsg("* string repetition requires integer count")
		}
	default:
		err = ErrArgumentMsg("invalid types for *")
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
			err = ErrArgumentMsg("- requires numeric arguments")
		}
	case float64:
		switch op2 := vm.Stack[top-1].(type) {
		case int:
			vm.Stack[top-1] = float64(op2) - op1
		case float64:
			vm.Stack[top-1] = op2 - op1
		default:
			err = ErrArgumentMsg("- requires numeric arguments")
		}
	default:
		err = ErrArgumentMsg("- requires numeric arguments")
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
				err = ErrArgumentMsg("division by zero")
			} else {
				vm.Stack[top-1] = op2 / op1
			}
		case float64:
			if op1 == 0 {
				err = ErrArgumentMsg("division by zero")
			} else {
				vm.Stack[top-1] = op2 / float64(op1)
			}
		default:
			err = ErrArgumentMsg("/ requires numeric arguments")
		}
	case float64:
		switch op2 := vm.Stack[top-1].(type) {
		case int:
			if op1 == 0 {
				err = ErrArgumentMsg("division by zero")
			} else {
				vm.Stack[top-1] = float64(op2) / op1
			}
		case float64:
			if op1 == 0 {
				err = ErrArgumentMsg("division by zero")
			} else {
				vm.Stack[top-1] = op2 / op1
			}
		default:
			err = ErrArgumentMsg("/ requires numeric arguments")
		}
	default:
		err = ErrArgumentMsg("/ requires numeric arguments")
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
			err = ErrArgumentMsg("sqrt requires non-negative argument")
		} else {
			vm.Stack[len(vm.Stack)-1] = math.Sqrt(float64(op))
		}
	case float64:
		if op < 0 {
			err = ErrArgumentMsg("sqrt requires non-negative argument")
		} else {
			vm.Stack[len(vm.Stack)-1] = math.Sqrt(op)
		}
	default:
		err = ErrArgumentMsg("sqrt requires numeric argument")
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
			err = ErrArgumentMsg("log requires positive argument")
		} else {
			vm.Stack[len(vm.Stack)-1] = math.Log(float64(op))
		}
	case float64:
		if op <= 0 {
			err = ErrArgumentMsg("log requires positive argument")
		} else {
			vm.Stack[len(vm.Stack)-1] = math.Log(op)
		}
	default:
		err = ErrArgumentMsg("log requires numeric argument")
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
			err = ErrArgumentMsg("log10 requires positive argument")
		} else {
			vm.Stack[len(vm.Stack)-1] = math.Log10(float64(op))
		}
	case float64:
		if op <= 0 {
			err = ErrArgumentMsg("log10 requires positive argument")
		} else {
			vm.Stack[len(vm.Stack)-1] = math.Log10(op)
		}
	default:
		err = ErrArgumentMsg("log10 requires numeric argument")
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
			err = ErrArgumentMsg("log2 requires positive argument")
		} else {
			vm.Stack[len(vm.Stack)-1] = math.Log2(float64(op))
		}
	case float64:
		if op <= 0 {
			err = ErrArgumentMsg("log2 requires positive argument")
		} else {
			vm.Stack[len(vm.Stack)-1] = math.Log2(op)
		}
	default:
		err = ErrArgumentMsg("log2 requires numeric argument")
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
			err = ErrArgumentMsg("max requires numeric arguments")
		}
	case float64:
		switch op2 := vm.Stack[top-1].(type) {
		case int:
			vm.Stack[top-1] = math.Max(float64(op2), op1)
		case float64:
			vm.Stack[top-1] = math.Max(op2, op1)
		default:
			err = ErrArgumentMsg("max requires numeric arguments")
		}
	default:
		err = ErrArgumentMsg("max requires numeric arguments")
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
			err = ErrArgumentMsg("min requires numeric arguments")
		}
	case float64:
		switch op2 := vm.Stack[top-1].(type) {
		case int:
			vm.Stack[top-1] = math.Min(float64(op2), op1)
		case float64:
			vm.Stack[top-1] = math.Min(op2, op1)
		default:
			err = ErrArgumentMsg("min requires numeric arguments")
		}
	default:
		err = ErrArgumentMsg("min requires numeric arguments")
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
		err = ErrArgumentMsg("sin requires numeric argument")
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
		err = ErrArgumentMsg("cos requires numeric argument")
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
		err = ErrArgumentMsg("tan requires numeric argument")
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
		err = ErrArgumentMsg("round requires numeric argument")
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
		err = ErrArgumentMsg("floor requires numeric argument")
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
		err = ErrArgumentMsg("ceil requires numeric argument")
	}
	return
}

// numWordsInit adds numeric core words to the VM
func numWordsInit(vm *VM) {
	vm.Define(Word{Name: "+", Run: add, Immediate: false})
	vm.Define(Word{Name: "-", Run: subtract, Immediate: false})
	vm.Define(Word{Name: "*", Run: multiply, Immediate: false})
	vm.Define(Word{Name: "/", Run: divide, Immediate: false})
	vm.Define(Word{Name: "sqrt", Run: sqrt, Immediate: false})
	vm.Define(Word{Name: "log", Run: log, Immediate: false})
	vm.Define(Word{Name: "log10", Run: log10, Immediate: false})
	vm.Define(Word{Name: "log2", Run: log2, Immediate: false})
	vm.Define(Word{Name: "max", Run: max, Immediate: false})
	vm.Define(Word{Name: "min", Run: min, Immediate: false})
	vm.Define(Word{Name: "sin", Run: sin, Immediate: false})
	vm.Define(Word{Name: "cos", Run: cos, Immediate: false})
	vm.Define(Word{Name: "tan", Run: tan, Immediate: false})
	vm.Define(Word{Name: "round", Run: round, Immediate: false})
	vm.Define(Word{Name: "floor", Run: floor, Immediate: false})
	vm.Define(Word{Name: "ceil", Run: ceil, Immediate: false})
}
