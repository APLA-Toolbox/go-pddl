package parser

import (
	"github.com/guilyx/go-pddl/src/models"
)

func (p *Parser) ParseDomain() *models.PddlError {
	name, err := p.DomainToolbox.parseDomainName()
	if err != nil {
		return p.DomainToolbox.NewPddlError("Failed to parse domain: %v", err.Error)
	}
	reqs, err := p.DomainToolbox.parseRequirements()
	if err != nil {
		return p.DomainToolbox.NewPddlError("Failed to parse domain: %v", err.Error)
	}
	typs, err := p.DomainToolbox.parseTypesDefinition()
	if err != nil {
		return p.DomainToolbox.NewPddlError("Failed to parse domain: %v", err.Error)
	}
	csts, err := p.DomainToolbox.parseConstantsDefinition()
	if err != nil {
		return p.DomainToolbox.NewPddlError("Failed to parse domain: %v", err.Error)
	}
	preds, err := p.DomainToolbox.parsePredicatesDefinition()
	if err != nil {
		return p.DomainToolbox.NewPddlError("Failed to parse domain: %v", err.Error)
	}
	funcs := p.DomainToolbox.parseFuncsDef()
	if err != nil {
		return p.DomainToolbox.NewPddlError("Failed to parse domain: %v", err.Error)
	}
	acts := p.DomainToolbox.parseActionsDef()
	if err != nil {
		return p.DomainToolbox.NewPddlError("Failed to parse domain: %v", err.Error)
	}
	d := &models.Domain{
		Name:         name,
		Actions:      acts,
		Constants:    csts,
		Functions:    funcs,
		Predicates:   preds,
		Requirements: reqs,
		Types:        typs,
	}
	p.Domain = d
	return nil
}
