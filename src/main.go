package main

import (
	"fmt"

	"github.com/guilyx/go-pddl/src/services"
)

func main() {
	_, err := services.Start()
	if err != nil {
		fmt.Printf("Initialization failed: %v", err)
		panic("Exit failure")
	}

	fmt.Println("Under construction, nothing to see here.")

	// Parse
	// ...

	// Plan
	//...
}
