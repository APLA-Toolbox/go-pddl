package models

import (
	"fmt"
)

var (
	AssignOps = map[string]bool{
		"=":        true,
		"assign":   true,
		"increase": true,
	}
)

type Formula interface {
	ToString(string) string
	ToJSON(string) string
}

type Node struct {
	Location *Location
}

type UnaryNode struct {
	Node    *Node
	Formula Formula
}

type BinaryNode struct {
	Node
	Left  Formula
	Right Formula
}

type MultiNode struct {
	Node
	Formula []Formula
}

type QuantNode struct {
	Variables []*TypedEntry
	UnaryNode *UnaryNode
}

type LiteralNode struct {
	Node       *Node
	Predicate  *Name
	Negative   bool
	Terms      []*Term
	IsEffect   bool
	Definition *Predicate
}

type AndNode struct {
	MultiNode *MultiNode
}

type OrNode struct {
	MultiNode *MultiNode
}

type NotNode struct {
	UnaryNode *UnaryNode
}

type ImplyNode struct {
	BinaryNode *BinaryNode
}

type ForAllNode struct {
	QuantNode *QuantNode
	IsEffect  bool
}

type ExistsNode struct {
	QuantNode *QuantNode
}

type WhenNode struct {
	Condition Formula
	UnaryNode *UnaryNode
}

type AssignNode struct {
	Node         *Node
	Operation    *Name
	AssignedTo   *FunctionInit
	IsNumber     bool
	Number       string
	FunctionInit *FunctionInit
	IsInit       bool
}

func (lit *LiteralNode) ToString(prefix string) string {
	var s string
	if lit.Negative {
		s += fmt.Sprintf("%s(not ", prefix)
		prefix = ""
	}
	s += fmt.Sprintf("%s(", prefix)
	s += lit.Predicate.Name
	for _, t := range lit.Terms {
		s += fmt.Sprintf(" %s", t.Name.Name)
	}
	s += ")"
	if lit.Negative {
		s += ")"
	}
	return s
}

func (lit *LiteralNode) ToJSON(prefix string) string {
	var s string
	if lit.Negative {
		s += "\"not\":{"
	}
	s += "\"" + lit.Predicate.Name + "\"" + ":{"
	for i, t := range lit.Terms {
		if i == len(lit.Terms) - 1 {
			s += fmt.Sprintf("\"%s\"", t.Name.Name)
		} else {
			s += fmt.Sprintf("\"%s\",", t.Name.Name)
		}
	}
	s += "},"
	if lit.Negative {
		s += "},"
	}
	return s
}

func (n *AndNode) ToString(prefix string) string {
	var s string
	s += fmt.Sprintf("%s(and", prefix)
	for _, f := range n.MultiNode.Formula {
		s += "\n"
		s += f.ToString(prefix + Indent(1))
	}
	s += ")"
	return s
}

func (n *AndNode) ToJSON(prefix string) string {
	var s string
	s += "\"and\":{"
	for i, f := range n.MultiNode.Formula {
		if i == len(n.MultiNode.Formula) - 1 {
			s += f.ToJSON("")
		} else {
			s += f.ToJSON("") + ","
		}
		
	}
	s += "}"
	return s
}

func (n *OrNode) ToString(prefix string) string {
	s := fmt.Sprintf("%s(or", prefix)
	for _, f := range n.MultiNode.Formula {
		s += "\n"
		s += f.ToString(prefix + Indent(1))
	}
	s += ")"
	return s
}

func (n *OrNode) ToJSON(prefix string) string {
	s := "\"or\":{"
	for _, f := range n.MultiNode.Formula {
		s += "\"" + f.ToString("") + "\","
	}
	s += "}"
	return s
}

func (n *NotNode) ToString(prefix string) string {
	s := fmt.Sprintf("%s(not", prefix)
	s += n.UnaryNode.Formula.ToString(prefix)
	s += ")"
	return s
}

func (n *NotNode) ToJSON(prefix string) string {
	s := "\"not\":{"
	s += n.UnaryNode.Formula.ToJSON("")
	s += "}"
	return s
}

func (n *ImplyNode) ToString(prefix string) string {
	s := fmt.Sprintf("%s(imply\n", prefix)
	s += n.BinaryNode.Left.ToString(prefix + Indent(1))
	s += "\n"
	s += n.BinaryNode.Right.ToString(prefix + Indent(1))
	s += ")"
	return s
}

func (n *ForAllNode) ToString(prefix string) string {
	s := fmt.Sprintf("%s(forall (", prefix)
	s += toStringTypedNames("", n.QuantNode.Variables)
	s += ")\n"
	s += n.QuantNode.UnaryNode.Formula.ToString(prefix + Indent(1))
	s += ")"
	return s
}

func (n *ExistsNode) ToString(prefix string) string {
	s := fmt.Sprintf("%s(exists (", prefix)
	s += toStringTypedNames("", n.QuantNode.Variables)
	s += ")\n"
	s += n.QuantNode.UnaryNode.Formula.ToString(prefix + Indent(1))
	s += ")"
	return s
}

func (n *ImplyNode) ToJSON(prefix string) string {
	s := "\"imply\":{"
	s += n.BinaryNode.Left.ToJSON("")
	s += "\n"
	s += n.BinaryNode.Right.ToJSON("")
	s += "},"
	return s
}

func (n *ForAllNode) ToJSON(prefix string) string {
	s := "\"forall\":{"
	s += "\"quant\":{" + toJSONTypedNames("", n.QuantNode.Variables) + "}\","
	s += "\"effect\":{"
	s += n.QuantNode.UnaryNode.Formula.ToJSON("")
	s += "}}"
	return s
}

func (n *ExistsNode) ToJSON(prefix string) string {
	s := "\"exists\":{"
	s += toJSONTypedNames("", n.QuantNode.Variables)
	s += "}"
	s += n.QuantNode.UnaryNode.Formula.ToJSON("")
	s += "},"
	return s
}

func (n *WhenNode) ToString(prefix string) string {
	s := fmt.Sprintf("%s(when\n", prefix)
	s += n.Condition.ToString(prefix + Indent(1))
	s += "\n"
	s += n.UnaryNode.Formula.ToString(prefix + Indent(1))
	s += ")"
	return s
}

func (n *WhenNode) ToJSON(prefix string) string {
	s := "\"when\":{"
	s += n.Condition.ToJSON("")
	s += n.UnaryNode.Formula.ToJSON("")
	s += "}"
	return s
}

func (n *AssignNode) ToString(prefix string) string {
	s := fmt.Sprintf("%s(%s ", prefix, n.Operation.Name)
	s += n.AssignedTo.ToString()
	if n.IsNumber {
		s += fmt.Sprintf(" %s", n.Number)
	} else {
		s += " "
		s += n.FunctionInit.ToString()
	}
	s += ")"
	return s
}

func (n *AssignNode) ToJSON(prefix string) string {
	s := "\"" + n.Operation.Name + "\":{"
	s += n.AssignedTo.ToJSON()
	if n.IsNumber {
		s += "\"" + n.Number + "\""
	} else {
		s += n.FunctionInit.ToJSON()
	}
	s += "}"
	return s
}

func (h *FunctionInit) ToString() string {
	var s string
	if len(h.Terms) == 0 {
		s += fmt.Sprintf("(%s)", h.Name.Name)
		return s
	}
	s += fmt.Sprintf("(%s", h.Name.Name)
	for _, t := range h.Terms {
		s += fmt.Sprintf(" %s", t.Name.Name)
	}
	s += ")"
	return s
}

func (h *FunctionInit) ToJSON() string {
	var s string
	if len(h.Terms) == 0 {
		s += fmt.Sprintf("\"%s\":{},", h.Name.Name)
		return s
	}
	s += fmt.Sprintf("\"%s\":{", h.Name.Name)
	for i := range h.Terms {
		if i == len(h.Terms) - 1 {
			s += fmt.Sprintf("%s", h.Terms[i].Name.Name)
		} else {
			s += fmt.Sprintf("%s,", h.Terms[i].Name.Name)
		}
	}
	s += "},"
	return s
}
