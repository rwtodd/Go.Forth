package forth

import (
	"bufio"
	"fmt"
	"io"
)

// define a few constant opcodes that are reliable
// so we don't have to look them up all the time
const (
	opReturn = iota
	opLitINT
	opLitUINT
	opCompileComma
	opBranch
	opBZR
)

// A Word in forth is an operation on the VM
type Word struct {
	Run       func(*VM) error
	Immediate bool
}

// VM is the forth virtual machine state, which all
// operations take
type VM struct {
	words []Word
	dict  map[string]uint16 // maps from names to indexes in `words'

	Stack  []interface{} // the data stack
	Rstack []interface{} // the return stack

	codeseg []uint16 // where the code for composite (user-defined) words go
	ip      int      // instruction pointer
	curdef  int      // the start-index of the word we are currently defining
	curname string   // the name of teh word we are defining

	Source *bufio.Reader // our input
	Sink   *bufio.Writer // out output

	marker uint16 // place to roll back to when we FORGET

	Compiling bool // are we compiling right now?
}

// Define adds a word to the VM
func (vm *VM) Define(name string, word Word) {
	vm.dict[name] = uint16(len(vm.words))
	vm.words = append(vm.words, word)
}

// Forget removes words from the VM up to the
// vm.marker.
func forget(vm *VM) error {
	if len(vm.words) < int(vm.marker) {
		return ErrBadState
	}

	for k, v := range vm.dict {
		if v >= vm.marker {
			delete(vm.dict, k)
		}
	}
	vm.words = vm.words[:vm.marker]
	return nil
}

// Mark sets the marker for a future call to Forget
func mark(vm *VM) error {
	vm.marker = uint16(len(vm.words))
	return nil
}

// Push a value onto the stack
func (vm *VM) Push(v interface{}) {
	vm.Stack = append(vm.Stack, v)
}

// debugPrint prints the codeseg...
func debugPrint(vm *VM) error {
	var revdict = make(map[uint16]string)
	for k, v := range vm.dict {
		revdict[v] = k
	}
	for i, v := range vm.codeseg {
		opcode, ok := revdict[v]
		if !ok {
			opcode = fmt.Sprintf("%d", int16(v))
		}
		fmt.Printf("%03d: %d (%s)\n", i, v, opcode)
	}
	return nil
}

// Pop a value from the stack, returning the value there
func (vm *VM) Pop() (v interface{}, err error) {
	l := len(vm.Stack) - 1
	if l < 0 {
		err = ErrUnderflow
		return
	}
	v = vm.Stack[l]
	vm.Stack = vm.Stack[:l]
	return
}

// RPush pushes a value onto the return stack
func (vm *VM) RPush(v interface{}) {
	vm.Rstack = append(vm.Rstack, v)
}

// RPop pops a value from the return stack, returning the value there
func (vm *VM) RPop() (v interface{}, err error) {
	l := len(vm.Rstack) - 1
	if l < 0 {
		err = ErrRStackUnderflow
		return
	}
	v = vm.Rstack[l]
	vm.Rstack = vm.Rstack[:l]
	return
}

// CreatePusher generates a word in the dictionary, and returns the
// index for the word.  No name is associated with the word.
func (vm *VM) CreatePusher(v interface{}) uint16 {
	vm.words = append(vm.words, Word{Run: func(fvm *VM) error { fvm.Push(v); return nil }, Immediate: false})
	return uint16(len(vm.words) - 1)
}

// NewVM returns a new Forth VM, initialized with the base
// wordset
func NewVM() *VM {
	ans := &VM{
		dict:      make(map[string]uint16),
		Compiling: true,
	}

	// SPECIAL... must be specific opcodes to match constants
	ans.Define("(RET)", Word{nil, false})
	ans.Define("(litINT)", Word{litINT, false})
	ans.Define("(litUINT)", Word{litUINT, false})
	ans.Define("compile,", Word{compileComma, false})
	ans.Define("(branch)", Word{branchUnconditional, false})
	ans.Define("(bzr)", Word{branchZero, false})
	// END SPECIALS

	branchWordsInit(ans)
	stackWordsInit(ans)
	ioWordsInit(ans)
	parseWordsInit(ans)
	numWordsInit(ans)

	// these come from this file...
	ans.Define("mark", Word{mark, false})
	ans.Define("forget", Word{forget, false})
	ans.Define("debug.", Word{debugPrint, false})
	return ans
}

// Run interprets an input stream 'r', writing output
// to an output stream 'w'
func (vm *VM) Run(r io.Reader, w io.Writer) error {
	vm.Source = bufio.NewReader(r)
	vm.Sink = bufio.NewWriter(w)
	vm.Compiling = true
	return interpret(vm)
}

// ResetState recovers from an error and puts us in
// a known state to restart the interpreter
func (vm *VM) ResetState() {
	vm.Stack = nil
	vm.Rstack = nil
	vm.Compiling = true
	vm.curdef = 0
	vm.curname = ""
	vm.ip = 0
}
