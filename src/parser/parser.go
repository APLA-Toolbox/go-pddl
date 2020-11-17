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
