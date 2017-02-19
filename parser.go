package forth

import (
	"bufio"
	"io"
	"strconv"
	"unicode"
)

// eatWhitespace skips whitespace and returns the next
// non-whitespace char
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

// readTilWhitespace reads from r until whitespace is found,
// filling the provided buf as it goes.
func readTilWhitespace(r *bufio.Reader, buf []rune) ([]rune, error) {
	var (
		ch  rune
		err error
	)

	for err == nil {
		ch, _, err = r.ReadRune()
		if (err != nil) || unicode.IsSpace(ch) {
			break
		}
		buf = append(buf, ch)
	}

	// EOF isn't really a problem for this function
	if err == io.EOF {
		err = nil
	}
	return buf, err
}

func nextToken(vm *VM, buf []rune) (string, error) {
	ch, err := eatWhitespace(vm.Source)
	if err != nil {
		return "", err
	}

	switch ch {
	case '"':
		str, err := vm.Source.ReadSlice('"')
		return string(str[:len(str)-1]), err
	default:
		buf = append(buf, ch)
		buf, err = readTilWhitespace(vm.Source, buf)
	}
	return string(buf), err
}

// processLiteral possibly turns a string into a number,
// and maybe other literal forms if I want to do so later.
func processLiteral(s string) interface{} {
	// try to make an integer...
	i, e := strconv.Atoi(s)
	if e == nil {
		return i
	}

	// try to make a float...
	f, e := strconv.ParseFloat(s, 64)
	if e == nil {
		return f
	}

	// just return the string...
	return s
}


// stopInterpret completes an interpretation and falls back to the compiler
// (assuming one was in play
func stopInterpret(vm *VM) error {
	if vm.Compiling {
		return ErrBadState
	}
	vm.Compiling = true
	return nil
}

// Interpret sets the compilation state of the VM to false, and
// reads words one at a time...
func interpret(vm *VM) (err error) {
	if !vm.Compiling {
		return ErrBadState
	}

	vm.Compiling = false
	buf := make([]rune, 0, 20)

	for (err == nil) && !vm.Compiling {
		var str string
		str, err = nextToken(vm, buf)
		if err != nil {
			if err == io.EOF {
				err = nil
			}
			return
		}

		// lookup the string in the dictionary
		idx, ok := vm.dict[str]

		// if it's not there, put it on the stack as a literal
		if !ok {
			vm.Push(processLiteral(str))
		} else {
			err = vm.words[idx].Run(vm)
		}
	}
	return
}

func compile(vm *VM) (err error) {
 	if vm.Compiling {
		return ErrBadState
	}

	vm.Compiling = true
		
	buf := make([]rune, 0, 20)

	for (err == nil) && !vm.Compiling {
		var str string
		str, err = nextToken(vm, buf)
		if err != nil {
			if err == io.EOF {
				err = nil
			}
			return
		}

		// lookup the string in the dictionary
		idx, ok := vm.dict[str]

		// if it's not there, put it on the stack as a literal
		if !ok {
			vm.Push(processLiteral(str))
		} else {
			err = vm.words[idx].Run(vm)
		}
	}
	return

}
