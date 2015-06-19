# Kiwi Language Reference

## BNF Grammar

                expr = relation [log-op expr] .
            relation = simple-expr [cmp-op relation] .
         simple-expr = term [add-op simple-expr] .
                term = factor [mul-op term] .
              factor = '(' expr ')' | expr-op expr | cast .
             expr-op = '~' | add-op .
              add-op = '+' | '-' .
                cast = terminal ['!' IDENT] .
            terminal = boolean | number | STRING | IDENT | func-call .
           func-call = IDENT func-call-args .
      func-call-args = '(' ')' | '(' expr-list ')' .
           expr-list = expr | expr-list ',' expr .
             boolean = 'true' | 'false' .
              number = (0-9)+ ['.' (0-9)+] .
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
         assign-stmt = IDENT ':=' expr '.' .
         return-stmt = 'return' [expr] '.' .
      func-call-stmt = func-call '.' .

## Language Constructs

### Comments

Both single and multiple-line comments are supported. Single-line comments
begin with `//` and span to the end of the current line. Multiple-line
comments open with `/*` and close with `*/`. It is possible to nest 
multiple-line comments.

    // this is a single-line comment

    /* this is a multiple-line comment
    that spans multiple lines. It also
    /* this is a nested comment */
    allows nested comments. */

### Data Types

### Operators

### Control Flow

### Functions

### Data Structures

