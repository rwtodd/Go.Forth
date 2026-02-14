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
			return ErrArgumentMsg("= requires comparable numeric or string types")
		}
	case float64:
		switch v2 := val2.(type) {
		case int:
			res = (v1 == float64(v2))
		case float64:
			res = (v1 == v2)
		default:
			return ErrArgumentMsg("= requires comparable numeric or string types")
		}
	case string:
		v2, ok := val2.(string)
		if !ok {
			return ErrArgumentMsg("= requires two strings")
		}
		res = (v1 == v2)
	default:
		return ErrArgumentMsg("= requires comparable numeric or string types")
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
			return ErrArgumentMsg("< requires numeric arguments")
		}
	case float64:
		switch v2 := val2.(type) {
		case int:
			res = (v1 < float64(v2))
		case float64:
			res = (v1 < v2)
		default:
			return ErrArgumentMsg("< requires numeric arguments")
		}
	case string:
		v2, ok := val2.(string)
		if !ok {
			return ErrArgumentMsg("< requires two strings")
		}
		res = (v1 < v2)
	default:
		return ErrArgumentMsg("< requires numeric or string type")
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
			return ErrArgumentMsg("> requires numeric arguments")
		}
	case float64:
		switch v2 := val2.(type) {
		case int:
			res = (v1 > float64(v2))
		case float64:
			res = (v1 > v2)
		default:
			return ErrArgumentMsg("> requires numeric arguments")
		}
	case string:
		v2, ok := val2.(string)
		if !ok {
			return ErrArgumentMsg("> requires two strings")
		}
		res = (v1 > v2)
	default:
		return ErrArgumentMsg("> requires numeric or string type")
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
			return ErrArgumentMsg("<= requires numeric arguments")
		}
	case float64:
		switch v2 := val2.(type) {
		case int:
			res = (v1 <= float64(v2))
		case float64:
			res = (v1 <= v2)
		default:
			return ErrArgumentMsg("<= requires numeric arguments")
		}
	case string:
		v2, ok := val2.(string)
		if !ok {
			return ErrArgumentMsg("<= requires two strings")
		}
		res = (v1 <= v2)
	default:
		return ErrArgumentMsg("<= requires numeric or string type")
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
			return ErrArgumentMsg(">= requires numeric arguments")
		}
	case float64:
		switch v2 := val2.(type) {
		case int:
			res = (v1 >= float64(v2))
		case float64:
			res = (v1 >= v2)
		default:
			return ErrArgumentMsg(">= requires numeric arguments")
		}
	case string:
		v2, ok := val2.(string)
		if !ok {
			return ErrArgumentMsg(">= requires two strings")
		}
		res = (v1 >= v2)
	default:
		return ErrArgumentMsg(">= requires numeric or string type")
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
		return ErrArgumentMsg("0= requires numeric type")
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
		return ErrArgumentMsg("0< requires numeric type")
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
		return ErrArgumentMsg("0> requires numeric type")
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
		return ErrArgumentMsg("and requires two integers")
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
		return ErrArgumentMsg("or requires two integers")
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
		return ErrArgumentMsg("xor requires two integers")
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
		return ErrArgumentMsg("invert requires an integer")
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
	vm.Define(Word{Name: "=", Run: equal, Immediate: false})
	vm.Define(Word{Name: "<", Run: lessThan, Immediate: false})
	vm.Define(Word{Name: ">", Run: greaterThan, Immediate: false})
	vm.Define(Word{Name: "<=", Run: lessThanOrEqual, Immediate: false})
	vm.Define(Word{Name: ">=", Run: greaterThanOrEqual, Immediate: false})
	vm.Define(Word{Name: "<>", Run: notEqual, Immediate: false})
	vm.Define(Word{Name: "0=", Run: zeroEqual, Immediate: false})
	vm.Define(Word{Name: "0<", Run: zeroLessThan, Immediate: false})
	vm.Define(Word{Name: "0>", Run: zeroGreaterThan, Immediate: false})
	vm.Define(Word{Name: "and", Run: andOp, Immediate: false})
	vm.Define(Word{Name: "or", Run: orOp, Immediate: false})
	vm.Define(Word{Name: "xor", Run: xorOp, Immediate: false})
	vm.Define(Word{Name: "invert", Run: invertOp, Immediate: false})
	vm.Define(Word{Name: "true", Run: trueConst, Immediate: false})
	vm.Define(Word{Name: "false", Run: falseConst, Immediate: false})
}
