package models

type Domain struct {
	Name         *Name
	Requirements []*Name
	Types        []*Type
	Constants    []*TypedEntry
	Predicates   []*Predicate
	Functions    []*Function
	Actions      []*Action
}
