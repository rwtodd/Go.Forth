# Random Extension Words Reference

This document lists the words available in the `random` extension. These words are available after activating the extension with `<<" random ">> <activate-extensions>`.

## Random Number Generation

- **randint**: randint ( max min -- n )
  Returns a random integer in the half-closed range [min, max).
  Example: `10 0 randint .` → (random integer between 0 and 9)

- **randfloat**: randfloat ( -- n )
  Returns a random float64 in the range [0.0, 1.0).
  Example: `randfloat .` → (random float e.g. 0.12345678)

- **randseed!**: randseed! ( n -- )
  Sets the random number generator seed using the provided integer.
  Example: `12345 randseed!`

## Array Operations

- **@shuffle**: @shuffle ( arr/var -- )
  Shuffles the elements of the given array (or an array stored in a variable) in-place.
  Example: `<< 1 2 3 4 5 >> <ints> variable x x @shuffle x @ .s`

- **@select**: @select ( arr/var -- value )
  Selects an element from the given array (or an array stored in a variable) at random and pushes it to the stack.
  Example: `<< "a" "b" "c" >> <strings> @select type`
