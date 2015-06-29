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
var builtins = map[string]func(*util.Stack, Params){

	// strlen - returns the length of a string
	"strlen": func(s *util.Stack, p Params) {
		s.Push(ValueEntry{
			Value: len(p[0].Value.(string)),
			Type:  NUMBER,
		})

	},
	// write - prints a value
	"write": func(s *util.Stack, p Params) {
		for i, _ := range p {
			fmt.Print(p[i].Value)
		}
	},
	// read - read a string
	"read": func(s *util.Stack, p Params) {
		str, err := rdStdin.ReadString('\n')
		if err != nil {
			panic(err)
		}
		s.Push(ValueEntry{
			Value: strings.TrimRight(str, "\n"),
			Type:  STRING,
		})
	},
}
