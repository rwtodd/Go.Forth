package main

import (
	"fmt"
	"os"

	"github.com/rwtodd/forth"
)

func main() {
	vm := forth.NewVM()

	for {
		vm.Run(os.Stdin, os.Stdout)
		if vm.Err == nil {
			break
		}
		fmt.Fprintf(os.Stderr, "Error: %v\n\n", vm.Err)
	}
}
