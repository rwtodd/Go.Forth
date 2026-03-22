# Go.Forth
An embeddable postfix mini-language for Go programs.

## Does it follow ANS FORTH?

No, and it intentionally diverges from ANS FORTH. Primarily because I did not want to simulate a raw memory space, which would
complicate interactions between the language and Go. Instead, the stack and all variables use garbage-collected Go `any` types, with overloads allowing operations like:

~~~~~~
: double dup + ;

" hi" double .  ( '+' works on strings )
hihi
~~~~~~

There are no words like `c,` to push raw data into a data segment.

Otherwise, it should feel pretty FORTHy, with immediate words
and `POSTPONE` letting you do compile-time programming.

## Is it fast?

Performance is secondary to ease of embedding and interaction with the host Go program. Anything that's annoyingly slow can
always be provided from the Go side of the wall.

## What is the status?

Go.Forth is a mature, embeddable Forth interpreter written in Go. It supports a comprehensive set of words for stack manipulation, arithmetic, control flow, arrays, dictionaries, variables, and more. See `Reference.md` for a complete word reference.

## Extensions

Go.Forth supports optional extensions that add specialized functionality:

- **random**: Random number generation and array shuffling (`randint`, `randfloat`, `randseed!`, `@shuffle`, `@select`)
- **regex**: Regular expression operations (`rx:`, `rx-match?`, `rx-gsub`, `rx-split`, etc.)
- **strings**: String manipulation utilities (`trim""`, `upper""`, `split""`, `replace""`, etc.)

Extensions can be activated with: `<<" extension-name ">> <activate-extensions>`

## Library Files

- `library.4th`: Higher-level utilities for array operations (`@sum`, `@prod`, `@map`, etc.)

## Words Implemented

Go.Forth supports a comprehensive set of words organized by category (see `Reference.md` for complete details):

### Stack Manipulation
dup drop swap over rot -rot pick roll -roll nip tuck >r r> r@ rdrop depth

### Variable Length Arguments
<< >> <<" <@push> <ints> <floats> <strings> <bytes> <things> @spread sprintf

### Arithmetic & Math
+ - * / ** sqrt log log10 log2 max min sin cos tan round floor ceil

### Comparison & Logic
= < > <= >= <> 0= 0< 0> and or xor invert true false

### Control Flow
if else then recur do loop +loop i j leave exit begin again until while repeat

### Input/Output
read skip " chr ord .s . type cr

### Compilation & Definition
: ; literal postpone immediate compile, ' ['] execute eval load

### Interpretation Control
[[ ]] [ ]

### Comments
\ ( \p

### Arrays
bytes ints floats strings things @ ! c@ c! @len @push @pop @shift @unshift

### Variables & Constants
variable variable-does constant @ !

### Local Variables
(| |)

### Quotations/Closures
[ ] execute

### Dictionaries
empty-dict d@ d! ddel dkeys d@| d@?

### Extensions
extension-list <activate-extensions>

### Internal/System
debug. mark forget read-token lookup-xt

You can define custom words, including immediate ("macro"-type) words using `postpone`.

Go.Forth supports advanced features like quotations/closures, local variables, dictionaries, and variable-length argument handling.

Example with quotations and locals:

~~~~~~
: make-adder (| n |) [ n + ] ;
5 make-adder [ @ execute ] " add5" variable-does
10 add5 .  \ prints 15
~~~~~~

Example with arrays and higher-order functions (using library.4th):

~~~~~~
<< 1 2 3 4 5 >> <ints> @sum .  \ prints 15
<< 1 2 3 4 5 >> <ints> [ 2 * ] @map @sum .  \ prints 30
~~~~~~

Go.Forth is ready for embedding in Go applications as a powerful scripting and control language.

## Prior Work

I had written a java-based FORTH interpreter a while back (it lives in
the [small-programs](https://github.com/rwtodd/small_programs) repo).
For that one, I __did__ simulate a flat memory space, and the result was
much closer to an ANS FORTH implementation -- or at least it could
eventually be one.  For this interpreter, the goals are different.

