package ast

import (
	"jpl/token"
)

type Node interface {
	Literal() string
}

type Expr interface {
	Node
	exprNode()
}

type Stmt interface {
	Node
	stmtNode()
}

type OperatorKind int

const (
	ADD = iota // 足し算
	SUB // 引き算
	MUL // 掛け算
	DIV // 割り算
	EXPONENT // べき乗
	MODULUS // 剰余算
	ASSIGN // 代入演算子
	GT // 超過
	GE // 以上
	EQ // 等号
	NOT_EQ // 不等号
	AND // 論理席
	OR // 論理和
	NOT // 否定
	PA // PLUS(ADD)+ASSIGN
	MA // MINUS(SUB)+ASSIGN
	SA // SLASH(DIV)+ASSIGN
	AA // ASTERISK(MUL)+ASSIGN
)

// expressions
type (
	Ident struct {
		Token *token.Token
		Name string
	}

	ArrayExpr struct {
		Token *token.Token
		Elements []Expr
	}

	IndexExpr struct {
		Token *token.Token
		Ident *Ident
		IndexList []Expr
	}

	Integer struct {
		Token *token.Token
		Value int64
	}

	Boolean struct {
		Token *token.Token
		Value bool
	}

	PrefixExpr struct {
		Token *token.Token
		Operator OperatorKind
		Right Expr
	}

	InfixExpr struct {
		Token *token.Token
		Left Expr
		Operator OperatorKind
		Right Expr
	}

	CallExpr struct {
		Token *token.Token
		Params []Expr
		Name string
	}
)

func (i *Ident) exprNode() {}
func (a *ArrayExpr) exprNode() {}
func (i *IndexExpr) exprNode() {}
func (i *Integer) exprNode() {}
func (b *Boolean) exprNode() {}
func (p *PrefixExpr) exprNode() {}
func (i *InfixExpr) exprNode() {}
func (c *CallExpr) exprNode() {}

func (i *Ident) Literal() string { return i.Token.Literal }
func (a *ArrayExpr) Literal() string { return a.Token.Literal }
func (i *IndexExpr) Literal() string { return i.Token.Literal }
func (i *Integer) Literal() string { return i.Token.Literal }
func (b *Boolean) Literal() string { return b.Token.Literal }
func (p *PrefixExpr) Literal() string { return p.Token.Literal }
func (i *InfixExpr) Literal() string { return i.Token.Literal }
func (c *CallExpr) Literal() string { return c.Token.Literal }

// statements
type (
	ExprStmt struct {
		Expr Expr
	}

	ReturnStmt struct {
		Token *token.Token
		Value Expr
	}

	BlockStmt struct {
		Token *token.Token
		List []Stmt
	}

	IfStmt struct {
		Token *token.Token
		Condition Expr
		Body *BlockStmt
		Else *BlockStmt
	}

	ForStmt struct {
		Token *token.Token
		Condition Expr
		Body *BlockStmt
	}

	FuncStmt struct {
		Token *token.Token
		Name string
		Params []*Ident
		Body *BlockStmt
	}
)

func (e *ExprStmt) stmtNode() {}
func (r *ReturnStmt) stmtNode() {}
func (b *BlockStmt) stmtNode() {}
func (i *IfStmt) stmtNode() {}
func (f *ForStmt) stmtNode() {}
func (f *FuncStmt) stmtNode() {}

func (e *ExprStmt) Literal() string { return e.Expr.Literal() }
func (r *ReturnStmt) Literal() string { return r.Token.Literal }
func (b *BlockStmt) Literal() string { return b.Token.Literal }
func (i *IfStmt) Literal() string { return i.Token.Literal }
func (f *ForStmt) Literal() string { return f.Token.Literal }
func (f *FuncStmt) Literal() string { return f.Token.Literal }
