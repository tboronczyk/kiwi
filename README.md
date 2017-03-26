# Kiwi Programming Langauge

[![Build Status](https://travis-ci.org/tboronczyk/kiwi.svg?branch=master)](https://travis-ci.org/tboronczyk/kiwi)

This is the README file for the Kiwi Programming Language distribution.

Kiwi is my on-going “forever project” to implement a programming language. Not
only is it a learning opportunity, it’s also a playground for me to explore
issues in language design.

Kiwi probably won’t be useful to anyone other than me, especially at the rate I
get to work on it, but my code is available as free software. You can do
whatever you want with it under the terms of the license provided in the
LICENSE file. Just don’t come after me if it melts down your computer.

## Code Examples

### Generate the Fibonacci Series

    func fibonacci n
    { 
        if n < 2 { 
            return n
        } 
        m := fibonacci(n - 1)
        p := fibonacci(n - 2)
        return m + p
    }

    i := 0
    while i < 10 { 
        i := i + 1
        fib := fibonacci(i)
        write(fib, "\n")
    }

### Fizz Buzz

    i := 0
    while i < 100 {
        i := i + 1
    
        if i % 15 = 0 {
            write("Fizz Buzz\n")
        }
        else i % 3 = 0 {
            write("Fizz\n")
        }
        else i % 5 = 0 {
            write("Buzz\n")
        }
        else {
            write(i, "\n")
        }
    }
