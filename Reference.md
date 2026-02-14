# FORTH Words Reference

This document lists all the FORTH words currently defined in the Go.Forth interpreter, along with their descriptions and examples.

Words are grouped by category.

## Stack Manipulation Words

- **dup**: : dup ( a -- a a ) <code>  
  Duplicates the top item on the stack.  
  Example: `5 dup .s` → stack: 5 5

- **over**: : over swap dup -rot ;  
  Copies the second item to the top.  
  Example: `1 2 over .s` → stack: 1 2 1

- **drop**: : drop ( a -- ) <code>  
  Removes the top item from the stack.  
  Example: `1 2 drop .s` → stack: 1

- **swap**: : swap ( a b -- b a ) <code>  
  Swaps the top two items.  
  Example: `1 2 swap .s` → stack: 2 1

- **rot**: : rot ( a b c -- b c a ) <code>  
  Rotates the top three items.  
  Example: `1 2 3 rot .s` → stack: 2 3 1

- **-rot**: : -rot rot rot ;  
  Rotates the top three items in the opposite direction.  
  Example: `1 2 3 -rot .s` → stack: 3 1 2

- **nip**: : nip swap drop ;  
  Removes the second item from the stack.  
  Example: `1 2 nip .s` → stack: 2

- **tuck**: : tuck swap over ;  
  Inserts a copy of the top item under the second item.  
  Example: `1 2 tuck .s` → stack: 2 1 2

- **>r**: >r push onto rstack  
  Moves the top item from the data stack to the return stack.  
  Example: `5 >r .s` → data stack empty, return stack: 5

- **r>**: r> pop from rstack  
  Moves the top item from the return stack to the data stack.  
  Example: `r> .s` → data stack: 5

- **r@**: r@ peek at rstack  
  Copies the top item from the return stack to the data stack.  
  Example: `r@ .s` → data stack: 5, return stack: 5

- **rdrop**: rdrop  
  Removes the top item from the return stack.  
  Example: `rdrop` → return stack empty

## Numeric Operations

- **+**: : + ( a b -- a+b ) <code>  
  Adds two numbers (int, float64, or strings).  
  Example: `3 4 + .` → 7

- **-**: : - ( a b -- a-b ) <code>  
  Subtracts two numbers.  
  Example: `10 3 - .` → 7

- **\*** : : * ( a b -- a*b ) <code>  
  Multiplies two numbers or repeats a string.  
  Example: `3 4 * .` → 12  
  Example: `"hello" 3 * .` → hellohellohello

- **/**: : / ( a b -- a/b ) <code>  
  Divides two numbers.  
  Example: `10 2 / .` → 5

- **sqrt**: : sqrt ( a -- sqrt(a) ) <code>  
  Computes the square root.  
  Example: `16 sqrt .` → 4

- **log**: : log ( a -- log(a) ) <code>  
  Computes the natural logarithm.  
  Example: `2.718 log .` → 1

- **log10**: : log10 ( a -- log10(a) ) <code>  
  Computes the base-10 logarithm.  
  Example: `100 log10 .` → 2

- **log2**: : log2 ( a -- log2(a) ) <code>  
  Computes the base-2 logarithm.  
  Example: `8 log2 .` → 3

- **max**: : max ( a b -- max(a,b) ) <code>  
  Returns the maximum of two numbers.  
  Example: `5 3 max .` → 5

- **min**: : min ( a b -- min(a,b) ) <code>  
  Returns the minimum of two numbers.  
  Example: `5 3 min .` → 3

- **sin**: : sin ( a -- sin(a) ) <code>  
  Computes the sine (input in radians).  
  Example: `0 sin .` → 0

- **cos**: : cos ( a -- cos(a) ) <code>  
  Computes the cosine (input in radians).  
  Example: `0 cos .` → 1

- **tan**: : tan ( a -- tan(a) ) <code>  
  Computes the tangent (input in radians).  
  Example: `0 tan .` → 0

- **round**: : round ( a -- round(a) ) <code>  
  Rounds to the nearest integer.  
  Example: `3.6 round .` → 4

- **floor**: : floor ( a -- floor(a) ) <code>  
  Rounds down to the nearest integer.  
  Example: `3.6 floor .` → 3

- **ceil**: : ceil ( a -- ceil(a) ) <code>  
  Rounds up to the nearest integer.  
  Example: `3.2 ceil .` → 4

## Comparison and Logical Words

System uses -1 for TRUE and 0 for FALSE. Comparisons work on ints, floats, and strings.

- **=**: : = ( a b -- bool ) <code>  
  Checks for equality.  
  Example: `1 1 = .` → -1

- **<**: : < ( a b -- bool ) <code>  
  Checks less than.  
  Example: `1 2 < .` → -1

- **>**: : > ( a b -- bool ) <code>  
  Checks greater than.  
  Example: `2 1 > .` → -1

- **<=**: : <= ( a b -- bool ) <code>  
  Checks less than or equal.  
  Example: `2 2 <= .` → -1

- **>=**: : >= ( a b -- bool ) <code>  
  Checks greater than or equal.  
  Example: `2 1 >= .` → -1

- **<>**: : <> ( a b -- bool ) <code>  
  Checks inequality.  
  Example: `1 2 <> .` → -1

- **0=**: : 0= ( a -- bool ) <code>  
  Checks if equal to zero (logical NOT).  
  Example: `0 0= .` → -1

- **0<**: : 0< ( a -- bool ) <code>  
  Checks if less than zero.  
  Example: `-1 0< .` → -1

- **0>**: : 0> ( a -- bool ) <code>  
  Checks if greater than zero.  
  Example: `1 0> .` → -1

- **and**: : and ( a b -- c ) <code>  
  Bitwise AND.  
  Example: `-1 1 and .` → 1

- **or**: : or ( a b -- c ) <code>  
  Bitwise OR.  
  Example: `1 2 or .` → 3

- **xor**: : xor ( a b -- c ) <code>  
  Bitwise XOR.  
  Example: `1 3 xor .` → 2

- **invert**: : invert ( a -- b ) <code>  
  Bitwise NOT.  
  Example: `0 invert .` → -1

- **true**: : true ( -- -1 ) <code>  
  Pushes constant -1 (TRUE).  

- **false**: : false ( -- 0 ) <code>  
  Pushes constant 0 (FALSE).

## Input/Output Words

- **read**: looks at the top of the stack, and tries to interpret it as a rune. If it can, then it reads until it finds that rune, and leaves the string it read at the top of the stack. A special case is when the delimiter is a space, in which case it reads until any whitespace is found.  
  Example: `32 read` (reads until space)

- **skip**: : skip ( delim -- ) read drop ;  
  Skips input until the delimiter.  
  Example: `32 skip` (skips until space)

- **"**: : " 34 read (compiling?) if postpone literal then ; immediate  
  Parses a string literal.  
  Example: `"hello" type` → hello

- **chr**: chrFromInt ('chr') takes an integer and makes a one-char string of it, interpreted as a rune  
  Converts int to character.  
  Example: `65 chr .` → A

- **ord**: ordFromStr ('ord') takes a one-character string and gives its rune value as an int. It is the inverse of 'chr'.  
  Converts character to int.  
  Example: `"A" ord .` → 65

- **.s**: printStack prints out the stack contents, without removing anything.  
  Prints the entire stack.  
  Example: `1 2 3 .s` → 3: 3 2: 2 1: 1

- **.**: printTop prints out the top element on the stack, removing it in the process. It puts a trailing space after the item.  
  Prints and removes top of stack.  
  Example: `42 .` → 42

- **type**: printOut ('type') prints out the top element on the stack, removing it in the process. It does not include a trailing space.  
  Prints top of stack without space.  
  Example: `42 type` → 42 (no trailing space)

- **cr**: cr simply prints a carriage return  
  Prints a newline.  
  Example: `cr` → (newline)

## Parser/Compilation Words

- **\\**: nlComment '\' skips until the next newline : \ '\n' skip ; immediate  
  Line comment.  
  Example: `\ This is a comment`

- **(**: parenComment '(' skips until the closing paren. : ( ')' skip ; immediate  
  Block comment.  
  Example: `( This is a comment )`

- **[**: Interpret sets the compilation state of the VM to false, and reads words one at a time...  
  Switches to interpretation mode.  
  Example: `[ 1 2 + ]` (interprets the code)

- **]**: stopInterpret completes an interpretation and falls back to the compiler  
  Switches back to compilation mode.  
  Example: `]` (after [ )

- **:**: compile (':') reads the name of a word to define, and then compiles the definition until ';' tells it to stop  
  Starts word definition.  
  Example: `: square dup * ;`

- **;**: stopCompile (';') terminates a compilation  
  Ends word definition.  
  Example: `;` (after : definition)

- **literal**: literal is an immediate word that reads an int from the stack and compiles it into the codestream  
  Compiles a literal value.  
  Example: `42 literal`

- **postpone**: postpone creates code that compiles code into the caller. For immediates, it creates code that calls code in the caller.  
  Postpones compilation of a word.  
  Example: `: test postpone + ; immediate`

- **immediate**: func makeImmediate ('immediate') makes the last defined word immediate  
  Makes the last defined word immediate.  
  Example: `: myword ... ; immediate`

- **compile,**: compileComma  
  Compiles an opcode (integer index) directly into the code stream.  
  Example: `10 compile,`

- **'**: tick ( ' name -- xt )  
  Finds the word in the dictionary and pushes its execution token (index) to the stack.  
  Example: `' dup`

- **[']**: bracketTick ( ['] name -- )  
  Compiles the execution token of the following word as a literal. Immediate.  
  Example: `: get-dup ['] dup ;`

- **execute**: execute ( xt -- )  
  Executes the word execution token on the stack.  
  Example: `' dup execute`

## Control Flow Words

- **if**: IF is an immediate word that stores a fixup address on the stack for ELSE / THEN to find, and stores a (bzr) with a dummy branch amount in the code stream.  
  Conditional execution.  
  Example: `1 if 42 . then`

- **else**: ELSE needs to issue a jump over the else-stuff, and then use opThen to fixup the IF to jump into the else-stuff. Finally, it needs to leave a fixup location on the stack for the final THEN.  
  Alternative branch in if-then.  
  Example: `0 if 1 . else 2 . then`

- **then**: THEN takes a fixup address from the stack, and inserts the right amount to jump over the IF (or ELSE) block. No new code is added to the codestream.  
  Ends if-then-else block.  
  Example: `then` (see if/else examples)

- **recur**: RECUR just jumps to the start of the current function  
  Recurses to start of current word.  
  Example: `: countdown dup . 1 - dup 0 > if recur then drop ;`

- **do**: limit start DO <body> LOOP/+LOOP defines a basic for-style loop. It needs to stash away the limit and current index on the R-stack prior to the loop proper. Then, at the start of the loop, it needs to test whether iteration should continue, or jump to the end of the loop: >r >r (test loop-body back-facing branch) rdrop rdrop  
  Starts a definite loop.  
  Example: `10 0 do i . loop`

- **loop**: opLoop  
  Ends a do-loop, incrementing by 1.  
  Example: `loop` (see do example)

- **+loop**: opLoopPlus  
  Ends a do-loop, incrementing by the top of stack.  
  Example: `2 +loop`

- **i**: getDoI gets the index of the current loop  
  Pushes current loop index.  
  Example: `i .`

- **j**: getDoJ gets the index of the outer loop in nested DO loops
  Pushes outer loop index in nested loops.
  Example: `j .`

## Array Words

Arrays are dynamic Go slices supporting bytes, ints, floats, and strings. `@` and `!` provide unified access with type coercion, while `c@` and `c!` are byte-specific with value wrapping.

- **bytes** (size -- []byte): Create a byte array of given size
  Example: `10 bytes` creates a 10-element byte array

- **ints** (size -- []int): Create an int array of given size
  Example: `5 ints` creates a 5-element int array

- **floats** (size -- []float64): Create a float64 array of given size
  Example: `3 floats` creates a 3-element float array

- **strings** (size -- []string): Create a string array of given size
  Example: `4 strings` creates a 4-element string array

- **@** (array index -- value): Get element at index with lossless type coercion
  Example: `myarray 0 @ .` gets first element

- **!** (value array index -- ): Set element at index with type checking - allows int→float64 but rejects float64→int
  Example: `42 myarray 0 !` sets first element to 42

- **c@** (array index -- byte): Get byte at index with auto-wrap (value & 0xff)
  Example: `mybytes 0 c@ .` gets byte value wrapped to 0-255

- **c!** (value array index -- ): Set byte at index with auto-wrap (value & 0xff)
  Example: `300 mybytes 0 c!` sets byte to 44 (300 & 0xff)

- **@push** (value array -- array'): Append value to array and return modified array
  Example: `99 myarray @push` appends 99 and returns new array

- **@pop** (array -- array' value): Remove and return last element and modified array
  Example: `myarray @pop .` pops last element and prints it

- **@shift** (array -- array' value): Remove and return first element and modified array
  Example: `myarray @shift .` shifts first element and prints it

- **@unshift** (value array -- array'): Prepend value to array and return modified array
  Example: `0 myarray @unshift` prepends 0 and returns new array

- **@len** (array -- length): Get array length
  Example: `myarray @len .` prints array length

## Variable Words

Variables provide named storage for any Go value, similar to traditional FORTH variables. They integrate seamlessly with array operations for automatic updates.

- **variable** ( "name" -- ): Create a new global variable initialized to 0
  Example: `variable x` creates variable x

- **@** (variable -- value): Get variable value, or get element from array in variable
  Example: `x @ .` prints the value of x
  Example: `x 0 @ .` gets element 0 from array in x

- **!** (value variable -- ): Set variable value, or set element in array in variable
  Example: `42 x !` sets x to 42
  Example: `99 x 1 !` sets element 1 in array in x to 99

### Variable Auto-Update with Arrays

When array operations are performed on variables containing slices, the variable is automatically updated with the modified array. This eliminates manual variable management:

- **@push** (value variable -- ): Append to array in variable
  Example: `99 x @push` appends 99 to array in x (x auto-updated)

- **@pop** (variable -- value): Pop from array in variable
  Example: `x @pop .` pops from array in x (x auto-updated, value printed)

- **@shift** (variable -- value): Shift from array in variable
  Example: `x @shift .` shifts from array in x (x auto-updated, value printed)

- **@unshift** (value variable -- ): Unshift to array in variable
  Example: `0 x @unshift` prepends 0 to array in x (x auto-updated)

- **@len** (variable -- length): Get length of array in variable
  Example: `x @len .` prints length of array in x

For manual control, extract the array first: `x @ @push` (standard array behavior).

## Dictionary Words

Dictionaries provide key-value storage with string keys, supporting any Go value types. They offer efficient lookups and modifications with comprehensive error handling.

- **empty-dict** ( -- dict ): Create a new empty dictionary
  Example: `empty-dict` creates an empty dictionary

- **d!** ( value dict key -- ): Set a value in the dictionary
  Example: `42 empty-dict " age" d!` sets "age" to 42

- **d@** ( dict key -- value ): Get a value from the dictionary
  Example: `my-dict " age" d@ .` prints the age value

- **ddel** ( dict key -- ): Delete a key from the dictionary
  Example: `my-dict " age" ddel` removes the "age" key

- **dkeys** ( dict -- keys ): Get all keys from the dictionary as an array
  Example: `my-dict dkeys @len .` prints number of keys

- **d@|** ( dict key default -- value ): Get value or default if key missing
  Example: `my-dict " missing" 0 d@| .` prints 0 if "missing" not found

- **d@?** ( dict key -- value -1 | 0 ): Query value with existence flag
  Example: `my-dict " age" d@?` returns value and -1 if found, or 0 if not found

- **mark**: mark  
  Sets a marker in the dictionary for a later `forget`.  
  Example: `mark`

- **forget**: forget  
  Forgets all definitions defined since the last `mark`.  
  Example: `forget`

## Internal Words

These are typically not used directly but are part of the implementation:

- **(setupDo)**: setupDo  
  Sets up the DO loop internals.

- **(testDo)**: testDo  
  Tests the DO loop condition.

- **(perfLoopPlus)**: (perfLoopPlus) ( amt -- )  
  Performs the loop increment.

- **debug.**: debugPrint  
  Prints the current code segment for debugging purposes.  
  Example: `debug.`