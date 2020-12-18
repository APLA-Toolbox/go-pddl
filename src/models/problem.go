package models

import (
	"fmt"
)

type Problem struct {
	Name              *Name
	Domain            *Name
	Requirements      []*Name
	Objects           []*TypedEntry
	InitialConditions []Formula
	Goal              Formula
}

func (p *Problem) PrintProblem() {
	var s string
	s += fmt.Sprintf("(define (problem %s)\n%s(:domain %s)\n",
		p.Name.Name, Indent(1), p.Domain.Name)
	s += toStringReqs(p.Requirements)
	s += toStringConsts(":objects", p.Objects)
	s += fmt.Sprintf("%s(:init", Indent(1))
	for _, f := range p.InitialConditions {
		s += "\n"
		s += f.ToString(Indent(2))
	}
	s += ")\n"
	s += fmt.Sprintf("%s(:goal\n", Indent(1))
	s += p.Goal.ToString(Indent(2))
	s += ")\n)\n"
	fmt.Println(s)
}

func (p *Problem) ToJSONProblem() {
	var s string
	s += "{"
	s += fmt.Sprintf("\"define\":{\"problem\":\"%s\",\"domain\":\"%s\"},", p.Name.Name, p.Domain.Name)
	s += toJSONReqs(p.Requirements)
	s += toJSONConsts("objects", p.Objects)
	s += "\"init\":{"
	for _, f := range p.InitialConditions {
		s += f.ToJSON("")
	}
	s += "}"
	s += "\"goal\":{"
	s += p.Goal.ToJSON("")
	s += "}}"
	fmt.Println(s)
}
