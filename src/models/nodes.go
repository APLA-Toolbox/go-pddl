package models

import (
	"fmt"
	"io"
)

var (
	AssignOps = map[string]bool{
		"=":        true,
		"assign":   true,
		"increase": true,
	}
)

type Formula interface {
	Print(io.Writer, string)
	ToString(string) string
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

func (lit *LiteralNode) Print(w io.Writer, prefix string) {
	if lit.Negative {
		fmt.Fprintf(w, "%s(not ", prefix)
		prefix = ""
	}
	fmt.Fprintf(w, "%s(", prefix)
	fmt.Fprint(w, lit.Predicate.Name)
	for _, t := range lit.Terms {
		fmt.Fprintf(w, " %s", t.Name.Name)
	}
	fmt.Fprint(w, ")")
	if lit.Negative {
		fmt.Fprint(w, ")")
	}
}

func (n *AndNode) Print(w io.Writer, prefix string) {
	fmt.Fprintf(w, "%s(and", prefix)
	for _, f := range n.MultiNode.Formula {
		fmt.Fprint(w, "\n")
		f.Print(w, prefix+Indent(1))
	}
	fmt.Fprint(w, ")")
}

func (n *OrNode) Print(w io.Writer, prefix string) {
	fmt.Fprintf(w, "%s(or", prefix)
	for _, f := range n.MultiNode.Formula {
		fmt.Fprint(w, "\n")
		f.Print(w, prefix+Indent(1))
	}
	fmt.Fprint(w, ")")
}

func (n *NotNode) Print(w io.Writer, prefix string) {
	fmt.Fprintf(w, "%s(not\n", prefix)
	n.UnaryNode.Formula.Print(w, prefix+Indent(1))
	fmt.Fprint(w, ")")
}

func (n *ImplyNode) Print(w io.Writer, prefix string) {
	fmt.Fprintf(w, "%s(imply\n", prefix)
	n.BinaryNode.Left.Print(w, prefix+Indent(1))
	fmt.Fprint(w, "\n")
	n.BinaryNode.Right.Print(w, prefix+Indent(1))
	fmt.Fprint(w, ")")
}

func (n *ForAllNode) Print(w io.Writer, prefix string) {
	fmt.Fprintf(w, "%s(forall (", prefix)
	fmt.Fprintf(w, "%s", toStringTypedNames("", n.QuantNode.Variables))
	fmt.Fprint(w, ")\n")
	n.QuantNode.UnaryNode.Formula.Print(w, prefix+Indent(1))
	fmt.Fprint(w, ")")
}

func (n *ExistsNode) Print(w io.Writer, prefix string) {
	fmt.Fprintf(w, "%s(exists (", prefix)
	fmt.Fprintf(w, "%s", toStringTypedNames("", n.QuantNode.Variables))
	fmt.Fprint(w, ")\n")
	n.QuantNode.UnaryNode.Formula.Print(w, prefix+Indent(1))
	fmt.Fprint(w, ")")
}

func (n *WhenNode) Print(w io.Writer, prefix string) {
	fmt.Fprintf(w, "%s(when\n", prefix)
	n.Condition.Print(w, prefix+Indent(1))
	fmt.Fprint(w, "\n")
	n.UnaryNode.Formula.Print(w, prefix+Indent(1))
	fmt.Fprint(w, ")")
}

func (n *AssignNode) Print(w io.Writer, prefix string) {
	fmt.Fprintf(w, "%s(%s ", prefix, n.Operation.Name)
	n.AssignedTo.Print(w)
	if n.IsNumber {
		fmt.Fprintf(w, " %s", n.Number)
	} else {
		fmt.Fprint(w, " ")
		n.FunctionInit.Print(w)
	}
	fmt.Fprintf(w, ")")
}

func (h *FunctionInit) Print(w io.Writer) {
	if len(h.Terms) == 0 {
		fmt.Fprintf(w, "(%s)", h.Name.Name)
		return
	}
	fmt.Fprintf(w, "(%s", h.Name.Name)
	for _, t := range h.Terms {
		fmt.Fprintf(w, " %s", t.Name.Name)
	}
	fmt.Fprint(w, ")")
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

func (n *AndNode) ToString(prefix string) string {
	var s string
	s += fmt.Sprintf("%s(and", prefix)
	for _, f := range n.MultiNode.Formula {
		s += "\n"
		s += f.ToString(prefix+Indent(1))
	}
	s += ")"
	return s
}

func (n *OrNode) ToString(prefix string) string {
	s := fmt.Sprintf("%s(or", prefix)
	for _, f := range n.MultiNode.Formula {
		s += "\n"
		s += f.ToString(prefix+Indent(1))
	}
	s += ")"
	return s
}

func (n *NotNode) ToString(prefix string) string {
	s := fmt.Sprintf("%s(not\n", prefix)
	s += n.UnaryNode.Formula.ToString(prefix+Indent(1))
	s += ")"
	return s
}

func (n *ImplyNode) ToString(prefix string) string {
	s := fmt.Sprintf("%s(imply\n", prefix)
	s += n.BinaryNode.Left.ToString(prefix+Indent(1))
	s += "\n"
	s += n.BinaryNode.Right.ToString(prefix+Indent(1))
	s += ")"
	return s
}

func (n *ForAllNode) ToString(prefix string) string {
	s := fmt.Sprintf("%s(forall (", prefix)
	s += toStringTypedNames("", n.QuantNode.Variables)
	s += ")\n"
	s += n.QuantNode.UnaryNode.Formula.ToString(prefix+Indent(1))
	s += ")"
	return s
}

func (n *ExistsNode) ToString(prefix string) string {
	s := fmt.Sprintf("%s(exists (", prefix)
	s += toStringTypedNames("", n.QuantNode.Variables)
	s += ")\n"
	s += n.QuantNode.UnaryNode.Formula.ToString(prefix+Indent(1))
	s += ")"
	return s
}

func (n *WhenNode) ToString(prefix string) string {
	s := fmt.Sprintf("%s(when\n", prefix)
	s += n.Condition.ToString(prefix+Indent(1))
	s += "\n"
	s += n.UnaryNode.Formula.ToString(prefix+Indent(1))
	s += ")"
	return s
}

func (n *AssignNode) ToString(prefix string) string {
	s := fmt.Sprintf( "%s(%s ", prefix, n.Operation.Name)
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
