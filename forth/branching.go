package forth

// (branch) branches unconditionally.
// The int16 relative move is the next word
// in the codeseg.  N.B. because of the way the interpreter
// runs, we actually compile code to jump to the
// target IP _minus_ _one_.  N.B. the jump amount is relative
// to the BRANCH instruction location, NOT the offset number's
// location.
func branchUnconditional(vm *VM) (err error) {
	num := int16(vm.codeseg[vm.ip+1])
	vm.ip += int(num)
	if vm.ip < -1 {
		err = ErrBadState
	}
	return
}

// (bzr) branches when the top of stack is zero. Otherwise
// it is a NOP.  The int16 relative move is the next word
// in the codeseg.  N.B. because of the way the interpreter
// runs, we actually compile code to jump to the
// target IP _minus_ _one_.
func branchZero(vm *VM) (err error) {
	var tos interface{}
	tos, err = vm.Pop()
	bval, ok := tos.(int)
	if ok && bval == 0 {
		err = branchUnconditional(vm)
	} else {
		vm.ip++
	}
	return
}

// IF is an immediate word that stores a fixup address
// on the stack for ELSE / THEN to find, and stores
// a (bzr) with a dummy branch amount in the code stream.
func opIf(vm *VM) (err error) {
	vm.Push(len(vm.codeseg) + 1)
	vm.codeseg = append(vm.codeseg, opBZR, 32768)
	return
}

// THEN takes a fixup address from the stack, and
// inserts the right amount to jump over the IF (or ELSE)
// block. No new code is added to the codestream.
func opThen(vm *VM) (err error) {
	var tos interface{}
	tos, err = vm.Pop()
	fixupLoc, ok := tos.(int)
	if ok {
		// 5    6     7       8   // fixupLoc == 6
		// BZR  FFFF  PRINT       // Right answer == 2  (8 - 6)
		vm.codeseg[fixupLoc] = uint16(len(vm.codeseg) - fixupLoc)
	} else {
		if err == nil {
			err = ErrBadState
		}
	}
	return
}

// ELSE needs to issue a jump over the else-stuff, and then
// use opThen to fixup the IF to jump into the else-stuff.
// Finally, it needs to leave a fixup location on the stack
// for the final THEN.
func opElse(vm *VM) (err error) {
	fupLoc := len(vm.codeseg) + 1
	vm.codeseg = append(vm.codeseg, opBranch, 32768)
	err = opThen(vm)
	vm.Push(fupLoc)
	return
}

// RECUR just jumps to the start of the current function
func recur(vm *VM) (err error) {
	// 5     6      7      8      // Start = 5  len(code) == 8
	// PRINT PRINT  PRINT  RECUR  // Right answer ==  -4 (5 - 8 - 1)
	distance := vm.curdef - len(vm.codeseg) - 1
	vm.codeseg = append(vm.codeseg, opBranch, uint16(distance))
    return
}

func branchWordsInit(vm *VM) {
	vm.Define("if", Word{opIf, true})
	vm.Define("else", Word{opElse, true})
	vm.Define("then", Word{opThen, true})
	vm.Define("recur", Word{recur, true})
}
