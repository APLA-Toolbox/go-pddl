package parser

import "github.com/guilyx/go-pddl/src/models"

func (p *Parser) ParseProblem() *models.PddlError {
	p.ProblemToolbox.Expects("(", "define")
	defer p.ProblemToolbox.Expects(")")
	tk, err := p.ProblemToolbox.PeekNth(2)
	if err != nil {
		return p.ProblemToolbox.NewPddlError("Failed to parse problem: %v", err.Error)
	}
	if tk.Text != "problem" {
		return p.ProblemToolbox.NewPddlError("Failed to parse problem: input file isn't a valid problem.")
	}
	name := p.ProblemToolbox.parseProbName()
	dom := p.ProblemToolbox.parseProbDomain()
	reqs, err := p.ProblemToolbox.parseRequirements()
	if err != nil {
		return p.ProblemToolbox.NewPddlError("Failed to parse problem: %v", err.Error)
	}
	obj := p.ProblemToolbox.parseObjsDecl()
	init := p.ProblemToolbox.parseInit()
	goal := p.ProblemToolbox.parseGoal()
	pb := &models.Problem {
		Domain: dom,
		Goal: goal,
		InitialConditions: init,
		Name: name,
		Objects: obj,
		Requirements: reqs,
	}
	p.Problem = pb
	return nil
}
