// SPDX-License-Identifier: MIT

package forth

import (
	"fmt"
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
// :h reads from stdin: category token
func help(vm *VM) error {
	// Read category token from input source
	var err error

	givenstr, err := nextToken(vm, nil)
	if err != nil {
		return err
	}
	givenstr = strings.ToLower(strings.TrimSpace(givenstr))
	if givenstr == "-" {
		categories := GetHelpCategories()
		for _, cat := range categories {
			fmt.Fprintln(vm.Sink, cat)
		}
		return nil
	}

	// Get help text
	helpText, exists := GetHelpText(givenstr)
	if !exists {
		return fmt.Errorf("help category '%s' not found (use category `-` to see all)", givenstr)
	}

	fmt.Fprintln(vm.Sink, helpText)
	return nil
}

func init() {
	// Register base help categories
	RegisterHelp("core", `swap ( a b -- b a ) exchanges the top two items on the stack.
dup ( a -- a a ) duplicates the top item on the stack.
drop ( a -- ) removes the top item from the stack.
over ( a b -- a b a ) copies the second item to the top.
rot ( a b c -- b c a ) rotates the top three items.
+ ( n1 n2 -- n3 ) adds two numbers.
- ( n1 n2 -- n3 ) subtracts n2 from n1.
* ( n1 n2 -- n3 ) multiplies two numbers.
/ ( n1 n2 -- n3 ) divides n1 by n2.
** ( n1 n2 -- n3 ) raises n1 to the power of n2.
sqrt ( n -- sqrt(n) ) computes square root.
log ( n -- log(n) ) computes natural logarithm.
log10 ( n -- log10(n) ) computes base-10 logarithm.
log2 ( n -- log2(n) ) computes base-2 logarithm.
max ( n1 n2 -- max ) returns the maximum of two numbers.
min ( n1 n2 -- min ) returns the minimum of two numbers.
sin ( n -- sin(n) ) computes sine.
cos ( n -- cos(n) ) computes cosine.
tan ( n -- tan(n) ) computes tangent.
round ( n -- round(n) ) rounds to nearest integer.
floor ( n -- floor(n) ) rounds down to integer.
ceil ( n -- ceil(n) ) rounds up to integer.
= ( a b -- bool ) tests if two values are equal.
< ( n1 n2 -- bool ) tests if n1 is less than n2.
> ( n1 n2 -- bool ) tests if n1 is greater than n2.
and ( n1 n2 -- n3 ) bitwise AND of two numbers.
or ( n1 n2 -- n3 ) bitwise OR of two numbers.
xor ( n1 n2 -- n3 ) bitwise XOR of two numbers.
not ( n -- ~n ) bitwise NOT of a number.
if ( bool -- ) begins a conditional block.
else ( -- ) provides alternative in conditional.
then ( -- ) ends a conditional block.
begin ( -- ) starts an indefinite loop.
until ( bool -- ) ends loop when condition is true.
do ( limit index -- ) starts a counted loop.
loop ( -- ) ends a counted loop.
i ( -- index ) gets current loop index.
j ( -- index ) gets outer loop index.
variable ( val name -- ) creates a named variable.
constant ( val name -- ) creates a named constant.
@ ( var -- val ) gets value from variable.
! ( val var -- ) sets value in variable.
execute ( xt -- ) executes an execution token.
catch? ( ?? n xt -- ?? bool ) executes xt catching errors.
throw ( str -- ) throws an exception with message.
mark ( -- ) marks current VM state.
forget ( -- ) forgets back to last mark.
debug. ( -- ) prints compiled code for debugging.`)

	RegisterHelp("containers", `things ( size -- arr ) creates an array of any values.
ints ( size -- arr ) creates an array of integers.
floats ( size -- arr ) creates an array of floats.
strings ( size -- arr ) creates an array of strings.
bytes ( size -- arr ) creates an array of bytes.
@push ( val arr -- arr ) adds value to end of array.
@pop ( arr -- arr val ) removes and returns last value.
@ ( arr idx -- val ) gets value at index from array or variable.
! ( val arr idx -- ) sets value at index in array or variable.
@len ( arr -- len ) gets array length.
empty-dict ( -- d ) creates a new empty dictionary.
d@ ( d key -- val ) gets value for key.
d! ( val d key -- ) sets value for key.
d@? ( d key -- val bool ) gets value and existence flag.
dkeys ( d -- keys ) gets all keys as array.
<< ( -- ) starts variable-length argument collection.
>> ( -- count ) ends collection, pushes item count.
<@push> ( arr item1...itemN count -- arr ) pushes multiple items to array.
@spread ( arr -- item1...itemN count ) spreads array to stack.`)

	RegisterHelp("io", `type ( str -- ) prints string without newline.
. ( val -- ) prints value with space.
cr ( -- ) prints a newline.
.s ( -- ) prints stack contents.
read-line ( -- str ) reads line from input.
read ( delim -- str ) reads until delimiter.
parse ( delim -- str ) parses from input source.
" ( -- str ) parses string literal.
chr ( n -- str ) converts number to character.
ord ( str -- n ) converts character to number.
skip-parse ( delim -- ) skips input until delimiter.`)

	RegisterHelp("parsing", `: ( name -- ) starts word definition.
; ( -- ) ends word definition.
immediate ( -- ) makes last word immediate.
[[ ( -- ) enters interpretation mode.
]] ( -- ) leaves interpretation mode.
literal ( val -- ) compiles literal value.
compile-xt ( xt -- ) compiles execution token.
postpone ( -- ) defers compilation of next word.
(| ( -- ) starts local variable declaration.
|) ( -- ) ends local variable declaration.
read-token ( -- str ) reads next token.
lookup-xt ( str -- xt ) finds execution token.
eval ( str -- ) evaluates string.
load ( filename -- ) loads and executes file.`)

	RegisterHelp("extensions", `extension-list ( -- arr ) gets list of registered extensions.
<activate-extensions> ( names... count -- ) activates extensions.
:h ( -- ) shows help system. Can give a category or "-" to list categories`)
}

func helpWordsInit(vm *VM) {
	vm.Define(&NativeWord{name: ":h", run: help, immediate: false})
	vm.Define(&NativeWord{name: "(help-text)", run: helpText, immediate: false})
}
