package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"strings"
)

// built-in functions, [name]func{implementation}
var builtins = map[string]func(*Stack, params, io.Reader, io.Writer, io.Writer){
	// strlen - returns the length of a string
	"strlen": func(s *Stack, p params, stdin io.Reader, stdout, stderr io.Writer) {
		s.Push(ScopeEntry{TypNumber, len(p[0].Value.(string))})
	},

	// write - prints a value
	"write": func(s *Stack, p params, stdin io.Reader, stdout, stderr io.Writer) {
		for i := range p {
			fmt.Fprint(stdout, p[i].Value)
		}
	},

	// read - read a string
	"read": func(s *Stack, p params, stdin io.Reader, stdout, stderr io.Writer) {
		b, err := ioutil.ReadAll(stdin)
		if err != nil {
			panic(err)
		}
		s.Push(ScopeEntry{TypString, strings.TrimRight(string(b), "\n")})
	},
}
