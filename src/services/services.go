package services

import (
	"fmt"
	"time"

	"github.com/guilyx/go-pddl/src/config"
	"github.com/guilyx/go-pddl/src/models"
	"github.com/guilyx/go-pddl/src/parser"
	"github.com/guilyx/go-pddl/src/planner"
)

type Pddl struct {
	Parser  *parser.Parser
	Planner *planner.Planner
}

func (p *Pddl) RegisterPlanner(d *models.Domain, pb *models.Problem) error {
	var conf *config.Config
	if p.Parser.DomainToolbox != nil {
		conf = p.Parser.DomainToolbox.Configuration
	} else if p.Parser.ProblemToolbox != nil {
		conf = p.Parser.ProblemToolbox.Configuration
	}
	planner, err := planner.NewPlanner(conf, d, pb)
	if err != nil {
		return fmt.Errorf("Failed to build planner: %v", err)
	}
	p.Planner = planner
	return nil
}

func Start() (*Pddl, error) {
	conf, err := config.NewConfig()
	if err != nil {
		return nil, fmt.Errorf("Failed to instantiate configuration: %v", err)
	}
	fmt.Println("Starting go-pddl... (v " + conf.Version + ", started at " + time.Now().String())

	// Parser Creation
	parser := parser.NewParser()
	err = parser.RegisterDomain(conf)
	if err != nil {
		return nil, err
	}
	err = parser.RegisterProblem(conf)
	if err != nil {
		return nil, err
	}

	return &Pddl{
		Parser:  parser,
	}, nil
}
