package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/tboronczyk/kiwi/ast"
	"github.com/tboronczyk/kiwi/parser"
	"github.com/tboronczyk/kiwi/runtime"
	"github.com/tboronczyk/kiwi/scanner"
)

var optAst = flag.Bool("ast", false, "Print parsed abstract syntax tree")

func main() {
	defer func() {
		if e := recover(); e != nil {
			log.Fatal(e)
		}
	}()

	flag.Parse()
	args := flag.Args()

	var in io.Reader
	switch len(args) {
	case 0:
		in = os.Stdin
		break
	case 1:
		fp, err := os.Open(args[0])
		if err != nil {
			panic(err)
		}
		in = fp
		break
	default:
		panic("junk arguments")
	}

	p := parser.New(scanner.New(bufio.NewReader(in)))
	r := runtime.New()
	v := ast.NewAstPrinter()

	n, err := p.Parse()
	if n == nil {
		if err != nil {
			fmt.Println(err)
		}
		return
	}
	if *optAst {
		n.Accept(v)
	} else {
		n.Accept(r)
	}
}
