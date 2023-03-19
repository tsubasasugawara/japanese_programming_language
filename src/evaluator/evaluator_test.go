package evaluator

import (
	"testing"

	"jpl/token"
	"jpl/parser"
	"jpl/object"
)

func TestEvaluator(t *testing.T) {
	tests := []struct {
		input string
		expectNum int
	} {
		{"5 + 5", 10},
		{"５＋１９", 24},
		{"6 - 3", 3},
		{"7 * 8", 56},
		{"９÷３", 3},
		{"9 * 9 * 0", 0},
		{"(9 + 9) * 7", 126},
	}

	for i, v := range tests {
		head := token.Tokenize(v.input)
		program, errors := parser.Parse(head)
		if len(errors) > 0 {
			t.Fatalf("Error.\n")
		}

		o := Eval(program.Nodes[0])
		if val := o.(*object.Integer).Value; val != v.expectNum {
			t.Fatalf("test%d : got=%d expect=%d", i, val, v.expectNum)
		}
	}
}

