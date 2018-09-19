package forth

import (
	"bufio"
	"fmt"
	"io"
	"unicode"
	"unicode/utf8"
)

// eatWhitespace eats whitespace and returns the next non-ws char
func eatWhitespace(r *bufio.Reader) (rune, error) {
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
func delimitedRead(source *bufio.Reader, delim rune, buf []rune) ([]rune, error) {
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
func delimitedWSRead(source *bufio.Reader, buf []rune) ([]rune, error) {
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
	case int:
		delim = rune(delimT)
	case string:
		var sz int
		delim, sz = utf8.DecodeRuneInString(delimT)
		// it needs to be a one-char string
		if sz != len(delimT) {
			return ErrArgument
		}
	default:
		return ErrArgument
	}

	buf := make([]rune, 0, 20)

	if delim == ' ' {
		buf, err = delimitedWSRead(vm.Source, buf)
	} else {
		buf, err = delimitedRead(vm.Source, delim, buf)
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

// : " 34 read (compiling?) if postpone literal then ; immediate
func openQuote(vm *VM) error {
	buf, err := delimitedRead(vm.Source, '"', nil)
	if err != nil {
		return err
	}
	str := string(buf)
	if vm.Compiling {
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
	chInt, ok := value.(int)
	if !ok {
		return ErrArgument
	}
	vm.Push(string([]rune{rune(chInt)}))
	return nil
}

// ordFromStr ('ord') takes a one-character string and gives its rune
// value as an int. It is the inverse of 'chr'.
func ordFromStr(vm *VM) error {
	value, err := vm.Pop()
	if err != nil {
		return err
	}
	chStr, ok := value.(string)
	if !ok {
		return ErrArgument
	}
	r, sz := utf8.DecodeRuneInString(chStr)
	// it needs to be a one-char string
	if sz != len(chStr) {
		return ErrArgument
	}
	vm.Push(int(r))

	return nil
}

// printStack prints out the stack contents, without removing
// anything.
func printStack(vm *VM) error {
	tot := len(vm.Stack)
	for i, v := range vm.Stack {
		fmt.Printf("%2d: %v\n", tot-i, v)
	}
	return nil
}

// printTop prints out the top element on the stack, removing
// it in the process. It puts a trailing space after the item.
func printTop(vm *VM) error {
	v, err := vm.Pop()
	if err == nil {
		fmt.Print(v, " ")
	}
	return err
}

// printOut ('type') prints out the top element on the stack, removing
// it in the process. It does not include a trailing space.
func printStr(vm *VM) error {
	v, err := vm.Pop()
	if err == nil {
		fmt.Print(v)
	}
	return err
}

// cr simply prints a carriage return
func printCR(vm *VM) error {
	fmt.Println()
	return nil
}

// ioWordsInit adds the io-related core words to the VM.
func ioWordsInit(vm *VM) {
	vm.Define("read", Word{read, false})
	vm.Define("skip", Word{skip, false})
	vm.Define("\"", Word{openQuote, true})
	vm.Define("chr", Word{chrFromInt, false})
	vm.Define("ord", Word{ordFromStr, false})
	vm.Define(".s", Word{printStack, false})
	vm.Define(".", Word{printTop, false})
	vm.Define("type", Word{printStr, false})
	vm.Define("cr", Word{printCR, false})
}
