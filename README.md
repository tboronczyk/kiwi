# Kiwi Programming Langauge

© 2012, 2015 Timothy Boronczyk

This is the README file for the Kiwi Programming Language distribution.

## Introduction

Kiwi is my on-going “forever project” to implement a programming language. Not
only is it a learning opportunity, it’s also a playground to explore issues in
language design.

Kiwi probably won’t be useful to anyone other than me, especially at the rate I
get to work on it, but the code is licensed as free software. You can do
whatever you want with it under the terms of the license provided in the
LICENSE file.

## Build

    sudo apt-get install build-essential autoconf pkg-config bison libicu-dev \
     check git
    git clone https://github.com/tboronczyk/Kiwi
    cd Kiwi
    autoreconf --install
    ./configure
    make

## Language Definition

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

### Statement Terminator

Statements are terminated with a period (i.e., `.`).

### Assignment, Equality, and Identity

The assignment operator is `:=`. The equality operator is `=`. The keyword
`is` is reserved for identity tests (strong equality) if such an operator is
needed as the language evolves.

The sections Declaring Variables and Boolean Logic have more information.

### Numeric Literals

Integers are written as bare literals (e.g., `42`). Floating-point numbers must
have at least one fractional digit (e.g., `42.0`).

Expressing integers in non-decimal bases is possible as well, such as base 2,
8, and 16. The radix and value are separated by `#` and the value can only
contain digits valid for the given radix (e.g., the binary representation of
42 is `2#101010` and the octal representation is `8#52`). The radix may be
omitted, in which base 16 is assumed (e.g., both `16#2A` and `#2A` are the
same).

### String Literals

String literals are enclosed by straight double-quotation marks (e.g.,
`"Hello World!"`). The following escape sequences are recognized within a
string:

 * `\n`  New line
 * `\r`  Carriage return
 * `\t`  Horizontal tab
 * `\\`  Backslash
 * `\"`  Double quote

### Identifiers

Identifiers may be named using any combination of letters, underscores, and numbers. If the identifier is the same as a reserved keyword, or if it begins
with a number, it must be prefixed with a backtick (e.g., `` `123_foo ``).

### Declaring Variables

Variables must be declared using the keyword `var` before they are used. 
Multiple variables may be declared by one `var` provided the identifiers are
separated by commas.

    var foo, bar.

### Conditional Branching

The keywords `if` and `else` are used for conditional branching. The
expression following `if` is the conditional’s test and the body is provided
surrounded by `{` and `}`.

    if hour < 12 {
        write("Good morning!").
    }
    else if hour < 18 {
        write("Good afternoon!").
    }
    else {
        write("Good night!").
    }

### Boolean Logic

Boolean operators test a relationship between two expressions and return
the values `True` or `False`. The operators are:

  * `&&`  Logical-and
  * `||`  Logical-or
  * `^^`  Logical-xor
  * `=`   Equality
  * `~=`  Inequality
  * `is`  Strong equality (if implemented)

The negation operator is `~`.

## Code Examples

### Generate the Nth Fibonacci Number

    func fibonacci n
    { 
        var m, p, r.
        if n < 2 { 
            return n.
        } 
        m := fibonacci(n - 1).
        p := fibonacci(n - 2).
        r := m + p.
    }

    var i := 1, fib.
    while i <= 10 { 
        fib := fibonacci(i).
        write(fib).
        i := i + 1.
    }

### Escaped Identifier, Binary Literal, and Nested Comment

    `var := 2#100.
    if `var > 10 {
        /* greater than 10 */
        write("value is greater than 10").
    }
    else if `var < 0 {
        /* negative */
        write("negative value").
    }
    /*
    else {
        /* less than 10 */
        write("value less or equal to 10").
    }
    */
