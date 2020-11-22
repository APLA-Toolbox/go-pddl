package models

type Problem struct {
	Name              *Name
	Domain            *Name
	Requirements      []*Name
	Objects           []*TypedEntry
	InitialConditions []Formula
	Goal              Formula
}
