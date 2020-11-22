package parser

import (
	"fmt"

	"github.com/guilyx/go-pddl/src/common"
	"github.com/guilyx/go-pddl/src/config"
	"github.com/guilyx/go-pddl/src/lexer"
	"github.com/guilyx/go-pddl/src/models"
)

type Parser struct {
	Domain         *models.Domain
	Problem        *models.Problem
	DomainToolbox  *ParserToolbox
	ProblemToolbox *ParserToolbox
}

func NewParser() *Parser {
	return &Parser{}
}

func (p *Parser) RegisterDomain(config *config.Config) error {
	if p == nil {
		return fmt.Errorf("Failed to register domain: parser is nil")
	}
	text, err := common.LoadFile(config.Domain)
	if err != nil {
		return fmt.Errorf("Failed to register domain: %v", err)
	}
	l, err := lexer.NewLexer(config.Domain, text)
	if err != nil {
		return fmt.Errorf("Failed to register domain: %v", err)
	}
	p.DomainToolbox, err = NewParserToolbox(config, l)
	if err != nil {
		return fmt.Errorf("Failed to register domain: %v", err)
	}
	return nil
}

func (p *Parser) RegisterProblem(config *config.Config) error {
	if p == nil {
		return fmt.Errorf("Failed to register problem: parser is nil")
	}
	text, err := common.LoadFile(config.Problem)
	if err != nil {
		return fmt.Errorf("Failed to register problem: %v", err)
	}
	l, err := lexer.NewLexer(config.Problem, text)
	if err != nil {
		return fmt.Errorf("Failed to register problem: %v", err)
	}
	p.ProblemToolbox, err = NewParserToolbox(config, l)
	if err != nil {
		return fmt.Errorf("Failed to register problem: %v", err)
	}
	return nil
}
