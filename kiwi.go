package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"

	"github.com/tboronczyk/kiwi/ast"
	"github.com/tboronczyk/kiwi/parser"
	"github.com/tboronczyk/kiwi/runtime"
	"github.com/tboronczyk/kiwi/scanner"
)

func main() {
	p := parser.New(scanner.New(bufio.NewReader(os.Stdin)))
	var v ast.NodeVisitor

	tree := flag.Bool("ast", false, "Display parsed abstract syntax tree")
	flag.Parse()
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
