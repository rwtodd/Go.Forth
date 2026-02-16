// SPDX-License-Identifier: MIT

package forth

// (branch) branches unconditionally.
// The int16 relative move is the next word
// in the codeseg.  N.B. because of the way the interpreter
// runs, we actually compile code to jump to the
// target IP _minus_ _one_.  N.B. the jump amount is relative
// to the BRANCH instruction location, NOT the offset number's
// target IP _minus_ _one_.  N.B. the jump amount is relative
// to the BRANCH instruction location, NOT the offset number's
// location.

type LoopType int

const (
	LoopDo LoopType = iota
	LoopBegin
)

type LoopCtx struct {
	StartIP   int
	BreakOps  []int
	Type      LoopType
	DoFixupIP int // for DO loops, the initial jump
}

type ConditionType int

const (
	CondIf ConditionType = iota
	CondElse
)

type ConditionCtx struct {
	FixupIP int
	Type    ConditionType
}

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
	vm.Push(&ConditionCtx{FixupIP: len(vm.codeseg) + 1, Type: CondIf})
	vm.codeseg = append(vm.codeseg, opBZR, 32768)
	return
}

// THEN takes a fixup address from the stack, and
// inserts the right amount to jump over the IF (or ELSE)
// block. No new code is added to the codestream.
func opThen(vm *VM) (err error) {
	var tos any
	tos, err = vm.Pop()
	if err != nil {
		return err
	}
	ctx, ok := tos.(*ConditionCtx)
	if ok {
		// 5    6     7       8   // fixupLoc == 6
		// BZR  FFFF  PRINT       // Right answer == 2  (8 - 6)
		vm.codeseg[ctx.FixupIP] = uint16(len(vm.codeseg) - ctx.FixupIP)
	} else {
		err = ErrBadState
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
	if err != nil {
		return err
	}
	vm.Push(&ConditionCtx{FixupIP: fupLoc, Type: CondElse})
	return
}

// RECUR just jumps to the start of the current function
func recur(vm *VM) (err error) {
	ctx := vm.CurrentCompCtx()
	if ctx == nil {
		return ErrBadStateMsg("recur used outside of definition")
	}

	if ctx.IsClosure {
		// Defer the decision until end of quotation.
		// Emit placeholder op and offset.
		// We use opRecurClosure as placeholder, but it might change to opCallOffset.
		ctx.RecurFixups = append(ctx.RecurFixups, len(vm.codeseg))
		vm.codeseg = append(vm.codeseg, opRecurClosure, 0)
	} else {
		vm.codeseg = append(vm.codeseg, uint16(ctx.WordIdx))
	}
	return
}

// TAILCALL jumps to the start of the current function, but
// bypassing the local variable creation.
func tailCall(vm *VM) (err error) {
	ctx := vm.CurrentCompCtx()
	if ctx == nil {
		return ErrBadStateMsg("tail-call used outside of definition")
	}

	if ctx.IsClosure {
		// Defer decision.
		// We emit a placeholder BRANCH.
		// We DO NOT emit opExitScope yet, because if we optimize for locals, we won't need it.
		// If we do need it (no locals optimization?), we will have to handle that?
		// Wait, if no locals, we just branch to start. No exit scope needed (scope stays same).

		// If there ARE locals, but we use TCO (recurse), we jump to AFTER enterScope.
		// So we reuse the scope. No exit scope needed.

		// So.. in ALL closure cases, we don't need explicit opExitScope for tail call!
		// Case 1: No locals. Scope is HeadScope. Jump to start. Scope remains HeadScope. Correct.
		// Case 2: Locals. TCO. Scope is LocalScope. Jump to Body. Scope remains LocalScope. Variables overwritten. Correct.

		// Is there any case where we needed opExitScope?
		// Only if we wanted to destroy the current scope and create a NEW one.
		// But TCO implies reusing the frame if possible.
		// Since Forth locals are just an array in a scope, reusing is fine.

		ctx.TailCallFixups = append(ctx.TailCallFixups, len(vm.codeseg))
		vm.codeseg = append(vm.codeseg, opBranch, 0)
	} else {
		// Standard word tail call
		offset := ctx.StartIP - 1 - len(vm.codeseg)
		vm.codeseg = append(vm.codeseg, opBranch, uint16(offset))
	}
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
	// We sort of assume start IP is right after the jump, but
	// for DO loops specifically, the backwards jump goes to the
	// FIXUP location (addr of 32768), not the instruction after it.
	// See calculate below: distToStart := ful - ...
	vm.Push(&LoopCtx{
		DoFixupIP: len(vm.codeseg) - 1,
		Type:      LoopDo,
		StartIP:   len(vm.codeseg), // Not really used for DO, but good to have
	})
	return
}

// opLoop ends a DO loop, incrementing the index by 1
func opLoop(vm *VM) error {
	return opLoopInternal(vm, true)
}

// opLoopPlus ends a DO loop, incrementing the index by the value on the stack
func opLoopPlus(vm *VM) error {
	return opLoopInternal(vm, false)
}

func opLoopInternal(vm *VM, pullVal bool) (err error) {
	opLoopPlus := vm.dict["(perfLoopPlus)"]
	opRAt := vm.dict["r@"]
	opRDrop := vm.dict["rdrop"]

	var loopCtx any
	loopCtx, err = vm.Pop()
	if err != nil {
		return
	}

	ctx, ok := loopCtx.(*LoopCtx)
	if !ok {
		return ErrBadState
	}

	ful := ctx.DoFixupIP
	distToEnd := len(vm.codeseg) + 3 - ful
	distToStart := ful - len(vm.codeseg) - 3

	if pullVal {
		vm.codeseg = append(vm.codeseg, opRAt)
		distToEnd++
		distToStart--
	}
	vm.codeseg[ful] = uint16(distToEnd)

	// Patch LEAVEs
	for _, breakIP := range ctx.BreakOps {
		// Target is distToEnd + whatever needed to reach end of loop from breakIP
		// Actually target is simply relative offset from breakIP to HERE (start of cleanups)
		// cleanups are about to be appended.
		// Current IP = len(vm.codeseg) (before append)
		// If we append loop logic, cleanups are at:
		// len(vm.codeseg) + [opRAt?] + opLoopPlus + opBranch + offset
		// = len + (1 if pullVal) + 3

		targetIdx := len(vm.codeseg) + 3

		offset := targetIdx - breakIP
		vm.codeseg[breakIP] = uint16(offset)
	}

	vm.codeseg = append(vm.codeseg, opLoopPlus,
		opBranch, uint16(distToStart),
		opRDrop, opRDrop, opRDrop)
	return
}

// LEAVE ( -- )
func opLeave(vm *VM) error {
	if vm.CurrentCompCtx() == nil {
		return ErrBadStateMsg("LEAVE used outside of definition")
	}
	var ctx *LoopCtx
	for i := len(vm.Stack) - 1; i >= 0; i-- {
		if c, ok := vm.Stack[i].(*LoopCtx); ok {
			ctx = c
			break
		}
	}
	if ctx == nil {
		return ErrBadStateMsg("LEAVE outside of loop")
	}

	vm.codeseg = append(vm.codeseg, opBranch, 0)
	ctx.BreakOps = append(ctx.BreakOps, len(vm.codeseg)-1)
	return nil
}

// EXIT ( -- )
func opExit(vm *VM) error {
	if vm.CurrentCompCtx() == nil {
		return ErrBadStateMsg("EXIT used outside of definition")
	}
	vm.codeseg = append(vm.codeseg, opReturn)
	return nil
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

// setupDo sets up a DO loop by pushing limit and start index to rstack and determining loop direction
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

// testDo checks if a DO loop should continue based on the loop direction and current index vs limit
func testDo(vm *VM) (err error) {
	rtop := len(vm.Rstack) - 1
	if rtop < 2 {
		return ErrUnderflow
	}
	rtest, rlim, ridx := vm.Rstack[rtop], vm.Rstack[rtop-1], vm.Rstack[rtop-2]
	testval, ok1 := rtest.(int)
	limval, ok2 := rlim.(int)
	ival, ok3 := ridx.(int)
	// fmt.Printf("rtop: %v  test: %v   limit: %v   idx: %v\n", rtop, testval, limval, ival)
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

// getDoI gets the index of the current DO loop
func getDoI(vm *VM) error {
	rlen := len(vm.Rstack)
	if rlen < 3 {
		return ErrUnderflow
	}
	vm.Push(vm.Rstack[rlen-3])
	return nil
}

// getDoJ gets the index of the outer loop in nested DO loops
func getDoJ(vm *VM) error {
	rlen := len(vm.Rstack)
	if rlen < 6 {
		return ErrUnderflow
	}
	vm.Push(vm.Rstack[rlen-6])
	return nil
}

func branchWordsInit(vm *VM) {
	vm.Define(Word{Name: "if", Run: opIf, Immediate: true})
	vm.Define(Word{Name: "else", Run: opElse, Immediate: true})
	vm.Define(Word{Name: "then", Run: opThen, Immediate: true})
	vm.Define(Word{Name: "recur", Run: recur, Immediate: true})
	vm.Define(Word{Name: "(tail-call)", Run: tailCall, Immediate: true})
	vm.Define(Word{Name: "do", Run: opDo, Immediate: true})
	vm.Define(Word{Name: "(setupDo)", Run: setupDo, Immediate: false})
	vm.Define(Word{Name: "(testDo)", Run: testDo, Immediate: false})
	vm.Define(Word{Name: "(perfLoopPlus)", Run: performLoopPlus, Immediate: false})
	vm.Define(Word{Name: "loop", Run: opLoop, Immediate: true})
	vm.Define(Word{Name: "+loop", Run: opLoopPlus, Immediate: true})
	vm.Define(Word{Name: "i", Run: getDoI, Immediate: false})
	vm.Define(Word{Name: "j", Run: getDoJ, Immediate: false})
	vm.Define(Word{Name: "leave", Run: opLeave, Immediate: true})
	vm.Define(Word{Name: "exit", Run: opExit, Immediate: true})
	vm.Define(Word{Name: "begin", Run: opBegin, Immediate: true})
	vm.Define(Word{Name: "again", Run: opAgain, Immediate: true})
	vm.Define(Word{Name: "until", Run: opUntil, Immediate: true})
	vm.Define(Word{Name: "while", Run: opWhile, Immediate: true})
	vm.Define(Word{Name: "repeat", Run: opRepeat, Immediate: true})
}

// BEGIN ( -- )
func opBegin(vm *VM) (err error) {
	vm.Push(&LoopCtx{
		StartIP:  len(vm.codeseg),
		Type:     LoopBegin,
		BreakOps: make([]int, 0),
	})
	return
}

// resolveBreaks patches any LEAVE/WHILE jumps to point to the current code end
func resolveBreaks(vm *VM, ctx *LoopCtx) {
	targetIdx := len(vm.codeseg)
	for _, breakIP := range ctx.BreakOps {
		offset := targetIdx - breakIP
		vm.codeseg[breakIP] = uint16(offset)
	}
}

// AGAIN ( -- )
func opAgain(vm *VM) (err error) {
	var val interface{}
	val, err = vm.Pop()
	if err != nil {
		return
	}
	ctx, ok := val.(*LoopCtx)
	if !ok || ctx.Type != LoopBegin {
		return ErrBadStateMsg("AGAIN without BEGIN")
	}

	// Unconditional jump back to start
	offset := ctx.StartIP - 1 - len(vm.codeseg)
	vm.codeseg = append(vm.codeseg, opBranch, uint16(offset))

	resolveBreaks(vm, ctx)
	return
}

// UNTIL ( flag -- )
func opUntil(vm *VM) (err error) {
	var val interface{}
	val, err = vm.Pop()
	if err != nil {
		return
	}
	ctx, ok := val.(*LoopCtx)
	if !ok || ctx.Type != LoopBegin {
		return ErrBadStateMsg("UNTIL without BEGIN")
	}

	// Conditional jump back if false (zero)
	offset := ctx.StartIP - 1 - len(vm.codeseg)
	vm.codeseg = append(vm.codeseg, opBZR, uint16(offset))

	resolveBreaks(vm, ctx)
	return
}

// WHILE ( flag -- )
func opWhile(vm *VM) (err error) {
	// Peek for LoopCtx
	var ctx *LoopCtx
	for i := len(vm.Stack) - 1; i >= 0; i-- {
		if c, ok := vm.Stack[i].(*LoopCtx); ok {
			ctx = c
			break
		}
	}
	if ctx == nil || ctx.Type != LoopBegin {
		return ErrBadStateMsg("WHILE outside of BEGIN loop")
	}

	// Branch if Zero (false) to exit
	// Append dummy offset, to be fixed by REPEAT/AGAIN/UNTIL
	vm.codeseg = append(vm.codeseg, opBZR, 0)
	ctx.BreakOps = append(ctx.BreakOps, len(vm.codeseg)-1)
	return
}

// REPEAT ( -- )
func opRepeat(vm *VM) error {
	return opAgain(vm)
}
