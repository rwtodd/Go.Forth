package regex

import (
	"regexp"
	"strings"
	"unicode"

	"github.com/rwtodd/Go.Forth/forth"
)

func init() {
	forth.RegisterExtensionWithHelp("regex", func(vm *forth.VM) error {
		vm.Define(forth.NewNativeWord("rx:", rxString))
		vm.Define(forth.NewNativeWord("rx-compile", rxCompile))
		vm.Define(forth.NewNativeWord("rx-match?", rxMatch))
		vm.Define(forth.NewNativeWord("rx-gsub", rxGsub))
		vm.Define(forth.NewNativeWord("rx-sub", rxSub))
		vm.Define(forth.NewNativeWord("rx-split", rxSplit))
		vm.Define(forth.NewNativeWord("rx-gmatch?", rxMatchXt))
		vm.Define(forth.NewNativeWord("rx-find", rxFind))
		vm.Define(forth.NewNativeWord("rx-gfind", rxGfind))
		vm.Define(forth.NewImmediateWord("rx-of", rxOf))

		// [rx:] validates compiling a Regex string. Equivalent to `[[ rx: /.../ ]] literal`
		return vm.Eval(`: [rx:] rx: [[ <<" literal ">> ]] <postpone> ; immediate`)
	}, `rx: ( -- str ) parses a regex string from the input stream.
rx-compile ( pattern -- rx ) compiles a regex pattern string into a regex object.
rx-match? ( string pattern -- False | groups start-index True ) matches string against pattern, returns groups array, start index, and true if matched.
rx-gsub ( string pattern replacement -- string ) globally replaces all matches of pattern in string with replacement.
rx-sub ( string pattern replacement -- string ) replaces the first match of pattern in string with replacement.
rx-split ( string pattern -- array ) splits string by pattern into an array of strings.
rx-gmatch? ( string pattern xt -- count ) executes xt for each match, passing index and groups array, returns match count.
rx-find ( string pattern -- string ) returns the first substring that matches the pattern.
rx-gfind ( string pattern -- array ) returns all substrings that match the pattern.
[rx:] ( -- ) immediate word that compiles a regex literal.
rx-of ( #of -- orig #of+1 / x -- ) immediate word for regex pattern matching in loops.`)
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

// rx-match? ( string pattern -- False | groups start-index True )
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

	matches := re.FindStringSubmatchIndex(str)
	if matches == nil {
		vm.Push(int64(0)) // False in Forth
		return nil
	}

	startIdx := matches[0]

	var groups []string
	for i := 0; i < len(matches); i += 2 {
		if matches[i] == -1 || matches[i+1] == -1 {
			groups = append(groups, "")
		} else {
			groups = append(groups, str[matches[i]:matches[i+1]])
		}
	}

	vm.Push(groups)
	vm.Push(int64(startIdx))
	vm.Push(int64(-1)) // True in Forth
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

// rx-gmatch? ( string pattern xt -- count )
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
			return forth.ErrArgumentMsg("rx-gmatch? expects an execution token")
		}
	}

	re, err := getPattern(vm, "rx-gmatch?")
	if err != nil {
		return err
	}

	strVal, err := vm.Pop()
	if err != nil {
		return err
	}
	str, ok := strVal.(string)
	if !ok {
		return forth.ErrArgumentMsg("rx-gmatch? expects a target string")
	}

	matches := re.FindAllStringSubmatchIndex(str, -1)
	if matches == nil {
		vm.Push(int64(0))
		return nil
	}

	for _, match := range matches {
		startIdx := match[0]

		// create array of subgroups
		var groups []string
		for i := 0; i < len(match); i += 2 {
			if match[i] == -1 || match[i+1] == -1 {
				groups = append(groups, "")
			} else {
				groups = append(groups, str[match[i]:match[i+1]])
			}
		}

		// push match info: index, groups array
		vm.Push(int64(startIdx))
		vm.Push(groups)

		// execute XT
		if err := xt.Run(vm); err != nil {
			return err
		}
	}

	vm.Push(int64(len(matches)))
	return nil
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

// rx-of ( #of -- orig #of+1 / x -- ) immediate
func rxOf(vm *forth.VM) error {
	val, err := vm.Pop()
	if err != nil {
		return err
	}
	count, ok := val.(int64)
	if !ok {
		return forth.ErrBadStateMsg("RX-OF expects a count on stack")
	}

	overTok, err := vm.LookupToken("over")
	if err != nil {
		return err
	}
	swapTok, err := vm.LookupToken("swap")
	if err != nil {
		return err
	}
	matchTok, err := vm.LookupToken("rx-match?")
	if err != nil {
		return err
	}
	ifTok, err := vm.LookupToken("if")
	if err != nil {
		return err
	}
	rotTok, err := vm.LookupToken("rot")
	if err != nil {
		return err
	}
	dropTok, err := vm.LookupToken("drop")
	if err != nil {
		return err
	}

	if err := overTok.Compile(vm); err != nil {
		return err
	}
	if err := swapTok.Compile(vm); err != nil {
		return err
	}
	if err := matchTok.Compile(vm); err != nil {
		return err
	}

	err = ifTok.Run(vm)
	if err != nil {
		return err
	}

	if err := rotTok.Compile(vm); err != nil {
		return err
	}
	if err := dropTok.Compile(vm); err != nil {
		return err
	}

	vm.Push(count + 1)
	return nil
}
