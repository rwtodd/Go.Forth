# Regex Extension Words Reference

This document lists the words available in the `regex` extension. These words are available after activating the extension with `<<" regex ">> <activate-extensions>`.

In general, the regex words will have a prefix `rx`

## Regex Strings

- **rx:** : rx: ( "string" -- str ) provides a way to give an
regex string without having to escape special characters.  The first characer read after whitespace is picked up as the delimiter, and perl-style delimiter pairs (e.g., `{` ... `}` are recognized and supported.  Any non-whitespace character is acceptable.

- **\[rx:\]** : [rx:] is an immediate version of `rx:`  which is 
equivalent to:  `[[ rx: /^regex/ ]] literal`.  This is the same
relationship as `'` to `[']`, and is likewise expected to be useful during compilation.

## Compiling a Regex

- **rx-compile** : rx-compile ( pattern -- rx ) 
  Compiles a regex pattern.  
  Example: `rx: /^hello/ rx-compile .` → compiled regex

## Matching a Regex

- **rx-match?** : rx-match? ( string pattern -- bool ) 
  Matches a regex pattern against a string.  The `pattern` 
  can be a string or a compiled regex.
  Example: `"hello" rx: /^hello/ rx-match? .` → True/False

- **rx-gsub** : rx-gsub ( string pattern replacement -- string ) 
  Replaces all occurrences of a regex pattern in a string with a replacement string.  The `pattern` can be a string or a compiled regex.  The `replacement` can contain backreferences to captured groups using Go-flavored `$1` syntax.

  Example: `" hello" rx: /l/ " x" rx-gsub .` → "hexxo"

- **rx-sub** : rx-sub ( string pattern replacement -- string ) 
  Replaces the first occurrence of a regex pattern in a string with a replacement string.  The `pattern` can be a string or a compiled regex.  The `replacement` can contain backreferences to captured groups using Go-flavored `$1` syntax.

  Example: `" hello" rx: /l/ " x" rx-sub .` → "hexlo"

- **rx-split** : rx-split ( string pattern -- array ) 
  Splits a string into an array of strings using a regex pattern as the delimiter.  The `pattern` can be a string or a compiled regex.

  Example: `" hello" rx: /l/ rx-split .` → ["he", "", "o"]

- **rx-match** : rx-match ( string pattern xt --  ) 
  Matches a regex pattern against a string.  The `pattern` 
  can be a string or a compiled regex.  If the pattern matches, 
  the xt is executed with the match information on the stack.  The xt should accept two arguments: the index into the string where the match starts, and an array of strings, where the first string is the entire match, and the subsequent strings are any captured groups.
  Example: `" hello" rx: /l/ [ nip 0 @ type cr ] rx-match` →  "l\nl\n"

- **rx-find** : rx-find ( string pattern -- string ) 
  Finds the first match of a regex pattern in a string and returns the matched substring. If no match is found, an empty string is returned. The `pattern` can be a string or a compiled regex.

  Example: `" hello" rx: /[A-Za-z]l/ rx-find .` → "el"

- **rx-gfind** : rx-gfind ( string pattern -- array ) 
  Finds all matches of a regex pattern in a string and returns an array of matched substrings. The `pattern` can be a string or a compiled regex.

  Example: `" hello lolo" rx: /lo/ rx-gfind .` → ["lo", "lo"]


  