package strings

import (
	"strings"
	"unicode"

	"github.com/rwtodd/Go.Forth/forth"
)

func init() {
	forth.RegisterExtensionWithHelp("strings", func(vm *forth.VM) error {
		vm.Define(forth.NewNativeWord(`blank""?`, isBlank))
		vm.Define(forth.NewNativeWord(`len""`, strLen))
		vm.Define(forth.NewNativeWord(`trim""`, trim))
		vm.Define(forth.NewNativeWord(`triml""`, trimL))
		vm.Define(forth.NewNativeWord(`trimr""`, trimR))
		vm.Define(forth.NewNativeWord(`starts""?`, startsWith))
		vm.Define(forth.NewNativeWord(`ends""?`, endsWith))
		vm.Define(forth.NewNativeWord(`upper""`, upper))
		vm.Define(forth.NewNativeWord(`lower""`, lower))
		vm.Define(forth.NewNativeWord(`split""`, split))
		vm.Define(forth.NewNativeWord(`@join""`, arrayJoin))
		vm.Define(forth.NewNativeWord(`<join"">`, stackJoin))
		vm.Define(forth.NewNativeWord(`sub""`, substring))
		vm.Define(forth.NewNativeWord(`replace""`, replace))
		vm.Define(forth.NewNativeWord(`index""`, indexStr))
		vm.Define(forth.NewNativeWord(`contains""?`, contains))
		return nil
	}, `blank""? ( str -- bool ) checks if a string is entirely whitespace or empty.
len"" ( str -- len ) returns the length of the string in bytes.
trim"" ( str -- str ) trims whitespace from both ends of the string.
triml"" ( str -- str ) trims leading whitespace from the string.
trimr"" ( str -- str ) trims trailing whitespace from the string.
starts""? ( str prefix -- bool ) checks if the target string starts with the prefix.
ends""? ( str suffix -- bool ) checks if the target string ends with the suffix.
upper"" ( str -- str ) converts the string to uppercase.
lower"" ( str -- str ) converts the string to lowercase.
split"" ( str sep -- arr ) splits the string by a separator into an array of strings.
@join"" ( arr sep -- str ) joins an array of strings by a separator.
<join""> ( str1 ... strn n sep -- str ) takes n strings from the stack and joins them by a separator.
sub"" ( str idx1 idx2 -- substr ) extracts the substring [idx1, idx2).
replace"" ( str old new -- str ) replaces all occurrences of old with new in the string.
index"" ( str substr -- index ) returns the byte index of the first instance of substr, or -1 if not found.
contains""? ( str substr -- bool ) checks if the target string contains the substring.`)
}

func pushBool(vm *forth.VM, b bool) {
	if b {
		vm.Push(int64(-1))
	} else {
		vm.Push(int64(0))
	}
}

// blank""? ( str -- bool )
func isBlank(vm *forth.VM) error {
	sVal, err := vm.Pop()
	if err != nil {
		return err
	}
	s, ok := sVal.(string)
	if !ok {
		return forth.ErrArgumentMsg("blank\"\"? expects a string")
	}
	pushBool(vm, strings.TrimSpace(s) == "")
	return nil
}

// len"" ( str -- len )
func strLen(vm *forth.VM) error {
	sVal, err := vm.Pop()
	if err != nil {
		return err
	}
	s, ok := sVal.(string)
	if !ok {
		return forth.ErrArgumentMsg("len\"\" expects a string")
	}
	vm.Push(int64(len(s)))
	return nil
}

// trim"" ( str -- str )
func trim(vm *forth.VM) error {
	sVal, err := vm.Pop()
	if err != nil {
		return err
	}
	s, ok := sVal.(string)
	if !ok {
		return forth.ErrArgumentMsg("trim\"\" expects a string")
	}
	vm.Push(strings.TrimSpace(s))
	return nil
}

// triml"" ( str -- str )
func trimL(vm *forth.VM) error {
	sVal, err := vm.Pop()
	if err != nil {
		return err
	}
	s, ok := sVal.(string)
	if !ok {
		return forth.ErrArgumentMsg("triml\"\" expects a string")
	}
	vm.Push(strings.TrimLeftFunc(s, unicode.IsSpace))
	return nil
}

// trimr"" ( str -- str )
func trimR(vm *forth.VM) error {
	sVal, err := vm.Pop()
	if err != nil {
		return err
	}
	s, ok := sVal.(string)
	if !ok {
		return forth.ErrArgumentMsg("trimr\"\" expects a string")
	}
	vm.Push(strings.TrimRightFunc(s, unicode.IsSpace))
	return nil
}

// starts""? ( str prefix -- bool )
func startsWith(vm *forth.VM) error {
	pVal, err := vm.Pop()
	if err != nil {
		return err
	}
	prefix, ok := pVal.(string)
	if !ok {
		return forth.ErrArgumentMsg("starts\"\"? expects a suffix string")
	}

	sVal, err := vm.Pop()
	if err != nil {
		return err
	}
	s, ok := sVal.(string)
	if !ok {
		return forth.ErrArgumentMsg("starts\"\"? expects a target string")
	}

	pushBool(vm, strings.HasPrefix(s, prefix))
	return nil
}

// ends""? ( str suffix -- bool )
func endsWith(vm *forth.VM) error {
	sfVal, err := vm.Pop()
	if err != nil {
		return err
	}
	suffix, ok := sfVal.(string)
	if !ok {
		return forth.ErrArgumentMsg("ends\"\"? expects a suffix string")
	}

	sVal, err := vm.Pop()
	if err != nil {
		return err
	}
	s, ok := sVal.(string)
	if !ok {
		return forth.ErrArgumentMsg("ends\"\"? expects a target string")
	}

	pushBool(vm, strings.HasSuffix(s, suffix))
	return nil
}

// upper"" ( str -- str )
func upper(vm *forth.VM) error {
	sVal, err := vm.Pop()
	if err != nil {
		return err
	}
	s, ok := sVal.(string)
	if !ok {
		return forth.ErrArgumentMsg("upper\"\" expects a string")
	}
	vm.Push(strings.ToUpper(s))
	return nil
}

// lower"" ( str -- str )
func lower(vm *forth.VM) error {
	sVal, err := vm.Pop()
	if err != nil {
		return err
	}
	s, ok := sVal.(string)
	if !ok {
		return forth.ErrArgumentMsg("lower\"\" expects a string")
	}
	vm.Push(strings.ToLower(s))
	return nil
}

// split"" ( str sep -- arr )
func split(vm *forth.VM) error {
	sepVal, err := vm.Pop()
	if err != nil {
		return err
	}
	sep, ok := sepVal.(string)
	if !ok {
		return forth.ErrArgumentMsg("split\"\" expects a separator string")
	}

	sVal, err := vm.Pop()
	if err != nil {
		return err
	}
	s, ok := sVal.(string)
	if !ok {
		return forth.ErrArgumentMsg("split\"\" expects a target string")
	}

	parts := strings.Split(s, sep)
	vm.Push(parts)
	return nil
}

// @join"" ( arr sep -- str )
func arrayJoin(vm *forth.VM) error {
	sepVal, err := vm.Pop()
	if err != nil {
		return err
	}
	sep, ok := sepVal.(string)
	if !ok {
		return forth.ErrArgumentMsg("@join\"\" expects a separator string")
	}

	arrVal, err := vm.Pop()
	if err != nil {
		return err
	}

	if v, ok := arrVal.(*forth.Variable); ok {
		arrVal = v.Value()
	}

	var parts []string
	switch a := arrVal.(type) {
	case []string:
		parts = a
	case []any:
		for _, item := range a {
			if s, ok := item.(string); ok {
				parts = append(parts, s)
			} else {
				return forth.ErrArgumentMsg("@join\"\" array contains non-strings")
			}
		}
	default:
		return forth.ErrArgumentMsg("@join\"\" expects an array of strings")
	}

	vm.Push(strings.Join(parts, sep))
	return nil
}

// <join""> ( str1 ... strn n sep -- str )
func stackJoin(vm *forth.VM) error {
	sepVal, err := vm.Pop()
	if err != nil {
		return err
	}
	sep, ok := sepVal.(string)
	if !ok {
		return forth.ErrArgumentMsg("<join\"\"> expects a separator string")
	}

	nVal, err := vm.Pop()
	if err != nil {
		return err
	}
	n, ok := nVal.(int64)
	if !ok {
		return forth.ErrArgumentMsg("<join\"\"> expects an integer count n")
	}
	if n < 0 {
		return forth.ErrArgumentMsg("<join\"\"> count cannot be negative")
	}

	parts := make([]string, n)
	for i := n - 1; i >= 0; i-- {
		sVal, err := vm.Pop()
		if err != nil {
			return err
		}
		str, ok := sVal.(string)
		if !ok {
			return forth.ErrArgumentMsg("<join\"\"> expected string on stack")
		}
		parts[i] = str
	}

	vm.Push(strings.Join(parts, sep))
	return nil
}

// sub"" ( str idx1 idx2 -- substr )
func substring(vm *forth.VM) error {
	idx2Val, err := vm.Pop()
	if err != nil {
		return err
	}
	idx2, ok := idx2Val.(int64)
	if !ok {
		return forth.ErrArgumentMsg("sub\"\" expects integer idx2")
	}

	idx1Val, err := vm.Pop()
	if err != nil {
		return err
	}
	idx1, ok := idx1Val.(int64)
	if !ok {
		return forth.ErrArgumentMsg("sub\"\" expects integer idx1")
	}

	sVal, err := vm.Pop()
	if err != nil {
		return err
	}
	s, ok := sVal.(string)
	if !ok {
		return forth.ErrArgumentMsg("sub\"\" expects a string")
	}

	if idx1 != 0 && idx2 == 0 {
		idx2 = int64(len(s))
	}

	if idx1 < 0 {
		idx1 = int64(len(s)) + idx1
	}
	if idx2 < 0 {
		idx2 = int64(len(s)) + idx2
	}

	if idx1 < 0 {
		idx1 = 0
	}
	if idx1 > int64(len(s)) {
		idx1 = int64(len(s))
	}

	if idx2 < 0 {
		idx2 = 0
	}
	if idx2 > int64(len(s)) {
		idx2 = int64(len(s))
	}

	if idx1 > idx2 {
		vm.Push("")
	} else {
		vm.Push(s[idx1:idx2])
	}
	return nil
}

// replace"" ( str old new -- str )
func replace(vm *forth.VM) error {
	newVal, err := vm.Pop()
	if err != nil {
		return err
	}
	newStr, ok := newVal.(string)
	if !ok {
		return forth.ErrArgumentMsg("replace\"\" expects string replacement")
	}

	oldVal, err := vm.Pop()
	if err != nil {
		return err
	}
	oldStr, ok := oldVal.(string)
	if !ok {
		return forth.ErrArgumentMsg("replace\"\" expects string old pattern")
	}

	sVal, err := vm.Pop()
	if err != nil {
		return err
	}
	s, ok := sVal.(string)
	if !ok {
		return forth.ErrArgumentMsg("replace\"\" expects a target string")
	}

	vm.Push(strings.ReplaceAll(s, oldStr, newStr))
	return nil
}

// index"" ( str substr -- index )
func indexStr(vm *forth.VM) error {
	subVal, err := vm.Pop()
	if err != nil {
		return err
	}
	sub, ok := subVal.(string)
	if !ok {
		return forth.ErrArgumentMsg("index\"\" expects a substring")
	}

	sVal, err := vm.Pop()
	if err != nil {
		return err
	}
	s, ok := sVal.(string)
	if !ok {
		return forth.ErrArgumentMsg("index\"\" expects a target string")
	}

	vm.Push(int64(strings.Index(s, sub)))
	return nil
}

// contains""? ( str substr -- bool )
func contains(vm *forth.VM) error {
	subVal, err := vm.Pop()
	if err != nil {
		return err
	}
	sub, ok := subVal.(string)
	if !ok {
		return forth.ErrArgumentMsg("contains\"\"? expects a substring")
	}

	sVal, err := vm.Pop()
	if err != nil {
		return err
	}
	s, ok := sVal.(string)
	if !ok {
		return forth.ErrArgumentMsg("contains\"\"? expects a target string")
	}

	pushBool(vm, strings.Contains(s, sub))
	return nil
}
