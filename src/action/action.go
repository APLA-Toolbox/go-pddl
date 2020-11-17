package Action

import (
	"github.com/guilyx/go-pddl/src/common"
)

type Action struct {
	Name                  string
	Parameters            common.StringSlice
	PositivePreconditions common.StringSlice
	NegativePreconditions common.StringSlice
	AddEffects            common.StringSlice
	DelEffects            common.StringSlice
}
