
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

\ map across input arrays into a new array (actually a flat-map if the xt returns multiple)
: @map-append ( dest arr xt ) (| arr xt |)  arr @len 0 do  arr i @ << xt execute >> 1 + <@push>  loop ; 
: @map ( arr xt -- result ) 0 things -rot @map-append ;
: @map>i ( arr xt -- result ) 0 ints -rot @map-append ;
: @map>f ( arr xt -- result ) 0 floats -rot @map-append ;

: @map2-append ( dest a1 a2 xt res ) (| a1 a2 xt |)
                 a1 @len a2 @len min 0 do i a1 over @ a2 rot @ << xt execute >> 2 + <@push> loop ; 
: @map2 ( a1 a2 xt  -- result)  0 things 3 -roll @map2-append ;
: @map2>i ( a1 a2 xt  -- result)  0 ints 3 -roll @map2-append ;
: @map2>f ( a1 a2 xt  -- result)  0 floats 3 -roll @map2-append ;

: @dot-prod ( a1 a2 -- prod ) ['] * @map2>f @sum ;


\ closures are nifty!
\ : make-counter (| n | -- counter ) [ n dup 1 + n! ] ;
\ or use variable-does 
: make-counter ( n "name" -- ) [ dup @ 1 + dup rot ! ] variable-does ;

: countdown ( n "name" -- ) [ dup @ dup 0= IF nip ELSE dup 1 - rot ! THEN ] variable-does ;

\ paragraph comments.. ignore until you see a blank line...
10 constant NL \ new line
: \p ( -- ) [[ NL ]] literal read   begin [[ NL ]] literal read " " = if exit then again ;

