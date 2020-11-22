package models

import (
	"fmt"
	"io"
)

type Domain struct {
	Name         *Name
	Requirements []*Name
	Types        []*Type
	Constants    []*TypedEntry
	Predicates   []*Predicate
	Functions    []*Function
	Actions      []*Action
}

func (d *Domain) PrintDomain(w io.Writer) {
	fmt.Fprintf(w, "(define (domain %s)\n", d.Name.Name)
	printReqsDef(w, d.Requirements)
	printTypesDef(w, d.Types)
	printConstsDef(w, ":constants", d.Constants)
	printPredsDef(w, d.Predicates)
	printFuncsDef(w, d.Functions)
	for _, act := range d.Actions {
		printAction(w, act)
	}
	fmt.Fprintln(w, ")")
}
