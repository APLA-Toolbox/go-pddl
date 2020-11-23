package models

import "fmt"

type Location struct {
	Path string
	Line int
}

func (l *Location) ToString() (string, error) {
	if l == nil {
		return "", fmt.Errorf("Failed to stringify location: location is nil")
	}
	return fmt.Sprintf("%s:%d", l.Path, l.Line), nil
}

type PddlError struct {
	Location *Location
	Error    error
}

func (pe *PddlError) ToError() error {
	if pe == nil {
		return fmt.Errorf("Failed to errorify pddl error")
	}
	loc, err := pe.Location.ToString() 
	if err != nil {
		return fmt.Errorf("Failed to errorify pddl error: %v", err)
	}
	return fmt.Errorf("%s: %v", loc, pe.Error)
}
