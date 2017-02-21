# Forth
An embeddable postfix mini-language for my Go programs.

## Is it a true ANS FORTH?

No, because I did not want to simulate a raw memory space, which would
complicate interactions between the language and Go.  So, instead,
I make the stack and all variables of garbage-collected type `interface{}`,
and provide overloads so you can say stuff like:

~~~~~~
4 1.2 * .  ( '*' can multiply an int with a float )
4.8

: double dup + ;

3 double . 
6

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
\ ( read skip " chr ord .s . type
[ ] : ; literal postpone immediate 
dup drop swap over rot -rot + * mark 
forget
~~~~~~

At this point, you can define custom words, which can include
immediate ("macro"-type words) which use `postpone`. 

The big missing feature is conditionals and loops.  Those will
be implemented next.  At some point I plan to incorporate a 
pre-parsed block of forth code to implement the rest of the high-level
features.  I believe this is a typical implementation technique.

## Prior Work

I had written a java-based FORTH interpreter a while back (it lives in
the [small-programs](https://github.com/rwtodd/small_programs) repo). For that
one, I __did__ simulate a flat memory space, and the result was much closer to
an ANS FORTH implementation -- or at least it could eventually be one.  For this 
interpreter, the goals are different.

