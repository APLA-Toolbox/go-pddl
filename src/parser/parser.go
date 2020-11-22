package parser

import (
	"fmt"
	"io"
	"io/ioutil"

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

func (p *Parser) RegisterDomain(file string, r io.Reader, config *config.Config) error {
	if p == nil {
		return fmt.Errorf("Failed to register domain: parser is nil")
	}
	text, err := ioutil.ReadAll(r)
	if err != nil {
		return fmt.Errorf("Failed to register domain: %v", err)
	}
	l, err := lexer.NewLexer(file, string(text))
	if err != nil {
		return fmt.Errorf("Failed to register domain: %v", err)
	}
	p.DomainToolbox, err = NewParserToolbox(config, l)
	if err != nil {
		return fmt.Errorf("Failed to register domain: %v", err)
	}
	return nil
}

func (p *Parser) RegisterProblem(file string, r io.Reader, config *config.Config) error {
	if p == nil {
		return fmt.Errorf("Failed to register problem: parser is nil")
	}
	text, err := ioutil.ReadAll(r)
	if err != nil {
		return fmt.Errorf("Failed to register problem: %v", err)
	}
	l, err := lexer.NewLexer(file, string(text))
	if err != nil {
		return fmt.Errorf("Failed to register problem: %v", err)
	}
	p.ProblemToolbox, err = NewParserToolbox(config, l)
	if err != nil {
		return fmt.Errorf("Failed to register problem: %v", err)
	}
	return nil
}
