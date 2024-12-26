package main

import (
	"fmt"
	"io"
	"strings"
)

// built-in functions, [name]func{implementation}
var builtins = map[string]func(*Stack, params, *RuntimeEnv){
	// strlen - returns the length of a string
	"strlen": func(s *Stack, p params, env *RuntimeEnv) {
		s.Push(ScopeEntry{TypNumber, len(p[0].Value.(string))})
	},

	// write - prints a value
	"write": func(s *Stack, p params, env *RuntimeEnv) {
		for i := range p {
			fmt.Fprint(env.stdout, p[i].Value)
		}
	},

	// read - read a string
	"read": func(s *Stack, p params, env *RuntimeEnv) {
		b, err := io.ReadAll(env.stdin)
		if err != nil {
			panic(err)
		}
		s.Push(ScopeEntry{TypString, strings.TrimRight(string(b), "\n")})
	},
}
