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

var optAst = flag.Bool("ast", false, "Print parsed abstract syntax tree")

func main() {
	flag.Parse()

	p := parser.New(scanner.New(bufio.NewReader(os.Stdin)))
	r := runtime.New()
	v := ast.NewAstPrinter()
	for {
		n, err := p.Parse()
		if n == nil {
			if err != nil {
				fmt.Println(err)
			}
			return
		}
		if *optAst {
			n.Accept(v)
		}
		n.Accept(r)
	}
}
