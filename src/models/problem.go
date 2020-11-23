package models

import (
	"fmt"
	"io"
)

type Problem struct {
	Name              *Name
	Domain            *Name
	Requirements      []*Name
	Objects           []*TypedEntry
	InitialConditions []Formula
	Goal              Formula
}

func (p *Problem) PrintProblem(w io.Writer) {
	fmt.Fprintf(w, "(define (problem %s)\n%s(:domain %s)\n",
		p.Name.Name, Indent(1), p.Domain.Name)
	printReqsDef(w, p.Requirements)
	printConstsDef(w, ":objects", p.Objects)

	fmt.Fprintf(w, "%s(:init", Indent(1))
	for _, f := range p.InitialConditions {
		fmt.Fprint(w, "\n")
		f.Print(w, Indent(2))
	}
	fmt.Fprint(w, ")\n")

	fmt.Fprintf(w, "%s(:goal\n", Indent(1))
	p.Goal.Print(w, Indent(2))

	fmt.Fprintln(w, ")\n)")
}
