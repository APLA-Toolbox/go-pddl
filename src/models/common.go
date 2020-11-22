package models

import "fmt"

type (
	defs struct {
		reqs   map[string]bool
		types  map[string]*Type
		consts map[string]*TypedEntry
		preds  map[string]*Predicate
		funcs  map[string]*Function
		vars   *varDefs
	}

	varDefs struct {
		up         *varDefs
		name       string
		definition *TypedEntry
	}
)

type Name struct {
	Name     string
	Location *Location
}

func (n *Name) ToString() (string, error) {
	if n == nil {
		return "", fmt.Errorf("Failed to stringify, name is nil")
	}
	return n.Name, nil
}

type TypeName struct {
	Name       *Name
	Definition *Type
}

type TypedEntry struct {
	Name  *Name
	Id    int
	Types []*TypeName
}

type Type struct {
	TypedEntry   *TypedEntry
	Predecessors []*Type
	Domain       []*TypedEntry
}

type FunctionInit struct {
	Name       *Name
	Terms      []*Term
	Definition *Function
}

type Term struct {
	Name       *Name
	IsConstant bool
	Definition *TypedEntry
}

type Function struct {
	Name   *Name
	Id     int
	Types  []*TypeName
	Params []*TypedEntry
}
