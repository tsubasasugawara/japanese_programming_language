package ast

type NodeKind int

const (
	NUMBER NodeKind = iota

	IDENT // 識別子

	ADD // 足し算
	SUB // 引き算
	MUL // 掛け算
	DIV // 割り算
	ASSIGN // 代入演算子

	GT // 超過
	GE // 以上
	EQ // 等号
	NOT_EQ // 不等号

	RETURN // 戻す(return)
	IF // もし
	ELSE // それ以外
	THEN // ならば
	FOR // 繰り返す

	FUNC // 関数
	CALL // 関数呼び出し
	BLOCK
)

type RuntimeType int

const (
	INTEGER RuntimeType = iota
)

type Node struct {
	NodeKind NodeKind
	Next     *Node

	Lhs *Node
	Rhs *Node

	Condition *Node
	Then *Node
	Else *Node

	Stmts []*Node
	
	Params []*Node
	Body *Node

	Type RuntimeType

	Num int // NUMBERの時に値を格納する
	Ident string // 識別子を格納する
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

func NewIdentNode(ident string) *Node {
	n := NewNode(IDENT)
	n.Ident = ident
	return n
}
