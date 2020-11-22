package main

import (
	"fmt"
	"os"

	"github.com/guilyx/go-pddl/src/services"
)

func main() {
	pddl, err := services.Start()
	if err != nil {
		fmt.Printf("Initialization failed: %v", err)
		panic("Exit failure")
	}

	// Parse
	errPddl := pddl.Parser.ParseDomain()
	if errPddl != nil {
		fmt.Println(*errPddl)
		panic("Failed to parse domain")
	}
	errPddl = pddl.Parser.ParseProblem()
	if errPddl != nil {
		fmt.Println("Failed to parse problem")
		panic("Failed to parse problem")
	}

	pddl.Parser.Domain.PrintDomain(os.Stdout)
	fmt.Printf("\n#################################################################\n")
	pddl.Parser.Problem.PrintProblem(os.Stdout)

	// Plan
	//...
}
