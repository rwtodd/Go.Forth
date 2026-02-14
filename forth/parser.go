// SPDX-License-Identifier: MIT

package forth

import (
	"fmt"
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
	// run the internal words
	var locals []any
	for {
		idx := vm.codeseg[vm.ip]
		switch idx {
		case opReturn:
			goto Done
		case opCreateLocals:
			vm.ip++
			count := vm.codeseg[vm.ip]
			locals = make([]any, count)
		case opLocalGet:
			vm.ip++
			lidx := vm.codeseg[vm.ip]
			vm.Push(locals[lidx])
		case opLocalSet:
			vm.ip++
			lidx := vm.codeseg[vm.ip]
			val, err := vm.Pop()
			if err != nil {
				return err
			}
			locals[lidx] = val

		default:
			if err := vm.words[idx].Run(vm); err != nil {
				return err
			}
		}
		vm.ip++
	}
Done:

	if len(vm.Rstack) < rstackLen {
		return ErrUnderflowMsg("return stack underflow")
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
	if buf, err = delimitedWSRead(vm.Source, buf); err == nil {
		// convert the buffer to lowercase
		for i, r := range buf {
			buf[i] = unicode.ToLower(r)
		}
	}

	return string(buf), err
}

// decodeLiteral possibly turns a string into a number,
// and maybe other literal forms if I want to do so later.
func decodeLiteral(s string) (interface{}, error) {
	// try to make an integer...
	i, e := strconv.Atoi(s)
	if e == nil {
		return i, e
	}

	// try to make a float...
	f, e := strconv.ParseFloat(s, 64)
	if e == nil {
		return f, e
	}

	// we can't tell what this token is!
	return nil, ErrArgumentMsg(fmt.Sprintf("unknown token <%s>", s))
}

// stopInterpret completes an interpretation and falls back to the compiler
// (assuming one was in play
func stopInterpret(vm *VM) error {
	if vm.Compiling {
		return ErrBadStateMsg("not compiling (stopInterpret called in compiler mode)")
	}
	vm.Compiling = true
	return nil
}

// Interpret sets the compilation state of the VM to false, and
// reads words one at a time...
func interpret(vm *VM) (err error) {
	if !vm.Compiling {
		return ErrBadStateMsg("already interpreting (interpret called in interpret mode)")
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
		if idx, ok := vm.dict[str]; ok {
			err = vm.words[idx].Run(vm)
		} else {
			// if it's not in the dict, put it on the stack as a literal
			var lit interface{}
			if lit, err = decodeLiteral(str); err == nil {
				vm.Push(lit)
			}
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
		return ErrBadStateMsg("not compiling (stopCompile called in interpret mode)")
	}

	// Deal with locals if we have them
	// Since we double-map, len(vm.CompileLocals) is not the count.
	// We need to find the max index.
	numLocals := 0
	if len(vm.CompileLocals) > 0 {
		maxIdx := -1
		for _, v := range vm.CompileLocals {
			if v > maxIdx {
				maxIdx = v
			}
		}
		numLocals = maxIdx + 1
	}

	if numLocals > 0 {
		// Insert opCreateLocals and the count at the beginning of the definition
		// The definition starts at vm.curdef
		// We need to shift everything from vm.curdef onwards by 2 spots

		// 1. Expand the slice
		vm.codeseg = append(vm.codeseg, 0, 0)
		copy(vm.codeseg[vm.curdef+2:], vm.codeseg[vm.curdef:])

		// 2. Insert the instruction
		vm.codeseg[vm.curdef] = opCreateLocals
		vm.codeseg[vm.curdef+1] = uint16(numLocals)

		// 3. Fixup recursion
		// Any branch that pointed to vm.curdef-1 (which is where recur points)
		// now points to vm.curdef+1 because of the shift.
		// However, we want it to point to vm.curdef-1 (the start of the function)
		// essentially, the relative offset needs to be decreased by 2.
		// NOTE: The target calculation is: TargetIP = InstIP + Offset
		// We want NewTargetIP == vm.curdef - 1
		for i := vm.curdef + 2; i < len(vm.codeseg); i++ {
			if vm.codeseg[i] == opBranch || vm.codeseg[i] == opBZR {
				offset := int16(vm.codeseg[i+1])
				targetIP := i + int(offset)
				// Check if it targets what used to be the start (now shifted)
				// The start was vm.curdef.
				// The RECUR logic calculated offset to hit vm.curdef - 1 (the IP before start)
				// Because Run loop increments IP.
				// Wait, let's re-verify RECUR logic.
				// recur: distance := vm.curdef - len(vm.codeseg) - 1
				// IP of branch is len(codeseg).
				// TargetIP = len(codeseg) + distance = curdef - 1.
				// NextIP = TargetIP + 1 = curdef. Correct.

				// So we are looking for TargetIP == vm.curdef - 1.
				if targetIP == vm.curdef+1 { // It shifted by 2
					// We want it to be vm.curdef - 1
					// NewTarget = OldTarget - 2
					// NewOffset = OldOffset - 2
					vm.codeseg[i+1] = uint16(offset - 2)
				}
			}
			// Skip arguments
			if vm.codeseg[i] == opBranch || vm.codeseg[i] == opBZR ||
				vm.codeseg[i] == opLitINT || vm.codeseg[i] == opLitUINT ||
				vm.codeseg[i] == opCreateLocals || vm.codeseg[i] == opLocalGet ||
				vm.codeseg[i] == opLocalSet {
				i++
			}
		}
	}

	vm.Compiling = false
	vm.codeseg = append(vm.codeseg, opReturn) // put a (RET)

	// create a composite word out of the current definition
	cw := CompositeWord{start: vm.curdef}
	vm.Define(vm.curname, Word{Run: cw.Run, Immediate: false})

	// Reset locals map
	vm.CompileLocals = nil
	return nil
}

// compileLocals ({:) parses local variable definitions
func compileLocals(vm *VM) error {
	if !vm.Compiling {
		return ErrBadStateMsg("interpret mode ({: called outside definition)")
	}

	if vm.CompileLocals == nil {
		vm.CompileLocals = make(map[string]int)
	}

	buf := make([]rune, 0, 20)
	var initList []int
	uninitMode := false

	for {
		str, err := nextToken(vm, buf)
		if err != nil {
			return err
		}

		if str == ":}" {
			break
		}
		if str == "|" {
			uninitMode = true
			continue
		}

		if str[len(str)-1] == '!' {
			return ErrArgumentMsg("local variable names cannot end in '!'")
		}

		// define the local
		// Check if it already exists? Assuming no redefinition for now or shadowing ok.
		// Use a simpler approach: finding max index to know next index.
		idx := 0
		if len(vm.CompileLocals) > 0 {
			maxIdx := -1
			for _, v := range vm.CompileLocals {
				if v > maxIdx {
					maxIdx = v
				}
			}
			idx = maxIdx + 1
		}

		vm.CompileLocals[str] = idx
		vm.CompileLocals[str+"!"] = idx

		if !uninitMode {
			initList = append(initList, idx)
		}
	}

	// Compile initialization code in reverse order
	for i := len(initList) - 1; i >= 0; i-- {
		vm.codeseg = append(vm.codeseg, opLocalSet, uint16(initList[i]))
	}

	return nil
}

// compile (':') reads the name of a word to define, and then compiles
// the definition until ';' tells it to stop
func compile(vm *VM) (err error) {
	if vm.Compiling {
		return ErrBadStateMsg("already compiling (compile called in compiler mode)")
	}

	vm.Compiling = true

	buf := make([]rune, 0, 20)

	// STEP 1: read the name
	var str string
	str, err = nextToken(vm, buf)
	vm.curname = str                        // remember the name of the definition
	vm.curdef = len(vm.codeseg)             // remember the start of the definition
	vm.CompileLocals = make(map[string]int) // reset locals

	for (err == nil) && vm.Compiling {
		str, err = nextToken(vm, buf)
		if err != nil {
			if err == io.EOF {
				err = nil
			}
			return
		}

		// Check locals first (shadowing)
		if idx, ok := vm.CompileLocals[str]; ok {
			if str[len(str)-1] == '!' {
				vm.codeseg = append(vm.codeseg, opLocalSet, uint16(idx))
			} else {
				vm.codeseg = append(vm.codeseg, opLocalGet, uint16(idx))
			}
			continue
		}

		// lookup the string in the dictionary
		if idx, ok := vm.dict[str]; ok {
			// compile in the word unless it's immediate
			if vm.words[idx].Immediate {
				err = vm.words[idx].Run(vm)
			} else {
				vm.codeseg = append(vm.codeseg, idx)
			}
		} else {
			// if it's not in the dict, compile it in as a literal
			var lit interface{}
			if lit, err = decodeLiteral(str); err == nil {
				compileLiteral(vm, lit)
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
		return ErrBadStateMsg("interpret mode (literal called outside definition)")
	}
	var value any
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
		return ErrArgumentMsg("compile, expects valid word index")
	}

	vm.codeseg = append(vm.codeseg, uint16(num))
	return nil
}

// postpone creates code that compiles code into the caller.  For
// immediates, it creates code that calls code in the caller.
func postpone(vm *VM) error {
	if !vm.Compiling {
		return ErrBadStateMsg("interpret mode (postpone called outside definition)")
	}

	buf := make([]rune, 0, 20)

	// STEP 1: read the name and look it up
	str, err := nextToken(vm, buf)
	if err != nil {
		return err
	}

	opcode, ok := vm.dict[str]
	if !ok {
		return ErrArgumentMsg(fmt.Sprintf("POSTPONE: no word <%s>", str))
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

// tick (') reads the next word and pushes its execution token (index)
// to the stack.
func tick(vm *VM) error {
	buf := make([]rune, 0, 20)
	str, err := nextToken(vm, buf)
	if err != nil {
		return err
	}

	idx, ok := vm.dict[str]
	if !ok {
		return ErrArgumentMsg(fmt.Sprintf("' : word <%s> not found", str))
	}

	vm.Push(int(idx))
	return nil
}

// bracketTick ([']) is an immediate word that reads the next word and
// compiles its execution token (index) as a literal.
func bracketTick(vm *VM) error {
	if !vm.Compiling {
		return ErrBadStateMsg("interpret mode (['] called outside definition)")
	}

	buf := make([]rune, 0, 20)
	str, err := nextToken(vm, buf)
	if err != nil {
		return err
	}

	idx, ok := vm.dict[str]
	if !ok {
		return ErrArgumentMsg(fmt.Sprintf("['] : word <%s> not found", str))
	}

	vm.codeseg = append(vm.codeseg, opLitUINT, idx)
	return nil
}

func parseWordsInit(vm *VM) {
	vm.Define("\\", Word{nlComment, true})
	vm.Define("(", Word{parenComment, true})
	vm.Define("[", Word{interpret, true})
	vm.Define("]", Word{stopInterpret, false})
	vm.Define(":", Word{compile, false})
	vm.Define(";", Word{stopCompile, true})
	vm.Define("literal", Word{literal, true})
	vm.Define("postpone", Word{postpone, true})
	vm.Define("immediate", Word{makeImmediate, false})
	vm.Define("'", Word{tick, false})
	vm.Define("[']", Word{bracketTick, true})
	vm.Define("{:", Word{compileLocals, true})
}
