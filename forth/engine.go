// SPDX-License-Identifier: MIT

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
	opCreateLocals
	opLocalGet
	opLocalSet

	opLitINT
	opLitUINT
	opCompileComma
	opBranch
	opBZR
)

// A Word in forth is an operation on the VM
type Word struct {
	Name      string
	Run       func(*VM) error
	Immediate bool
	PrevIdx   uint16 // Previous definition index, 0 means no previous
}

// Variable represents a FORTH variable
type Variable struct {
	value any
}

// VM is the forth virtual machine state, which all
// operations take
type VM struct {
	words []Word
	dict  map[string]uint16 // maps from names to indexes in `words'

	Stack  []any // the data stack
	Rstack []any // the return stack

	codeseg    []uint16 // where the code for composite (user-defined) words go
	ip         int      // instruction pointer
	curdef     int      // the start-index of the word we are currently defining
	curwordidx int      // the index in words of the word we are currently defining

	Source *bufio.Reader // our input
	Sink   io.Writer     // our output

	marker     uint16 // place to roll back to when we FORGET for words
	codeMarker int    // place to roll back to when we FORGET for codeseg

	Compiling     bool           // are we compiling right now?
	CompileLocals map[string]int // locals defined in the current word
}

// Define adds a word to the VM
func (vm *VM) Define(word Word) {
	vm.words = append(vm.words, word)
	vm.AddToDict(uint16(len(vm.words) - 1))
}

// AddToDict adds a word at the given index to the dictionary,
// handling redefinition logic
func (vm *VM) AddToDict(idx uint16) {
	word := &vm.words[idx]
	if existingIdx, exists := vm.dict[word.Name]; exists {
		if existingIdx < 10 {
			// Cannot redefine special opcodes!
			panic("cannot redefine special opcodes!")
		}
		word.PrevIdx = existingIdx
	} else {
		word.PrevIdx = 0
	}
	vm.dict[word.Name] = idx
}

// Forget removes words from the VM up to the
// vm.marker.
func forget(vm *VM) error {
	if len(vm.words) < int(vm.marker) {
		return ErrBadStateMsg("cannot forget below marker")
	}
	if len(vm.codeseg) < vm.codeMarker {
		return ErrBadStateMsg("cannot forget below code marker")
	}

	// Traverse words backwards from the end and restore previous definitions
	for i := len(vm.words) - 1; i >= int(vm.marker); i-- {
		word := &vm.words[i]
		if word.PrevIdx > 0 {
			vm.dict[word.Name] = word.PrevIdx
		} else {
			delete(vm.dict, word.Name)
		}
		// Nil out the Run function to help GC
		word.Run = nil
		word.Name = ""
	}
	vm.words = vm.words[:vm.marker]
	vm.codeseg = vm.codeseg[:vm.codeMarker]
	return nil
}

// Mark sets the marker for a future call to Forget
func mark(vm *VM) error {
	vm.marker = uint16(len(vm.words))
	vm.codeMarker = len(vm.codeseg)
	return nil
}

// variable creates a new FORTH variable
func variable(vm *VM) error {
	buf := make([]rune, 0, 20)
	str, err := nextToken(vm, buf)
	if err != nil {
		return err
	}
	varObj := &Variable{value: 0}
	idx := vm.CreatePusher(varObj)
	vm.dict[str] = idx
	return nil
}

// Push a value onto the stack
func (vm *VM) Push(v any) {
	vm.Stack = append(vm.Stack, v)
}

// debugPrint prints the codeseg with smart opcode decoding
func debugPrint(vm *VM) error {
	for i := 0; i < len(vm.codeseg); i++ {
		v := vm.codeseg[i]
		switch v {
		case opCreateLocals:
			if i+1 < len(vm.codeseg) {
				fmt.Printf("%03d: %d (createLocals %d)\n", i, v, vm.codeseg[i+1])
				i++ // skip the data
			} else {
				fmt.Printf("%03d: %d (createLocals ???)\n", i, v)
			}
		case opLocalGet:
			if i+1 < len(vm.codeseg) {
				fmt.Printf("%03d: %d (localGet %d)\n", i, v, vm.codeseg[i+1])
				i++ // skip the data
			} else {
				fmt.Printf("%03d: %d (localGet ???)\n", i, v)
			}
		case opLocalSet:
			if i+1 < len(vm.codeseg) {
				fmt.Printf("%03d: %d (localSet %d)\n", i, v, vm.codeseg[i+1])
				i++ // skip the data
			} else {
				fmt.Printf("%03d: %d (localSet ???)\n", i, v)
			}
		case opLitINT:
			if i+1 < len(vm.codeseg) {
				fmt.Printf("%03d: %d (litINT %d)\n", i, v, int16(vm.codeseg[i+1]))
				i++ // skip the data
			} else {
				fmt.Printf("%03d: %d (litINT ???)\n", i, v)
			}
		case opLitUINT:
			if i+1 < len(vm.codeseg) {
				fmt.Printf("%03d: %d (litUINT %d)\n", i, v, vm.codeseg[i+1])
				i++ // skip the data
			} else {
				fmt.Printf("%03d: %d (litUINT ???)\n", i, v)
			}
		case opBranch:
			if i+1 < len(vm.codeseg) {
				fmt.Printf("%03d: %d (branch %d)\n", i, v, int16(vm.codeseg[i+1]))
				i++ // skip the data
			} else {
				fmt.Printf("%03d: %d (branch ???)\n", i, v)
			}
		case opBZR:
			if i+1 < len(vm.codeseg) {
				fmt.Printf("%03d: %d (bzr %d)\n", i, v, int16(vm.codeseg[i+1]))
				i++ // skip the data
			} else {
				fmt.Printf("%03d: %d (bzr ???)\n", i, v)
			}
		default:
			// Regular word call - look up the name
			if int(v) < len(vm.words) {
				fmt.Printf("%03d: %d (%s)\n", i, v, vm.words[v].Name)
			} else {
				fmt.Printf("%03d: %d (INVALID INDEX %d)\n", i, v, v)
			}
		}
	}
	return nil
}

// Pop a value from the stack, returning the value there
func (vm *VM) Pop() (v any, err error) {
	l := len(vm.Stack) - 1
	if l < 0 {
		err = ErrUnderflow // Standard data stack underflow
		return
	}
	v = vm.Stack[l]
	vm.Stack = vm.Stack[:l]
	return
}

// RPush pushes a value onto the return stack
func (vm *VM) RPush(v any) {
	vm.Rstack = append(vm.Rstack, v)
}

// RPop pops a value from the return stack, returning the value there
func (vm *VM) RPop() (v any, err error) {
	l := len(vm.Rstack) - 1
	if l < 0 {
		err = ErrUnderflowMsg("return stack underflow")
		return
	}
	v = vm.Rstack[l]
	vm.Rstack[l] = nil
	vm.Rstack = vm.Rstack[:l]
	return
}

var iAmAPusher = "-- Value Pusher --" // used for all Value Pushers...

// CreatePusher generates a word in the dictionary, and returns the
// index for the word.  No name is associated with the word.
func (vm *VM) CreatePusher(v any) uint16 {
	vm.words = append(vm.words, Word{Name: iAmAPusher, Run: func(fvm *VM) error { fvm.Push(v); return nil }, Immediate: false})
	return uint16(len(vm.words) - 1)
}

func execute(vm *VM) error {
	val, err := vm.Pop()
	if err != nil {
		return err
	}
	idx, ok := val.(int)
	if !ok {
		return ErrArgumentMsg("execute requires an integer index")
	}
	if idx < 0 || idx >= len(vm.words) {
		return ErrArgumentMsg("invalid word index")
	}
	return vm.words[idx].Run(vm)
}

// NewVM returns a new Forth VM, initialized with the base
// wordset
func NewVM() *VM {
	ans := &VM{
		dict:          make(map[string]uint16),
		Compiling:     true,
		CompileLocals: make(map[string]int),
	}

	// SPECIAL... must be specific opcodes to match constants
	ans.Define(Word{Name: "(RET)", Run: nil, Immediate: false})
	ans.Define(Word{Name: "(createLocals)", Run: nil, Immediate: false})
	ans.Define(Word{Name: "(localGet)", Run: nil, Immediate: false})
	ans.Define(Word{Name: "(localSet)", Run: nil, Immediate: false})

	ans.Define(Word{Name: "(litINT)", Run: litINT, Immediate: false})
	ans.Define(Word{Name: "(litUINT)", Run: litUINT, Immediate: false})
	ans.Define(Word{Name: "compile,", Run: compileComma, Immediate: false})
	ans.Define(Word{Name: "(branch)", Run: branchUnconditional, Immediate: false})
	ans.Define(Word{Name: "(bzr)", Run: branchZero, Immediate: false})
	// END SPECIALS

	branchWordsInit(ans)
	stackWordsInit(ans)
	ioWordsInit(ans)
	parseWordsInit(ans)
	numWordsInit(ans)
	arrayWordsInit(ans)
	dictWordsInit(ans)
	comparisonWordsInit(ans)

	// these come from this file...
	ans.Define(Word{Name: "mark", Run: mark, Immediate: false})
	ans.Define(Word{Name: "forget", Run: forget, Immediate: false})
	ans.Define(Word{Name: "debug.", Run: debugPrint, Immediate: false})
	ans.Define(Word{Name: "variable", Run: variable, Immediate: false})
	ans.Define(Word{Name: "execute", Run: execute, Immediate: false})
	return ans
}

// Run interprets an input stream 'r', writing output
// to an output stream 'w'
func (vm *VM) Run(r io.Reader, w io.Writer) error {
	vm.Source = bufio.NewReader(r)
	vm.Sink = w
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
	vm.curwordidx = -1
	vm.ip = 0
}
