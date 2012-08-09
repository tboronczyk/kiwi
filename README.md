Copyright (c) 2012, Timothy Boronczyk

This is the file README for the Kiwi Programming Language distribution.

Kiwi is a modern, lightweight programming language. It is free software. You
can redistribute it and/or modify it under the terms of the license provided
under the name LICENSE.

Comments
========
Kiwi has single-line and multi-line comment support. Single-line comments start 
with `//` and go until the end of the line. Multi-line comments start with `/*`
and end with `*/`. Nested multi-line comments are supported. Both styles of
comments are ignored by the parser and so they may appear anywhere in a source
file.

Assignment, Equality, and Identity
==================================
The assignment operator is `:`. The equality operator is `=`.  The  keyword
`is` is used for identity tests (strong equality). See the Variable Declaration
and Boolean Logic sections for more information.

Numeric Literals
================
Numeric literals may be specified in a variety of styles. Standard integers are 
written as bare literals, for example the integer value forty-two is `42`.
Floats must have at least one decimal place (forty-two is written `42.0`).

Expressing integers in non-decimal bases is possible as well, such as base 2,
8, and 16. The radix is written first, followed by `#`, and then the value
observing valid digits for the selected radix. The radix portion may be omitted
in wich case hexadecimal is assumed. The integer value forty-two can thus be
written in binary as `2#101010`, in octal as `8#52`, and headecimal as either
`16#2A` or `#2A`.

String Literals
===============
String literals are enclosed by quotation marks. Within a string, the following
escape sequences are recognized:

 - `\n`          newline
 - `\r`          carriage return
 - `\t`          horizontal tab
 - `\\`          backslash
 - `\"`          double quote


Identifiers
===========
Valid identifiers start with a letter or underscore and may contain a mix of
letters, underscores, and numbers. If the programmer wishes to use an identifier
that is the same as a reserved word he may preface the identifier with a
backtick.

Variable Declaration
====================
Variables must be declared before they are used. They are declared using the
keyword `var`. Kiwi allows multiple variables to be declared with one `var` with
the variables being separated by commas.

...

Conditional Branching
=====================
The keywords `if` and `else` are used to write conditional branches of code. 
...

Boolean Logic
=============
The operators `&&`, `||`, and `^^` are logical-and, logical-or, and logical-xor
respectively. They return the Boolean values `true` or `false` and their behavior
is short-circuiting. `~` is logical-not, and `~=` is not-equals.

The `=` operator tests for equality between values, and `is` is used for
identity tests. They too return Boolean `true` or `false`.

...

Examples
========
    // generate the Nth Fibonacci number 
    func fibonacci: n { 
        var m, p, r.
        if n < 2 { 
            return n.
        } 
        m: fibonacci(n - 1).
        p: fibonacci(n - 2).
        r: m + p.
        return r.
    }
    var i: 1, fib.
    while i <= 10 { 
        fib: fibonacci(i).
        print(fib).
        i: i + 1.
    }


    // escaped identifier, binary literal, and nested comment
    var `var: 2#100.
    if `var > 10 {
        /* greater than 10 */
        print("value is greater than 10").
    }
    else if `var < 0 {
        /* negative */
        print("negative value").
    }
    /*
    else {
        /* less than 10 */
        print("value less or equal to 10").
    }
    */
