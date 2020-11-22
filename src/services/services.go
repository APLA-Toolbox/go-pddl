package services

import (
	"fmt"
	"time"

	"github.com/guilyx/go-pddl/src/config"
	"github.com/guilyx/go-pddl/src/parser"
	"github.com/guilyx/go-pddl/src/planner"
)

type Pddl struct {
	Parser  *parser.Parser
	Planner *planner.Planner
}

func Start() (*Pddl, error) {
	conf, err := config.NewConfig()
	if err != nil {
		return nil, fmt.Errorf("Failed to instantiate configuration: %v", err)
	}
	parser, err := parser.NewParser(conf)
	if err != nil {
		return nil, fmt.Errorf("Failed to build parser: %v", err)
	}
	planner, err := planner.NewPlanner(conf)
	if err != nil {
		return nil, fmt.Errorf("Failed to build planner: %v", err)
	}
	fmt.Println("Starting go-pddl... (v" + conf.Version + ", started at " + time.Now().String())

	return &Pddl{
		Parser:  parser,
		Planner: planner,
	}, nil
}
