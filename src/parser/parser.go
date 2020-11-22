package parser

import "github.com/guilyx/go-pddl/src/config"

type Parser struct {
	Configuration *config.Config
}

func NewParser(config *config.Config) (*Parser, error) {
	return &Parser{Configuration: config}, nil
}
