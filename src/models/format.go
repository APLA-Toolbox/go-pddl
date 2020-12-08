package models

import (
	"fmt"
)

func toStringReqs(reqs []*Name) string {
	if len(reqs) == 0 {
		return ""
	}
	s := fmt.Sprintf("%s(:requirements\n", Indent(1))
	for i, r := range reqs {
		sTemp := r.Name
		if i == len(reqs)-1 {
			sTemp += ")"
		}
		s += fmt.Sprintf("%s %s\n", Indent(2), s)
	}
	return s
}

func toStringTypesDef(ts []*Type) string {
	if len(ts) == 0 {
		return ""
	}
	s := fmt.Sprintf("%s(:types", Indent(1))
	ids := []*TypedEntry{}
	for _, t := range ts {
		if t.TypedEntry.Name.Location.Line == 0 {
			// Skip undeclared implicit types like object.
			continue
		}
		ids = append(ids, t.TypedEntry)
	}
	s += toStringTypedNames("\n"+Indent(2), ids)
	s += ")\n"
	return s
}

func toStringConsts(def string, cs []*TypedEntry) string {
	if len(cs) == 0 {
		return ""
	}
	var s string
	s += fmt.Sprintf("%s(%s", Indent(1), def)
	s += toStringTypedNames("\n"+Indent(2), cs)
	s += ")\n"
	return s
}

func toStringPredicates(ps []*Predicate) string {
	var s string
	if len(ps) == 0 {
		return ""
	}
	s += fmt.Sprintf("%s(:predicates\n", Indent(1))
	for i, p := range ps {
		if p.Name.Location.Line == 0 {
			continue
		}
		s += fmt.Sprintf("%s(%s", Indent(2), p.Name.Name)
		s += toStringTypedNames(" ", p.Parameters)
		s += ")"
		if i < len(ps)-1 {
			s += "\n"
		}
	}
	s += ")\n"
	return s
}

func toStringFunctions(fs []*Function) string {
	var s string
	if len(fs) == 0 {
		return ""
	}
	s += fmt.Sprintf("%s(:functions\n", Indent(1))
	for i, f := range fs {
		s += fmt.Sprintf("%s(%s", Indent(2), f.Name.Name)
		s += toStringTypedNames(" ", f.Params)
		s += ")"
		if len(f.Types) > 0 {
			s += fmt.Sprintf(" - %s\n", toStringType(f.Types))
		}
		if i < len(fs)-1 {
			s += "\n"
		}
	}
	s += fmt.Sprintf(")")
	return s
}

func toStringAction(act *Action) string {
	var s string
	s += fmt.Sprintf("%s(:action %s\n", Indent(1), act.Name.Name)
	s += fmt.Sprintf("%s:parameters (", Indent(2))
	s += toStringTypedNames("", act.Params)
	s += ")"
	if act.Precondition != nil {
		s += "\n"
		s += fmt.Sprintf("%s:precondition\n", Indent(2))
		s += act.Precondition.ToString(Indent(3))
	}
	if act.Effect != nil {
		s += "\n"
		s += fmt.Sprintf("%s:effect\n", Indent(2))
		s += act.Effect.ToString(Indent(3))
	}
	s += ")\n"
	return s
}

func toStringTypedNames(prefix string, ns []*TypedEntry) string {
	var s string
	if len(ns) == 0 {
		return ""
	}
	tprev := toStringType(ns[0].Types)
	sep := prefix
	for _, n := range ns {
		tcur := toStringType(n.Types)
		if tcur != tprev {
			if tprev == "" {
				// Should be impossible.
				str, _ := n.Name.Location.ToString()
				panic(str + ": untyped declarations in the middle of a typed list")
			}
			s += fmt.Sprintf(" - %s", tprev)
			tprev = tcur
			sep = prefix
			if sep == "" {
				sep = " "
			}
		}
		s += fmt.Sprintf("%s%s", sep, n.Name.Name)
		sep = " "
	}
	if tprev != "" {
		s += fmt.Sprintf(" - %s", tprev)
	}
	return s
}

func toStringType(t []*TypeName) string {
	var str string
	switch len(t) {
	case 0:
		break
	case 1:
		if t[0].Name.Location.Line == 0 {
			break
		}
		str = t[0].Name.Name
	default:
		str = "(either"
		for _, n := range t {
			str += " " + n.Name.Name
		}
		str += ")"
	}
	return str
}
