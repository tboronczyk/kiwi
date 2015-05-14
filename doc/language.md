# Kiwi Language Reference

## BNF Grammar

           expr = relation [log-op expr] .
       relation = simple-expr [cmp-op relation] .
    simple-expr = term [add-op simple-expr] .
           term = factor [mul-op term] .
         factor = '(' expr ')' | '~' expr | add-op expr | terminal .
      func-call = IDENT '(' expr-list ')' .
          addop = '+' | '-' .
       terminal = func-call | boolean | NUMBER | STRING | IDENT .
      func-call = IDENT '(' [expr-list] ')' .
        boolean = 'true' | 'false' .
      expr-list = expr | expr-list ',' expr .
         mul-op = '*' | '/' | '%' .
         cmp-op = '=' | '~=' | '<' | '<=' | '>' | '>=' .
         log-op = '&&' | '||' .

           stmt = assign-stmt | if-stmt | while-stmt | fcall-stmt .
      stmt-list = stmt | stmt-list stmt .
    assign-stmt = IDENT ':=' expr ';' .
        if-stmt = 'if' expr '{' stmt-list '}' .
     while-stmt = 'while' expr '{' stmt-list '}' .
     fcall-stmt = func-call . ';' .
