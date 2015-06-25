package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"

	"github.com/tboronczyk/kiwi/analyzer"
	"github.com/tboronczyk/kiwi/ast"
	"github.com/tboronczyk/kiwi/parser"
	"github.com/tboronczyk/kiwi/runtime"
	"github.com/tboronczyk/kiwi/scanner"
)

func main() {
	defer func() {
		if e := recover(); e != nil {
			fmt.Println(e)
		}
	}()

	p := parser.New(scanner.New(bufio.NewReader(os.Stdin)))
	a := analyzer.New()

	tree := flag.Bool("ast", false, "Display parsed abstract syntax tree")
	flag.Parse()

	var v ast.NodeVisitor
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
		n.Accept(a)
		n.Accept(v)
	}
}
