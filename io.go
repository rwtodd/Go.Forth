package forth

import (
	"fmt"
)

func printStack(vm *VM) error {
	tot := len(vm.Stack)
	for i, v := range vm.Stack {
		fmt.Printf("%2d: %v\n", tot-i, v)
	}
	return nil
}

func printTop(vm *VM) error {
	v, err := vm.Pop()
	if err == nil {
		fmt.Print(v, " ")
	}
	return err
}
