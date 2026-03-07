// SPDX-License-Identifier: MIT

package forth

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"
)

// sourceLevel holds information about a single input source
type sourceLevel struct {
	reader  io.RuneReader
	closer  io.Closer
	context string
}

// define a few constant opcodes that are reliable
// so we don't have to look them up all the time
const (
	opReturn = iota
	opEnterScope
	opExitScope
	opLocalGet
	opLocalSet
	opPushClosure

	opLitINT
	opLitUINT
	opLitFloat16
	opLitDecimal
	opCompileComma
	opBranch
	opBZR
	opCallOffset
	opRecurClosure
	opLit
)

// VmMark represents a checkpoint of the VM state for forget operations
type VmMark struct {
	WordCount    uint16
	CodeSegCount int
	LiteralCount int
}

// Word is an interface for any operation in the VM
type Word interface {
	Name() string
	IsImmediate() bool
	SetImmediate(bool)
	Run(*VM) error
	Previous() uint16
	SetPrevious(uint16)
}

// NativeWord is a Word implemented in Go
type NativeWord struct {
	name      string
	run       func(*VM) error
	immediate bool
	prevIdx   uint16
}

func (w *NativeWord) Name() string {
	return w.name
}

// NewNativeWord creates a new native word
func NewNativeWord(name string, run func(*VM) error) *NativeWord {
	return &NativeWord{name: name, run: run}
}

// NewImmediateWord creates a new immediate native word
func NewImmediateWord(name string, run func(*VM) error) *NativeWord {
	return &NativeWord{name: name, run: run, immediate: true}
}

func (w *NativeWord) IsImmediate() bool {
	return w.immediate
}

func (w *NativeWord) SetImmediate(b bool) {
	w.immediate = b
}

func (w *NativeWord) Run(vm *VM) error {
	err := w.run(vm)
	if err != nil {
		return fmt.Errorf("%w\n at `%s`", err, w.Name())
	}
	return nil
}

func (w *NativeWord) Previous() uint16 {
	return w.prevIdx
}

func (w *NativeWord) SetPrevious(idx uint16) {
	w.prevIdx = idx
}

// CompilationCtx holds the state for the word currently being defined
type CompilationCtx struct {
	StartIP        int            // Instruction pointer where this definition starts
	WordIdx        int            // Index of the word being defined
	CompileLocals  map[string]int // Locals defined in the current word
	IsClosure      bool           // Whether this context is a closure
	RecurFixups    []int          // IP locations that need recur fixup
	TailCallFixups []int          // IP locations that need tail-call fixup
	ExitFixups     []int          // IP locations that need exit fixup
	Suspended      bool           // true if we are running interpret
}

// Scope represents a runtime environment for variables
type Scope struct {
	Locals []any
	Parent *Scope
}

// ExecutionToken is an interface for something that can be executed
// ExecutionToken is an interface for something that can be executed
type ExecutionToken interface {
	Run(*VM) error
	Compile(*VM) error
	Name() string
}

// WordToken is an ExecutionToken that wraps a dictionary word index
type WordToken struct {
	Token uint16
}

func (wt WordToken) Name() string {
	return fmt.Sprintf("Token(%d)", wt.Token)
}

func (wt WordToken) Run(vm *VM) error {
	if int(wt.Token) >= len(vm.words) {
		return ErrArgumentMsg("invalid word index in token")
	}
	err := vm.words[wt.Token].Run(vm)
	if err != nil {
		return fmt.Errorf("%w\n at `%s`", err, wt.Name())
	}
	return nil
}

func (wt WordToken) Compile(vm *VM) error {
	if int(wt.Token) >= len(vm.words) {
		return ErrArgumentMsg("invalid word index in token")
	}
	vm.codeseg = append(vm.codeseg, wt.Token)
	return nil
}

// Closure represents a captured execution context
type Closure struct {
	StartIP int
	Env     *Scope
}

func (c Closure) Name() string {
	return fmt.Sprintf("Closure(@%d)", c.StartIP)
}

func (c Closure) Run(vm *VM) error {
	oldHead := vm.HeadScope
	vm.HeadScope = c.Env // Restore captured environment
	err := vm.RunAt(c.StartIP)
	vm.HeadScope = oldHead // Restore previous environment
	if err != nil {
		return fmt.Errorf("%w\n at `%s`", err, c.Name())
	}
	return nil
}

func (c Closure) Compile(vm *VM) error {
	// Compile the closure as a literal
	vm.literals = append(vm.literals, c)
	vm.codeseg = append(vm.codeseg, opLit, uint16(len(vm.literals)-1))

	// Look up EXECUTE
	execIdx, ok := vm.dict["execute"]
	if !ok {
		return ErrBadStateMsg("execute word not found")
	}
	// Compile call to EXECUTE
	vm.codeseg = append(vm.codeseg, execIdx)
	return nil
}

// Variable represents a FORTH variable
type Variable struct {
	value any
}

// Value returns the current value of the variable
func (v *Variable) Value() any {
	return v.value
}

// SetValue sets the value of the variable
func (v *Variable) SetValue(val any) {
	v.value = val
}

// VariableWord is a Word that holds a Variable and an optional ExecutionToken
type VariableWord struct {
	name string
	val  *Variable
	xt   ExecutionToken
	// standard Word impl details
	immediate bool
	prevIdx   uint16
}

func (w *VariableWord) Name() string {
	return w.name
}

func (w *VariableWord) IsImmediate() bool {
	return w.immediate
}

func (w *VariableWord) SetImmediate(b bool) {
	w.immediate = b
}

func (w *VariableWord) Previous() uint16 {
	return w.prevIdx
}

func (w *VariableWord) SetPrevious(idx uint16) {
	w.prevIdx = idx
}

func (w *VariableWord) Run(vm *VM) error {
	vm.Push(w.val)
	if w.xt != nil {
		err := w.xt.Run(vm)
		if err != nil {
			return fmt.Errorf("%w\n at `%s`", err, w.Name())
		}
	}
	return nil
}

// VM is the forth virtual machine state, which all
// operations take
type VM struct {
	words []Word            // Polymorphic list of words
	dict  map[string]uint16 // maps from names to indexes in `words'

	Stack  []any // the data stack
	Rstack []any // the return stack

	literals []any          // Pool of literals to reference by index
	strMap   map[string]int // Map for string interning

	codeseg    []uint16 // where the code for composite (user-defined) words go
	ip         int      // instruction pointer
	curdef     int      // the start-index of the word we are currently defining
	curwordidx int      // the index in words of the word we are currently defining

	CodeSource  io.RuneReader
	CodeContext string
	InputSource io.Reader
	Sink        io.Writer // our output

	marker VmMark // place to roll back to when we FORGET

	CompStack []CompilationCtx // Stack of compilation contexts
	HeadScope *Scope           // Current variable scope

	ActivatedExtensions map[string]bool // Set of activated extensions
}

// CurrentCompCtx returns the current compilation context, or nil if not compiling
// Note: This returns a pointer to the slice element, so modifications affect the stack.
func (vm *VM) CurrentCompCtx() *CompilationCtx {
	if len(vm.CompStack) == 0 {
		return nil
	}
	ctx := &vm.CompStack[len(vm.CompStack)-1]
	if ctx.Suspended {
		return nil
	}
	return ctx
}

// PushCompCtx pushes a new compilation context
func (vm *VM) PushCompCtx(startIP, wordIdx int) {
	vm.CompStack = append(vm.CompStack, CompilationCtx{
		StartIP:        startIP,
		WordIdx:        wordIdx,
		CompileLocals:  make(map[string]int),
		IsClosure:      false,
		RecurFixups:    make([]int, 0),
		TailCallFixups: make([]int, 0),
		ExitFixups:     make([]int, 0),
		Suspended:      false,
	})
}

// PopCompCtx pops the current compilation context
func (vm *VM) PopCompCtx() {
	if len(vm.CompStack) > 0 {
		vm.CompStack = vm.CompStack[:len(vm.CompStack)-1]
	}
}

// PushScope pushes a new variable scope
func (vm *VM) PushScope(size int) {
	vm.HeadScope = &Scope{
		Locals: make([]any, size),
		Parent: vm.HeadScope,
	}
}

// PopScope pops the current variable scope
func (vm *VM) PopScope() {
	if vm.HeadScope != nil {
		vm.HeadScope = vm.HeadScope.Parent
	}
}

// Define adds a word to the VM
func (vm *VM) Define(word Word) {
	vm.words = append(vm.words, word)
	vm.AddToDict(uint16(len(vm.words) - 1))
}

// AddToDict adds a word at the given index to the dictionary,
// handling redefinition logic
func (vm *VM) AddToDict(idx uint16) {
	word := vm.words[idx]
	if existingIdx, exists := vm.dict[word.Name()]; exists {
		if existingIdx < 10 {
			// Cannot redefine special opcodes!
			panic("cannot redefine special opcodes!")
		}
		word.SetPrevious(existingIdx)
	} else {
		word.SetPrevious(0)
	}
	vm.dict[word.Name()] = idx
}

// Forget removes words from the VM up to the
// vm.marker.
func forget(vm *VM) error {
	if len(vm.words) < int(vm.marker.WordCount) {
		return ErrBadStateMsg("cannot forget below marker")
	}
	if len(vm.codeseg) < vm.marker.CodeSegCount {
		return ErrBadStateMsg("cannot forget below code marker")
	}

	// Traverse words backwards from the end and restore previous definitions
	for i := len(vm.words) - 1; i >= int(vm.marker.WordCount); i-- {
		word := vm.words[i]
		if word.Previous() > 0 {
			vm.dict[word.Name()] = word.Previous()
		} else {
			delete(vm.dict, word.Name())
		}
	}
	clear(vm.words[vm.marker.WordCount:])
	vm.words = vm.words[:vm.marker.WordCount]
	vm.codeseg = vm.codeseg[:vm.marker.CodeSegCount]
	if len(vm.literals) >= int(vm.marker.LiteralCount) {
		// Remove interned strings that are being forgotten
		for k, v := range vm.strMap {
			if v >= int(vm.marker.LiteralCount) {
				delete(vm.strMap, k)
			}
		}
		clear(vm.literals[vm.marker.LiteralCount:])
		vm.literals = vm.literals[:vm.marker.LiteralCount]
	}
	return nil
}

// Mark sets the marker for a future call to Forget
func mark(vm *VM) error {
	vm.marker = VmMark{
		WordCount:    uint16(len(vm.words)),
		CodeSegCount: len(vm.codeseg),
		LiteralCount: len(vm.literals),
	}
	return nil
}

// constant creates a new FORTH constant
func constant(vm *VM) error {
	strVal, err := vm.Pop()
	if err != nil {
		return err
	}
	str, ok := strVal.(string)
	if !ok {
		return ErrArgumentMsg("constant expects a string name")
	}

	val, err := vm.Pop()
	if err != nil {
		return err
	}

	return createConstant(vm, str, val)
}

func createConstant(vm *VM, name string, val any) error {
	// Look up @
	atIdx, ok := vm.dict["@"]
	if !ok {
		return ErrBadStateMsg("@ word not found")
	}

	// Create VariableWord with @ as XT
	varObj := &Variable{value: val}
	word := &VariableWord{
		name: strings.ToLower(strings.TrimSpace(name)),
		val:  varObj,
		xt:   WordToken{Token: atIdx},
	}

	vm.Define(word)
	return nil
}

// variableDoes creates a new FORTH word that behaves like a variable with custom behavior
// ( val xt name -- )
func variableDoes(vm *VM) error {
	strVal, err := vm.Pop()
	if err != nil {
		return err
	}
	str, ok := strVal.(string)
	if !ok {
		return ErrArgumentMsg("variable-does expects a string name")
	}

	xtVal, err := vm.Pop()
	if err != nil {
		return err
	}

	var xt ExecutionToken
	if tok, ok := xtVal.(ExecutionToken); ok {
		xt = tok
	} else if idx, ok := xtVal.(int64); ok {
		xt = WordToken{Token: uint16(idx)}
	} else {
		return ErrArgumentMsg("variable-does expects execution token or word index")
	}

	val, err := vm.Pop()
	if err != nil {
		return err
	}

	return createVariableDoes(vm, str, val, xt)
}

func createVariableDoes(vm *VM, name string, val any, xt ExecutionToken) error {
	varObj := &Variable{value: val}
	word := &VariableWord{
		name: strings.ToLower(strings.TrimSpace(name)),
		val:  varObj,
		xt:   xt,
	}

	vm.Define(word)
	return nil
}

// variable creates a new FORTH variable
// ( val name -- )
func variable(vm *VM) error {
	strVal, err := vm.Pop()
	if err != nil {
		return err
	}
	str, ok := strVal.(string)
	if !ok {
		return ErrArgumentMsg("variable expects a string name")
	}

	val, err := vm.Pop()
	if err != nil {
		return err
	}

	return createVariable(vm, str, val)
}

func createVariable(vm *VM, name string, val any) error {
	varObj := &Variable{value: val}
	// Optimization: nil XT avoids overhead
	word := &VariableWord{
		name: strings.ToLower(strings.TrimSpace(name)),
		val:  varObj,
		xt:   nil,
	}
	vm.Define(word)
	return nil
}

// Push a value onto the stack
func (vm *VM) Push(v any) {
	vm.Stack = append(vm.Stack, v)
}

// debugPrint prints the codeseg with smart opcode decoding
func debugPrint(vm *VM) error {
	opTestDo, ok := vm.dict["(testDo)"]
	if !ok {
		return ErrBadStateMsg("testDo not found")
	}
	for i := 0; i < len(vm.codeseg); i++ {
		v := vm.codeseg[i]
		switch v {
		case opEnterScope:
			if i+1 < len(vm.codeseg) {
				fmt.Printf("%03d: %d (enterScope %d)\n", i, v, vm.codeseg[i+1])
				i++ // skip the data
			} else {
				fmt.Printf("%03d: %d (enterScope ???)\n", i, v)
			}
		case opExitScope:
			fmt.Printf("%03d: %d (exitScope)\n", i, v)
		case opLocalGet:
			if i+1 < len(vm.codeseg) {
				packed := vm.codeseg[i+1]
				fmt.Printf("%03d: %d (localGet %d %d)\n", i, v, packed>>10, packed&0x3FF)
				i++ // skip the data
			} else {
				fmt.Printf("%03d: %d (localGet ???)\n", i, v)
			}
		case opLocalSet:
			if i+1 < len(vm.codeseg) {
				packed := vm.codeseg[i+1]
				fmt.Printf("%03d: %d (localSet %d %d)\n", i, v, packed>>10, packed&0x3FF)
				i++ // skip the data
			} else {
				fmt.Printf("%03d: %d (localSet ???)\n", i, v)
			}
		case opPushClosure:
			if i+1 < len(vm.codeseg) {
				fmt.Printf("%03d: %d (pushClosure %d)\n", i, v, int16(vm.codeseg[i+1]))
				i++ // skip the data
			} else {
				fmt.Printf("%03d: %d (pushClosure ???)\n", i, v)
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
		case opLitFloat16:
			if i+1 < len(vm.codeseg) {
				val := float16ToFloat64(vm.codeseg[i+1])
				fmt.Printf("%03d: %d (litFloat16 %g)\n", i, v, val)
				i++ // skip the data
			} else {
				fmt.Printf("%03d: %d (litFloat16 ???)\n", i, v)
			}
		case opLitDecimal:
			if i+1 < len(vm.codeseg) {
				val := decimal16ToFloat64(vm.codeseg[i+1])
				fmt.Printf("%03d: %d (litDecimal %g)\n", i, v, val)
				i++ // skip the data
			} else {
				fmt.Printf("%03d: %d (litDecimal ???)\n", i, v)
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
		case opCallOffset:
			if i+1 < len(vm.codeseg) {
				fmt.Printf("%03d: %d (callOffset %d)\n", i, v, int16(vm.codeseg[i+1]))
				i++ // skip the data
			} else {
				fmt.Printf("%03d: %d (callOffset ???)\n", i, v)
			}
		case opTestDo:
			if i+1 < len(vm.codeseg) {
				fmt.Printf("%03d: %d (testDo %d)\n", i, v, int16(vm.codeseg[i+1]))
				i++ // skip the data
			} else {
				fmt.Printf("%03d: %d (testDo ???)\n", i, v)
			}
		case opRecurClosure:
			if i+1 < len(vm.codeseg) {
				fmt.Printf("%03d: %d (recurClosure %d)\n", i, v, int16(vm.codeseg[i+1]))
				i++ // skip the data
			} else {
				fmt.Printf("%03d: %d (recurClosure ???)\n", i, v)
			}
		case opLit:
			if i+1 < len(vm.codeseg) {
				idx := vm.codeseg[i+1]
				var val string
				if int(idx) < len(vm.literals) {
					val = fmt.Sprintf("%v", vm.literals[idx])
				} else {
					val = fmt.Sprintf("INVALID INDEX %d", idx)
				}
				fmt.Printf("%03d: %d (lit #%d = %s)\n", i, v, idx, val)
				i++ // skip the data
			} else {
				fmt.Printf("%03d: %d (lit ???)\n", i, v)
			}
		default:
			// Regular word call - look up the name
			if int(v) < len(vm.words) {
				fmt.Printf("%03d: %d (%s)\n", i, v, vm.words[v].Name())
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

func execute(vm *VM) error {
	val, err := vm.Pop()
	if err != nil {
		return err
	}
	switch v := val.(type) {
	case ExecutionToken:
		return v.Run(vm)
	case int64:
		// Support legacy int index execution for now? Or convert?
		// Plan implies full switch. But raw ints might be on stack?
		// Let's support ints as word indexes for backward comp if needed,
		// but ideally everything should be Token.
		// For now, allow int for safety.
		idx := v
		if idx < 0 || int(idx) >= len(vm.words) {
			return ErrArgumentMsg("invalid word index")
		}
		return vm.words[idx].Run(vm)
	default:
		return ErrArgumentMsg("execute requires an execution token or word index")
	}
}

// enterScope implements the opEnterScope opcode
func enterScope(vm *VM) error {
	vm.ip++
	count := vm.codeseg[vm.ip]
	vm.PushScope(int(count))
	return nil
}

// exitScope implements the opExitScope opcode
func exitScope(vm *VM) error {
	vm.PopScope()
	return nil
}

// localGet implements the opLocalGet opcode
func localGet(vm *VM) error {
	packed := vm.codeseg[vm.ip+1]
	depth := packed >> 10
	lidx := packed & 0x3FF
	vm.ip++

	// find the scope at depth
	scope := vm.HeadScope
	for i := 0; i < int(depth); i++ {
		if scope == nil {
			return ErrBadStateMsg("local variable depth too deep")
		}
		scope = scope.Parent
	}
	if scope == nil {
		return ErrBadStateMsg("local variable depth too deep (nil scope)")
	}
	vm.Push(scope.Locals[lidx])
	return nil
}

// localSet implements the opLocalSet opcode
func localSet(vm *VM) error {
	packed := vm.codeseg[vm.ip+1]
	depth := packed >> 10
	lidx := packed & 0x3FF
	vm.ip++

	val, err := vm.Pop()
	if err != nil {
		return err
	}

	// find the scope at depth
	scope := vm.HeadScope
	for i := 0; i < int(depth); i++ {
		if scope == nil {
			return ErrBadStateMsg("local variable depth too deep")
		}
		scope = scope.Parent
	}
	if scope == nil {
		return ErrBadStateMsg("local variable depth too deep (nil scope)")
	}
	scope.Locals[lidx] = val
	return nil
}

// pushClosure implements the opPushClosure opcode
func pushClosure(vm *VM) error {
	vm.ip++
	offset := int16(vm.codeseg[vm.ip])
	targetIP := vm.ip + int(offset)
	vm.Push(Closure{StartIP: int(targetIP), Env: vm.HeadScope})
	return nil
}

// callOffset implements the opCallOffset opcode
// It calls a function at a relative offset from the current IP.
func callOffset(vm *VM) error {
	vm.ip++
	offset := int16(vm.codeseg[vm.ip])
	targetIP := vm.ip + int(offset)
	return vm.RunAt(targetIP)
}

// litFloat16 implements the opLitFloat16 opcode
// It reads the next 16-bits from the codeseg, interprets as Float16, and pushes float64.
func litFloat16(vm *VM) error {
	vm.ip++
	val := float16ToFloat64(vm.codeseg[vm.ip])
	vm.Stack = append(vm.Stack, val)
	return nil
}

// litDecimal implements the opLitDecimal opcode
// It reads the next 16-bits from the codeseg, interprets as Decimal, and pushes float64.
func litDecimal(vm *VM) error {
	vm.ip++
	val := decimal16ToFloat64(vm.codeseg[vm.ip])
	vm.Stack = append(vm.Stack, val)
	return nil
}

// recurClosure implements the opRecurClosure opcode
// It calls a function at a relative offset, but temporarily restores the parent scope
// to ensure captured variables are accessed correctly during the recursive call.
func recurClosure(vm *VM) error {
	vm.ip++
	offset := int16(vm.codeseg[vm.ip])
	targetIP := vm.ip + int(offset)

	oldHead := vm.HeadScope
	// We need to unwrap ONE level of scope (the local scope of the current closure)
	// so that the new instance is called with the captured environment as HeadScope.
	if oldHead != nil {
		vm.HeadScope = oldHead.Parent
	}

	err := vm.RunAt(targetIP)

	vm.HeadScope = oldHead
	vm.HeadScope = oldHead
	return err
}

// lit implements the opLit opcode
// It pushes a value from the literals pool to the stack
func lit(vm *VM) error {
	vm.ip++
	idx := vm.codeseg[vm.ip]
	if int(idx) >= len(vm.literals) {
		return ErrBadStateMsg("literal index out of bounds")
	}
	vm.Push(vm.literals[idx])
	return nil
}

// opThrow implements the throw word ( str -- )
func opThrow(vm *VM) error {
	strVal, err := vm.Pop()
	if err != nil {
		return err
	}
	str, ok := strVal.(string)
	if !ok {
		return ErrArgumentMsg("throw expects a string message")
	}
	return ErrUserMsg(str)
}

// opCatch implements the catch? word ( ?? n xt -- ?? bool )
func opCatch(vm *VM) error {
	xtVal, err := vm.Pop()
	if err != nil {
		return err
	}

	nVal, err := vm.Pop()
	if err != nil {
		return err
	}

	n, ok := nVal.(int64)
	if !ok {
		return ErrArgumentMsg("catch? expects an integer number of stack items")
	}

	if n < 0 || int(n) > len(vm.Stack) {
		return ErrArgumentMsg("catch? invalid number of stack items")
	}

	var xt ExecutionToken
	if tok, ok := xtVal.(ExecutionToken); ok {
		xt = tok
	} else if idx, ok := xtVal.(int64); ok {
		xt = WordToken{Token: uint16(idx)}
	} else {
		return ErrArgumentMsg("catch? expects execution token or word index")
	}

	// Isolate the working stack
	originalStack := vm.Stack[:len(vm.Stack)-int(n)]
	workingStack := make([]any, int(n))
	copy(workingStack, vm.Stack[len(vm.Stack)-int(n):])

	vm.Stack = workingStack

	// Save state to restore on error
	rStackLen := len(vm.Rstack)
	oldScope := vm.HeadScope

	execErr := xt.Run(vm)

	if execErr != nil {
		// On error: restore original stack, discard working stack
		vm.Stack = originalStack

		// Clean up return stack from aborted run
		if len(vm.Rstack) > rStackLen {
			clear(vm.Rstack[rStackLen:])
			vm.Rstack = vm.Rstack[:rStackLen]
		}
		// Clean up scope
		vm.HeadScope = oldScope

		// Push error string and -1 (true)
		vm.Push(execErr.Error())
		vm.Push(int64(-1))
	} else {
		// On success: append working stack, push 0 (false)
		vm.Stack = append(originalStack, vm.Stack...)
		vm.Push(int64(0)) // 0 represents false logic in Go.Forth usually
	}

	return nil
}

// RunAt runs the code segment starting at the given IP
func (vm *VM) RunAt(startIP int) error {
	rstackLen := len(vm.Rstack)
	oldIP := vm.ip
	vm.ip = startIP

	for {
		idx := vm.codeseg[vm.ip]
		if idx == opReturn {
			break
		}
		if err := vm.words[idx].Run(vm); err != nil {
			vm.ip = oldIP
			return err
		}
		vm.ip++
	}

	if len(vm.Rstack) != rstackLen {
		if len(vm.Rstack) < rstackLen {
			return ErrUnderflowMsg("return stack underflow")
		}
		// clean up the rstack and exit
		clear(vm.Rstack[rstackLen:])
		vm.Rstack = vm.Rstack[:rstackLen]
	}

	vm.ip = oldIP
	return nil
}

// NewVM returns a new Forth VM, initialized with the base
// wordset and standard input/output readers
func NewVM(in io.Reader, out io.Writer) *VM {
	ans := &VM{
		dict:                make(map[string]uint16),
		strMap:              make(map[string]int),
		ActivatedExtensions: make(map[string]bool),
		InputSource:         in,
		Sink:                out,
	}

	// SPECIAL... must be specific opcodes to match constants
	ans.Define(&NativeWord{name: "(RET)", run: nil, immediate: false})
	ans.Define(&NativeWord{name: "(enterScope)", run: enterScope, immediate: false})
	ans.Define(&NativeWord{name: "(exitScope)", run: exitScope, immediate: false})
	ans.Define(&NativeWord{name: "(localGet)", run: localGet, immediate: false})
	ans.Define(&NativeWord{name: "(localSet)", run: localSet, immediate: false})
	ans.Define(&NativeWord{name: "(pushClosure)", run: pushClosure, immediate: false})

	ans.Define(&NativeWord{name: "(litINT)", run: litINT, immediate: false})
	ans.Define(&NativeWord{name: "(litUINT)", run: litUINT, immediate: false})
	ans.Define(&NativeWord{name: "(litFloat16)", run: litFloat16, immediate: false})
	ans.Define(&NativeWord{name: "(litDecimal)", run: litDecimal, immediate: false})
	ans.Define(&NativeWord{name: "compile-xt", run: compileComma, immediate: true})
	ans.Define(&NativeWord{name: "(branch)", run: branchUnconditional, immediate: false})
	ans.Define(&NativeWord{name: "(bzr)", run: branchZero, immediate: false})
	ans.Define(&NativeWord{name: "(call-offset)", run: callOffset, immediate: false})
	ans.Define(&NativeWord{name: "(recur-closure)", run: recurClosure, immediate: false})
	ans.Define(&NativeWord{name: "(lit)", run: lit, immediate: false})
	// END SPECIALS

	branchWordsInit(ans)
	stackWordsInit(ans)
	ioWordsInit(ans)
	parseWordsInit(ans)
	numWordsInit(ans)
	arrayWordsInit(ans)
	dictWordsInit(ans)
	comparisonWordsInit(ans)
	varLenWordsInit(ans)
	extensionsWordsInit(ans)

	// these come from this file...
	ans.Define(&NativeWord{name: "mark", run: mark, immediate: false})
	ans.Define(&NativeWord{name: "forget", run: forget, immediate: false})
	ans.Define(&NativeWord{name: "debug.", run: debugPrint, immediate: false})
	ans.Define(&NativeWord{name: "variable", run: variable, immediate: false})
	ans.Define(&NativeWord{name: "variable-does", run: variableDoes, immediate: false})
	ans.Define(&NativeWord{name: "constant", run: constant, immediate: false})
	ans.Define(&NativeWord{name: "execute", run: execute, immediate: false})
	ans.Define(&NativeWord{name: "catch?", run: opCatch, immediate: false})
	ans.Define(&NativeWord{name: "throw", run: opThrow, immediate: false})

	_ = mark(ans) // give the vm an initial mark after all the core words are added
	return ans
}

// SetInput changes the VM's input source for interactive reads
func (vm *VM) SetInput(in io.Reader) {
	vm.InputSource = in
}

// SetOutput changes the VM's standard output stream
func (vm *VM) SetOutput(out io.Writer) {
	vm.Sink = out
}

// wrapError prepends the current source context to an error
func (vm *VM) wrapError(err error) error {
	if err == nil || vm.CodeContext == "" {
		return err
	}
	return fmt.Errorf("%s: %w", vm.CodeContext, err)
}

// ReadRune implements io.RuneReader for the VM, pulling from the current CodeSource.
func (vm *VM) ReadRune() (rune, int, error) {
	if vm.CodeSource == nil {
		return 0, 0, io.EOF
	}
	return vm.CodeSource.ReadRune()
}

// Run executes the given code source, tracking the context.
// It replaces the previous code source for the duration of the interpretation.
func (vm *VM) Run(codeSource io.Reader, context string) error {
	oldSource := vm.CodeSource
	oldContext := vm.CodeContext

	// We wrap in a bufio.Reader to ensure we have RuneReader capabilities and a trailing newline
	inner := io.MultiReader(codeSource, strings.NewReader("\n"))
	vm.CodeSource = bufio.NewReader(inner)
	vm.CodeContext = context

	err := interpret(vm)

	if closer, ok := codeSource.(io.Closer); ok {
		_ = closer.Close()
	}

	vm.CodeSource = oldSource
	vm.CodeContext = oldContext

	return err
}

// Load executes a file by name, setting the context to the filename.
func (vm *VM) Load(filename string) error {
	f, err := os.Open(filename)
	if err != nil {
		return err
	}
	return vm.Run(bufio.NewReader(f), filename)
}

// Eval interprets a string completely and returns any errors.
// This allows evaluation of strings at any time without consuming the rest of the current script.
func (vm *VM) Eval(code string) error {
	ctxName := "eval"
	if vm.CodeContext != "" {
		ctxName = fmt.Sprintf("eval (%s)", vm.CodeContext)
	}
	return vm.Run(strings.NewReader(code), ctxName)
}

// ResetState recovers from an error and puts us in
// a known state to restart the interpreter
func (vm *VM) ResetState() {
	vm.Stack = nil
	vm.Rstack = nil
	vm.CompStack = nil
	vm.ip = 0
}

// ClearResetState resets the execution state.
func (vm *VM) ClearResetState() {
	vm.ResetState()
}
