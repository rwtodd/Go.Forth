// SPDX-License-Identifier: MIT

package main

import (
	"flag"
	"fmt"
	"os"
	"strings"

	_ "github.com/rwtodd/Go.Forth/extensions/random"
	_ "github.com/rwtodd/Go.Forth/extensions/regex"
	_ "github.com/rwtodd/Go.Forth/extensions/strings"
	"github.com/rwtodd/Go.Forth/forth"
)

func main() {
	var extensions []string
	var preExprs []string
	var postExprs []string
	var interactive bool

	flag.Func("x", "Comma-separated list of extensions to activate", func(s string) error {
		exts := strings.Split(s, ",")
		for _, ext := range exts {
			ext = strings.TrimSpace(ext)
			if ext != "" {
				extensions = append(extensions, ext)
			}
		}
		return nil
	})

	flag.Func("e", "Evaluate expression before loading files", func(s string) error {
		preExprs = append(preExprs, s)
		return nil
	})

	flag.Func("ee", "Evaluate expression after loading files", func(s string) error {
		postExprs = append(postExprs, s)
		return nil
	})

	flag.BoolVar(&interactive, "i", false, "Run interactively after processing args")
	flag.Parse()

	files := flag.Args()

	vm := forth.NewVM(os.Stdin, os.Stdout)

	// 1. Activate extensions
	if len(extensions) > 0 {
		for _, ext := range extensions {
			vm.Push(ext)
		}
		vm.Push(int64(len(extensions)))
		if err := vm.Run(strings.NewReader("<activate-extensions>"), "cmdline"); err != nil {
			fmt.Fprintf(os.Stderr, "Error activating extensions: %v\n", err)
			os.Exit(1)
		}
	}

	// 2. Pre-expressions
	for _, expr := range preExprs {
		if err := vm.Run(strings.NewReader(expr), "cmdline"); err != nil {
			fmt.Fprintf(os.Stderr, "Error evaluating -e '%s': %v\n", expr, err)
			os.Exit(1)
		}
	}

	// 3. Load files
	for _, file := range files {
		if err := vm.Load(file); err != nil {
			fmt.Fprintf(os.Stderr, "Error loading %s: %v\n", file, err)
			os.Exit(1)
		}
	}

	// 4. Post-expressions
	for _, expr := range postExprs {
		if err := vm.Run(strings.NewReader(expr), "cmdline"); err != nil {
			fmt.Fprintf(os.Stderr, "Error evaluating -ee '%s': %v\n", expr, err)
			os.Exit(1)
		}
	}

	// 5. Interactive prompt
	hasArgs := len(files) > 0 || len(preExprs) > 0 || len(postExprs) > 0 || len(extensions) > 0
	if interactive || !hasArgs {
		for {
			err := vm.Run(os.Stdin, "stdin")
			if err == nil {
				break
			}
			fmt.Fprintf(os.Stderr, "Error: %v\n\n", err)
			vm.ResetState()
		}
	}
}
