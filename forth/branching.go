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
	// fmt.Printf("Branch to %v\n", vm.ip + 1)
	if vm.ip < -1 || vm.ip >= len(vm.codeseg) {
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

// limit start DO <body> LOOP/+LOOP defines a basic for-style loop.
// It needs to stash away the limit and current index on the R-stack
// prior to the loop proper. Then, at the start of the loop, it needs to
// test whether iteration should continue, or jump to the end
// of the loop:
// >r >r (test loop-body back-facing branch) rdrop rdrop
func opDo(vm *VM) (err error) {
	opSetup := vm.dict["(setupDo)"]
	opTest := vm.dict["(testDo)"]
	vm.codeseg = append(vm.codeseg, opSetup, opTest, 32768)
	vm.Push(len(vm.codeseg) - 1)
	return
}

func opLoop(vm *VM) error {
	return opLoopInternal(vm, true)
}
func opLoopPlus(vm *VM) error {
	return opLoopInternal(vm, false)
}

func opLoopInternal(vm *VM, pullVal bool) (err error) {
	opLoopPlus := vm.dict["(perfLoopPlus)"]
	opRAt := vm.dict["r@"]
	opRDrop := vm.dict["rdrop"]

	var fixUpLoc interface{}
	fixUpLoc, err = vm.Pop()
	if err != nil {
		return
	}

	ful, ok := fixUpLoc.(int)
	if !ok {
		err = ErrBadState
	}

	distToEnd := len(vm.codeseg) + 3 - ful
	distToStart := ful - len(vm.codeseg) - 3

	if pullVal {
		vm.codeseg = append(vm.codeseg, opRAt)
		distToEnd++
		distToStart--
	}
	vm.codeseg[ful] = uint16(distToEnd)
	vm.codeseg = append(vm.codeseg, opLoopPlus,
		opBranch, uint16(distToStart),
		opRDrop, opRDrop, opRDrop)
	return
}

// (perfLoopPlus) ( amt -- )
func performLoopPlus(vm *VM) (err error) {
	rtop := len(vm.Rstack) - 1
	if rtop < 2 {
		return ErrUnderflow
	}
	ridx := vm.Rstack[rtop-2]
	iidx, ok := ridx.(int)

	var amt interface{}
	amt, err = vm.Pop()
	if err != nil {
		return err
	}

	iamt, ok2 := amt.(int)
	if ok && ok2 {
		vm.Rstack[rtop-2] = (iamt + iidx)
	} else {
		err = ErrBadState
	}
	return
}

func setupDo(vm *VM) (err error) {
	err = toR(vm)
	if err != nil {
		return err
	}
	err = toR(vm)
	if err != nil {
		return err
	}
	rtop := len(vm.Rstack) - 1
	rlim, ridx := vm.Rstack[rtop], vm.Rstack[rtop-1]
	limval, ok1 := rlim.(int)
	ival, ok2 := ridx.(int)
	if ok1 && ok2 {
		switch {
		case limval > ival:
			vm.RPush(1)
		case limval < ival:
			vm.RPush(-1)
		default:
			vm.RPush(0)
		}
	} else {
		err = ErrBadState
	}
	return
}

func testDo(vm *VM) (err error) {
	rtop := len(vm.Rstack) - 1
	if rtop < 2 {
		return ErrUnderflow
	}
	rtest, rlim, ridx := vm.Rstack[rtop], vm.Rstack[rtop-1], vm.Rstack[rtop-2]
	testval, ok1 := rtest.(int)
	limval, ok2 := rlim.(int)
	ival, ok3 := ridx.(int)
	// fmt.Printf("rtop: %v  test: %v   limit: %v   idx: %v\n",rtop, testval, limval, ival);
	noloop := true
	if ok1 && ok2 && ok3 {
		switch testval {
		case 0:
			noloop = true
		case 1:
			noloop = ival >= limval
		case -1:
			noloop = ival <= limval
		}
	} else {
		err = ErrBadState
	}
	if noloop {
		err = branchUnconditional(vm)
	} else {
		vm.ip++
	}
	return
}

func getDoI(vm *VM) error {
	rlen := len(vm.Rstack)
	if rlen < 3 {
		return ErrUnderflow
	}
	vm.Push(vm.Rstack[rlen-3])
	return nil
}

func getDoJ(vm *VM) error {
	rlen := len(vm.Rstack)
	if rlen < 6 {
		return ErrUnderflow
	}
	vm.Push(vm.Rstack[rlen-6])
	return nil
}

func branchWordsInit(vm *VM) {
	vm.Define("if", Word{opIf, true})
	vm.Define("else", Word{opElse, true})
	vm.Define("then", Word{opThen, true})
	vm.Define("recur", Word{recur, true})
	vm.Define("do", Word{opDo, true})
	vm.Define("(setupDo)", Word{setupDo, false})
	vm.Define("(testDo)", Word{testDo, false})
	vm.Define("(perfLoopPlus)", Word{performLoopPlus, false})
	vm.Define("loop", Word{opLoop, true})
	vm.Define("+loop", Word{opLoopPlus, true})
	vm.Define("i", Word{getDoI, false})
	vm.Define("j", Word{getDoJ, false})
}
