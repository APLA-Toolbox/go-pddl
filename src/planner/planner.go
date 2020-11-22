package planner

import "github.com/guilyx/go-pddl/src/config"

type Planner struct {
	Configuration *config.Config
}

func NewPlanner(config *config.Config) (*Planner, error) {
	return &Planner{Configuration: config}, nil
}
