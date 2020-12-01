package planner

import (
	"fmt"

	"github.com/guilyx/go-pddl/src/config"
	"github.com/guilyx/go-pddl/src/models"
)

type BreadthFirstSearch struct {
	Configuration *config.Config
	Domain        *models.Domain
	Problem       *models.Problem
}

func NewBFS(config *config.Config, dom *models.Domain, pb *models.Problem) (Planner, error) {
	if config == nil || dom == nil || pb == nil {
		return nil, fmt.Errorf("Failed to create new planner: entities are nil")
	}
	return &BreadthFirstSearch{
		Configuration: config,
		Domain:        dom,
		Problem:       pb,
	}, nil
}

func (dfs *BreadthFirstSearch) Initialize() (interface{}, error) {
	panic("unimplemented")
}

func (dfs *BreadthFirstSearch) Search() (models.Plan, error) {
	panic("unimplemented")
}
