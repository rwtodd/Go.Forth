package regex

import (
	"regexp"
	"strings"
	"unicode"

	"github.com/rwtodd/Go.Forth/forth"
)

func init() {
	forth.RegisterExtension("regex", func(vm *forth.VM) error {
		vm.Define(forth.NewNativeWord("rx:", rxString))
		vm.Define(forth.NewNativeWord("rx-compile", rxCompile))
		vm.Define(forth.NewNativeWord("rx-match?", rxMatch))
		vm.Define(forth.NewNativeWord("rx-gsub", rxGsub))
		vm.Define(forth.NewNativeWord("rx-sub", rxSub))
		vm.Define(forth.NewNativeWord("rx-split", rxSplit))
		vm.Define(forth.NewNativeWord("rx-match", rxMatchXt))
		vm.Define(forth.NewNativeWord("rx-find", rxFind))
		vm.Define(forth.NewNativeWord("rx-gfind", rxGfind))

		// [rx:] validates compiling a Regex string. Equivalent to `[[ rx: /.../ ]] literal`
		return vm.Eval(`: [rx:] rx: [[ <<" literal ">> ]] <postpone> ; immediate`)
	})
}

func parseRegexString(vm *forth.VM) (string, error) {
	// skip whitespace
	var delim rune
	for {
		ch, _, err := vm.ReadRune()
		if err != nil {
			return "", err
		}
		if !unicode.IsSpace(ch) {
			delim = ch
			break
		}
	}

	// figure out end delimiter if it's a paired character
	endDelim := delim
	switch delim {
	case '(':
		endDelim = ')'
	case '[':
		endDelim = ']'
	case '{':
		endDelim = '}'
	case '<':
		endDelim = '>'
	}

	var buf strings.Builder
	for {
		ch, _, err := vm.ReadRune()
		if err != nil {
			return "", err
		}
		if ch == endDelim {
			break
		}
		buf.WriteRune(ch)
	}
	return buf.String(), nil
}

// rx: parses a regex string from the input stream
func rxString(vm *forth.VM) error {
	str, err := parseRegexString(vm)
	if err != nil {
		return err
	}
	vm.Push(str)
	return nil
}

// getPattern helper extracts pattern from stack. Returns *regexp.Regexp
func getPattern(vm *forth.VM, wordName string) (*regexp.Regexp, error) {
	patVal, err := vm.Pop()
	if err != nil {
		return nil, err
	}

	switch p := patVal.(type) {
	case *regexp.Regexp:
		return p, nil
	case string:
		return regexp.Compile(p)
	default:
		return nil, forth.ErrArgumentMsg(wordName + " expects a string or compiled regex")
	}
}

// rx-compile ( pattern -- rx )
func rxCompile(vm *forth.VM) error {
	patVal, err := vm.Pop()
	if err != nil {
		return err
	}
	pattern, ok := patVal.(string)
	if !ok {
		return forth.ErrArgumentMsg("rx-compile expects a pattern string")
	}

	re, err := regexp.Compile(pattern)
	if err != nil {
		return err
	}
	vm.Push(re)
	return nil
}

// rx-match? ( string pattern -- bool )
func rxMatch(vm *forth.VM) error {
	re, err := getPattern(vm, "rx-match?")
	if err != nil {
		return err
	}

	strVal, err := vm.Pop()
	if err != nil {
		return err
	}
	str, ok := strVal.(string)
	if !ok {
		return forth.ErrArgumentMsg("rx-match? expects a target string")
	}

	if re.MatchString(str) {
		vm.Push(int64(-1)) // True in Forth
	} else {
		vm.Push(int64(0)) // False in Forth
	}
	return nil
}

// rx-gsub ( string pattern replacement -- string )
func rxGsub(vm *forth.VM) error {
	replVal, err := vm.Pop()
	if err != nil {
		return err
	}
	repl, ok := replVal.(string)
	if !ok {
		return forth.ErrArgumentMsg("rx-gsub expects a replacement string")
	}

	re, err := getPattern(vm, "rx-gsub")
	if err != nil {
		return err
	}

	strVal, err := vm.Pop()
	if err != nil {
		return err
	}
	str, ok := strVal.(string)
	if !ok {
		return forth.ErrArgumentMsg("rx-gsub expects a target string")
	}

	vm.Push(re.ReplaceAllString(str, repl))
	return nil
}

// rx-sub ( string pattern replacement -- string )
func rxSub(vm *forth.VM) error {
	replVal, err := vm.Pop()
	if err != nil {
		return err
	}
	repl, ok := replVal.(string)
	if !ok {
		return forth.ErrArgumentMsg("rx-sub expects a replacement string")
	}

	re, err := getPattern(vm, "rx-sub")
	if err != nil {
		return err
	}

	strVal, err := vm.Pop()
	if err != nil {
		return err
	}
	str, ok := strVal.(string)
	if !ok {
		return forth.ErrArgumentMsg("rx-sub expects a target string")
	}

	loc := re.FindStringSubmatchIndex(str)
	if loc == nil {
		vm.Push(str) // no match, return original string
		return nil
	}

	var result []byte
	result = re.ExpandString(result, repl, str, loc)

	// Create the final substituted string
	finalStr := str[:loc[0]] + string(result) + str[loc[1]:]
	vm.Push(finalStr)
	return nil
}

// rx-split ( string pattern -- array )
func rxSplit(vm *forth.VM) error {
	re, err := getPattern(vm, "rx-split")
	if err != nil {
		return err
	}

	strVal, err := vm.Pop()
	if err != nil {
		return err
	}
	str, ok := strVal.(string)
	if !ok {
		return forth.ErrArgumentMsg("rx-split expects a target string")
	}

	parts := re.Split(str, -1)
	vm.Push(parts)
	return nil
}

// rx-match ( string pattern xt -- )
func rxMatchXt(vm *forth.VM) error {
	xtVal, err := vm.Pop()
	if err != nil {
		return err
	}
	xt, ok := xtVal.(forth.ExecutionToken)
	if !ok {
		// support integer index fallback like in engine.go
		if idx, isInt := xtVal.(int64); isInt {
			xt = forth.WordToken{Token: uint16(idx)}
		} else {
			return forth.ErrArgumentMsg("rx-match expects an execution token")
		}
	}

	re, err := getPattern(vm, "rx-match")
	if err != nil {
		return err
	}

	strVal, err := vm.Pop()
	if err != nil {
		return err
	}
	str, ok := strVal.(string)
	if !ok {
		return forth.ErrArgumentMsg("rx-match expects a target string")
	}

	matches := re.FindStringSubmatchIndex(str)
	if matches == nil {
		return nil // no match, do nothing
	}

	startIdx := matches[0]

	// create array of subgroups
	var groups []string
	for i := 0; i < len(matches); i += 2 {
		if matches[i] == -1 || matches[i+1] == -1 {
			groups = append(groups, "")
		} else {
			groups = append(groups, str[matches[i]:matches[i+1]])
		}
	}

	// push match info: index, groups array
	vm.Push(int64(startIdx))
	vm.Push(groups)

	// execute XT
	return xt.Run(vm)
}

// rx-find ( string pattern -- string )
func rxFind(vm *forth.VM) error {
	re, err := getPattern(vm, "rx-find")
	if err != nil {
		return err
	}

	strVal, err := vm.Pop()
	if err != nil {
		return err
	}
	str, ok := strVal.(string)
	if !ok {
		return forth.ErrArgumentMsg("rx-find expects a target string")
	}

	vm.Push(re.FindString(str))
	return nil
}

// rx-gfind ( string pattern -- array )
func rxGfind(vm *forth.VM) error {
	re, err := getPattern(vm, "rx-gfind")
	if err != nil {
		return err
	}

	strVal, err := vm.Pop()
	if err != nil {
		return err
	}
	str, ok := strVal.(string)
	if !ok {
		return forth.ErrArgumentMsg("rx-gfind expects a target string")
	}

	var res []string
	matches := re.FindAllString(str, -1)
	if matches != nil {
		res = matches
	} else {
		res = []string{}
	}
	vm.Push(res)
	return nil
}
