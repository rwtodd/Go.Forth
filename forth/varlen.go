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

// endVarLenSprintf (>>sprintf) calculates count and calls (sprintf)
// In compile mode, compiles "depth r> - (sprintf)"
func endVarLenSprintf(vm *VM) error {
	if vm.CurrentCompCtx() != nil {
		// Compiling
		depthIdx, ok1 := vm.dict["depth"]
		fromRIdx, ok2 := vm.dict["r>"]
		subIdx, ok3 := vm.dict["-"]
		sprintfIdx, ok4 := vm.dict["(sprintf)"]

		if !ok1 || !ok2 || !ok3 || !ok4 {
			return ErrBadStateMsg("system words for >>sprintf missing")
		}

		vm.codeseg = append(vm.codeseg, depthIdx, fromRIdx, subIdx, sprintfIdx)
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

		if err := sprintf(vm); err != nil {
			return err
		}
	}
	return nil
}

// varLenWordsInit adds variable-length argument words to the VM
func varLenWordsInit(vm *VM) {
	vm.Define(&NativeWord{name: "(sprintf)", run: sprintf, immediate: false})
	vm.Define(&NativeWord{name: "<<", run: startVarLen, immediate: true})
	vm.Define(&NativeWord{name: ">>sprintf", run: endVarLenSprintf, immediate: true})
}
