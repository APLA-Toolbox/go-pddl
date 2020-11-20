package main

import (
	"fmt"
	"github.com/guilyx/go-pddl/src/parser"
	"github.com/guilyx/go-pddl/src/common"
	"github.com/akamensky/argparse"
)

func main() {
	// Handle input
	parser := argparse.NewParser("PDDL Plarser", "Parses PDDL a domain/problem combinaison and 
								 computes an optimal plan")
	dm := parser.String("dm", "domain", &argparse.Options{Required: true, Help: "Path to domain.pddl file"})
	pb := parser.String("pb", "problem", &argparse.Options{Required: true, Help: "Path to the problem.pddl file"})
	err := parser.Parse(os.Args)
	if err != nil {
		common.CheckError(err)
	}
	
	// Parse
	p := parser.Parser{}
	p.DomainName = "toto"
	fmt.Println(p)

	// Plan
}
