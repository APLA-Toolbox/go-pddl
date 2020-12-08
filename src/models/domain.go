package models

import (
	"fmt"
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

func (d *Domain) PrintDomain() {
	var s string
	if d == nil {
		panic("Domain is nil, can't print")
	}
	s += fmt.Sprintf("(define (domain %s)\n", d.Name.Name)
	s += toStringReqs(d.Requirements)
	s += toStringTypesDef(d.Types)
	s += toStringConsts(":constants", d.Constants)
	s += toStringPredicates(d.Predicates)
	s += toStringFunctions(d.Functions)
	for _, act := range d.Actions {
		s += toStringAction(act)
	}
	s += ")\n"
	fmt.Println(s)
}
