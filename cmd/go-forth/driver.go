// SPDX-License-Identifier: MIT

package main

import (
	"fmt"
	"os"

	_ "github.com/rwtodd/Go.Forth/extensions/random"
	_ "github.com/rwtodd/Go.Forth/extensions/regex"
	_ "github.com/rwtodd/Go.Forth/extensions/strings"
	"github.com/rwtodd/Go.Forth/forth"
)

func main() {
	vm := forth.NewVM(os.Stdin, os.Stdout)

	// Parse arguments as files to load before interactive prompt
	if len(os.Args) > 1 {
		for _, arg := range os.Args[1:] {
			if err := vm.Load(arg); err != nil {
				fmt.Fprintf(os.Stderr, "Error loading %s: %v\n", arg, err)
				os.Exit(1)
			}
		}
	}

	for {
		err := vm.Run(os.Stdin, "stdin")
		if err == nil {
			break
		}
		fmt.Fprintf(os.Stderr, "Error: %v\n\n", err)
		vm.ResetState()
	}
}
