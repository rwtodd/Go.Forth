// SPDX-License-Identifier: MIT

package forth

import (
	"fmt"
)

// (sprintf) ( format args... count -- result )
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

// (vpush) ( array item1 ... itemN count -- result )
func vpush(vm *VM) error {
	// 1. Pop count
	countVal, err := vm.Pop()
	if err != nil {
		return err
	}
	count, ok := countVal.(int)
	if !ok || count < 0 {
		return ErrArgumentMsg("vpush count must be non-negative integer")
	}

	currentLen := len(vm.Stack)
	// We need array + count items on stack.
	if currentLen < count+1 {
		return ErrUnderflowMsg("vpush stack underflow")
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
					return ErrArgumentMsg("lossy conversion from float to int in vpush")
				}
				a = append(a, int(v))
			default:
				return ErrArgumentMsg("int array vpush expects numeric values")
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
				return ErrArgumentMsg("float array vpush expects numeric values")
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
				return ErrArgumentMsg("string array vpush expects string values")
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
				return ErrArgumentMsg("byte array vpush expects integer values")
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
		return ErrArgumentMsg("vpush expects typed array or variable on stack")
	}

	return nil
}

// (vmake-int) ( item1 ... itemN count -- []int )
func vmakeInt(vm *VM) error {
	countVal, err := vm.Pop()
	if err != nil {
		return err
	}
	count, ok := countVal.(int)
	if !ok || count < 0 {
		return ErrArgumentMsg("vmake-int count must be non-negative integer")
	}

	currentLen := len(vm.Stack)
	if currentLen < count {
		return ErrUnderflowMsg("vmake-int stack underflow")
	}

	items := vm.Stack[currentLen-count:]
	res := make([]int, 0, count)

	for _, item := range items {
		switch v := item.(type) {
		case int:
			res = append(res, v)
		case float64:
			if v != float64(int(v)) {
				return ErrArgumentMsg("lossy conversion from float to int in vmake-int")
			}
			res = append(res, int(v))
		default:
			return ErrArgumentMsg("vmake-int expects numeric values")
		}
	}

	vm.Stack = vm.Stack[:currentLen-count]
	vm.Push(res)
	return nil
}

// (vmake-float) ( item1 ... itemN count -- []float64 )
func vmakeFloat(vm *VM) error {
	countVal, err := vm.Pop()
	if err != nil {
		return err
	}
	count, ok := countVal.(int)
	if !ok || count < 0 {
		return ErrArgumentMsg("vmake-float count must be non-negative integer")
	}

	currentLen := len(vm.Stack)
	if currentLen < count {
		return ErrUnderflowMsg("vmake-float stack underflow")
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
			return ErrArgumentMsg("vmake-float expects numeric values")
		}
	}

	vm.Stack = vm.Stack[:currentLen-count]
	vm.Push(res)
	return nil
}

// (vmake-string) ( item1 ... itemN count -- []string )
func vmakeString(vm *VM) error {
	countVal, err := vm.Pop()
	if err != nil {
		return err
	}
	count, ok := countVal.(int)
	if !ok || count < 0 {
		return ErrArgumentMsg("vmake-string count must be non-negative integer")
	}

	currentLen := len(vm.Stack)
	if currentLen < count {
		return ErrUnderflowMsg("vmake-string stack underflow")
	}

	items := vm.Stack[currentLen-count:]
	res := make([]string, 0, count)

	for _, item := range items {
		v, ok := item.(string)
		if !ok {
			return ErrArgumentMsg("vmake-string expects string values")
		}
		res = append(res, v)
	}

	vm.Stack = vm.Stack[:currentLen-count]
	vm.Push(res)
	return nil
}

// (vmake-byte) ( item1 ... itemN count -- []byte )
func vmakeByte(vm *VM) error {
	countVal, err := vm.Pop()
	if err != nil {
		return err
	}
	count, ok := countVal.(int)
	if !ok || count < 0 {
		return ErrArgumentMsg("vmake-byte count must be non-negative integer")
	}

	currentLen := len(vm.Stack)
	if currentLen < count {
		return ErrUnderflowMsg("vmake-byte stack underflow")
	}

	items := vm.Stack[currentLen-count:]
	res := make([]byte, 0, count)

	for _, item := range items {
		v, ok := item.(int)
		if !ok {
			return ErrArgumentMsg("vmake-byte expects integer values")
		}
		res = append(res, byte(v))
	}

	vm.Stack = vm.Stack[:currentLen-count]
	vm.Push(res)
	return nil
}

// (vmake-any) ( item1 ... itemN count -- []any )
func vmakeAny(vm *VM) error {
	countVal, err := vm.Pop()
	if err != nil {
		return err
	}
	count, ok := countVal.(int)
	if !ok || count < 0 {
		return ErrArgumentMsg("vmake-any count must be non-negative integer")
	}

	currentLen := len(vm.Stack)
	if currentLen < count {
		return ErrUnderflowMsg("vmake-any stack underflow")
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

// endVarLenHelper compiles or runs logic to calculate count and call target
func endVarLenHelper(vm *VM, targetName string) error {
	if vm.CurrentCompCtx() != nil {
		// Compiling
		depthIdx, ok1 := vm.dict["depth"]
		fromRIdx, ok2 := vm.dict["r>"]
		subIdx, ok3 := vm.dict["-"]
		targetIdx, ok4 := vm.dict[targetName]

		if !ok1 || !ok2 || !ok3 || !ok4 {
			return ErrBadStateMsg(fmt.Sprintf("system words for >>%s missing", targetName))
		}

		vm.codeseg = append(vm.codeseg, depthIdx, fromRIdx, subIdx, targetIdx)
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

		// Call target
		targetIdx, ok := vm.dict[targetName]
		if !ok {
			return ErrBadStateMsg(fmt.Sprintf("target word %s missing", targetName))
		}
		return vm.words[targetIdx].Run(vm)
	}
	return nil
}

// endVarLenSprintf (>>sprintf)
func endVarLenSprintf(vm *VM) error {
	return endVarLenHelper(vm, "(sprintf)")
}

// endVarLenPush (>>@push)
func endVarLenPush(vm *VM) error {
	return endVarLenHelper(vm, "(vpush)")
}

// endVarLenInts (>>@i)
func endVarLenInts(vm *VM) error {
	return endVarLenHelper(vm, "(vmake-int)")
}

// endVarLenFloats (>>@f)
func endVarLenFloats(vm *VM) error {
	return endVarLenHelper(vm, "(vmake-float)")
}

// endVarLenStrings (>>@s)
func endVarLenStrings(vm *VM) error {
	return endVarLenHelper(vm, "(vmake-string)")
}

// endVarLenBytes (>>@b)
func endVarLenBytes(vm *VM) error {
	return endVarLenHelper(vm, "(vmake-byte)")
}

// endVarLenAny (>>)
func endVarLenAny(vm *VM) error {
	return endVarLenHelper(vm, "(vmake-any)")
}

// varLenWordsInit adds variable-length argument words to the VM
func varLenWordsInit(vm *VM) {
	vm.Define(&NativeWord{name: "(sprintf)", run: sprintf, immediate: false})
	vm.Define(&NativeWord{name: "(vpush)", run: vpush, immediate: false})
	vm.Define(&NativeWord{name: "(vmake-int)", run: vmakeInt, immediate: false})
	vm.Define(&NativeWord{name: "(vmake-float)", run: vmakeFloat, immediate: false})
	vm.Define(&NativeWord{name: "(vmake-string)", run: vmakeString, immediate: false})
	vm.Define(&NativeWord{name: "(vmake-byte)", run: vmakeByte, immediate: false})
	vm.Define(&NativeWord{name: "(vmake-any)", run: vmakeAny, immediate: false})

	vm.Define(&NativeWord{name: "<<", run: startVarLen, immediate: true})
	vm.Define(&NativeWord{name: ">>sprintf", run: endVarLenSprintf, immediate: true})
	vm.Define(&NativeWord{name: ">>@push", run: endVarLenPush, immediate: true})
	vm.Define(&NativeWord{name: ">>@i", run: endVarLenInts, immediate: true})
	vm.Define(&NativeWord{name: ">>@f", run: endVarLenFloats, immediate: true})
	vm.Define(&NativeWord{name: ">>@s", run: endVarLenStrings, immediate: true})
	vm.Define(&NativeWord{name: ">>@b", run: endVarLenBytes, immediate: true})
	vm.Define(&NativeWord{name: ">>", run: endVarLenAny, immediate: true})
}
