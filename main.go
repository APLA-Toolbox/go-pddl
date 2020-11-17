package main

import (
	"github.com/guilyx/go-pddl/src/parser"
	"fmt"
)

func main() {
	p := parser.Parser{}
	p.DomainName = "toto"
	fmt.Println(p)
}