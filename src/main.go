package main

import (
	"fmt"

	"github.com/guilyx/go-pddl/src/services"
)

func main() {
	pddl, err := services.Start()
	if err != nil {
		fmt.Printf("Initialization failed: %v", err)
		panic("Exit failure")
	}

	// Parse
	d, errPddl := pddl.Parser.ParseDomain()
	if errPddl != nil {
		fmt.Println(errPddl.ToError())
		return
	}
	fmt.Println("Domain successfully parsed...")
	pb, errPddl := pddl.Parser.ParseProblem()
	if errPddl != nil {
		fmt.Println(errPddl.ToError())
		panic("Failed to parse problem")
	}
	fmt.Println("Problem successfully parsed...")

	d.ToJSONDomain()
	fmt.Printf("\n\n")
	pb.ToJSONProblem()

	// Plan
	// err = pddl.RegisterPlanner(d, pb)
	// if err != nil {
	// 	panic("Exit failure")
	// }
}
