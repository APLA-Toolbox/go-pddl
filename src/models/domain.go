package models

type Domain struct {
	Name         *Name
	Requirements []*Name
	Types        []*Type
	Constants    []*TypedEntry
	Predicate    []*Predicate
	Functions    []*Function
	Actions      []*Action
}
