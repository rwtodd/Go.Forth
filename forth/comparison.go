// SPDX-License-Identifier: MIT

package forth

const (
	forthTrue  = -1
	forthFalse = 0
)

// BoolToForth converts a Go bool to a Forth boolean int
func BoolToForth(b bool) int {
	if b {
		return forthTrue
	}
	return forthFalse
}

// : = ( a b -- bool )
func equal(vm *VM) error {
	val2, err := vm.Pop()
	if err != nil {
		return err
	}
	val1, err := vm.Pop()
	if err != nil {
		return err
	}

	res := false
	switch v1 := val1.(type) {
	case int:
		switch v2 := val2.(type) {
		case int:
			res = (v1 == v2)
		case float64:
			res = (float64(v1) == v2)
		default:
			return ErrArgument
		}
	case float64:
		switch v2 := val2.(type) {
		case int:
			res = (v1 == float64(v2))
		case float64:
			res = (v1 == v2)
		default:
			return ErrArgument
		}
	case string:
		v2, ok := val2.(string)
		if !ok {
			return ErrArgument
		}
		res = (v1 == v2)
	default:
		return ErrArgument
	}
	vm.Push(BoolToForth(res))
	return nil
}

// : < ( a b -- bool )
func lessThan(vm *VM) error {
	val2, err := vm.Pop()
	if err != nil {
		return err
	}
	val1, err := vm.Pop()
	if err != nil {
		return err
	}

	res := false
	switch v1 := val1.(type) {
	case int:
		switch v2 := val2.(type) {
		case int:
			res = (v1 < v2)
		case float64:
			res = (float64(v1) < v2)
		default:
			return ErrArgument
		}
	case float64:
		switch v2 := val2.(type) {
		case int:
			res = (v1 < float64(v2))
		case float64:
			res = (v1 < v2)
		default:
			return ErrArgument
		}
	case string:
		v2, ok := val2.(string)
		if !ok {
			return ErrArgument
		}
		res = (v1 < v2)
	default:
		return ErrArgument
	}
	vm.Push(BoolToForth(res))
	return nil
}

// : > ( a b -- bool )
func greaterThan(vm *VM) error {
	val2, err := vm.Pop()
	if err != nil {
		return err
	}
	val1, err := vm.Pop()
	if err != nil {
		return err
	}

	res := false
	switch v1 := val1.(type) {
	case int:
		switch v2 := val2.(type) {
		case int:
			res = (v1 > v2)
		case float64:
			res = (float64(v1) > v2)
		default:
			return ErrArgument
		}
	case float64:
		switch v2 := val2.(type) {
		case int:
			res = (v1 > float64(v2))
		case float64:
			res = (v1 > v2)
		default:
			return ErrArgument
		}
	case string:
		v2, ok := val2.(string)
		if !ok {
			return ErrArgument
		}
		res = (v1 > v2)
	default:
		return ErrArgument
	}
	vm.Push(BoolToForth(res))
	return nil
}

// : <= ( a b -- bool )
func lessThanOrEqual(vm *VM) error {
	val2, err := vm.Pop()
	if err != nil {
		return err
	}
	val1, err := vm.Pop()
	if err != nil {
		return err
	}

	res := false
	switch v1 := val1.(type) {
	case int:
		switch v2 := val2.(type) {
		case int:
			res = (v1 <= v2)
		case float64:
			res = (float64(v1) <= v2)
		default:
			return ErrArgument
		}
	case float64:
		switch v2 := val2.(type) {
		case int:
			res = (v1 <= float64(v2))
		case float64:
			res = (v1 <= v2)
		default:
			return ErrArgument
		}
	case string:
		v2, ok := val2.(string)
		if !ok {
			return ErrArgument
		}
		res = (v1 <= v2)
	default:
		return ErrArgument
	}
	vm.Push(BoolToForth(res))
	return nil
}

// : >= ( a b -- bool )
func greaterThanOrEqual(vm *VM) error {
	val2, err := vm.Pop()
	if err != nil {
		return err
	}
	val1, err := vm.Pop()
	if err != nil {
		return err
	}

	res := false
	switch v1 := val1.(type) {
	case int:
		switch v2 := val2.(type) {
		case int:
			res = (v1 >= v2)
		case float64:
			res = (float64(v1) >= v2)
		default:
			return ErrArgument
		}
	case float64:
		switch v2 := val2.(type) {
		case int:
			res = (v1 >= float64(v2))
		case float64:
			res = (v1 >= v2)
		default:
			return ErrArgument
		}
	case string:
		v2, ok := val2.(string)
		if !ok {
			return ErrArgument
		}
		res = (v1 >= v2)
	default:
		return ErrArgument
	}
	vm.Push(BoolToForth(res))
	return nil
}

// : <> ( a b -- bool )
func notEqual(vm *VM) error {
	if err := equal(vm); err != nil {
		return err
	}
	// logical invert the result
	val, _ := vm.Pop() // error checked in equal
	if val.(int) == forthFalse {
		vm.Push(forthTrue)
	} else {
		vm.Push(forthFalse)
	}
	return nil
}

// : 0= ( a -- bool )
func zeroEqual(vm *VM) error {
	val, err := vm.Pop()
	if err != nil {
		return err
	}

	res := false
	switch v := val.(type) {
	case int:
		res = (v == 0)
	case float64:
		res = (v == 0.0)
	default:
		return ErrArgument
	}
	vm.Push(BoolToForth(res))
	return nil
}

// : 0< ( a -- bool )
func zeroLessThan(vm *VM) error {
	val, err := vm.Pop()
	if err != nil {
		return err
	}

	res := false
	switch v := val.(type) {
	case int:
		res = (v < 0)
	case float64:
		res = (v < 0.0)
	default:
		return ErrArgument
	}
	vm.Push(BoolToForth(res))
	return nil
}

// : 0> ( a -- bool )
func zeroGreaterThan(vm *VM) error {
	val, err := vm.Pop()
	if err != nil {
		return err
	}

	res := false
	switch v := val.(type) {
	case int:
		res = (v > 0)
	case float64:
		res = (v > 0.0)
	default:
		return ErrArgument
	}
	vm.Push(BoolToForth(res))
	return nil
}

// : and ( a b -- c )
func andOp(vm *VM) error {
	val2, err := vm.Pop()
	if err != nil {
		return err
	}
	val1, err := vm.Pop()
	if err != nil {
		return err
	}

	v1, ok1 := val1.(int)
	v2, ok2 := val2.(int)
	if !ok1 || !ok2 {
		return ErrArgument
	}

	vm.Push(v1 & v2)
	return nil
}

// : or ( a b -- c )
func orOp(vm *VM) error {
	val2, err := vm.Pop()
	if err != nil {
		return err
	}
	val1, err := vm.Pop()
	if err != nil {
		return err
	}

	v1, ok1 := val1.(int)
	v2, ok2 := val2.(int)
	if !ok1 || !ok2 {
		return ErrArgument
	}

	vm.Push(v1 | v2)
	return nil
}

// : xor ( a b -- c )
func xorOp(vm *VM) error {
	val2, err := vm.Pop()
	if err != nil {
		return err
	}
	val1, err := vm.Pop()
	if err != nil {
		return err
	}

	v1, ok1 := val1.(int)
	v2, ok2 := val2.(int)
	if !ok1 || !ok2 {
		return ErrArgument
	}

	vm.Push(v1 ^ v2)
	return nil
}

// : invert ( a -- b )
func invertOp(vm *VM) error {
	val, err := vm.Pop()
	if err != nil {
		return err
	}

	v, ok := val.(int)
	if !ok {
		return ErrArgument
	}

	vm.Push(^v)
	return nil
}

// : true ( -- -1 )
func trueConst(vm *VM) error {
	vm.Push(forthTrue)
	return nil
}

// : false ( -- 0 )
func falseConst(vm *VM) error {
	vm.Push(forthFalse)
	return nil
}

func comparisonWordsInit(vm *VM) {
	vm.Define("=", Word{equal, false})
	vm.Define("<", Word{lessThan, false})
	vm.Define(">", Word{greaterThan, false})
	vm.Define("<=", Word{lessThanOrEqual, false})
	vm.Define(">=", Word{greaterThanOrEqual, false})
	vm.Define("<>", Word{notEqual, false})
	vm.Define("0=", Word{zeroEqual, false})
	vm.Define("0<", Word{zeroLessThan, false})
	vm.Define("0>", Word{zeroGreaterThan, false})
	vm.Define("and", Word{andOp, false})
	vm.Define("or", Word{orOp, false})
	vm.Define("xor", Word{xorOp, false})
	vm.Define("invert", Word{invertOp, false})
	vm.Define("true", Word{trueConst, false})
	vm.Define("false", Word{falseConst, false})
}
