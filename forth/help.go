// SPDX-License-Identifier: MIT

package forth

import (
	"fmt"
	"regexp"
	"sort"
	"strings"
)

// RegisterHelp registers help text for a category
func RegisterHelp(category, helpText string) {
	registryLock.Lock()
	defer registryLock.Unlock()
	helpRegistry[category] = helpText
}

// GetHelpCategories returns a sorted list of all help categories
func GetHelpCategories() []string {
	registryLock.RLock()
	defer registryLock.RUnlock()

	categories := make([]string, 0, len(helpRegistry))
	for category := range helpRegistry {
		categories = append(categories, category)
	}
	sort.Strings(categories)
	return categories
}

// GetHelpText returns the help text for a category
func GetHelpText(category string) (string, bool) {
	registryLock.RLock()
	defer registryLock.RUnlock()

	text, exists := helpRegistry[category]
	return text, exists
}

// helpText implements the (help-text) word
// ( category -- help-text )
func helpText(vm *VM) error {
	catVal, err := vm.Pop()
	if err != nil {
		return err
	}
	category, ok := catVal.(string)
	if !ok {
		return ErrArgumentMsg("(help-text) expects a string category")
	}

	helpText, exists := GetHelpText(strings.ToLower(category))
	if !exists {
		return fmt.Errorf("help category '%s' not found", category)
	}

	vm.Push(helpText)
	return nil
}

// help implements the :h word
// :h reads from stdin: category token, then rest of line as regex
func help(vm *VM) error {
	// Read category token from input source
	var err error

	vm.Push("What category (leave empty for the list): ")
	err = printStr(vm)
	if err != nil {
		return err
	}
	err = readLine(vm)
	if err != nil {
		return err
	}

	givenstrV, err := vm.Pop()
	if err != nil {
		return err
	}

	givenstr, ok := givenstrV.(string)
	if !ok {
		return ErrArgumentMsg(":h requires a string!")
	}

	givenstr = strings.TrimSpace(givenstr)
	if len(givenstr) == 0 {
		categories := GetHelpCategories()
		for _, cat := range categories {
			fmt.Fprintln(vm.Sink, cat)
		}
		return nil
	}

	// Split givenstr into category and optional regex
	parts := strings.Fields(givenstr)
	category := strings.ToLower(parts[0])
	var regexStr string
	if len(parts) > 1 {
		regexStr = strings.Join(parts[1:], " ")
	}

	// Get help text
	helpText, exists := GetHelpText(category)
	if !exists {
		return fmt.Errorf("help category '%s' not found", category)
	}

	if regexStr == "" {
		// No regex: show all paragraphs
		fmt.Println(vm.Sink, helpText)
	} else {
		paragraphs := strings.Split(helpText, "\n")

		// Has regex: filter paragraphs
		re, err := regexp.Compile(regexStr)
		if err != nil {
			return fmt.Errorf("invalid regex pattern: %v", err)
		}

		found := false
		for _, para := range paragraphs {
			para = strings.TrimSpace(para)
			if para != "" && re.MatchString(para) {
				fmt.Fprintln(vm.Sink, para)
				found = true
			}
		}

		if !found {
			fmt.Fprintf(vm.Sink, "No help text matches pattern '%s'\n", regexStr)
		}
	}

	return nil
}

func helpWordsInit(vm *VM) {
	vm.Define(&NativeWord{name: ":h", run: help, immediate: false})
	vm.Define(&NativeWord{name: "(help-text)", run: helpText, immediate: false})
}
