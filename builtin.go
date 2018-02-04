package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

var rdStdin = bufio.NewReader(os.Stdin)

// built-in functions, [name]func{implementation}
var builtins = map[string]func(*Stack, params){

	// strlen - returns the length of a string
	"strlen": func(s *Stack, p params) {
		s.Push(Entry{
			Value:    len(p[0].Value.(string)),
			DataType: NUMBER,
		})

	},
	// write - prints a value
	"write": func(s *Stack, p params) {
		for i, _ := range p {
			fmt.Print(p[i].Value)
		}
	},
	// read - read a string
	"read": func(s *Stack, p params) {
		str, err := rdStdin.ReadString('\n')
		if err != nil {
			panic(err)
		}
		s.Push(Entry{
			Value:    strings.TrimRight(str, "\n"),
			DataType: STRING,
		})
	},
}
