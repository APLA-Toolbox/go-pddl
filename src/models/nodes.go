package models

import "io"

var (
	AssignOps = map[string]bool{
		"=":        true,
		"assign":   true,
		"increase": true,
	}
)

type Formula interface {
	Print(io.Writer, string)
	Check(defs, []*error)
}

type Node struct {
	Location *Location
}

type UnaryNode struct {
	Node    *Node
	Formula *Formula
}

type BinaryNode struct {
	Node
	Left  *Formula
	Right *Formula
}

type MultiNode struct {
	Node
	Formula []*Formula
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
	Condition *Formula
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
