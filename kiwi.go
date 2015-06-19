package main

import (
	"bufio"
	"fmt"
	"github.com/tboronczyk/kiwi/runtime"
	"github.com/tboronczyk/kiwi/parser"
	"github.com/tboronczyk/kiwi/scanner"
	"github.com/tboronczyk/kiwi/symtable"
	"os"
)

func main() {
	r := bufio.NewReader(os.Stdin)
	s := scanner.New(r)
	p := parser.New()
	p.InitScanner(s)

	varTable := symtable.New()
	funTable := symtable.New()
	runtime.LoadBuiltins(funTable)
	for {
		n, err := p.Parse()
		if n == nil {
			if err != nil {
				fmt.Println(err)
			}
			return
		}
		runtime.Eval(n, varTable, funTable)
	}
}
