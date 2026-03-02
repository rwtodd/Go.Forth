package regex

import (
	"regexp"

	"github.com/rwtodd/Go.Forth/forth"
)

func init() {
	forth.RegisterExtension("regex", func(vm *forth.VM) error {
		vm.Define(forth.NewNativeWord("match?", match))
		return nil
	})
}

// match? ( string pattern -- bool )
func match(vm *forth.VM) error {
	patVal, err := vm.Pop()
	if err != nil {
		return err
	}
	pattern, ok := patVal.(string)
	if !ok {
		return forth.ErrArgumentMsg("match? expects a pattern string")
	}

	strVal, err := vm.Pop()
	if err != nil {
		return err
	}
	str, ok := strVal.(string)
	if !ok {
		return forth.ErrArgumentMsg("match? expects a target string")
	}

	matched, err := regexp.MatchString(pattern, str)
	if err != nil {
		return err
	}

	if matched {
		vm.Push(int64(-1)) // True in Forth
	} else {
		vm.Push(int64(0)) // False in Forth
	}
	return nil
}
