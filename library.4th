
\ a for-loop over an array
: @for-{ ( arr -- ) postpone dup postpone @len 0 postpone literal postpone do 
                    postpone dup postpone i postpone @  postpone swap postpone >r ; immediate
: }-@for postpone r> postpone loop postpone drop ; immediate   

\ example use: array sum and product
\ : @sum  ( arr -- sum )  0 swap @for-{ + }-@for ;
\ : @prod ( arr -- prod ) 1 swap @for-{ * }-@for ;


: @for-each (| arr xt | -- ?? ) arr @len 0 do arr i @ xt execute loop ;
: @sum  ( arr -- sum )  0 swap ['] + @for-each ;
: @prod ( arr -- prod ) 1 swap ['] * @for-each ;

\ map across 2 input arrays into a third new array 
: @map (| arr xt |)  0 things   arr @len 0 do   arr i @ xt execute   swap @push   loop ; 
: @map2 (| a1 a2 xt |)  0 things a1 @len a2 @len min 0 do i a1 over @ a2 rot @ xt execute swap @push loop ; 

: @dot-prod ( a1 a2 -- prod ) ['] * @map2 @sum ;


\ closures are nifty!
\ : make-counter (| n | -- counter ) [ n dup 1 + n! ] ;
\ or use variable-does 
: make-counter ( n "name" -- ) [ dup @ 1 + dup rot ! ] variable-does ;

: countdown ( n "name" -- ) [ dup @ dup 0= IF nip ELSE dup 1 - rot ! THEN ] variable-does ;

\ paragraph comments.. ignore until you see a blank line...
10 constant NL \ new line
: \p ( -- ) [[ NL ]] literal read   begin [[ NL ]] literal read " " = if exit then again ;

