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
func (c CompositeWord) Run(vm *VM) error {
	return vm.RunAt(c.start)
}

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

// stopInterpret complete an interpretation and falls back to the compiler
// (assuming one was in play
func stopInterpret(vm *VM) error {
	if vm.CurrentCompCtx() == nil {
		return ErrBadStateMsg("not compiling (stopInterpret called in compiler mode)")
	}
	return nil
}

func interpret(vm *VM) (err error) {
	// If we are called recursively from compile (via [ ... ]), that's fine.
	// We just loop until we hit ] or EOF.

	// In original code:
	// if !vm.Compiling { return ErrBadStateMsg(...) }
	// vm.Compiling = false

	// Here, we just run the loop. The "mode switch" is implicit by running this function.

	buf := make([]rune, 0, 20)

	for err == nil {
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
			if vm.words[idx].Name == "]]" {
				// Return to compilation if we are compiling
				if vm.CurrentCompCtx() != nil {
					return nil
				}
				// If not compiling, ]] is an error or no-op?
				// In standard forth ] enters compilation.
				// But our ]] says "stopInterpret".
				// Let's assume for now it just breaks the loop if we were inside [[ ... ]]
				// If we are top level, it might be an error.
				// For now, let's return and let caller decide.
				return nil
			}
			err = vm.words[idx].Run(vm)
		} else {
			// if it's not in the dict, put it on the stack as a literal
			var lit any
			if lit, err = decodeLiteral(str); err == nil {
				vm.Push(lit)
			}
		}
	}
	return
}

// func makeImmediate ('immediate') makes the last defined word immediate
func makeImmediate(vm *VM) error {
	ctx := vm.CurrentCompCtx()
	if ctx == nil {
		return ErrBadStateMsg("immediate called without compiling a word!")
	}
	vm.words[ctx.WordIdx].Immediate = true
	return nil
}

// resolveLocal searches for a local variable in the compilation stack
func resolveLocal(vm *VM, name string) (depth int, idx int, found bool) {
	depth = 0
	// iterate backwards through CompStack
	for i := len(vm.CompStack) - 1; i >= 0; i-- {
		ctx := &vm.CompStack[i]
		if idx, ok := ctx.CompileLocals[name]; ok {
			return depth, idx, true
		}
		if len(ctx.CompileLocals) > 0 {
			depth++
		}
	}
	return 0, 0, false
}

// stopCompile (';') terminates a compilation
func stopCompile(vm *VM) error {
	ctx := vm.CurrentCompCtx()
	if ctx == nil {
		return ErrBadStateMsg("not compiling (stopCompile called in interpret mode)")
	}

	// Deal with locals if we have them
	// Since we double-map, len(vm.CompileLocals) is actually 2x count.
	numLocals := 0
	if len(ctx.CompileLocals) > 0 {
		numLocals = len(ctx.CompileLocals) / 2
	}

	if numLocals > 0 {
		// Insert opEnterScope and the count at the beginning of the definition
		// The definition starts at ctx.StartIP
		// We need to shift everything from ctx.StartIP onwards by 2 spots

		// 1. Expand the slice
		vm.codeseg = append(vm.codeseg, 0, 0)
		copy(vm.codeseg[ctx.StartIP+2:], vm.codeseg[ctx.StartIP:])

		// 2. Insert the instruction
		vm.codeseg[ctx.StartIP] = opEnterScope
		vm.codeseg[ctx.StartIP+1] = uint16(numLocals)

		// AND emit opExitScope before opReturn
		vm.codeseg = append(vm.codeseg, opExitScope)
	}

	// We pop the context BEFORE finishing, or effectively we are done compiling THIS word.
	vm.codeseg = append(vm.codeseg, opReturn) // put a (RET)

	// create a composite word out of the current definition and finish setting it up
	cw := CompositeWord{start: ctx.StartIP}
	vm.words[ctx.WordIdx].Run = cw.Run
	vm.AddToDict(uint16(ctx.WordIdx))

	vm.PopCompCtx()
	return nil
}

// compileLocals ((|) parses local variable definitions
func compileLocals(vm *VM) error {
	ctx := vm.CurrentCompCtx()
	if ctx == nil {
		return ErrBadStateMsg("interpret mode ((| called outside definition)")
	}

	if ctx.CompileLocals == nil {
		ctx.CompileLocals = make(map[string]int)
	}

	buf := make([]rune, 0, 20)
	var initList []int
	uninitMode := false

	for {
		str, err := nextToken(vm, buf)
		if err != nil {
			return err
		}

		if !uninitMode {
			if str == "|)" {
				break
			} else if str == "|" {
				uninitMode = true
				continue
			} else if str == "--" || str == ")" {
				return ErrArgumentMsg("invalid token in locals list: " + str)
			}
		} else {
			if str == ")" {
				break
			} else if str == "--" {
				// skip comment tokens until )
				for {
					str2, err2 := nextToken(vm, buf)
					if err2 != nil {
						return err2
					}
					if str2 == ")" {
						break
					}
				}
				break
			} else if str == "|)" || str == "|" {
				return ErrArgumentMsg("invalid token in locals list: " + str)
			}
		}

		if str[len(str)-1] == '!' {
			return ErrArgumentMsg("local variable names cannot end in '!'")
		}

		// define the local
		idx := 0
		if len(ctx.CompileLocals) > 0 {
			idx = len(ctx.CompileLocals) / 2
		}

		if _, exists := ctx.CompileLocals[str]; exists {
			return ErrArgumentMsg("duplicate local variable: " + str)
		}
		ctx.CompileLocals[str] = idx
		ctx.CompileLocals[str+"!"] = idx

		if !uninitMode {
			initList = append(initList, idx)
		}
	}

	// Compile initialization code in reverse order
	for i := len(initList) - 1; i >= 0; i-- {
		vm.codeseg = append(vm.codeseg, opLocalSet, 0, uint16(initList[i]))
	}

	return nil
}

// compileLoop manages the compilation loop
func compileLoop(vm *VM) error {
	buf := make([]rune, 0, 20)
	startDepth := len(vm.CompStack)

	for {
		// Check if we finished the current context (moved up the stack)
		if len(vm.CompStack) < startDepth {
			break
		}

		str, err := nextToken(vm, buf)
		if err != nil {
			if err == io.EOF {
				return nil
			}
			return err
		}

		// Check locals first (shadowing) with depth resolution
		if depth, idx, ok := resolveLocal(vm, str); ok {
			if str[len(str)-1] == '!' {
				vm.codeseg = append(vm.codeseg, opLocalSet, uint16(depth), uint16(idx))
			} else {
				vm.codeseg = append(vm.codeseg, opLocalGet, uint16(depth), uint16(idx))
			}
			continue
		}

		// lookup the string in the dictionary
		if idx, ok := vm.dict[str]; ok {
			if vm.words[idx].Immediate {
				err = vm.words[idx].Run(vm)
				if err != nil {
					return err
				}
				// The immediate word might have popped the context (e.g. ;)
			} else {
				vm.codeseg = append(vm.codeseg, idx)
			}
		} else {
			// if it's not in the dict, compile it in as a literal
			var lit any
			if lit, err = decodeLiteral(str); err == nil {
				compileLiteral(vm, lit)
			} else {
				return err
			}
		}
	}
	return nil
}

// compile (':') reads the name of a word to define, and then compiles
// the definition until ';' tells it to stop
func compile(vm *VM) (err error) {
	if vm.CurrentCompCtx() != nil {
		return ErrBadStateMsg("already compiling (compile called in compiler mode)")
	}

	buf := make([]rune, 0, 20)

	// STEP 1: read the name and reserve a slot in vm.words
	var str string
	str, err = nextToken(vm, buf)
	if err != nil {
		return err
	}

	// Reserve a slot for this word
	wordIdx := len(vm.words)
	vm.words = append(vm.words, Word{Name: str})

	startIP := len(vm.codeseg)
	vm.PushCompCtx(startIP, wordIdx)

	return compileLoop(vm)
}

// quotationStart ( [: ) begins a quotation
func quotationStart(vm *VM) error {
	if vm.CurrentCompCtx() != nil {
		// Nested compilation
		// skip over the body
		vm.codeseg = append(vm.codeseg, opBranch, 0)
		jumpIdx := len(vm.codeseg) - 1

		// Push the jump index to the stack so ;] can resolve it
		vm.Push(jumpIdx) // Stack: jumpIdx

		// Start new context
		parentIdx := vm.CurrentCompCtx().WordIdx
		vm.PushCompCtx(len(vm.codeseg), parentIdx)
		vm.CurrentCompCtx().IsClosure = true

		// Capture stack depth to clean up any partial compilation artifacts (like IF/ELSE contexts) on error
		stackDepth := len(vm.Stack)

		if err := compileLoop(vm); err != nil {
			vm.PopCompCtx()
			// Restore stack to remove any compilation artifacts
			if len(vm.Stack) > stackDepth {
				vm.Stack = vm.Stack[:stackDepth]
			}
			vm.Pop() // pop the jumpIdx
			return err
		}
	} else {
		// Top-Level compilation (intepreted quotation)
		// We need to act like interpreted execution.
		// We start compiling into codeseg (which is persistent).
		// We use a dummy WordIdx? Or -1?
		vm.PushCompCtx(len(vm.codeseg), -1)
		vm.CurrentCompCtx().IsClosure = true
		// And we need to enter the compile loop!
		return compileLoop(vm)
	}
	return nil
}

// quotationEnd ends a quotation definition
func quotationEnd(vm *VM) error {
	ctx := vm.CurrentCompCtx()
	if !ctx.IsClosure {
		// Should be impossible given how quotationStart works
		return ErrBadStateMsg("quotationEnd called on non-closure")
	}

	// Before we finish, we might need to insert locals handling code
	// at the BEGINNING of the closure.
	hasLocals := len(ctx.CompileLocals) > 0
	shift := 0
	if hasLocals {
		// We need to inject opEnterScope <count> at ctx.StartIP
		// This will shift all subsequent code by 2 words.
		shift = 2

		// Find max local index
		maxIdx := -1
		for _, idx := range ctx.CompileLocals {
			if idx > maxIdx {
				maxIdx = idx
			}
		}
		numLocals := maxIdx + 1

		// Inject code!
		// We need to shift vm.codeseg from StartIP onwards.
		vm.codeseg = append(vm.codeseg, 0, 0)
		copy(vm.codeseg[ctx.StartIP+2:], vm.codeseg[ctx.StartIP:])
		vm.codeseg[ctx.StartIP] = opEnterScope
		vm.codeseg[ctx.StartIP+1] = uint16(numLocals)

		// Also append opExitScope at the END
		vm.codeseg = append(vm.codeseg, opExitScope)
	}

	// FIXUPS!
	// Now we know if we have locals, and where the code starts.
	// We need to fixup all RECUR and TAIL-CALL placeholders.

	// Fixup RECUR
	for _, ip := range ctx.RecurFixups {
		// ip is the location of the placeholder op (opRecurClosure).
		// If we shifted code, ip needs to be shifted too?
		// YES, because the fixups were recorded RELATIVE TO START OF SEGMENT.
		// Since we inserted at StartIP, everything after it shifted.
		// All fixups are inside the closure, so they are after StartIP.
		realIP := ip + shift

		targetIP := ctx.StartIP // New start (opEnterScope if locals, or first instr)

		// Calculate offset: Target - (AddressOfArgument)
		// opRecurClosure/opCallOffset implementations do:
		// vm.ip++ (points to argument)
		// target = vm.ip + offset
		// So: offset = target - vm.ip
		// At runtime, vm.ip will be realIP + 1.
		offset := targetIP - (realIP + 1)

		// Patch opcode and offset
		vm.codeseg[realIP+1] = uint16(offset)

		if hasLocals {
			vm.codeseg[realIP] = opRecurClosure
		} else {
			vm.codeseg[realIP] = opCallOffset
		}
	}

	// Fixup TAIL-CALL
	for _, ip := range ctx.TailCallFixups {
		realIP := ip + shift
		targetIP := ctx.StartIP

		if hasLocals {
			// Optimization: Skip opEnterScope (index 0 and 1). Jump to index 2.
			// This reuses the current scope frame.
			targetIP += 2
		}

		// Tail call is always opBranch
		vm.codeseg[realIP] = opBranch

		// Calculate offset involves the -1 because branch is relative to instruction?
		// See branchUnconditional: vm.ip += int(num) -> next instruction is ip+1+num.
		// We want next instruction to be TargetIP.
		// TargetIP = (realIP + 1) + 1 + num  (Wait, ip is pointing to opcode?)
		// When executing branch: ip points to opcode.
		// num = vm.codeseg[ip+1]
		// vm.ip += num.
		// Next loop: vm.ip++ (so effective next is ip + num + 1)
		// We want effective next to be TargetIP.
		// TargetIP = ip + num + 1
		// num = TargetIP - ip - 1

		offset := targetIP - realIP - 1
		vm.codeseg[realIP+1] = uint16(offset)
	}

	vm.codeseg = append(vm.codeseg, opReturn)
	vm.PopCompCtx()

	// Post-processing
	if vm.CurrentCompCtx() != nil {
		// We were nested.
		// 1. Resolve the branch
		val, err := vm.Pop() // pop jumpIdx
		if err != nil {
			return err
		}
		jumpIdx := val.(int)

		offset := len(vm.codeseg) - jumpIdx
		vm.codeseg[jumpIdx] = uint16(offset) // Patch opBranch

		// 2. Emit opPushClosure
		// We use a relative offset offsets relative to current instruction.
		// opPushClosure (here) + offset = StartIP.
		// vm.ip points to offset argument.
		// StartIP = (len + 1) + offset.
		// offset = StartIP - (len + 1).
		offset = ctx.StartIP - (len(vm.codeseg) + 1)
		vm.codeseg = append(vm.codeseg, opPushClosure, uint16(offset))
	} else {
		// We were top-level.
		// We just finished compiling the closure.
		// We need to push the Closure object to the stack.
		vm.Push(Closure{StartIP: ctx.StartIP, Env: vm.HeadScope})
	}

	return nil
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
	if vm.CurrentCompCtx() == nil {
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
	if vm.CurrentCompCtx() == nil {
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
	if vm.CurrentCompCtx() == nil {
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
	vm.Define(Word{Name: "\\", Run: nlComment, Immediate: true})
	vm.Define(Word{Name: "(", Run: parenComment, Immediate: true})
	vm.Define(Word{Name: "[[", Run: interpret, Immediate: true})
	vm.Define(Word{Name: "]]", Run: stopInterpret, Immediate: false})
	vm.Define(Word{Name: ":", Run: compile, Immediate: false})
	vm.Define(Word{Name: ";", Run: stopCompile, Immediate: true})
	vm.Define(Word{Name: "literal", Run: literal, Immediate: true})
	vm.Define(Word{Name: "postpone", Run: postpone, Immediate: true})
	vm.Define(Word{Name: "immediate", Run: makeImmediate, Immediate: false})
	vm.Define(Word{Name: "'", Run: tick, Immediate: false})
	vm.Define(Word{Name: "[']", Run: bracketTick, Immediate: true})
	vm.Define(Word{Name: "(|", Run: compileLocals, Immediate: true})
	vm.Define(Word{Name: "[", Run: quotationStart, Immediate: true})
	vm.Define(Word{Name: "]", Run: quotationEnd, Immediate: true})
}
