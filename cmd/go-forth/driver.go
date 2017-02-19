package main

import (
	"fmt"
	"os"

	"github.com/rwtodd/forth"
)

func main() {
	vm := forth.NewVM()

	for {
		err := vm.Run(os.Stdin, os.Stdout)
		if err == nil {
			break
		}
		fmt.Fprintf(os.Stderr, "Error: %v\n\n", err)
		vm.ResetState()
	}
}
