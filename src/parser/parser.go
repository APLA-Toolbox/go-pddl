package parser

import (
	"github.com/guilyx/go-pddl/src/common"
)

type Parser struct {
	DomainName   string
	Requirements common.StringSlice
	Types        common.StringSlice
	Actions      common.StringSlice
	Predicates   common.StringMap
}

func (p *Parser) LoadDomain(filename string) (string, error) {
	s, err := common.LoadFile(filename)
	if err != nil {
		return "", err
	}
	return s, nil
}

func (p *Parser) LoadProblem(filename string) (string, error) {
	s, err := common.LoadFile(filename)
	if err != nil {
		return "", err
	}
	return s, nil
}
