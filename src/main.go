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
		fmt.Println(errPddl.ToError())
		return
	}
	fmt.Println("Domain successfully parsed...")
	errPddl = pddl.Parser.ParseProblem()
	if errPddl != nil {
		fmt.Println(errPddl.ToError())
		panic("Failed to parse problem")
	}
	fmt.Println("Problem successfully parsed...")

	if pddl.Parser.DomainToolbox.Configuration.PrintPddl {
		fmt.Printf("\n#################################################################")
		fmt.Printf("\n####################### D O M A I N #############################")
		fmt.Printf("\n#################################################################\n\n")
		pddl.Parser.Domain.PrintDomain(os.Stdout)
		fmt.Printf("\n#################################################################")
		fmt.Printf("\n###################### P R O B L E M ############################")
		fmt.Printf("\n#################################################################\n\n")
		pddl.Parser.Problem.PrintProblem(os.Stdout)
	}

	// Plan
	//...
}
