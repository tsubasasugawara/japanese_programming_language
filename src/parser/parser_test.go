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
