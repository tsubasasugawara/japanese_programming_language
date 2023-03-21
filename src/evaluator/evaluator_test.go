package evaluator

import (
	"testing"

	"jpl/token"
	"jpl/parser"
	"jpl/object"
)

func TestCalc(t *testing.T) {
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
		{"-9 * (-8)", 72},
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

func TestComparisonOperators(t *testing.T) {
	tests := []struct {
		input string
		expect bool
	} {
		{"5 < 7", true},
		{"5 <= 8", true},
		{"5 <= 5", true},
		{"7 > 5", true},
		{"7 >= 5", true},
		{"7 >= 7", true},
		{"7 == 7", true},
		{"8 != 9", true},
		{"7 < 5", false},
		{"5 < 5", false},
		{"7 <= 5", false},
		{"5 > 7", false},
		{"5 > 5", false},
		{"5 >= 7", false},
		{"8 == 9", false},
		{"9 != 9", false},
	}
	for i, v := range tests {
		head := token.Tokenize(v.input)
		program, errors := parser.Parse(head)
		if len(errors) > 0 {
			t.Fatalf("Error.\n")
		}

		o := Eval(program.Nodes[0])
		if val := o.(*object.Boolean).Value; val != v.expect {
			t.Fatalf("test%d : got=%t expect=%t", i, val, v.expect)
		}
	}
}
