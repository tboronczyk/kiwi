package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
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

	p := NewParser(NewScanner(bufio.NewReader(in)))
	r := NewRuntime(os.Stdin, os.Stdout, os.Stderr)
	v := NewAstPrinter()

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
