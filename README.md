Copyright (c) 2012, Timothy Boronczyk

This is the file README for the Kiwi Programming Language distribution.

Kiwi is a modern, lightweight programming language. It is free software. You
can redistribute it and/or modify it under the terms of the license provided
under the name LICENSE.

Comments
========
Kiwi has single-line and multi-line comment support. Single-line comments start 
with `//` and go until the end of the line. Multi-line comments start with `/*`
and end with `*/`. Nested multi-line comments are supported.  

Assignment, Equality, and Identity
==================================
The assignment operator is `:=`. The equality operator is `==`.  The  keyword
`is` is used for identity tests (strong equality). See the Variable Declaration
and Boolean Logic sections for more information.

Numeric Literals
================
Numeric literals may be specified in a variety of styles. Standard integers are 
written as bare literals, for example the integer value seven is 7, and the 
value forty-two is 42. Kiwi makes available a unique notation for expressing 
integers in non-decimal bases as well, such as base 2 (binary), base 8 (octal), 
and base 16 (hexadecimal) -- the radix is written, followed by `#`, and then 
the value observing valid digits for the selected radix.  The radix portion may 
be omitted in which case the notation will default to hexadecimal. The integer 
value forty-two is thus written in binary as 2#101010 (valid digits are 0 and 
1), in octal as 8#52 (valid digits are 0 through 7), and hexadecimal as either 
16#2A or #2A (valid digits are 0 through 9, and case-insensitive A through F).  

String Literals
===============
String literals are enclosed by double quotation marks. Within a string, the
following escape sequences are recognized:

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

The wildcard identifier `_` is used when one wishes to specify the value is 
unknown or to be discarded. In Boolean tests its value is considered true.

Variable Declaration
====================
Variables must be declared before they are used. They are declared using the
keyword `var`. Kiwi allows multiple variables to be declared with one `var` with
the variables being separated by commas.

...

Conditional Branching
=====================
The keywords `if` and `else` are used to write conditional branches of code. Both
require a condition and body-- the condition are evaluated in order and when the
statements that make up the associated body to the first one determined true are
executed. The first branch is always marked with `if`, and any number of
alternates are marked with `else`. The `_` wildcard may be used to denote a default
else case. Consider the following example:

    if i > 10 {
        // this block will be executed if the value stored
        // in i is greater than 10
    }
    else i < 0 {
        // this block will be executed if the value is 
        // negative
    }
    else _ {
        // this block will be executed if the value is
        // positive and less than or equal to 10
    }

Boolean Logic
=============
The operators `&&`, `||`, and `^^` are logical-and, logical-or, and logical-xor
respectively. They return the Boolean values `true` or `false` and their behavior
is short-circuiting. `~` is logical-not, and `~=` is not-equals.

The `==` operator tests for equality between values, and `is` is used for
identity tests. They too return Boolean `true` or `false`.

...
