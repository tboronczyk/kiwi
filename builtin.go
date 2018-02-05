package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

var rStdin = bufio.NewReader(os.Stdin)

// built-in functions, [name]func{implementation}
var builtins = map[string]func(*Stack, params){
	// strlen - returns the length of a string
	"strlen": func(s *Stack, p params) {
		s.Push(ScopeEntry{TypNumber, len(p[0].Value.(string))})
	},

	// write - prints a value
	"write": func(s *Stack, p params) {
		for i, _ := range p {
			fmt.Print(p[i].Value)
		}
	},

	// read - read a string
	"read": func(s *Stack, p params) {
		str, err := rStdin.ReadString('\n')
		if err != nil {
			panic(err)
		}
		s.Push(ScopeEntry{TypString, strings.TrimRight(str, "\n")})
	},
}
