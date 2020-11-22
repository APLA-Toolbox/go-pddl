package parser

import (
	"github.com/guilyx/go-pddl/src/models"
)

func (p *Parser) ParseDomain() *models.PddlError {
	d := &models.Domain{}
	p.Domain = d
	return nil
}
