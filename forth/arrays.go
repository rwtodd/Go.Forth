// SPDX-License-Identifier: MIT

package forth

// Array creation words

// bytes creates a byte array of given size
func bytes(vm *VM) error {
	size, err := vm.Pop()
	if err != nil {
		return err
	}
	sz, ok := size.(int)
	if !ok || sz < 0 {
		return ErrArgument
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
		return ErrArgument
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
		return ErrArgument
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
		return ErrArgument
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
		return ErrArgument
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
			return ErrIndexOutOfBounds
		}
		vm.Push(int(a[index]))
	case []int:
		if index < 0 || index >= len(a) {
			return ErrIndexOutOfBounds
		}
		vm.Push(a[index])
	case []float64:
		if index < 0 || index >= len(a) {
			return ErrIndexOutOfBounds
		}
		vm.Push(a[index])
	case []string:
		if index < 0 || index >= len(a) {
			return ErrIndexOutOfBounds
		}
		vm.Push(a[index])
	default:
		return ErrArgument
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
		return ErrArgument
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
			return ErrIndexOutOfBounds
		}
		v, ok := val.(int)
		if !ok {
			return ErrArgument
		}
		a[index] = byte(v)
		if varPtr != nil {
			varPtr.value = a
		}
	case []int:
		if index < 0 || index >= len(a) {
			return ErrIndexOutOfBounds
		}
		switch v := val.(type) {
		case int:
			a[index] = v
		case float64:
			if v != float64(int(v)) {
				return ErrArgument // lossy conversion
			}
			a[index] = int(v)
		default:
			return ErrArgument
		}
		if varPtr != nil {
			varPtr.value = a
		}
	case []float64:
		if index < 0 || index >= len(a) {
			return ErrIndexOutOfBounds
		}
		switch v := val.(type) {
		case int:
			a[index] = float64(v)
		case float64:
			a[index] = v
		default:
			return ErrArgument
		}
		if varPtr != nil {
			varPtr.value = a
		}
	case []string:
		if index < 0 || index >= len(a) {
			return ErrIndexOutOfBounds
		}
		v, ok := val.(string)
		if !ok {
			return ErrArgument
		}
		a[index] = v
		if varPtr != nil {
			varPtr.value = a
		}
	default:
		return ErrArgument
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
		return ErrArgument
	}
	a, ok := arr.([]byte)
	if !ok {
		return ErrArgument
	}
	if index < 0 || index >= len(a) {
		return ErrIndexOutOfBounds
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
		return ErrArgument
	}
	a, ok := arr.([]byte)
	if !ok {
		return ErrArgument
	}
	if index < 0 || index >= len(a) {
		return ErrIndexOutOfBounds
	}
	v, ok := val.(int)
	if !ok {
		return ErrArgument
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
			return ErrArgument
		}
		newArr := append(a, byte(v))
		if varPtr != nil {
			varPtr.value = newArr
		} else {
			vm.Push(newArr)
		}
	case []int:
		switch v := val.(type) {
		case int:
			newArr := append(a, v)
			if varPtr != nil {
				varPtr.value = newArr
			} else {
				vm.Push(newArr)
			}
		case float64:
			if v != float64(int(v)) {
				return ErrArgument
			}
			newArr := append(a, int(v))
			if varPtr != nil {
				varPtr.value = newArr
			} else {
				vm.Push(newArr)
			}
		default:
			return ErrArgument
		}
	case []float64:
		switch v := val.(type) {
		case int:
			newArr := append(a, float64(v))
			if varPtr != nil {
				varPtr.value = newArr
			} else {
				vm.Push(newArr)
			}
		case float64:
			newArr := append(a, v)
			if varPtr != nil {
				varPtr.value = newArr
			} else {
				vm.Push(newArr)
			}
		default:
			return ErrArgument
		}
	case []string:
		v, ok := val.(string)
		if !ok {
			return ErrArgument
		}
		newArr := append(a, v)
		if varPtr != nil {
			varPtr.value = newArr
		} else {
			vm.Push(newArr)
		}
	default:
		return ErrArgument
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
			return ErrArgument
		}
		val := a[len(a)-1]
		newArr := a[:len(a)-1]
		if varPtr != nil {
			varPtr.value = newArr
		} else {
			vm.Push(newArr)
		}
		vm.Push(int(val))
	case []int:
		if len(a) == 0 {
			return ErrArgument
		}
		val := a[len(a)-1]
		newArr := a[:len(a)-1]
		if varPtr != nil {
			varPtr.value = newArr
		} else {
			vm.Push(newArr)
		}
		vm.Push(val)
	case []float64:
		if len(a) == 0 {
			return ErrArgument
		}
		val := a[len(a)-1]
		newArr := a[:len(a)-1]
		if varPtr != nil {
			varPtr.value = newArr
		} else {
			vm.Push(newArr)
		}
		vm.Push(val)
	case []string:
		if len(a) == 0 {
			return ErrArgument
		}
		val := a[len(a)-1]
		newArr := a[:len(a)-1]
		if varPtr != nil {
			varPtr.value = newArr
		} else {
			vm.Push(newArr)
		}
		vm.Push(val)
	default:
		return ErrArgument
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
			return ErrArgument
		}
		val := a[0]
		newArr := a[1:]
		if varPtr != nil {
			varPtr.value = newArr
		} else {
			vm.Push(newArr)
		}
		vm.Push(int(val))
	case []int:
		if len(a) == 0 {
			return ErrArgument
		}
		val := a[0]
		newArr := a[1:]
		if varPtr != nil {
			varPtr.value = newArr
		} else {
			vm.Push(newArr)
		}
		vm.Push(val)
	case []float64:
		if len(a) == 0 {
			return ErrArgument
		}
		val := a[0]
		newArr := a[1:]
		if varPtr != nil {
			varPtr.value = newArr
		} else {
			vm.Push(newArr)
		}
		vm.Push(val)
	case []string:
		if len(a) == 0 {
			return ErrArgument
		}
		val := a[0]
		newArr := a[1:]
		if varPtr != nil {
			varPtr.value = newArr
		} else {
			vm.Push(newArr)
		}
		vm.Push(val)
	default:
		return ErrArgument
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
			return ErrArgument
		}
		newArr := append([]byte{byte(v)}, a...)
		if varPtr != nil {
			varPtr.value = newArr
		} else {
			vm.Push(newArr)
		}
	case []int:
		switch v := val.(type) {
		case int:
			newArr := append([]int{v}, a...)
			if varPtr != nil {
				varPtr.value = newArr
			} else {
				vm.Push(newArr)
			}
		case float64:
			if v != float64(int(v)) {
				return ErrArgument
			}
			newArr := append([]int{int(v)}, a...)
			if varPtr != nil {
				varPtr.value = newArr
			} else {
				vm.Push(newArr)
			}
		default:
			return ErrArgument
		}
	case []float64:
		switch v := val.(type) {
		case int:
			newArr := append([]float64{float64(v)}, a...)
			if varPtr != nil {
				varPtr.value = newArr
			} else {
				vm.Push(newArr)
			}
		case float64:
			newArr := append([]float64{v}, a...)
			if varPtr != nil {
				varPtr.value = newArr
			} else {
				vm.Push(newArr)
			}
		default:
			return ErrArgument
		}
	case []string:
		v, ok := val.(string)
		if !ok {
			return ErrArgument
		}
		newArr := append([]string{v}, a...)
		if varPtr != nil {
			varPtr.value = newArr
		} else {
			vm.Push(newArr)
		}
	default:
		return ErrArgument
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
		return ErrArgument
	}
	return nil
}

// arrayWordsInit adds array-related words to the VM
func arrayWordsInit(vm *VM) {
	vm.Define("bytes", Word{bytes, false})
	vm.Define("ints", Word{ints, false})
	vm.Define("floats", Word{floats, false})
	vm.Define("strings", Word{stringArray, false})
	vm.Define("@", Word{arrayGet, false})
	vm.Define("!", Word{arraySet, false})
	vm.Define("c@", Word{byteGet, false})
	vm.Define("c!", Word{byteSet, false})
	vm.Define("@push", Word{arrayPush, false})
	vm.Define("@pop", Word{arrayPop, false})
	vm.Define("@shift", Word{arrayShift, false})
	vm.Define("@unshift", Word{arrayUnshift, false})
	vm.Define("@len", Word{arrayLen, false})
}
