package runtime

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/tboronczyk/kiwi/util"
)

var rdStdin = bufio.NewReader(os.Stdin)

// built-in functions, [name]func{implementation}
var builtins = map[string]func(*util.Stack, params){

	// strlen - returns the length of a string
	"strlen": func(s *util.Stack, p params) {
		s.Push(valueEntry{
			value: len(p[0].value.(string)),
			dtype: NUMBER,
		})

	},
	// write - prints a value
	"write": func(s *util.Stack, p params) {
		for i, _ := range p {
			fmt.Print(p[i].value)
		}
	},
	// read - read a string
	"read": func(s *util.Stack, p params) {
		str, err := rdStdin.ReadString('\n')
		if err != nil {
			panic(err)
		}
		s.Push(valueEntry{
			value: strings.TrimRight(str, "\n"),
			dtype: STRING,
		})
	},
}
