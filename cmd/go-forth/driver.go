// SPDX-License-Identifier: MIT

package main

import (
	"fmt"
	"os"
	"strings"

	_ "github.com/rwtodd/Go.Forth/extensions/random"
	_ "github.com/rwtodd/Go.Forth/extensions/regex"
	"github.com/rwtodd/Go.Forth/forth"
)

func main() {
	vm := forth.NewVM()
	vm.PushSource(os.Stdin, "stdin")

	// Parse arguments as files to load before interactive prompt
	if len(os.Args) > 1 {
		var sb strings.Builder
		for _, arg := range os.Args[1:] {
			fmt.Fprintf(&sb, "\" %s\" load\n", strings.ReplaceAll(arg, "\"", "\\\""))
		}
		vm.PushSource(strings.NewReader(sb.String()), "cmdline args")
	}

	for {
		err := vm.Run(os.Stdout)
		if err == nil {
			break
		}
		fmt.Fprintf(os.Stderr, "Error: %v\n\n", err)
		vm.ResetState()
	}
}
