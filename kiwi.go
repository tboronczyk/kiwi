package main

import (
	"bufio"
	"flag"
	"fmt"
	"github.com/tboronczyk/kiwi/ast"
	"github.com/tboronczyk/kiwi/parser"
	"github.com/tboronczyk/kiwi/runtime"
	"github.com/tboronczyk/kiwi/scanner"
	"github.com/tboronczyk/kiwi/symtable"
	"os"
)

func main() {
	tree := flag.Bool("ast", false, "Display parsed abstract syntax tree")
	flag.Parse()

	r := bufio.NewReader(os.Stdin)
	s := scanner.New(r)
	p := parser.New()
	p.InitScanner(s)

	if *tree {
		print(p)
		return
	}
	exec(p)
}

func print(p parser.Parser) {
	v := ast.NewAstPrinter()
	for {
		n, err := p.Parse()
		if n == nil {
			if err != nil {
				fmt.Println(err)
			}
			return
		}
		n.Accept(v)
	}
}

func exec(p parser.Parser) {
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
