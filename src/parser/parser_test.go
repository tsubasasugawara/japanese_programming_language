package parser

import (
	"testing"

	"jpl/ast"
	"jpl/token"
)

func TestParseer(t *testing.T) {
	tests := []struct {
		input    string
		nodeKind ast.NodeKind
		lhs      int
		rhs      int
	}{
		{
			"５ ＋ 6",
			ast.ADD,
			5,
			6,
		},
	}

	for i, v := range tests {
		token := token.Tokenize(v.input)
		node := Parse(token)

		if node.NodeKind != v.nodeKind {
			t.Fatalf("test%d : got=%d expect=%d\n", i, v.nodeKind, node.NodeKind)
		}

		if node.Lhs.Num != v.lhs {
			t.Fatalf("test%d : got=%d expect=%d\n", i, v.lhs, node.Lhs.Num)
		}

		if node.Rhs.Num != v.rhs {
			t.Fatalf("test%d : got=%d expect=%d\n", i, v.rhs, node.Rhs.Num)
		}
	}
}
