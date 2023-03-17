package ast

type NodeKind int

const (
	NUMBER NodeKind = iota

	ADD // 足し算
	SUB // 引き算
	MUL // 掛け算
	DIV // 割り算
)

type Node struct {
	NodeKind NodeKind
	Next     *Node

	Lhs *Node
	Rhs *Node

	Num int // NUMBERの時に値を格納する
}

func NewNode(nodeKind NodeKind) *Node {
	n := &Node{NodeKind: nodeKind}
	return n
}

func NewNumberNode(num int) *Node {
	n := NewNode(NUMBER)
	n.Num = num
	return n
}

func NewNodeBinop(nodeKind NodeKind, lhs *Node, rhs *Node) *Node {
	n := NewNode(nodeKind)
	n.Lhs = lhs
	n.Rhs = rhs
	return n
}