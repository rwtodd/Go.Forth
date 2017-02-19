package forth

import (
	"bufio"
	"io"
	"strconv"
	"unicode"
)

// CompositeWord represents a word made up of opcodes for other defined words
type CompositeWord struct {
	start int
}

// Run on a composite word:
// don't use the typical return stack... use
// Go's stack instead... this can make the RStack
// an auto-cleaned scratch space, which doesn't have
// to remain balanced like a typical FORTH.
// The only downside is you can't play return-address games
// to force double exits or delayed tail calls.  But, from
// what I've seen on c.l.f, that kind of behavior doesn't
// work on all FORTHS anyway.
func (c CompositeWord) Run(vm *VM) error {
	// setup the composite word
	rstackLen := len(vm.Rstack)
	oldIP := vm.ip
	vm.ip = c.start

	// run the internal words
	for {
		idx := vm.codeseg[vm.ip]
		if idx == 0 {
			break
		}
		vm.words[idx].Run(vm)
		vm.ip++
	}

	if len(vm.Rstack) < rstackLen {
		return ErrRStackUnderflow
	}

	// clean up the rstack and exit
	vm.Rstack = vm.Rstack[:rstackLen]
	vm.ip = oldIP
	return nil
}

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

// func makeImmediate ('immediate') makes the last defined word immediate
func makeImmediate(vm *VM) error {
	vm.words[len(vm.words)-1].Immediate = true
	return nil
}

// stopCompile (';') terminates a compilation
func stopCompile(vm *VM) error {
	if !vm.Compiling {
		return ErrBadState
	}
	vm.Compiling = false
	vm.codeseg = append(vm.codeseg, 0) // put a (RET)

	// create a composite word out of the current definition
	cw := CompositeWord{start: vm.curdef}
	vm.Define(vm.curname, Word{Run: cw.Run, Immediate: false})
	return nil
}

// compile (':') reads the name of a word to define, and then compiles
// the definition until ';' tells it to stop
func compile(vm *VM) (err error) {
	if vm.Compiling {
		return ErrBadState
	}

	vm.Compiling = true

	buf := make([]rune, 0, 20)

	// STEP 1: read the name
	var str string
	str, err = nextToken(vm, buf)
	vm.curname = str            // remember the name of the definition
	vm.curdef = len(vm.codeseg) // remember the start of the definition

	for (err == nil) && vm.Compiling {
		str, err = nextToken(vm, buf)
		if err != nil {
			if err == io.EOF {
				err = nil
			}
			return
		}

		// lookup the string in the dictionary
		idx, ok := vm.dict[str]

		if !ok {
			// if it's not there, compile it in as a literal
			vm.codeseg = append(vm.codeseg, vm.CreatePusher(processLiteral(str)))
		} else {
			// otherwise, compile in the word unless it's immediate
			if vm.words[idx].Immediate {
				err = vm.words[idx].Run(vm)
			} else {
				vm.codeseg = append(vm.codeseg, idx)
			}
		}
	}
	return
}
