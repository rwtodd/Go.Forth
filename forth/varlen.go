// SPDX-License-Identifier: MIT

package forth

import (
	"fmt"
	"io"
)

// sprintf ( format args... count -- result )
func sprintf(vm *VM) error {
	// 1. Pop count
	countVal, err := vm.Pop()
	if err != nil {
		return err
	}
	count, ok := countVal.(int)
	if !ok {
		return ErrArgumentMsg("sprintf count must be an integer")
	}

	currentLen := len(vm.Stack)
	if currentLen < count+1 {
		return ErrUnderflowMsg("sprintf stack underflow")
	}

	// 2. Identify args and format
	// Stack: ... format arg1 ... argN
	// positions:
	// format: currentLen - count - 1
	// args: currentLen - count ... currentLen

	formatIdx := currentLen - count - 1
	formatVal := vm.Stack[formatIdx]
	format, ok := formatVal.(string)
	if !ok {
		return ErrArgumentMsg("sprintf format must be a string")
	}

	args := vm.Stack[formatIdx+1:]

	// 3. Format
	res := fmt.Sprintf(format, args...)

	// 4. Cleanup and Push result
	// We remove format and args.
	// New length = formatIdx.
	vm.Stack = vm.Stack[:formatIdx]
	vm.Push(res)

	return nil
}

// <@push> ( array item1 ... itemN count -- result )
func varLenPush(vm *VM) error {
	// 1. Pop count
	countVal, err := vm.Pop()
	if err != nil {
		return err
	}
	count, ok := countVal.(int)
	if !ok || count < 0 {
		return ErrArgumentMsg("<@push> count must be non-negative integer")
	}

	currentLen := len(vm.Stack)
	// We need array + count items on stack.
	if currentLen < count+1 {
		return ErrUnderflowMsg("<@push> stack underflow")
	}

	// 2. Identify array and items
	// Stack: ... array item1 ... itemN
	arrIdx := currentLen - count - 1
	arrVal := vm.Stack[arrIdx]
	items := vm.Stack[arrIdx+1:]

	var varPtr *Variable
	if v, ok := arrVal.(*Variable); ok {
		varPtr = v
		arrVal = v.value
	}

	// 3. Append based on type
	switch a := arrVal.(type) {
	case []int:
		for _, item := range items {
			switch v := item.(type) {
			case int:
				a = append(a, v)
			case float64:
				if v != float64(int(v)) {
					return ErrArgumentMsg("lossy conversion from float to int in <@push>")
				}
				a = append(a, int(v))
			default:
				return ErrArgumentMsg("int array <@push> expects numeric values")
			}
		}
		if varPtr != nil {
			varPtr.value = a
			vm.Stack = vm.Stack[:arrIdx] // remove array and items
			vm.Push(varPtr)
		} else {
			vm.Stack = vm.Stack[:arrIdx] // remove array and items
			vm.Push(a)
		}
	case []float64:
		for _, item := range items {
			switch v := item.(type) {
			case int:
				a = append(a, float64(v))
			case float64:
				a = append(a, v)
			default:
				return ErrArgumentMsg("float array <@push> expects numeric values")
			}
		}
		if varPtr != nil {
			varPtr.value = a
			vm.Stack = vm.Stack[:arrIdx]
			vm.Push(varPtr)
		} else {
			vm.Stack = vm.Stack[:arrIdx]
			vm.Push(a)
		}
	case []string:
		for _, item := range items {
			v, ok := item.(string)
			if !ok {
				return ErrArgumentMsg("string array <@push> expects string values")
			}
			a = append(a, v)
		}
		if varPtr != nil {
			varPtr.value = a
			vm.Stack = vm.Stack[:arrIdx]
			vm.Push(varPtr)
		} else {
			vm.Stack = vm.Stack[:arrIdx]
			vm.Push(a)
		}
	case []byte:
		for _, item := range items {
			v, ok := item.(int)
			if !ok {
				return ErrArgumentMsg("byte array <@push> expects integer values")
			}
			a = append(a, byte(v))
		}
		if varPtr != nil {
			varPtr.value = a
			vm.Stack = vm.Stack[:arrIdx]
			vm.Push(varPtr)
		} else {
			vm.Stack = vm.Stack[:arrIdx]
			vm.Push(a)
		}
	case []any:
		a = append(a, items...)
		if varPtr != nil {
			varPtr.value = a
			vm.Stack = vm.Stack[:arrIdx]
			vm.Push(varPtr)
		} else {
			vm.Stack = vm.Stack[:arrIdx]
			vm.Push(a)
		}
	default:
		return ErrArgumentMsg("<@push> expects typed array or variable on stack")
	}

	return nil
}

// <ints> ( item1 ... itemN count -- []int )
func varLenInts(vm *VM) error {
	countVal, err := vm.Pop()
	if err != nil {
		return err
	}
	count, ok := countVal.(int)
	if !ok || count < 0 {
		return ErrArgumentMsg("<ints> count must be non-negative integer")
	}

	currentLen := len(vm.Stack)
	if currentLen < count {
		return ErrUnderflowMsg("<ints> stack underflow")
	}

	items := vm.Stack[currentLen-count:]
	res := make([]int, 0, count)

	for _, item := range items {
		switch v := item.(type) {
		case int:
			res = append(res, v)
		case float64:
			if v != float64(int(v)) {
				return ErrArgumentMsg("lossy conversion from float to int in <ints>")
			}
			res = append(res, int(v))
		default:
			return ErrArgumentMsg("<ints> expects numeric values")
		}
	}

	vm.Stack = vm.Stack[:currentLen-count]
	vm.Push(res)
	return nil
}

// <floats> ( item1 ... itemN count -- []float64 )
func varLenFloats(vm *VM) error {
	countVal, err := vm.Pop()
	if err != nil {
		return err
	}
	count, ok := countVal.(int)
	if !ok || count < 0 {
		return ErrArgumentMsg("<floats> count must be non-negative integer")
	}

	currentLen := len(vm.Stack)
	if currentLen < count {
		return ErrUnderflowMsg("<floats> stack underflow")
	}

	items := vm.Stack[currentLen-count:]
	res := make([]float64, 0, count)

	for _, item := range items {
		switch v := item.(type) {
		case int:
			res = append(res, float64(v))
		case float64:
			res = append(res, v)
		default:
			return ErrArgumentMsg("<floats> expects numeric values")
		}
	}

	vm.Stack = vm.Stack[:currentLen-count]
	vm.Push(res)
	return nil
}

// <strings> ( item1 ... itemN count -- []string )
func varLenStrings(vm *VM) error {
	countVal, err := vm.Pop()
	if err != nil {
		return err
	}
	count, ok := countVal.(int)
	if !ok || count < 0 {
		return ErrArgumentMsg("<strings> count must be non-negative integer")
	}

	currentLen := len(vm.Stack)
	if currentLen < count {
		return ErrUnderflowMsg("<strings> stack underflow")
	}

	items := vm.Stack[currentLen-count:]
	res := make([]string, 0, count)

	for _, item := range items {
		v, ok := item.(string)
		if !ok {
			return ErrArgumentMsg("<strings> expects string values")
		}
		res = append(res, v)
	}

	vm.Stack = vm.Stack[:currentLen-count]
	vm.Push(res)
	return nil
}

// <bytes> ( item1 ... itemN count -- []byte )
func varLenBytes(vm *VM) error {
	countVal, err := vm.Pop()
	if err != nil {
		return err
	}
	count, ok := countVal.(int)
	if !ok || count < 0 {
		return ErrArgumentMsg("<bytes> count must be non-negative integer")
	}

	currentLen := len(vm.Stack)
	if currentLen < count {
		return ErrUnderflowMsg("<bytes> stack underflow")
	}

	items := vm.Stack[currentLen-count:]
	res := make([]byte, 0, count)

	for _, item := range items {
		v, ok := item.(int)
		if !ok {
			return ErrArgumentMsg("<bytes> expects integer values")
		}
		res = append(res, byte(v))
	}

	vm.Stack = vm.Stack[:currentLen-count]
	vm.Push(res)
	return nil
}

// <things> ( item1 ... itemN count -- []any )
func varLenAny(vm *VM) error {
	countVal, err := vm.Pop()
	if err != nil {
		return err
	}
	count, ok := countVal.(int)
	if !ok || count < 0 {
		return ErrArgumentMsg("<things> count must be non-negative integer")
	}

	currentLen := len(vm.Stack)
	if currentLen < count {
		return ErrUnderflowMsg("<things> stack underflow")
	}

	// Copy the slice from stack
	src := vm.Stack[currentLen-count:]
	res := make([]any, count)
	copy(res, src)

	vm.Stack = vm.Stack[:currentLen-count]
	vm.Push(res)
	return nil
}

// startVarLen (<<) pushes current stack depth to return stack
// In compile mode, compiles "depth >r"
func startVarLen(vm *VM) error {
	if vm.CurrentCompCtx() != nil {
		// Compiling
		depthIdx, ok1 := vm.dict["depth"]
		toRIdx, ok2 := vm.dict[">r"]
		if !ok1 || !ok2 {
			return ErrBadStateMsg("system words depth or >r missing")
		}
		vm.codeseg = append(vm.codeseg, depthIdx, toRIdx)
	} else {
		// Interpreting
		if err := depth(vm); err != nil {
			return err
		}
		if err := toR(vm); err != nil {
			return err
		}
	}
	return nil
}

// endVarLen (>>) calculates count of items since <<
// logic: depth r> -
func endVarLen(vm *VM) error {
	if vm.CurrentCompCtx() != nil {
		// Compiling
		depthIdx, ok1 := vm.dict["depth"]
		fromRIdx, ok2 := vm.dict["r>"]
		subIdx, ok3 := vm.dict["-"]

		if !ok1 || !ok2 || !ok3 {
			return ErrBadStateMsg("system words for >> missing")
		}

		vm.codeseg = append(vm.codeseg, depthIdx, fromRIdx, subIdx)
	} else {
		// Interpreting
		if err := depth(vm); err != nil {
			return err
		}
		if err := fromR(vm); err != nil {
			return err
		}
		// manual subtraction
		rVal, err := vm.Pop()
		if err != nil {
			return err
		}
		dVal, err := vm.Pop()
		if err != nil {
			return err
		}
		rInt, ok1 := rVal.(int)
		dInt, ok2 := dVal.(int)
		if !ok1 || !ok2 {
			return ErrBadStateMsg("stack depth calculation failed")
		}
		vm.Push(dInt - rInt)
	}
	return nil
}

// varLenQuote ( string1 ... stringN count -- )
// Parses whitespace-separated tokens until ">>". (Interpreting/Compiling)
func varLenQuote(vm *VM) error {
	buf := make([]rune, 0, 20)
	count := 0
	for {
		// 1. Eat whitespace to start next token
		ch, err := eatWhitespace(vm)
		if err != nil {
			if err == io.EOF {
				return ErrArgumentMsg("unexpected EOF inside <<\" ... \">>")
			}
			return err
		}

		// 2. Read token until whitespace
		// Prepare buf with 'ch'
		buf = buf[:0]
		buf = append(buf, ch)
		// Read subsequent chars
		buf, err = delimitedWSRead(vm, buf)
		if err != nil && err != io.EOF {
			return err
		}

		token := string(buf)

		if token == "\">>" {
			break
		}

		if vm.CurrentCompCtx() != nil {
			compileLiteral(vm, token)
		} else {
			vm.Push(token)
		}
		count++

		if err == io.EOF {
			return ErrArgumentMsg("unexpected EOF inside <<\" ... >>")
		}
	}

	if vm.CurrentCompCtx() != nil {
		compileLiteral(vm, count)
	} else {
		vm.Push(count)
	}

	return nil
}

// varLenWordsInit adds variable-length argument words to the VM
func varLenWordsInit(vm *VM) {
	vm.Define(&NativeWord{name: "sprintf", run: sprintf, immediate: false})
	vm.Define(&NativeWord{name: "<@push>", run: varLenPush, immediate: false})
	vm.Define(&NativeWord{name: "<ints>", run: varLenInts, immediate: false})
	vm.Define(&NativeWord{name: "<floats>", run: varLenFloats, immediate: false})
	vm.Define(&NativeWord{name: "<strings>", run: varLenStrings, immediate: false})
	vm.Define(&NativeWord{name: "<bytes>", run: varLenBytes, immediate: false})
	vm.Define(&NativeWord{name: "<things>", run: varLenAny, immediate: false})

	vm.Define(&NativeWord{name: "<<", run: startVarLen, immediate: true})
	vm.Define(&NativeWord{name: ">>", run: endVarLen, immediate: true})
	vm.Define(&NativeWord{name: "<<\"", run: varLenQuote, immediate: true})
}
