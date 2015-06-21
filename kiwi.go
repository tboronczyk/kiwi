package main

import (
	"bufio"
	"flag"
	"fmt"
	"github.com/tboronczyk/kiwi/ast"
	"github.com/tboronczyk/kiwi/parser"
	"github.com/tboronczyk/kiwi/runtime"
	"github.com/tboronczyk/kiwi/scanner"
//	"github.com/tboronczyk/kiwi/symtable"
	"os"
)

func main() {
	tree := flag.Bool("ast", false, "Display parsed abstract syntax tree")
	flag.Parse()

	p := parser.New()
	p.InitScanner(scanner.New(bufio.NewReader(os.Stdin)))

	var v ast.Visitor
	if *tree {
		v = ast.NewAstPrinter()
	} else {
		v = runtime.New()
	}
	
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
