package forth

import (
	"fmt"
	"io"
	"strconv"
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
		if idx == opReturn {
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

// parenComment '(' skips until the closing paren.
// : ( ')' skip ; immediate
func parenComment(vm *VM) error {
	vm.Push(int(')'))
	return skip(vm)
}

// nlComment '\' skips until the next newline
// : \ '\n' skip ; immediate
func nlComment(vm *VM) error {
	vm.Push(int('\n'))
	return skip(vm)
}

func nextToken(vm *VM, buf []rune) (string, error) {
	ch, err := eatWhitespace(vm.Source)
	if err != nil {
		return "", err
	}

	buf = append(buf, ch)
	buf, err = delimitedWSRead(vm.Source, buf)
	return string(buf), err
}

// decodeLiteral possibly turns a string into a number,
// and maybe other literal forms if I want to do so later.
func decodeLiteral(s string) interface{} {
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
			vm.Push(decodeLiteral(str))
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
	vm.codeseg = append(vm.codeseg, opReturn) // put a (RET)

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
			compileLiteral(vm, decodeLiteral(str))
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

// (litINT) reads the next 16-bits from the codeseg and pushes that number on the stack as an int
// The 16 bits are considered signed
func litINT(vm *VM) error {
	vm.ip++
	num := int16(vm.codeseg[vm.ip])
	vm.Stack = append(vm.Stack, int(num))
	return nil
}

// (litUINT) reads the next 16-bits from the codeseg and pushes that number on the stack as an int
// The 16 bits are considered unsigned
func litUINT(vm *VM) error {
	vm.ip++
	num := vm.codeseg[vm.ip]
	vm.Stack = append(vm.Stack, int(num))
	return nil
}

// compileLiteral is a helper function to put a literal into the compiled
// codestream. This will be the one place we'll have to add code to have more
// special types that don't just go to CreatePusher()
func compileLiteral(vm *VM, value interface{}) {
	switch num := value.(type) {
	case int:
		switch {
		case (num >= -32768) && (num < 32768):
			vm.codeseg = append(vm.codeseg, opLitINT, uint16(num))
		case (num >= 0) && (num < 65536):
			vm.codeseg = append(vm.codeseg, opLitUINT, uint16(num))
		default:
			vm.codeseg = append(vm.codeseg, vm.CreatePusher(num))
		}
	default:
		vm.codeseg = append(vm.codeseg, vm.CreatePusher(value))
	}
}

// literal is an immediate word that reads an int from the stack and compiles it into the codestream
// if possible, and uses a pusher if necessary.
func literal(vm *VM) (err error) {
	if !vm.Compiling {
		return ErrBadState
	}
	var value interface{}
	value, err = vm.Pop()
	if err != nil {
		return
	}
	compileLiteral(vm, value)
	return
}

// compileComma takes the top of the stack and puts that opcode literally
// into the code sequence.
func compileComma(vm *VM) error {
	value, err := vm.Pop()
	if err != nil {
		return err
	}

	num, ok := value.(int)
	if !ok || (num < 0) || (num > len(vm.words)) {
		return ErrArgument
	}

	vm.codeseg = append(vm.codeseg, uint16(num))
	return nil
}

// postpone creates code that compiles code into the caller.  For
// immediates, it creates code that calls code in the caller.
func postpone(vm *VM) error {
	if !vm.Compiling {
		return ErrBadState
	}

	buf := make([]rune, 0, 20)

	// STEP 1: read the name and look it up
	str, err := nextToken(vm, buf)
	if err != nil {
		return err
	}

	opcode, ok := vm.dict[str]
	if !ok {
		return fmt.Errorf("POSTPONE: no word <%s>", str)
	}

	// STEP 2: generate the code
	if vm.words[opcode].Immediate {
		// just call the immediate in the caller when it runs
		vm.codeseg = append(vm.codeseg, opcode)
	} else {
		// need to compile a sequence to compile the opcode into the caller's caller
		vm.codeseg = append(vm.codeseg, opLitUINT, opcode, opCompileComma)
	}

	return nil
}
