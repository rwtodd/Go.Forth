# Forth
An embeddable postfix mini-language for my Go programs.

## Is it a true ANS FORTH?

No, because I did not want to simulate a raw memory space, which would
complicate interactions between the language and Go.  So, instead, I
make the stack and all variables of garbage-collected type
`interface{}`, and provide overloads so you can say stuff like:

~~~~~~
: double dup + ;

" hi" double .  ( '+' works on strings )
hihi
~~~~~~

Similarly, I won't have words like `c,` to push raw data into a data segment.

Otherwise, though, it should feel pretty FORTHy, with immediate words 
and `POSTPONE` letting you do compile-time programming.

## Is it fast?

It's too early to tell how fast or slow it will be, but I'm focused
first and formost on making it easy to embed and interact with the host
Go program.  Speed is secondary, becuse anything that's annoyingly slow can
always be provided from the Go side of the wall.

## What is the status?

This is just preliminary work.  Words implemented:

~~~~~~
\ ( read skip " chr ord .s . type cr
[ ] : ; literal postpone immediate 
dup drop swap over rot -rot + * mark 
forget if else then recur  >r r> r@ rdrop
do loop +loop i j
~~~~~~

At this point, you can define custom words, which can include
immediate ("macro"-type words) which use `postpone`. 

As of Sept 2018, we have IF/ELSE/THEN, RECUR, and DO loops.

~~~~~~
: block ( size -- ) 
  0 swap tuck 0 DO over over DO j type i .  LOOP cr LOOP drop drop ;
: blocks ( num -- ) 
  0 DO   i block  LOOP ;
5 blocks
00
00 01
10 11
00 01 02
10 11 12
20 21 22
00 01 02 03
10 11 12 13
20 21 22 23
30 31 32 33

: regreet ( num -- ) " HELLO! " . -1 + dup IF recur ELSE drop THEN ;
5 regreet
HELLO!  HELLO!  HELLO!  HELLO!  HELLO! 
~~~~~~

At this point, it's actually starting to be useful to embed in things as a basic
control language.  I need to flesh out the math functions, and make it easy to 
deal with Go arrays.

## Prior Work

I had written a java-based FORTH interpreter a while back (it lives in
the [small-programs](https://github.com/rwtodd/small_programs) repo).
For that one, I __did__ simulate a flat memory space, and the result was
much closer to an ANS FORTH implementation -- or at least it could
eventually be one.  For this interpreter, the goals are different.

