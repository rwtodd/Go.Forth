# FORTH Words Reference

This document lists all the FORTH words currently defined in the Go.Forth interpreter, along with their descriptions and examples.

Words are grouped by category.

## Stack Manipulation Words

- **dup**: : dup ( a -- a a )   
  Duplicates the top item on the stack.  
  Example: `5 dup .s` → stack: 5 5

- **over**: : over swap dup -rot ;  
  Copies the second item to the top.  
  Example: `1 2 over .s` → stack: 1 2 1

- **drop**: : drop ( a -- )   
  Removes the top item from the stack.  
  Example: `1 2 drop .s` → stack: 1

- **swap**: : swap ( a b -- b a )   
  Swaps the top two items.  
  Example: `1 2 swap .s` → stack: 2 1

- **rot**: : rot ( a b c -- b c a )   
  Rotates the top three items.  
  Example: `1 2 3 rot .s` → stack: 2 3 1

- **-rot**: : -rot rot rot ;  
  Rotates the top three items in the opposite direction.  
  Example: `1 2 3 -rot .s` → stack: 3 1 2

- **pick**: pick ( xu ... x1 x0 u -- xu ... x1 x0 xu )
  Copies the u-th item (0-indexed) to the top of the stack. `0 pick` is `dup`, `1 pick` is `over`.
  Example: `10 20 30 0 pick .s` → stack: 10 20 30 30

- **roll**: roll ( xu ... x1 x0 u -- ... x1 x0 xu )
  Rotates the u-th item to the top of the stack. `2 roll` is `rot`, `1 roll` is `swap`.
  Example: `10 20 30 2 roll .s` → stack: 20 30 10

- **-roll**: -roll ( xu ... x1 x0 u -- x0 xu ... x1 )
  Rotates the top of the stack to the u-th position. `2 -roll` is `-rot`.
  Example: `10 20 30 2 -roll .s` → stack: 30 10 20

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

- **depth**: depth ( -- n )  
  Pushes the number of items currently on the stack.  
  Example: `1 2 3 depth .` → 3

## Variable Length Arguments

Variable length arguments allow words to take an arbitrary number of items from the stack.
The system uses `<<` to mark the start, and `>>` to calculate the number of items that have been placed on the stack since the `<<`.

- **<<**: startVarLen ( -- )  
  Marks the start of a variable-length argument list by pushing the current stack depth to the return stack.  
  Example: `<< 1 2 3 ...`

- **>>**: endVarLen ( -- count )  
  Calculates the number of items pushed to the stack since the last `<<`.  
  Example: `<< 10 20 >> .` → 2

- **<<"**: varLenQuote ( -- string1 ... stringN count )  
  Parses whitespace-separated tokens until the `">>` token is found.
  Treats each token as a string literal.
  Compatible with interpretation and compilation.
  Example: `<<" one two three ">> .s` → stack: one two three 3

- **sprintf**: sprintf ( fmt args... count -- result )  
  Formats a string using the provided format string and arguments. Can be used with `<< ... >>`.  
  Example: `"val: %d" << 42 >> sprintf .` → val: 42

- **<@push>**: varLenPush ( array item1...itemN count -- array' )  
  Appends items to the array (or variable containing an array) provided before `<<`.  
  Example: `0 ints << 1 2 3 >> <@push> .s` → stack: [1 2 3]

- **<ints>**: varLenInts ( item1...itemN count -- []int )  
  Creates a new integer array from the stack items.  
  Example: `<< 10 20 >> <ints> .s` → stack: [10 20]

- **<floats>**: varLenFloats ( item1...itemN count -- []float64 )  
  Creates a new float array from the stack items.  
  Example: `<< 1.5 2.5 >> <floats> .s` → stack: [1.5 2.5]

- **<strings>**: varLenStrings ( item1...itemN count -- []string )  
  Creates a new string array from the stack items.  
  Example: `<< "a" "b" >> <strings> .s` → stack: [a b]

- **<bytes>**: varLenBytes ( item1...itemN count -- []byte )  
  Creates a new byte array from the stack items.  
  Example: `<< 65 66 >> <bytes> .s` → stack: [65 66]

- **<things>**: varLenAny ( item1...itemN count -- []any )  
  Creates a new generic array from the stack items.  
  Example: `<< 1 " a" >> <things> .s` → stack: [1 "a"]

## Numeric Operations

- **+**: : + ( a b -- a+b )   
  Adds two numbers (int, float64, or strings).  
  Example: `3 4 + .` → 7

- **-**: : - ( a b -- a-b )   
  Subtracts two numbers.  
  Example: `10 3 - .` → 7

- **\*** : : * ( a b -- a*b )   
  Multiplies two numbers or repeats a string.  
  Example: `3 4 * .` → 12  
  Example: `"hello" 3 * .` → hellohellohello

- **\*\*** : : \*\* ( a b -- a\*\*b )   
  raises a to the power of b
  Example: `2 3 \*\* .` →  8.0

- **/**: : / ( a b -- a/b )   
  Divides two numbers.  
  Example: `10 2 / .` → 5

- **sqrt**: : sqrt ( a -- sqrt(a) )   
  Computes the square root.  
  Example: `16 sqrt .` → 4

- **log**: : log ( a -- log(a) )   
  Computes the natural logarithm.  
  Example: `2.718 log .` → 1

- **log10**: : log10 ( a -- log10(a) )   
  Computes the base-10 logarithm.  
  Example: `100 log10 .` → 2

- **log2**: : log2 ( a -- log2(a) )   
  Computes the base-2 logarithm.  
  Example: `8 log2 .` → 3

- **max**: : max ( a b -- max(a,b) )   
  Returns the maximum of two numbers.  
  Example: `5 3 max .` → 5

- **min**: : min ( a b -- min(a,b) )   
  Returns the minimum of two numbers.  
  Example: `5 3 min .` → 3

- **sin**: : sin ( a -- sin(a) )   
  Computes the sine (input in radians).  
  Example: `0 sin .` → 0

- **cos**: : cos ( a -- cos(a) )   
  Computes the cosine (input in radians).  
  Example: `0 cos .` → 1

- **tan**: : tan ( a -- tan(a) )   
  Computes the tangent (input in radians).  
  Example: `0 tan .` → 0

- **round**: : round ( a -- round(a) )   
  Rounds to the nearest integer.  
  Example: `3.6 round .` → 4

- **floor**: : floor ( a -- floor(a) )   
  Rounds down to the nearest integer.  
  Example: `3.6 floor .` → 3

- **ceil**: : ceil ( a -- ceil(a) )   
  Rounds up to the nearest integer.  
  Example: `3.2 ceil .` → 4

## Comparison and Logical Words

System uses -1 for TRUE and 0 for FALSE. Comparisons work on ints, floats, and strings.

- **=**: : = ( a b -- bool )   
  Checks for equality.  
  Example: `1 1 = .` → -1

- **<**: : < ( a b -- bool )   
  Checks less than.  
  Example: `1 2 < .` → -1

- **>**: : > ( a b -- bool )   
  Checks greater than.  
  Example: `2 1 > .` → -1

- **<=**: : <= ( a b -- bool )   
  Checks less than or equal.  
  Example: `2 2 <= .` → -1

- **>=**: : >= ( a b -- bool )   
  Checks greater than or equal.  
  Example: `2 1 >= .` → -1

- **<>**: : <> ( a b -- bool )   
  Checks inequality.  
  Example: `1 2 <> .` → -1

- **0=**: : 0= ( a -- bool )   
  Checks if equal to zero (logical NOT).  
  Example: `0 0= .` → -1

- **0<**: : 0< ( a -- bool )   
  Checks if less than zero.  
  Example: `-1 0< .` → -1

- **0>**: : 0> ( a -- bool )   
  Checks if greater than zero.  
  Example: `1 0> .` → -1

- **and**: : and ( a b -- c )   
  Bitwise AND.  
  Example: `-1 1 and .` → 1

- **or**: : or ( a b -- c )   
  Bitwise OR.  
  Example: `1 2 or .` → 3

- **xor**: : xor ( a b -- c )   
  Bitwise XOR.  
  Example: `1 3 xor .` → 2

- **invert**: : invert ( a -- b )   
  Bitwise NOT.  
  Example: `0 invert .` → -1

- **true**: : true ( -- -1 )   
  Pushes constant -1 (TRUE).  

- **false**: : false ( -- 0 )   
  Pushes constant 0 (FALSE).

## Input/Output Words

- **read**: looks at the top of the stack, and tries to interpret it as a rune. If it can, then it reads until it finds that rune, and leaves the string it read at the top of the stack. A special case is when the delimiter is a space, in which case it reads until any whitespace is found.  
  Example: `32 read` (reads until space)

- **skip**: : skip ( delim -- ) read drop ;  
  Skips input until the delimiter.  
  Example: `32 skip` (skips until space)

- **"**: : " 34 escapedRead (compiling?) if postpone literal then ; immediate  
  Parses a string literal. Supports C-style escape sequences: `\n`, `\t`, `\r`, `\\`, `\"`.  
  Example: `"hello\nworld" type`

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

- **[[**: Interpret sets the compilation state of the VM to false, and reads words one at a time...  
  Switches to interpretation mode.  
  Example: `[[ 1 2 + ]]` (interprets the code)

- **]]**: stopInterpret completes an interpretation and falls back to the compiler  
  Switches back to compilation mode.  
  Example: `]]` (after [[ )

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

- **do**: limit start DO body LOOP/+LOOP defines a basic for-style loop. It needs to stash away the limit and current index on the R-stack prior to the loop proper. Then, at the start of the loop, it needs to test whether iteration should continue, or jump to the end of the loop: >r >r (test loop-body back-facing branch) rdrop rdrop  
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

- **leave**: leave
  Immediately exits the current loop (DO/LOOP). Execution continues after the loop.
  Example: `10 0 do i 5 = if leave then loop`

- **exit**: exit
  Immediately returns from the current word definition.
  Example: `: test 1 exit 2 ;` -> test pushes 1.

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

- **things** (size -- []any): Create a generic array of given size
  Example: `3 things` creates a 3-element generic array

- **@** (array index -- value): Get element at index with lossless type coercion
  Example: `myarray 0 @ .` gets first element

- **!** (value array index -- ): Set element at index with type checking - allows int→float64 but rejects float64→int
  Example: `42 myarray 0 !` sets first element to 42

- **c@** (array index -- byte): Get byte at index with auto-wrap (value & 0xff)
  Example: `mybytes 0 c@ .` gets byte value wrapped to 0-255

- **c!** (value array index -- ): Set byte at index with auto-wrap (value & 0xff)
  Example: `300 mybytes 0 c!` sets byte to 44 (300 & 0xff)

- **@push** (array value -- array'): Append value to array and return modified array
  Example: `myarray 99 @push` appends 99 and returns new array

- **@pop** (array -- array' value): Remove and return last element and modified array
  Example: `myarray @pop .` pops last element and prints it

- **@shift** (array -- array' value): Remove and return first element and modified array
  Example: `myarray @shift .` shifts first element and prints it

- **@unshift** (array value -- array'): Prepend value to array and return modified array
  Example: `myarray 0 @unshift` prepends 0 and returns new array

- **@len** (array -- length): Get array length
  Example: `myarray @len .` prints array length

## Variable Words

Variables provide named storage for any Go value, similar to traditional FORTH variables. They integrate seamlessly with array operations for automatic updates.

- **variable** ( value "name" -- ): Create a new global variable initialized to the value on the stack.
  Example: `0 variable x` creates variable x with initial value 0
  Example: `5 variable y` creates variable y with initial value 5

- **variable-does** ( value xt "name" -- ): Create a new word that behaves like a variable but executes the xt when called.
  The execution token (xt) receives the variable's value (address) on the stack.
  Example: `: constant ['] @ variable-does ;` defines constant
  Example: `3.14159 constant PI` uses the constant defined above

- **constant** ( value "name" -- ): Create a constant.
  Example: `42 constant ANSWER`
  Example: `ANSWER .` prints 42

- **@** (variable -- value): Get variable value, or get element from array in variable
  Example: `x @ .` prints the value of x
  Example: `x 0 @ .` gets element 0 from array in x

- **!** (value variable -- ): Set variable value, or set element in array in variable
  Example: `42 x !` sets x to 42
  Example: `99 x 1 !` sets element 1 in array in x to 99

### Variable Auto-Update with Arrays

When array operations are performed on variables containing slices, the variable is automatically updated with the modified array and the variable is returned on the stack for consistent stack pictures:

- **@push** (variable value -- variable): Append to array in variable
  Example: `x 99 @push` appends 99 to array in x (x auto-updated, x returned)

- **@pop** (variable -- variable value): Pop from array in variable
  Example: `x @pop drop` pops from array in x (x auto-updated, value discarded, x remains)

- **@shift** (variable -- variable value): Shift from array in variable
  Example: `x @shift drop` shifts from array in x (x auto-updated, value discarded, x remains)

- **@unshift** (variable value -- variable): Unshift to array in variable
  Example: `x 0 @unshift` prepends 0 to array in x (x auto-updated, x returned)

- **@len** (variable -- length): Get length of array in variable
  Example: `x @len .` prints length of array in x

For manual control, extract the array first: `x @ @push` (standard array behavior).

## Local Variables

Local variables can be defined within a word using the `(| ... )` syntax. They are scoped to the word execution and are automatically cleaned up when the word exits.

- **(|** ... **|)**: Define local variables.
  - Takes names from the code stream.
  - Initialized from the stack in reverse order of definition (e.g., `(| a b |)` initializes `b` from TOS, then `a`).
  - Supports `|` separator for uninitialized locals.
  - Optionally, `--` followed by comment tokens until `)`.
  - If no middle `|`, closing is `|)`.
  - Locals shadow global dictionary words.
  - Usage inside the word:
    - `name` pushes the local's value to the stack.
    - `name!` pops a value from the stack into the local.

  Example Basic:
  ```forth
  : add3 (| a b c |) a b c + + ;
  1 2 3 add3 . \ prints 6
  ```

  Example Uninitialized:
  ```forth
  : swap-locals (| a b | temp ) a temp! b a! temp b! a b ;
  ```

  Example Recursion:
  Recursion using `recur` works correctly with locals, creating a new stack frame for variables.
  ```forth
  : fact (| n |) n 1 <= IF 1 ELSE n n 1 - recur * THEN ;
  ```

  Example with comments:
  ```forth
  : test (| a b | c -- result ) a b + c * ;
  ```

## Quotations/Closures

Quotations (also called closures) allow capturing code blocks for later execution. They are created using the `[ ... ]` syntax and return a closure object that can be executed with `execute`.

- **[** ... **]**: quotationStart  
  Creates a quotation (closure) containing the compiled code between the brackets.  
  Example: `[ 1 2 + ] execute` → pushes 3 to the stack

- **execute**: execute ( closure -- )  
  Executes a quotation/closure.  
  Example: `[ 42 ] execute .` → prints 42

Quotations capture the current environment and can access local variables from their defining scope. They are particularly useful for creating anonymous functions, implementing control structures, and functional programming patterns.

Example with locals:
```forth
: make-adder (| n |) [ n + ] ;
5 make-adder constant add5
10 add5 execute . \ prints 15
```

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
