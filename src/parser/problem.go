package parser

import "github.com/guilyx/go-pddl/src/models"

func (p *Parser) ParseProblem() *models.PddlError {
	pb := &models.Problem{}
	p.Problem = pb
	return nil
}
