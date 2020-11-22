package models

import (
	"fmt"
	"strings"
)

const (
	// objectTypeName is the name of the default
	// object type.
	objectTypeName = "object"

	// totalCostName is the name of the total-cost
	// function.
	totalCostName = "total-cost"
)

// Check returns a slice of all semantic errors in the domain.
//
// If the problem is nil then only the domain is Checked.  The domain must not be nil.
func Check(d *Domain, p *Problem) []error {
	var errs errors
	defs := CheckDomain(d, &errs)
	if p == nil {
		return errs
	}
	if p.Domain.Name != d.Name.Name {
		errs.errorf("problem %s expects domain %s, but got %s",
			p.Name.Name, p.Domain.Name, d.Name.Name)
	}
	CheckReqsDef(defs, p.Requirements, &errs)
	CheckConstsDef(defs, p.Objects, &errs)
	for i := range p.InitialConditions {
		p.InitialConditions[i].Check(defs, &errs)
	}
	p.Goal.Check(defs, &errs)
	// Check the metric
	return errs
}

func CheckDomain(d *Domain, errs *errors) defs {
	defs := defs{
		reqs:   make(map[string]bool),
		types:  make(map[string]*Type),
		consts: make(map[string]*TypedEntry),
		preds:  make(map[string]*Predicate),
		funcs:  make(map[string]*Function),
	}
	CheckReqsDef(defs, d.Requirements, errs)
	CheckTypesDef(defs, d, errs)
	CheckConstsDef(defs, d.Constants, errs)
	CheckPredsDef(defs, d, errs)
	CheckFuncsDef(defs, d.Functions, errs)
	for i := range d.Actions {
		CheckActionDef(defs, d.Actions[i], errs)
	}
	return defs
}

func CheckReqsDef(defs defs, rs []*Name, errs *errors) {
	for _, r := range rs {
		req := strings.ToLower(r.Name)
		if !supportedReqs[req] {
			errs.add(r.Location, "requirement %s is not supported", r)
			continue
		}
		if defs.reqs[req] {
			errs.multipleDefs(r, "requirement")
		}
		defs.reqs[req] = true
	}
	if defs.reqs[":adl"] {
		defs.reqs[":strips"] = true
		defs.reqs[":typing"] = true
		defs.reqs[":negative-preconditions"] = true
		defs.reqs[":disjunctive-preconditions"] = true
		defs.reqs[":equality"] = true
		defs.reqs[":quantified-preconditions"] = true
		defs.reqs[":conditional-effects"] = true
	}
	if defs.reqs[":quantified-preconditions"] {
		defs.reqs[":existential-preconditions"] = true
		defs.reqs[":universal-preconditions"] = true
	}
}

// CheckTypesDef Checks a list of type definitions, maps type names to their definitions, and
// builds the list of all super types of each type.  If the implicit object type was not defined
// then  it is added.
func CheckTypesDef(defs defs, d *Domain, errs *errors) {
	if len(d.Types) > 0 && !defs.reqs[":typing"] {
		errs.badReq(d.Types[0].TypedEntry.Name.Location, ":types", ":typing")
	}
	// Ensure that object is defined
	if !objectDefined(d.Types) {
		d.Types = append(d.Types, &Type{
			TypedEntry: &TypedEntry{
				Name: &Name{Name: objectTypeName},
			},
		})
	}

	// Map type names to their definitions
	for i, t := range d.Types {
		if len(t.TypedEntry.Types) > 1 {
			errs.add(t.TypedEntry.Name.Location, "either super types are not semantically defined")
			continue
		}
		if defs.types[strings.ToLower(t.TypedEntry.Name.Name)] != nil {
			errs.multipleDefs(t.TypedEntry.Name, "type")
			continue
		}
		d.Types[i].TypedEntry.Id = len(defs.types)
		defs.types[strings.ToLower(t.TypedEntry.Name.Name)] = d.Types[i]
	}

	// Link parent types to their definitions
	for i := range d.Types {
		CheckTypeNames(defs, d.Types[i].TypedEntry.Types, errs)
	}

	// Build super type lists
	for i := range d.Types {
		d.Types[i].Predecessors = superTypes(defs, d.Types[i])
		if len(d.Types[i].Predecessors) <= 0 {
			panic("no predecessors!")
		}
	}
}

// ObjectDefined returns true if the object type is in the list of defined types.
func objectDefined(ts []*Type) bool {
	for _, t := range ts {
		if t.TypedEntry.Name.Name == objectTypeName {
			return true
		}
	}
	return false
}

// SuperTypes returns a slice of the parent types of the given type, including the type itself.
func superTypes(defs defs, t *Type) (supers []*Type) {
	seen := make([]bool, len(defs.types))
	stk := []*Type{t}
	for len(stk) > 0 {
		t := stk[len(stk)-1]
		stk = stk[:len(stk)-1]
		if seen[t.TypedEntry.Id] {
			continue
		}
		seen[t.TypedEntry.Id] = true
		supers = append(supers, t)
		for _, s := range t.TypedEntry.Types {
			if s.Definition != nil {
				stk = append(stk, s.Definition)
			}
		}
	}
	if obj := defs.types[objectTypeName]; !seen[obj.TypedEntry.Id] {
		supers = append(supers, obj)
	}
	return
}

// CheckConstsDef Checks a list of constant or object definitions and maps names to their definitions.
func CheckConstsDef(defs defs, objs []*TypedEntry, errs *errors) {
	for i, obj := range objs {
		if defs.consts[strings.ToLower(obj.Name.Name)] != nil {
			errs.multipleDefs(obj.Name, "object")
			continue
		}
		objs[i].Id = len(defs.consts)
		defs.consts[strings.ToLower(obj.Name.Name)] = objs[i]
	}
	CheckTypedEntries(defs, objs, errs)

	// Add the object to the list of objects for its type
	for i := range objs {
		obj := objs[i]
		for _, t := range obj.Types {
			if t.Definition == nil {
				continue
			}
			for _, s := range t.Definition.Predecessors {
				s.addToDomain(obj)
			}
		}
	}
}

// AddToDomain adds an object to the list of all objects of the given type.  If the object has
// already been added then it is not added again.
func (t *Type) addToDomain(obj *TypedEntry) {
	for _, o := range t.Domain {
		if o == obj {
			return
		}
	}
	t.Domain = append(t.Domain, obj)
}

// CheckPredsDef Checks a list of predicate definitions and maps their names to their definitions.
// If :equality is required and the implicit = predicate was not defined then it is added.
func CheckPredsDef(defs defs, d *Domain, errs *errors) {
	if defs.reqs[":equality"] && !equalDefined(d.Predicates) {
		d.Predicates = append(d.Predicates, &Predicate{
			Name: &Name{
				Name: "=",
			},
			Id: len(defs.preds),
			Parameters: []*TypedEntry{
				{Name: &Name{Name: "?x"}},
				{Name: &Name{Name: "?y"}},
			},
		})
	}
	for i, p := range d.Predicates {
		if defs.preds[strings.ToLower(p.Name.Name)] != nil {
			errs.multipleDefs(p.Name, "predicate")
			continue
		}
		CheckTypedEntries(defs, p.Parameters, errs)
		counts := make(map[string]int, len(p.Parameters))
		for _, param := range p.Parameters {
			if counts[param.Name.Name] > 0 {
				errs.multipleDefs(param.Name, "parameter")
			}
			counts[param.Name.Name]++
		}
		d.Predicates[i].Id = len(defs.preds)
		defs.preds[strings.ToLower(p.Name.Name)] = d.Predicates[i]
	}
}

// EqualDefined returns true if the = predicate is in the list of defined predicates.
func equalDefined(ps []*Predicate) bool {
	for _, p := range ps {
		if p.Name.Name == "=" {
			return true
		}
	}
	return false
}

// CheckFuncsDef Checks a list of function definitions and maps their names to their definitions.
func CheckFuncsDef(defs defs, fs []*Function, errs *errors) {
	if len(fs) > 0 && !defs.reqs[":action-costs"] {
		errs.badReq(fs[0].Name.Location, ":functions", ":action-costs")
	}
	for i, f := range fs {
		if defs.funcs[strings.ToLower(f.Name.Name)] != nil {
			errs.multipleDefs(f.Name, "function")
			continue
		}
		CheckTypedEntries(defs, f.Params, errs)
		counts := make(map[string]int, len(f.Params))
		for _, param := range f.Params {
			if counts[param.Name.Name] > 0 {
				errs.multipleDefs(param.Name, "parameter")
			}
			counts[param.Name.Name]++
		}
		fs[i].Id = len(defs.funcs)
		defs.funcs[strings.ToLower(f.Name.Name)] = fs[i]
	}
}

func CheckActionDef(defs defs, act *Action, errs *errors) {
	CheckTypedEntries(defs, act.Params, errs)
	counts := make(map[string]int, len(act.Params))
	for i, param := range act.Params {
		if counts[param.Name.Name] > 0 {
			errs.multipleDefs(param.Name, "parameter")
		}
		counts[param.Name.Name]++
		defs.vars = defs.vars.push(act.Params[i])
	}
	if act.Precondition != nil {
		act.Precondition.Check(defs, errs)
	}
	if act.Effect != nil {
		act.Effect.Check(defs, errs)
	}
	for _ = range act.Params {
		defs.vars.pop()
	}
}

// Push returns a new varDefs with the given definitions defined.
func (v *varDefs) push(d *TypedEntry) *varDefs {
	return &varDefs{
		up:         v,
		name:       d.Name.Name,
		definition: d,
	}
}

// CheckTypedEntries ensures that the types of a list of typed indentifiers are valid.  If they
// are valid then they are linked to their type definitions.  All identifiers that have no declared
// type are linked to the object type.
func CheckTypedEntries(defs defs, lst []*TypedEntry, errs *errors) {
	for i := range lst {
		CheckTypeNames(defs, lst[i].Types, errs)
		if len(lst[i].Types) == 0 {
			lst[i].Types = []*TypeName{{
				Name:       &Name{Name: objectTypeName},
				Definition: defs.types[objectTypeName],
			}}
		}
	}
}

// CheckTypeNames Checks that all of the type names are defined.  Each defined type name
// is linked to its type definition.
func CheckTypeNames(defs defs, ts []*TypeName, errs *errors) {
	if len(ts) > 0 && !defs.reqs[":typing"] {
		errs.badReq(ts[0].Name.Location, "types", ":typing")
	}
	for j, t := range ts {
		switch def := defs.types[strings.ToLower(t.Name.Name)]; def {
		case nil:
			errs.undefined(t.Name, "type")
		default:
			ts[j].Definition = def
		}
	}
}

func (u *UnaryNode) Check(defs defs, errs *errors) {
	u.Formula.Check(defs, errs)
}

func (b *BinaryNode) Check(defs defs, errs *errors) {
	b.Left.Check(defs, errs)
	b.Right.Check(defs, errs)
}

func (m *MultiNode) Check(defs defs, errs *errors) {
	for i := range m.Formula {
		m.Formula[i].Check(defs, errs)
	}
}

func (q *QuantNode) Check(defs defs, errs *errors) {
	CheckTypedEntries(defs, q.Variables, errs)
	counts := make(map[string]int, len(q.Variables))
	for i, v := range q.Variables {
		if counts[v.Name.Name] > 0 {
			errs.multipleDefs(v.Name, "variable")
		}
		counts[v.Name.Name]++
		defs.vars = defs.vars.push(q.Variables[i])
	}
	q.UnaryNode.Check(defs, errs)
	for _ = range q.Variables {
		defs.vars = defs.vars.pop()
	}
}

// Pop returns a varDefs with the latest definition removed.
func (v *varDefs) pop() *varDefs {
	return v.up
}

func (n *OrNode) Check(defs defs, errs *errors) {
	if !defs.reqs[":disjunctive-preconditions"] {
		errs.badReq(n.MultiNode.Location, "or", ":disjunctive-preconditions")
	}
	n.MultiNode.Check(defs, errs)
}

func (i *ImplyNode) Check(defs defs, errs *errors) {
	if !defs.reqs[":disjunctive-preconditions"] {
		errs.badReq(i.BinaryNode.Location, "imply", ":disjunctive-preconditions")
	}
	i.BinaryNode.Check(defs, errs)
}

func (f *ForAllNode) Check(defs defs, errs *errors) {
	switch {
	case !f.IsEffect && !defs.reqs[":universal-preconditions"]:
		errs.badReq(f.QuantNode.UnaryNode.Node.Location, "forall", ":universal-preconditions")
	case f.IsEffect && !defs.reqs[":conditional-effects"]:
		errs.badReq(f.QuantNode.UnaryNode.Node.Location, "forall", ":conditional-effects")
	}
	f.QuantNode.Check(defs, errs)
}

func (e *ExistsNode) Check(defs defs, errs *errors) {
	if !defs.reqs[":existential-preconditions"] {
		errs.badReq(e.QuantNode.UnaryNode.Node.Location, "exists", ":existential-preconditions")
	}
	e.QuantNode.Check(defs, errs)
}

func (w *WhenNode) Check(defs defs, errs *errors) {
	if !defs.reqs[":conditional-effects"] {
		errs.badReq(w.UnaryNode.Node.Location, "when", ":conditional-effects")
	}
	w.Condition.Check(defs, errs)
	w.UnaryNode.Check(defs, errs)
}

func (lit *LiteralNode) Check(defs defs, errs *errors) {
	if lit.Definition = defs.preds[strings.ToLower(lit.Predicate.Name)]; lit.Definition == nil {
		errs.undefined(lit.Predicate, "predicate")
		return
	}
	if lit.IsEffect {
		if lit.Negative {
			lit.Definition.NegEffect = true
		} else {
			lit.Definition.PosEffect = true
		}
	}
	CheckInst(defs, lit.Predicate, lit.Terms, lit.Definition.Parameters, errs)
}

func CheckInst(defs defs, n *Name, args []*Term, params []*TypedEntry, errs *errors) {
	if len(args) != len(params) {
		var argName = "arguments"
		if len(params) == 1 {
			argName = argName[:len(argName)-1]
		}
		errs.add(n.Location, "%s requires %d %s", n, len(params), argName)
	}

	for i := range args {
		kind := "constant"
		args[i].Definition = defs.consts[strings.ToLower(args[i].Name.Name)]
		if args[i].IsVariable {
			args[i].Definition = defs.vars.find(args[i].Name.Name)
			kind = "variable"
		}
		if args[i].Definition == nil {
			errs.undefined(args[i].Name, kind)
			return
		}
		if !compatTypes(params[i].Types, args[i].Definition.Types) {
			errs.add(args[i].Name.Location,
				"%s [type %s] is incompatible with parameter %s [type %s] of %s",
				args[i], typeString(args[i].Definition.Types),
				params[i], typeString(params[i].Types), n)
		}
	}
}

func (v *varDefs) find(n string) *TypedEntry {
	if v == nil {
		return nil
	}
	if strings.ToLower(v.name) == strings.ToLower(n) {
		return v.definition
	}
	return v.up.find(n)
}

func compatTypes(left, right []*TypeName) bool {
	for _, r := range right {
		if r.Definition == nil {
			return true
		}
		ok := false
		for _, l := range left {
			if l.Definition == nil {
				return true
			}
			for _, s := range r.Definition.Predecessors {
				if s == l.Definition {
					ok = true
					break
				}
			}
			if ok {
				break
			}
		}
		if !ok {
			return false
		}
	}
	return true
}

func (a *AndNode) Check(defs defs, errs *errors) {
	fmt.Println("AndNodes check isn't implemented yet")
}

func (a *AssignNode) Check(defs defs, errs *errors) {
	if !defs.reqs[":action-costs"] {
		errs.badReq(a.Node.Location, a.Operation.Name, ":action-costs")
	}
	a.AssignedTo.Check(defs, errs)
	if a.IsNumber {
		if negative(a.Number) {
			errs.add(a.Node.Location, "assigned value must not be negative with :action-costs")
		}
	} else {
		a.FunctionInit.Check(defs, errs)
	}

	if !a.IsInit {
		if a.AssignedTo.Definition != nil && !a.AssignedTo.Definition.isTotalCost() {
			errs.add(a.AssignedTo.Name.Location, "assignment target must be a 0-ary total-cost function with :action-costs")
		}
		if !a.IsNumber && a.FunctionInit.Definition != nil && a.FunctionInit.Definition.isTotalCost() {
			errs.add(a.FunctionInit.Name.Location, "assigned value must not be total-cost with :action-costs")
		}
	}
}

func (f *Function) isTotalCost() bool {
	return f.Name.Name == totalCostName && len(f.Params) == 0
}

// Negative returns true if the string is a negative number.
func negative(n string) bool {
	neg := false
	for _, s := range n {
		if s != '-' {
			break
		}
		neg = !neg
	}
	return neg
}

func (h *FunctionInit) Check(defs defs, errs *errors) {
	if h.Definition = defs.funcs[strings.ToLower(h.Name.Name)]; h.Definition == nil {
		errs.undefined(h.Name, " function")
		return
	}
	CheckInst(defs, h.Name, h.Terms, h.Definition.Params, errs)
}

// Errors wraps a slice of errors.
type errors []error

func (es *errors) add(l *Location, f string, vs ...interface{}) {
	ls, _ := l.ToString()
	f = ls + ": " + f
	*es = append(*es, fmt.Errorf(fmt.Sprintf(f, vs...)))
}

func (es *errors) errorf(f string, vs ...interface{}) {
	*es = append(*es, fmt.Errorf(f, vs...))
}

func (es *errors) undefined(name *Name, kind string) {
	es.add(name.Location, "undefined %s %s", kind, name.Name)
}

func (es *errors) multipleDefs(name *Name, kind string) {
	es.add(name.Location, "%s %s defined multiple times", kind, name.Name)
}

func (es *errors) badReq(l *Location, used, reqd string) {
	*es = append(*es, MissingRequirementError{
		Location:    l,
		Cause:       used,
		Requirement: reqd,
	})
}

type MissingRequirementError struct {
	Location    *Location
	Cause       string
	Requirement string
}

func (r MissingRequirementError) Error() string {
	l, _ := r.Location.ToString()
	return l + ": " + r.Cause + " requires " + r.Requirement
}
