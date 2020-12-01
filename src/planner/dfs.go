package planner

import (
	"fmt"

	"github.com/guilyx/go-pddl/src/config"
	"github.com/guilyx/go-pddl/src/models"
)

type DepthFirstSearch struct {
	Configuration *config.Config
	Domain        *models.Domain
	Problem       *models.Problem
}

func NewDFS(config *config.Config, dom *models.Domain, pb *models.Problem) (Planner, error) {
	if config == nil || dom == nil || pb == nil {
		return nil, fmt.Errorf("Failed to create new planner: entities are nil")
	}
	return &DepthFirstSearch{
		Configuration: config,
		Domain:        dom,
		Problem:       pb,
	}, nil
}

func (dfs *DepthFirstSearch) Initialize() (interface{}, error) {
	panic("unimplemented")
}

func (dfs *DepthFirstSearch) Search() (models.Plan, error) {
	panic("unimplemented")
}
