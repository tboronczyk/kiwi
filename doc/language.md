# Kiwi Language Reference

## BNF Grammar

                expr = relation [log-op expr] .
            relation = simple-expr [cmp-op relation] .
         simple-expr = term [add-op simple-expr] .
                term = factor [mul-op term] .
              factor = '(' expr ')' | expr-op expr | terminal .
             expr-op = '~' | add-op .
              add-op = '+' | '-' .
            terminal = boolean | NUMBER | STRING | IDENT | func-call .
           func-call = IDENT func-call-args .
      func-call-args = '(' ')' | '(' expr-list ')' .
           expr-list = expr | expr-list ',' expr .
             boolean = 'true' | 'false' .
              mul-op = '*' | '/' | '%' .
              cmp-op = '=' | '~=' | '>' | '>=' | '<' | '<=' .
              log-op = '&&' | '||' .

                stmt = if-stmt | while-stmt | func-def-stmt | assign-stmt |
                       return-stmt | func-call-stmt .
             if-stmt = 'if' expr braced-stmt-list .
    braced-stmt-list = '{' stmt-list '}' .
           stmt-list = stmt | stmt-list stmt .
          while-stmt = 'while' expr braced-stmt-list .
       func-def-stmt = 'func' ident-list braced-stmt-list .
          ident-list = IDENT | ident-list ',' IDENT .
         assign-stmt = IDENT ':=' expr ';' .
         return-stmt = 'return' [expr] ';' .
      func-call-stmt = func-call ';' .
