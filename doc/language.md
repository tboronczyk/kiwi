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

Kiwi is a dynamic, strongly-typed language. The fundamental data types are:

Type   | Name    | Examples
-------|---------|--------------
`bool` | boolean | true, false
`num`  | number  | 42, 3.1415
`str`  | string  | "Hello world"

A variable’s type is derived from the type of the literal or expression
value assigned to it.

    foo := true // foo takes type from boolean literal and is type bool
    bar := 42   // bar takes type from numeric literal and is type number
    baz := boo  // baz takes type from value of foo (bool) and is type bool

A variable may only be assigned a value of the same data type.

    baz := bar  // fails because baz is type bool and value of bar is number
    baz := bar:bool  // casting the value of bar permits the assignment

The type casting behavior is show in the following table:

<table>
  <tr>
    <th colspan="2">&nbsp;</th>
    <th>bool</th><th>num</th><th>str</th>
  </tr><tr>
    <th rowspan="2">bool</th>
    <td>true</td><td>true</td><td>1</td><td>"true"</td>
  </tr><tr>
    <td>false</td><td>false</td><td>0</td><td>"false"</td>
  </tr><tr>
    <th rowspan="3">num</th>
    <td>0</td><td>false</td><td>0</td><td>"0"</td>
  </tr><tr>
    <td>1</td><td>true</td><td>1</td><td>"1"</td>
  </tr><tr>
    <td>42</td><td>true</td><td>42</td><td>"42"</td>
  </tr><tr>
    <th rowspan="7">str</th>
    <td>"" (empty)</td><td>false</td><td>0</td><td>""</td>
  </tr><tr>
    <td>"true"</td><td>true</td><td>1</td><td>"true"</td>
  </tr><tr>
    <td>"false"</td><td>false</td><td>0</td><td>"false"</td>
  </tr><tr>
    <td>"0"</td><td>false</td><td>0</td><td>"0"</td>
  </tr><tr>
    <td>"1"</td><td>true</td><td>1</td><td>"1"</td>
  </tr><tr>
    <td>"42"</td><td>true</td><td>42</td><td>"42"</td>
  </tr><tr>
    <td>"42foo"</td><td>true</td><td>1</td><td>"42foo"</td>
  </tr>
</table>

### Operators

Operators are listed here in order of decreasing precedence. Operators with
higher precedence are evaluated first.

Prec. | Type          | Operators
------|---------------|----------------------------
 5    | Unary         | `+` `-` `~`
 4    | Multiplcation | `*` `/` `%`
 3    | Addition      | `+` `-`
 2    | Comparison    | `=` `~=` `>` `>=` `<` `<=`
 1    | Logic         | `&&` `||`

### Control Flow

### Functions

### Data Structures

## ABNF Grammar

    ; RFC5243 App. B defines ALPHA, CHAR, DIGIT, DQUOTE, and LF

    expr            = term [expr-op expr]
    expr-op         = add-op / mul-op / cmp-op / log-op
    add-op          = "+" / "-"
    mul-op          = "*" / "/" / "%"
    cmp-op          = "=" / "~=" / ">" / ">=" / "<" / "<="
    log-op          = "&&" / "||"
    term            = "(" expr ")" / term-op expr / cast
    term-op         = "~" / "+" / "-"
    cast            = terminal [":" ident]
    terminal        = boolean / number / string / ident / func-call
    paren-expr-list = "(" [expr *("," expr)] ")"
    boolean         = "true" / "false"
    func-call       = ident paren-expr-list
    stmt            = if-stmt / while-stmt / func-def / return-stmt /
                      assign-stmt / func-call
    if-stmt         = "if" expr brace-stmt-list [else-clause]
    brace-stmt-list = "{" *stmt "}"
    else-clause     = "else" (brace-stmt-list / expr brace-stmt-list
                      else-clause)
    while-stmt      = "while" expr brace-stmt-list
    func-def        = "func" [ident-list] brace-stmt-list
    ident-list      = ident *("," ident)
    return-stmt     = "return" [expr] LF
    assign-stmt     = ident ":=" expr LF
    ident           = ALPHA *(ALPHA / DIGIT / "_")
    number          = DIGIT ["." 1*DIGIT]
    string          = DQUOTE *CHAR DQUOTE
