package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/jawher/mow.cli"
)

func main() {
	defer func() {
		if e := recover(); e != nil {
			log.Fatal(e)
		}
	}()

	app := cli.App("kiwi", "the kiwi language interpreter")
	app.Spec = "[-t] [FILE]"

	tree := app.BoolOpt("t tree", false, "print out syntax tree")
	file := app.StringArg("FILE", "", "source file")

	app.Action = func() {
		var fp io.Reader
		if *file == "" {
			fp = os.Stdin
		} else {
			var err error
			fp, err = os.Open(*file)
			if err != nil {
				panic(err)
			}
		}

		p := NewParser(NewScanner(bufio.NewReader(fp)))

		n, err := p.Parse()
		if err != nil {
			fmt.Println(err)
		}
		if n == nil {
			return
		}

		if *tree {
			n.Accept(NewAstPrinter())
		} else {
			n.Accept(NewRuntime(os.Stdin, os.Stdout, os.Stderr))
		}
	}
	app.Run(os.Args)
}
