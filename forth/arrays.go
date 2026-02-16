// SPDX-License-Identifier: MIT

package forth

import "fmt"

// Array creation words

// bytes creates a byte array of given size
func bytes(vm *VM) error {
	size, err := vm.Pop()
	if err != nil {
		return err
	}
	sz, ok := size.(int)
	if !ok || sz < 0 {
		return ErrArgumentMsg("bytes expects non-negative integer size")
	}
	arr := make([]byte, sz)
	vm.Push(arr)
	return nil
}

// ints creates an int array of given size
func ints(vm *VM) error {
	size, err := vm.Pop()
	if err != nil {
		return err
	}
	sz, ok := size.(int)
	if !ok || sz < 0 {
		return ErrArgumentMsg("ints expects non-negative integer size")
	}
	arr := make([]int, sz)
	vm.Push(arr)
	return nil
}

// floats creates a float64 array of given size
func floats(vm *VM) error {
	size, err := vm.Pop()
	if err != nil {
		return err
	}
	sz, ok := size.(int)
	if !ok || sz < 0 {
		return ErrArgumentMsg("floats expects non-negative integer size")
	}
	arr := make([]float64, sz)
	vm.Push(arr)
	return nil
}

// stringArray creates a string array of given size
func stringArray(vm *VM) error {
	size, err := vm.Pop()
	if err != nil {
		return err
	}
	sz, ok := size.(int)
	if !ok || sz < 0 {
		return ErrArgumentMsg("strings expects non-negative integer size")
	}
	arr := make([]string, sz)
	vm.Push(arr)
	return nil
}

// Array access words

// @ gets element at index from array or value from variable
func arrayGet(vm *VM) error {
	idx, err := vm.Pop()
	if err != nil {
		return err
	}
	if varPtr, ok := idx.(*Variable); ok {
		// Variable access: variable @
		vm.Push(varPtr.value)
		return nil
	}
	// Array access: array index @
	arr, err := vm.Pop()
	if err != nil {
		return err
	}
	index, ok := idx.(int)
	if !ok {
		return ErrArgumentMsg("array index must be an integer")
	}
	var slice any
	if varPtr, ok := arr.(*Variable); ok {
		slice = varPtr.value
	} else {
		slice = arr
	}
	switch a := slice.(type) {
	case []byte:
		if index < 0 || index >= len(a) {
			return ErrIndexOutOfBoundsMsg(fmt.Sprintf("index %d out of bounds (len %d)", index, len(a)))
		}
		vm.Push(int(a[index]))
	case []int:
		if index < 0 || index >= len(a) {
			return ErrIndexOutOfBoundsMsg(fmt.Sprintf("index %d out of bounds (len %d)", index, len(a)))
		}
		vm.Push(a[index])
	case []float64:
		if index < 0 || index >= len(a) {
			return ErrIndexOutOfBoundsMsg(fmt.Sprintf("index %d out of bounds (len %d)", index, len(a)))
		}
		vm.Push(a[index])
	case []string:
		if index < 0 || index >= len(a) {
			return ErrIndexOutOfBoundsMsg(fmt.Sprintf("index %d out of bounds (len %d)", index, len(a)))
		}
		vm.Push(a[index])
	default:
		return ErrArgumentMsg("expected array or variable")
	}
	return nil
}

// ! sets element at index in array or value in variable
func arraySet(vm *VM) error {
	idx, err := vm.Pop()
	if err != nil {
		return err
	}
	if varPtr, ok := idx.(*Variable); ok {
		// Variable assignment: value variable !
		val, err := vm.Pop()
		if err != nil {
			return err
		}
		varPtr.value = val
		return nil
	}
	// Array assignment: value array index !
	arr, err := vm.Pop()
	if err != nil {
		return err
	}
	val, err := vm.Pop()
	if err != nil {
		return err
	}
	index, ok := idx.(int)
	if !ok {
		return ErrArgumentMsg("array index must be an integer")
	}
	var slice any
	var varPtr *Variable
	if v, ok := arr.(*Variable); ok {
		varPtr = v
		slice = v.value
	} else {
		slice = arr
	}
	switch a := slice.(type) {
	case []byte:
		if index < 0 || index >= len(a) {
			return ErrIndexOutOfBoundsMsg(fmt.Sprintf("index %d out of bounds (len %d)", index, len(a)))
		}
		v, ok := val.(int)
		if !ok {
			return ErrArgumentMsg("byte array requires integer value")
		}
		a[index] = byte(v)
		if varPtr != nil {
			varPtr.value = a
		}
	case []int:
		if index < 0 || index >= len(a) {
			return ErrIndexOutOfBoundsMsg(fmt.Sprintf("index %d out of bounds (len %d)", index, len(a)))
		}
		switch v := val.(type) {
		case int:
			a[index] = v
		case float64:
			if v != float64(int(v)) {
				return ErrArgumentMsg("lossy conversion from float to int")
			}
			a[index] = int(v)
		default:
			return ErrArgumentMsg("int array requires numeric value")
		}
		if varPtr != nil {
			varPtr.value = a
		}
	case []float64:
		if index < 0 || index >= len(a) {
			return ErrIndexOutOfBoundsMsg(fmt.Sprintf("index %d out of bounds (len %d)", index, len(a)))
		}
		switch v := val.(type) {
		case int:
			a[index] = float64(v)
		case float64:
			a[index] = v
		default:
			return ErrArgumentMsg("float array requires numeric value")
		}
		if varPtr != nil {
			varPtr.value = a
		}
	case []string:
		if index < 0 || index >= len(a) {
			return ErrIndexOutOfBoundsMsg(fmt.Sprintf("index %d out of bounds (len %d)", index, len(a)))
		}
		v, ok := val.(string)
		if !ok {
			return ErrArgumentMsg("string array requires string value")
		}
		a[index] = v
		if varPtr != nil {
			varPtr.value = a
		}
	default:
		return ErrArgumentMsg("expected array or variable")
	}
	return nil
}

// c@ gets byte at index with wrapping
func byteGet(vm *VM) error {
	idx, err := vm.Pop()
	if err != nil {
		return err
	}
	arr, err := vm.Pop()
	if err != nil {
		return err
	}
	index, ok := idx.(int)
	if !ok {
		return ErrArgumentMsg("c@ expects integer index")
	}
	a, ok := arr.([]byte)
	if !ok {
		return ErrArgumentMsg("c@ expects byte array")
	}
	if index < 0 || index >= len(a) {
		return ErrIndexOutOfBoundsMsg(fmt.Sprintf("index %d out of bounds (len %d)", index, len(a)))
	}
	vm.Push(int(a[index] & 0xff))
	return nil
}

// c! sets byte at index with wrapping
func byteSet(vm *VM) error {
	idx, err := vm.Pop()
	if err != nil {
		return err
	}
	arr, err := vm.Pop()
	if err != nil {
		return err
	}
	val, err := vm.Pop()
	if err != nil {
		return err
	}
	index, ok := idx.(int)
	if !ok {
		return ErrArgumentMsg("c! expects integer index")
	}
	a, ok := arr.([]byte)
	if !ok {
		return ErrArgumentMsg("c! expects byte array")
	}
	if index < 0 || index >= len(a) {
		return ErrIndexOutOfBoundsMsg(fmt.Sprintf("index %d out of bounds (len %d)", index, len(a)))
	}
	v, ok := val.(int)
	if !ok {
		return ErrArgumentMsg("c! expects integer value")
	}
	a[index] = byte(v & 0xff)
	return nil
}

// Array manipulation words

// @push appends value to array and returns modified array
func arrayPush(vm *VM) error {
	arr, err := vm.Pop()
	if err != nil {
		return err
	}
	val, err := vm.Pop()
	if err != nil {
		return err
	}
	var varPtr *Variable
	if v, ok := arr.(*Variable); ok {
		varPtr = v
		arr = v.value
	}
	switch a := arr.(type) {
	case []byte:
		v, ok := val.(int)
		if !ok {
			return ErrArgumentMsg("byte array push expects integer")
		}
		newArr := append(a, byte(v))
		if varPtr != nil {
			varPtr.value = newArr
			vm.Push(varPtr)
		} else {
			vm.Push(newArr)
		}
	case []int:
		switch v := val.(type) {
		case int:
			newArr := append(a, v)
			if varPtr != nil {
				varPtr.value = newArr
				vm.Push(varPtr)
			} else {
				vm.Push(newArr)
			}
		case float64:
			if v != float64(int(v)) {
				return ErrArgumentMsg("lossy conversion from float to int")
			}
			newArr := append(a, int(v))
			if varPtr != nil {
				varPtr.value = newArr
				vm.Push(varPtr)
			} else {
				vm.Push(newArr)
			}
		default:
			return ErrArgumentMsg("int array push expects numeric value")
		}
	case []float64:
		switch v := val.(type) {
		case int:
			newArr := append(a, float64(v))
			if varPtr != nil {
				varPtr.value = newArr
				vm.Push(varPtr)
			} else {
				vm.Push(newArr)
			}
		case float64:
			newArr := append(a, v)
			if varPtr != nil {
				varPtr.value = newArr
				vm.Push(varPtr)
			} else {
				vm.Push(newArr)
			}
		default:
			return ErrArgumentMsg("float array push expects numeric value")
		}
	case []string:
		v, ok := val.(string)
		if !ok {
			return ErrArgumentMsg("string array push expects string")
		}
		newArr := append(a, v)
		if varPtr != nil {
			varPtr.value = newArr
			vm.Push(varPtr)
		} else {
			vm.Push(newArr)
		}
	default:
		return ErrArgumentMsg("expected array or variable")
	}
	return nil
}

// @pop removes and returns last element and modified array
func arrayPop(vm *VM) error {
	arr, err := vm.Pop()
	if err != nil {
		return err
	}
	var varPtr *Variable
	if v, ok := arr.(*Variable); ok {
		varPtr = v
		arr = v.value
	}
	switch a := arr.(type) {
	case []byte:
		if len(a) == 0 {
			return ErrArgumentMsg("cannot pop empty array")
		}
		val := a[len(a)-1]
		newArr := a[:len(a)-1]
		if varPtr != nil {
			varPtr.value = newArr
			vm.Push(varPtr)
		} else {
			vm.Push(newArr)
		}
		vm.Push(int(val))
	case []int:
		if len(a) == 0 {
			return ErrArgumentMsg("cannot pop empty array")
		}
		val := a[len(a)-1]
		newArr := a[:len(a)-1]
		if varPtr != nil {
			varPtr.value = newArr
			vm.Push(varPtr)
		} else {
			vm.Push(newArr)
		}
		vm.Push(val)
	case []float64:
		if len(a) == 0 {
			return ErrArgumentMsg("cannot pop empty array")
		}
		val := a[len(a)-1]
		newArr := a[:len(a)-1]
		if varPtr != nil {
			varPtr.value = newArr
			vm.Push(varPtr)
		} else {
			vm.Push(newArr)
		}
		vm.Push(val)
	case []string:
		if len(a) == 0 {
			return ErrArgumentMsg("cannot pop empty array")
		}
		val := a[len(a)-1]
		newArr := a[:len(a)-1]
		if varPtr != nil {
			varPtr.value = newArr
			vm.Push(varPtr)
		} else {
			vm.Push(newArr)
		}
		vm.Push(val)
	default:
		return ErrArgumentMsg("expected array or variable")
	}
	return nil
}

// @shift removes and returns first element and modified array
func arrayShift(vm *VM) error {
	arr, err := vm.Pop()
	if err != nil {
		return err
	}
	var varPtr *Variable
	if v, ok := arr.(*Variable); ok {
		varPtr = v
		arr = v.value
	}
	switch a := arr.(type) {
	case []byte:
		if len(a) == 0 {
			return ErrArgumentMsg("cannot shift empty array")
		}
		val := a[0]
		newArr := a[1:]
		if varPtr != nil {
			varPtr.value = newArr
			vm.Push(varPtr)
		} else {
			vm.Push(newArr)
		}
		vm.Push(int(val))
	case []int:
		if len(a) == 0 {
			return ErrArgumentMsg("cannot shift empty array")
		}
		val := a[0]
		newArr := a[1:]
		if varPtr != nil {
			varPtr.value = newArr
			vm.Push(varPtr)
		} else {
			vm.Push(newArr)
		}
		vm.Push(val)
	case []float64:
		if len(a) == 0 {
			return ErrArgumentMsg("cannot shift empty array")
		}
		val := a[0]
		newArr := a[1:]
		if varPtr != nil {
			varPtr.value = newArr
			vm.Push(varPtr)
		} else {
			vm.Push(newArr)
		}
		vm.Push(val)
	case []string:
		if len(a) == 0 {
			return ErrArgumentMsg("cannot shift empty array")
		}
		val := a[0]
		newArr := a[1:]
		if varPtr != nil {
			varPtr.value = newArr
			vm.Push(varPtr)
		} else {
			vm.Push(newArr)
		}
		vm.Push(val)
	default:
		return ErrArgumentMsg("expected array or variable")
	}
	return nil
}

// @unshift prepends value to array and returns modified array
func arrayUnshift(vm *VM) error {
	arr, err := vm.Pop()
	if err != nil {
		return err
	}
	val, err := vm.Pop()
	if err != nil {
		return err
	}
	var varPtr *Variable
	if v, ok := arr.(*Variable); ok {
		varPtr = v
		arr = v.value
	}
	switch a := arr.(type) {
	case []byte:
		v, ok := val.(int)
		if !ok {
			return ErrArgumentMsg("byte array unshift expects integer")
		}
		newArr := append([]byte{byte(v)}, a...)
		if varPtr != nil {
			varPtr.value = newArr
			vm.Push(varPtr)
		} else {
			vm.Push(newArr)
		}
	case []int:
		switch v := val.(type) {
		case int:
			newArr := append([]int{v}, a...)
			if varPtr != nil {
				varPtr.value = newArr
				vm.Push(varPtr)
			} else {
				vm.Push(newArr)
			}
		case float64:
			if v != float64(int(v)) {
				return ErrArgumentMsg("lossy conversion from float to int")
			}
			newArr := append([]int{int(v)}, a...)
			if varPtr != nil {
				varPtr.value = newArr
				vm.Push(varPtr)
			} else {
				vm.Push(newArr)
			}
		default:
			return ErrArgumentMsg("int array unshift expects numeric value")
		}
	case []float64:
		switch v := val.(type) {
		case int:
			newArr := append([]float64{float64(v)}, a...)
			if varPtr != nil {
				varPtr.value = newArr
				vm.Push(varPtr)
			} else {
				vm.Push(newArr)
			}
		case float64:
			newArr := append([]float64{v}, a...)
			if varPtr != nil {
				varPtr.value = newArr
				vm.Push(varPtr)
			} else {
				vm.Push(newArr)
			}
		default:
			return ErrArgumentMsg("float array unshift expects numeric value")
		}
	case []string:
		v, ok := val.(string)
		if !ok {
			return ErrArgumentMsg("string array unshift expects string")
		}
		newArr := append([]string{v}, a...)
		if varPtr != nil {
			varPtr.value = newArr
			vm.Push(varPtr)
		} else {
			vm.Push(newArr)
		}
	default:
		return ErrArgumentMsg("expected array or variable")
	}
	return nil
}

// @len returns the length of the array
func arrayLen(vm *VM) error {
	arr, err := vm.Pop()
	if err != nil {
		return err
	}
	if varPtr, ok := arr.(*Variable); ok {
		arr = varPtr.value
	}
	switch a := arr.(type) {
	case []byte:
		vm.Push(len(a))
	case []int:
		vm.Push(len(a))
	case []float64:
		vm.Push(len(a))
	case []string:
		vm.Push(len(a))
	default:
		return ErrArgumentMsg("expected array or variable")
	}
	return nil
}

// arrayWordsInit adds array-related words to the VM
func arrayWordsInit(vm *VM) {
	vm.Define(&NativeWord{name: "bytes", run: bytes, immediate: false})
	vm.Define(&NativeWord{name: "ints", run: ints, immediate: false})
	vm.Define(&NativeWord{name: "floats", run: floats, immediate: false})
	vm.Define(&NativeWord{name: "strings", run: stringArray, immediate: false})
	vm.Define(&NativeWord{name: "@", run: arrayGet, immediate: false})
	vm.Define(&NativeWord{name: "!", run: arraySet, immediate: false})
	vm.Define(&NativeWord{name: "c@", run: byteGet, immediate: false})
	vm.Define(&NativeWord{name: "c!", run: byteSet, immediate: false})
	vm.Define(&NativeWord{name: "@push", run: arrayPush, immediate: false})
	vm.Define(&NativeWord{name: "@pop", run: arrayPop, immediate: false})
	vm.Define(&NativeWord{name: "@shift", run: arrayShift, immediate: false})
	vm.Define(&NativeWord{name: "@unshift", run: arrayUnshift, immediate: false})
	vm.Define(&NativeWord{name: "@len", run: arrayLen, immediate: false})
}
