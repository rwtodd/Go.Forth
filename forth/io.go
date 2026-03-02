// SPDX-License-Identifier: MIT

package forth

import (
	"fmt"
	"io"
	"unicode"
	"unicode/utf8"
)

// eatWhitespace eats whitespace and returns the next non-ws char
func eatWhitespace(r io.RuneReader) (rune, error) {
	var (
		ch  rune
		err error
	)

	for err == nil {
		ch, _, err = r.ReadRune()
		if (err != nil) || !unicode.IsSpace(ch) {
			return ch, err
		}
	}

	return 'X', nil
}

// delimitedRead reads from the `source` until the delimiter (a rune)
// is found.  It will use the provided `buf` to
// avoid allocation, if one is provided.
func delimitedRead(source io.RuneReader, delim rune, buf []rune) ([]rune, error) {
	var (
		ch  rune
		err error
	)

	for err == nil {
		ch, _, err = source.ReadRune()
		if (err != nil) || (ch == delim) {
			break
		}
		buf = append(buf, ch)
	}

	if err == io.EOF {
		err = nil
	}
	return buf, err
}

// delimitedWSRead reads from the `source` until whitespace
// is found.  It will use the provided `buf` to
// avoid allocation, if one is provided.
func delimitedWSRead(source io.RuneReader, buf []rune) ([]rune, error) {
	var (
		ch  rune
		err error
	)

	for err == nil {
		ch, _, err = source.ReadRune()
		if (err != nil) || unicode.IsSpace(ch) {
			break
		}
		buf = append(buf, ch)
	}
	if err == io.EOF {
		err = nil
	}

	return buf, err
}

// read looks at the top of the stack, and tries to interpret it
// as a rune.  If it can, then it reads until it finds that rune,
// and leaves the string it read at the top of the stack.
//
// A special case is when the delimiter is a space, in which case
// it reads until any whitespace is found.
func read(vm *VM) error {
	var (
		delim rune
		err   error
	)

	delimStack, err := vm.Pop()
	if err != nil {
		return err
	}

	switch delimT := delimStack.(type) {
	case int64:
		delim = rune(delimT)
	case string:
		var sz int
		delim, sz = utf8.DecodeRuneInString(delimT)
		// it needs to be a one-char string
		if sz != len(delimT) {
			return ErrArgumentMsg("read requires a single-character delimiter")
		}
	default:
		return ErrArgumentMsg("read requires an integer or string delimiter")
	}

	buf := make([]rune, 0, 20)

	if delim == ' ' {
		buf, err = delimitedWSRead(vm, buf)
	} else {
		buf, err = delimitedRead(vm, delim, buf)
	}
	vm.Push(string(buf))
	return err
}

// : skip ( delim -- ) read drop ;
func skip(vm *VM) error {
	err := read(vm)
	if err != nil {
		return err
	}
	_, err = vm.Pop()
	return err
}

// escapedRead reads from the `source` until the delimiter (a rune) is found,
// processing escape sequences along the way.
func escapedRead(source io.RuneReader, delim rune, buf []rune) ([]rune, error) {
	var (
		ch  rune
		err error
	)

	for err == nil {
		ch, _, err = source.ReadRune()
		if err != nil {
			break
		}
		if ch == delim {
			break
		}
		if ch == '\\' {
			// Escape sequence
			ch, _, err = source.ReadRune()
			if err != nil {
				break
			}
			switch ch {
			case 'n':
				ch = '\n'
			case 't':
				ch = '\t'
			case 'r':
				ch = '\r'
			case '\\':
				ch = '\\'
			case '"':
				ch = '"'
			default:
				// Unknown escape, keep literal backslash and char?
				// Or strict error? C usually keeps literal char.
				// user said "interpret common escape codes like \n \t \\ and \"".
				// Let's keep the backslash if it's not a known escape?
				// But simpler to just append 'ch' if it's not special?
				// Actually, if it's unknown, usually it's just the char.
				// e.g. \a -> a.
			}
		}
		buf = append(buf, ch)
	}

	if err == io.EOF {
		err = nil
	}
	return buf, err
}

// : " 34 read (compiling?) if postpone literal then ; immediate
func openQuote(vm *VM) error {
	buf, err := escapedRead(vm, '"', nil)
	if err != nil {
		return err
	}
	str := string(buf)
	if vm.CurrentCompCtx() != nil {
		compileLiteral(vm, str)
	} else {
		vm.Push(str)
	}
	return nil
}

// chrFromInt ('chr') takes an integer and makes a one-char string of it, interpreted
// as a rune
func chrFromInt(vm *VM) error {
	value, err := vm.Pop()
	if err != nil {
		return err
	}
	chInt, ok := value.(int64)
	if !ok {
		return ErrArgumentMsg("chr requires an integer")
	}
	vm.Push(string([]rune{rune(chInt)}))
	return nil
}

// ordFromStr ('ord') takes the first character of a string and gives its rune
// value as an int. It is the inverse of 'chr'.
func ordFromStr(vm *VM) error {
	value, err := vm.Pop()
	if err != nil {
		return err
	}
	chStr, ok := value.(string)
	if !ok {
		return ErrArgumentMsg("ord requires a string")
	}

	if len(chStr) == 0 {
		return ErrArgumentMsg("ord requires a non-empty string")
	}

	r, _ := utf8.DecodeRuneInString(chStr)
	vm.Push(int64(r))

	return nil
}

// printStack prints out the stack contents, without removing
// anything.
func printStack(vm *VM) error {
	tot := len(vm.Stack)
	for i, v := range vm.Stack {
		fmt.Fprintf(vm.Sink, "%2d: %v\n", tot-i, v)
	}
	return nil
}

// printTop prints out the top element on the stack, removing
// it in the process. It puts a trailing space after the item.
func printTop(vm *VM) error {
	v, err := vm.Pop()
	if err == nil {
		fmt.Fprint(vm.Sink, v, " ")
	}
	return err
}

// printOut ('type') prints out the top element on the stack, removing
// it in the process. It does not include a trailing space.
func printStr(vm *VM) error {
	v, err := vm.Pop()
	if err == nil {
		fmt.Fprint(vm.Sink, v)
	}
	return err
}

// cr simply prints a carriage return
func printCR(vm *VM) error {
	fmt.Fprintln(vm.Sink)
	return nil
}

// ioWordsInit adds the io-related core words to the VM.
func ioWordsInit(vm *VM) {
	vm.Define(&NativeWord{name: "read", run: read, immediate: false})
	vm.Define(&NativeWord{name: "skip", run: skip, immediate: false})
	vm.Define(&NativeWord{name: "\"", run: openQuote, immediate: true})
	vm.Define(&NativeWord{name: "chr", run: chrFromInt, immediate: false})
	vm.Define(&NativeWord{name: "ord", run: ordFromStr, immediate: false})
	vm.Define(&NativeWord{name: ".s", run: printStack, immediate: false})
	vm.Define(&NativeWord{name: ".", run: printTop, immediate: false})
	vm.Define(&NativeWord{name: "type", run: printStr, immediate: false})
	vm.Define(&NativeWord{name: "cr", run: printCR, immediate: false})
}
