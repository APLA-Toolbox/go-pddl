package models

type Predicate struct {
	Name       *Name
	Id         int
	Parameters []*TypedEntry
	PosEffect  bool
	NegEffect  bool
}

type Action struct {
	Name         *Name
	Params       []*TypedEntry
	Precondition Formula
	Effect       Formula
}
