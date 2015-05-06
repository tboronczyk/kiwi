# Kiwi Language Reference

## BNF Grammar

           expr = relation [log-op expr] .
       relation = simple-expr [cmp-op relation] .
    simple-expr = term [add-op simple-expr] .
           term = factor [mul-op term] .
         factor = '(' expr ')' | '~' expr | add-op expr | terminal .
          addop = '+' | '-' .
       terminal = boolean | NUMBER | STRING | IDENT .
         mul-op = '*' | '/' | '%' .
        boolean = 'true' | 'false' .
         cmp-op = '=' | '~=' | '<' | '<=' | '>' | '>=' .
         log-op = '&&' | '||' .
