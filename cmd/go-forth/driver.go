// SPDX-License-Identifier: MIT

package main

import (
	"fmt"
	"os"

	_ "github.com/rwtodd/Go.Forth/extensions/regex"
	"github.com/rwtodd/Go.Forth/forth"
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
