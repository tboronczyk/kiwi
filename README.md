# Kiwi Programming Langauge

© 2012, 2015 Timothy Boronczyk

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
            return n.
        } 
        m := fibonacci(n - 1).
        p := fibonacci(n - 2).
        r := m + p.
    }

    i := 1.
    while i <= 10 { 
        fib := fibonacci(i).
        write(fib).
        i := i + 1.
    }
