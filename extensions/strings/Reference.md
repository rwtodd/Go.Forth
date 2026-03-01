# Strings Extension Words Reference

This document lists the words available in the `strings` extension. These words are available after activating the extension with `<<" strings ">> <activate-extensions>`.

## String Queries

- **blank""?**: blank""? ( str -- bool )
  Checks if a string is entirely whitespace or empty.
  Example: `"   " blank""? .` → -1
  
- **len""**: len"" ( str -- len )
  Returns the length of the string in bytes.
  Example: `" hello" len"" .` → 5

- **starts""?**: starts""? ( str prefix -- bool )
  Checks if the target string starts with the prefix.
  Example: `" hello world" " hello" starts""? .` → -1
  
- **ends""?**: ends""? ( str suffix -- bool )
  Checks if the target string ends with the suffix.
  Example: `" hello world" " rld" ends""? .` → -1

- **contains""?**: contains""? ( str substr -- bool )
  Checks if the target string contains the substring.
  Example: `" hello world" " lo " contains""? .` → -1

- **index""**: index"" ( str substr -- index )
  Returns the byte index of the first instance of `substr`, or -1 if not found.
  Example: `" hello world" " world" index"" .` → 6

## String Manipulations

- **trim""**: trim"" ( str -- str )
  Trims whitespace from both ends of the string.
  Example: `"  hello  " trim"" type`

- **triml""**: triml"" ( str -- str )
  Trims leading whitespace from the string.
  
- **trimr""**: trimr"" ( str -- str )
  Trims trailing whitespace from the string.
  
- **upper""**: upper"" ( str -- str )
  Converts the string to uppercase.

- **lower""**: lower"" ( str -- str )
  Converts the string to lowercase.

- **replace""**: replace"" ( str old new -- str )
  Replaces all occurrences of `old` with `new` in the string.
  Example: `" hello world" " o" " a" replace"" type`

- **sub""**: sub"" ( str idx1 idx2 -- substr )
  Extracts the substring `[idx1, idx2)`. Negative indices count from the end of the string. If `idx1 != 0` and `idx2 == 0`, `idx2` is treated as the end of the string.
  Example: `" hello" 1 -1 sub"" type` → ell

## Delimiting and Formatting

- **split""**: split"" ( str sep -- arr )
  Splits the string by a separator into an array of strings.
  Example: `" a,b,c" " ," split"" <strings> @len .` → 3
  
- **@join""**: @join"" ( arr sep -- str )
  Joins an array (or a variable pointing to an array) of strings by a separator.
  Example: `<< " a" " b" " c" >> <strings> " ," @join"" type` → a,b,c

- **<join"">**: <join""> ( str1 ... strn n sep -- str )
  Takes `n` strings from the stack and joins them by a separator.
  Example: `" a" " b" " c" 3 " ," <join""> type` → a,b,c
