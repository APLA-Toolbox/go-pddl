package planner

import "github.com/guilyx/go-pddl/src/models"

type Planner interface {
	Initialize() (interface{}, error)
	Search() (models.Plan, error)
}
