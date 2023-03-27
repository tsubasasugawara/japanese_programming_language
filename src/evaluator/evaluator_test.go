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

		env := object.NewEnvironment()
		o := Eval(program.Nodes[0], env)
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

		env := object.NewEnvironment()
		o := Eval(program.Nodes[0], env)
		if val := o.(*object.Boolean).Value; val != v.expect {
			t.Fatalf("test%d : got=%t expect=%t\n", i, val, v.expect)
		}
	}
}


func TestIdentifier(t *testing.T) {
	tests := []struct{
		input string
		expect int
	}{
		{"a = 5 a", 5},
		{"test=10 test", 10},
		{"test1=10 test1", 10},
		{"こんにちは＝１００ こんにちは", 100},
		{"世界 ＝ ２３８ 世界", 238},
		{"ワールド ＝ ２３５ ワールド", 235},
	}

	for i, v := range tests {
		head := token.Tokenize(v.input)
		program, errors := parser.Parse(head)
		if len(errors) > 0 {
			t.Fatalf("Error.\n")
		}

		env := object.NewEnvironment()
		Eval(program.Nodes[0], env)
		e2 := Eval(program.Nodes[1], env)
		if val := e2.(*object.Integer).Value; val != v.expect {
			t.Fatalf("test%d : got=%d expect=%d\n", i, val, v.expect)
		}
	}
}

func TestReturnStatement(t *testing.T) {
	tests := []struct {
		input string
		expect int
	} {
		{"5+5 戻す", 10},
	}

	for i, v := range tests {
		head := token.Tokenize(v.input)
		program, errors := parser.Parse(head)
		if len(errors) > 0 {
			t.Fatalf("Error.\n")
		}

		env := object.NewEnvironment()
		e := Eval(program.Nodes[0], env)
		if val := e.(*object.Integer).Value; val != v.expect {
			t.Fatalf("test%d : got=%d expect=%d\n", i, val, v.expect)
		}
	}
}

func TestIfStatement(t *testing.T) {
	tests := []struct {
		input string
		expect int
	} {
		{"もし 5==5 ならば 10 戻す", 10},
		{"もし 5!=5 10 戻す それ以外 15 戻す", 15},
	}
	
	for i, v := range tests {
		head := token.Tokenize(v.input)
		program, errors := parser.Parse(head)
		if len(errors) > 0 {
			t.Fatalf("Error.\n")
		}

		env := object.NewEnvironment()
		e := Eval(program.Nodes[0], env)
		if val := e.(*object.Integer).Value; val != v.expect {
			t.Fatalf("test%d : got=%d expect=%d\n", i, val, v.expect)
		}
	}
}

func TestIfStatements(t *testing.T) {
	input := `
		a = 1
		もし a==1 ならば
			a = a + 10
		それ以外
			a = a - 10
		a 戻す
		`
	head := token.Tokenize(input)
	program, errors := parser.Parse(head)
	if len(errors) > 0 {
		t.Fatalf("Error\n")
	}

	env := object.NewEnvironment()
	Eval(program.Nodes[0], env)
	Eval(program.Nodes[1], env)
	v := Eval(program.Nodes[2], env)
	if val := v.(*object.Integer).Value; val != 11 {
		t.Fatalf("got=%d expect=%d\n", val, 11)
	}
}

func TestForStatement(t *testing.T) {
	input := `
	a = 1
	a < 5 ならば 繰り返す a = a + 1
	a 戻す
	`
	head := token.Tokenize(input)
	program, errors := parser.Parse(head)
	if len(errors) > 0 {
		t.Fatalf("Error,\n")
	}

	env := object.NewEnvironment()
	Eval(program.Nodes[0], env)
	Eval(program.Nodes[1], env)
	v := Eval(program.Nodes[2], env)
	if val := v.(*object.Integer).Value; val != 5 {
		t.Fatalf("got=%d expect=%d\n", val, 5)
	}
}