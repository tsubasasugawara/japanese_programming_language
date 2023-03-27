package parser

import (
	"testing"

	"jpl/ast"
	"jpl/token"
)

func TestOperator(t *testing.T) {
	tests := []struct {
		input    string
		nodeKind ast.NodeKind
		lhs      int
		rhs      int
	}{
		{ "6５+6", ast.ADD, 65, 6,},
		{ "５ ＋ 6", ast.ADD, 5, 6,},
		{ "5５ - 5", ast.SUB, 55, 5,},
		{"６ー6", ast.SUB, 6, 6,},
		{"５*5", ast.MUL, 5, 5},
		{"５＊5", ast.MUL, 5, 5},
		{"2５×5", ast.MUL, 25, 5},
		{"５/5", ast.DIV, 5, 5},
		{"５／5", ast.DIV, 5, 5},
		{"５÷45", ast.DIV, 5, 45},
	}

	for i, v := range tests {
		token := token.Tokenize(v.input)
		program, _ := Parse(token)

		for _, node := range program.Nodes {
			if node.NodeKind != v.nodeKind {
				t.Fatalf("test%d(kind) : got=%d expect=%d\n", i, node.NodeKind, v.nodeKind)
			}
			if node.Lhs == nil || node.Lhs.Num != v.lhs {
				t.Fatalf("test%d(lhs) : got=%d expect=%d\n", i, node.Lhs.Num, v.lhs)
			}
			if node.Rhs == nil || node.Rhs.Num != v.rhs {
				t.Fatalf("test%d(rhs) : got=%d expect=%d\n", i, node.Rhs.Num, v.rhs)
			}
		}
	}
}

func TestUnaryOperator(t *testing.T) {
	tests := []struct {
		input string
		nodeKind ast.NodeKind
		lhs int
		rhs int
	} {
		{"+5", ast.ADD, 0, 5},
		{"-5", ast.SUB, 0, 5},
	}

	for i, v := range tests {
		token := token.Tokenize(v.input)
		program, _ := Parse(token)

		for _, node := range program.Nodes {
			if node.NodeKind != v.nodeKind {
				t.Fatalf("test%d(kind) : got=%d expect=%d\n", i, node.NodeKind, v.nodeKind)
			}
			if node.Lhs == nil || node.Lhs.Num != v.lhs {
				t.Fatalf("test%d(lhs) : got=%d expect=%d\n", i, node.Lhs.Num, v.lhs)
			}
			if node.Rhs == nil || node.Rhs.Num != v.rhs {
				t.Fatalf("test%d(rhs) : got=%d expect=%d\n", i, node.Rhs.Num, v.rhs)
			}
		}
	}
}

func TestComparisonOperators(t *testing.T) {
	tests := []struct {
		input string
		nodeKind ast.NodeKind
		lhs int
		rhs int
	} {
		{"5 < 9", ast.GT, 5, 9},
		{"５＜９", ast.GT, 5, 9},
		{"5 <= 9", ast.GE, 5, 9},
		{"５＜＝９", ast.GE, 5, 9},
		{"9>5", ast.GT, 5, 9},
		{"９＞５", ast.GT, 5, 9},
		{"9>=5", ast.GE, 5, 9},
		{"9＞＝5", ast.GE, 5, 9},
		{"5==5", ast.EQ, 5, 5},
		{"５＝＝５", ast.EQ, 5, 5},
		{"5!=9", ast.NOT_EQ, 5, 9},
		{"５！＝９", ast.NOT_EQ, 5, 9},
	}

	for i, v := range tests {
		token := token.Tokenize(v.input)
		program, _ := Parse(token)

		for _, node := range program.Nodes {
			if node.NodeKind != v.nodeKind {
				t.Fatalf("test%d(kind) : got=%d expect=%d\n", i, node.NodeKind, v.nodeKind)
			}
			if node.Lhs == nil || node.Lhs.Num != v.lhs {
				t.Fatalf("test%d(lhs) : got=%d expect=%d\n", i, node.Lhs.Num, v.lhs)
			}
			if node.Rhs == nil || node.Rhs.Num != v.rhs {
				t.Fatalf("test%d(rhs) : got=%d expect=%d\n", i, node.Rhs.Num, v.rhs)
			}
		}
	}
}

func TestIdentifier(t *testing.T) {
	tests := []struct {
		input string
		nodeKind ast.NodeKind
		lhs string
		rhs int
	} {
		{"こ=5", ast.ASSIGN, "こ", 5},
		{"a＝10", ast.ASSIGN, "a", 10},
	}

	for i, v := range tests {
		token := token.Tokenize(v.input)
		program, _ := Parse(token)

		for _, node := range program.Nodes {
			if node.NodeKind != v.nodeKind {
				t.Fatalf("test%d(kind) : got=%d expect=%d\n", i, node.NodeKind, v.nodeKind)
			}
			if node.Lhs == nil || node.Lhs.Ident != v.lhs{
				t.Fatalf("test%d(lhs) : got=%s expect=%s\n", i, node.Lhs.Ident, v.lhs)
			}
			if node.Rhs == nil || node.Rhs.Num != v.rhs {
				t.Fatalf("test%d(rhs) : got=%d expect=%d\n", i, node.Rhs.Num, v.rhs)
			}
		}
	}
}

func TestReturnStatement(t *testing.T) {
	tests := []struct {
		input string
		nodeKind ast.NodeKind
		lhs ast.NodeKind
	} {
		{"5+5 戻す", ast.RETURN, ast.ADD},
	}

	for i, v := range tests {
		token := token.Tokenize(v.input)
		program, _ := Parse(token)

		for _, node := range program.Nodes {
			if node.NodeKind != v.nodeKind {
				t.Fatalf("test%d(kind) : got=%d expect=%d\n", i, node.NodeKind, v.nodeKind)
			}
			if node.Lhs == nil || node.Lhs.NodeKind != v.lhs{
				t.Fatalf("test%d(lhs) : got=%d expect=%d\n", i, node.Lhs.NodeKind, v.lhs)
			}
		}
	}
}

func TestIfStatement(t *testing.T) {
	input := "もし 5 == 5 ならば 10 戻す"
	token := token.Tokenize(input)
	program, _ := Parse(token)

	node := program.Nodes[0]
	if node.NodeKind != ast.IF{
		t.Fatalf("got=%d expect=%d\n", node.NodeKind, ast.IF)
	}
	if node.Condition.NodeKind != ast.EQ {
		t.Fatalf("got=%d expect=%d\n", node.Condition.NodeKind, ast.EQ)
	}
	if node.Then.NodeKind != ast.RETURN {
		t.Fatalf("got=%d expect=%d\n", node.Then.NodeKind, ast.RETURN)
	}
}

func TestIfElseStatement(t *testing.T) {
	input := "もし 5 != 5 ならば 10 戻す それ以外 15 戻す"
	token := token.Tokenize(input)
	program, _ := Parse(token)

	node := program.Nodes[0]
	if node.NodeKind != ast.IF{
		t.Fatalf("got=%d expect=%d\n", node.NodeKind, ast.IF)
	}
	if node.Condition.NodeKind != ast.NOT_EQ {
		t.Fatalf("got=%d expect=%d\n", node.Condition.NodeKind, ast.EQ)
	}
	if node.Then.NodeKind != ast.RETURN {
		t.Fatalf("got=%d expect=%d\n", node.Then.NodeKind, ast.RETURN)
	}
	if node.Else == nil && node.Else.NodeKind != ast.RETURN {
		t.Fatalf("got=%d expect=%d\n", node.Else.NodeKind, ast.RETURN)
	}
}