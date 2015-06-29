# Kiwi Language Reference

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

## BNF Grammar

               expr := relation (log-op expr)? .
           relation := simple-expr (cmp-op relation)? .
        simple-expr := term (add-op simple-expr)? .
               term := factor (mul-op term)? .
             factor := '(' expr ')' | expr-op expr | cast .
            expr-op := '~' | '+' | '-' .
               cast := terminal (':' IDENT)? .
           terminal := boolean | number | STRING | IDENT | func-call .
    paren-expr-list := '(' ')' | '(' expr (',' expr)* ')' .
            boolean := 'true' | 'false' .
             number := (0-9)+ ['.' (0-9)+] .
          func-call := IDENT paren-expr-list .
             mul-op := '*' | '/' | '%' .
             add-op := '+' | '-' .
             cmp-op := '=' | '~=' | '>' | '>=' | '<' | '<=' .
             log-op := '&&' | '||' .
               stmt := if-stmt | while-stmt | func-def | return-stmt | 
                       assign-stmt | func-call .
            if-stmt := 'if' expr braced-stmt-list (else-clause)? .
    brace-stmt-list := '{' (stmt)* '}'
        else-clause := 'else' (brace-stmt-list | expr brace-stmt-list else-clause) .
         while-stmt := 'while' expr brace-stmt-list .
           func-def := 'func' (ident-list)? brace-stmt-list .
         ident-list := IDENT (',' IDENT)? .
        return-stmt := 'return' (expr)? NL .
        assign-stmt := IDENT ':=' expr NL .
