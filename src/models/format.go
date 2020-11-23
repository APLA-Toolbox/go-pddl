package models

import (
	"fmt"
	"io"
)

func printReqsDef(w io.Writer, reqs []*Name) {
	if len(reqs) == 0 {
		return
	}
	fmt.Fprintf(w, "%s(:requirements\n", Indent(1))
	for i, r := range reqs {
		s := r.Name
		if i == len(reqs)-1 {
			s += ")"
		}
		fmt.Fprintln(w, Indent(2), s)
	}
}

func printTypesDef(w io.Writer, ts []*Type) {
	if len(ts) == 0 {
		return
	}
	fmt.Fprintf(w, "%s(:types", Indent(1))
	ids := []*TypedEntry{}
	for _, t := range ts {
		if t.TypedEntry.Name.Location.Line == 0 {
			// Skip undeclared implicit types like object.
			continue
		}
		ids = append(ids, t.TypedEntry)
	}
	printTypedNames(w, "\n"+Indent(2), ids)
	fmt.Fprintln(w, ")")
}

// PrintConstsDef prints a constant definition with the given definition name
// (should be either :constants or :objects).
func printConstsDef(w io.Writer, def string, cs []*TypedEntry) {
	if len(cs) == 0 {
		return
	}
	fmt.Fprintf(w, "%s(%s", Indent(1), def)
	printTypedNames(w, "\n"+Indent(2), cs)
	fmt.Fprintln(w, ")")
}

func printPredsDef(w io.Writer, ps []*Predicate) {
	if len(ps) == 0 {
		return
	}
	fmt.Fprintf(w, "%s(:predicates\n", Indent(1))
	for i, p := range ps {
		if p.Name.Location.Line == 0 {
			// Skip undefined implicit predicates like =.
			continue
		}
		fmt.Fprintf(w, "%s(%s", Indent(2), p.Name.Name)
		printTypedNames(w, " ", p.Parameters)
		fmt.Fprint(w, ")")
		if i < len(ps)-1 {
			fmt.Fprint(w, "\n")
		}
	}
	fmt.Fprintln(w, ")")
}

func printFuncsDef(w io.Writer, fs []*Function) {
	if len(fs) == 0 {
		return
	}
	fmt.Fprintf(w, "%s(:functions\n", Indent(1))
	for i, f := range fs {
		fmt.Fprintf(w, "%s(%s", Indent(2), f.Name.Name)
		printTypedNames(w, " ", f.Params)
		fmt.Fprint(w, ")")
		if len(f.Types) > 0 {
			fmt.Fprint(w, " - ", typeString(f.Types))
		}
		if i < len(fs)-1 {
			fmt.Fprint(w, "\n")
		}
	}
	fmt.Fprintln(w, ")")
}

func printAction(w io.Writer, act *Action) {
	fmt.Fprintf(w, "%s(:action %s\n", Indent(1), act.Name.Name)
	fmt.Fprintf(w, "%s:parameters (", Indent(2))
	printTypedNames(w, "", act.Params)
	fmt.Fprint(w, ")")
	if act.Precondition != nil {
		fmt.Fprint(w, "\n")
		fmt.Fprintf(w, "%s:precondition\n", Indent(2))
		act.Precondition.Print(w, Indent(3))
	}
	if act.Effect != nil {
		fmt.Fprint(w, "\n")
		fmt.Fprintf(w, "%s:effect\n", Indent(2))
		act.Effect.Print(w, Indent(3))
	}
	fmt.Fprintln(w, ")")
}

// DeclGroup is a group of declarators along with their type.
type declGroup struct {
	typ  string
	ents []string
}

// DeclGroups implements sort.Interface, sorting the list of typed declarations by their type name.
type declGroups []declGroup

func (t declGroups) Len() int {
	return len(t)
}

func (t declGroups) Less(i, j int) bool {
	return t[i].typ < t[j].typ
}

func (t declGroups) Swap(i, j int) {
	t[i], t[j] = t[j], t[i]
}

// PrintTypedNames prints a slice of TypedNames. Adjacent items with the same type are
// all printed in a group.  Each group is preceeded by the prefix.
func printTypedNames(w io.Writer, prefix string, ns []*TypedEntry) {
	if len(ns) == 0 {
		return
	}
	tprev := typeString(ns[0].Types)
	sep := prefix
	for _, n := range ns {
		tcur := typeString(n.Types)
		if tcur != tprev {
			if tprev == "" {
				// Should be impossible.
				str, _ := n.Name.Location.ToString()
				panic(str + ": untyped declarations in the middle of a typed list")
			}
			fmt.Fprintf(w, " - %s", tprev)
			tprev = tcur
			sep = prefix
			if sep == "" {
				sep = " "
			}
		}
		fmt.Fprintf(w, "%s%s", sep, n.Name.Name)
		sep = " "
	}
	if tprev != "" {
		fmt.Fprintf(w, " - %s", tprev)
	}
}

// TypeString returns the string representation of a type.
func typeString(t []*TypeName) (str string) {
	switch len(t) {
	case 0:
		break
	case 1:
		if t[0].Name.Location.Line == 0 {
			// Use the empty string for undeclared
			// implicit types (such as object).
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
	return
}
