package forth

import (
	"fmt"
)

func printStack(vm *VM) {
	tot := len(vm.Stack)
	for i, v := range vm.Stack {
		fmt.Printf("%2d: %v\n", tot-i, v)
	}
}

func printTop(vm *VM) {
	fmt.Print(vm.Pop(), " ")
}
