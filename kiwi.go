package main

import (
	"bufio"
	"fmt"
	"github.com/tboronczyk/kiwi/ast"
	"github.com/tboronczyk/kiwi/parser"
	"github.com/tboronczyk/kiwi/scanner"
	"os"
)

func main() {
	r := bufio.NewReader(os.Stdin)
	s := scanner.NewScanner(r)
	p := parser.NewParser()
	p.InitScanner(s)

	for {
		n, err := p.Parse()
		if n == nil {
			if err != nil {
				fmt.Println(err)
			}
			return
		}
		ast.Print(n)
	}
}
