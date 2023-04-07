package evaluator

import (
	"fmt"
	"testing"

	"jpl/lexer"
	"jpl/parser"
	"jpl/object"
)

func TestCalc(t *testing.T) {
	tests := []struct {
		input string
		expectNum int64
	} {
		{"5 + 5", 10},
		{"５＋１９", 24},
		{"6 - 3", 3},
		{"7 * 8", 56},
		{"９÷３", 3},
		{"9 * 9 * 0", 0},
		{"(9 + 9) * 7", 126},
		{"-9 * (-8)", 72},
		{"9^2", 81},
		{"(-9)^3", -729},
		{"5%3", 2},
		{"７％５", 2},
	}

	for i, v := range tests {
		head := lexer.Tokenize(v.input)
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
		head := lexer.Tokenize(v.input)
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
		expect int64
	}{
		{"a = 5 a", 5},
		{"test=10 test", 10},
		{"test1=10 test1", 10},
		{"こんにちは＝１００ こんにちは", 100},
		{"世界 ＝ ２３８ 世界", 238},
		{"ワールド ＝ ２３５ ワールド", 235},
	}

	for i, v := range tests {
		head := lexer.Tokenize(v.input)
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
		expect string
	} {
		{"5+5 戻す", "10"},
	}

	for i, v := range tests {
		head := lexer.Tokenize(v.input)
		program, errors := parser.Parse(head)
		if len(errors) > 0 {
			t.Fatalf("Error.\n")
		}

		env := object.NewEnvironment()
		e := Eval(program.Nodes[0], env)
		if val := e.Inspect(); val != v.expect {
			t.Fatalf("test%d : got=%s expect=%s\n", i, val, v.expect)
		}
	}
}

func TestIfStatements(t *testing.T) {
	input := `
		a = 1
		もし a==1 ならば {
			a = a + 10
		} それ以外 {
			a = a - 10
		}
		a 戻す
		`
	head := lexer.Tokenize(input)
	program, errors := parser.Parse(head)
	if len(errors) > 0 {
		t.Fatalf("Error\n")
	}

	env := object.NewEnvironment()
	Eval(program.Nodes[0], env)
	Eval(program.Nodes[1], env)
	v := Eval(program.Nodes[2], env)
	if val := v.Inspect(); val != "11" {
		t.Fatalf("got=%s expect=%s\n", val, "11")
	}
}

func TestForStatement(t *testing.T) {
	input := `
	a = 1
	a < 5 ならば 繰り返す {
		a = a + 1
	}
	a 戻す
	`
	head := lexer.Tokenize(input)
	program, errors := parser.Parse(head)
	if len(errors) > 0 {
		t.Fatalf("Error,\n")
	}

	env := object.NewEnvironment()
	Eval(program.Nodes[0], env)
	Eval(program.Nodes[1], env)
	v := Eval(program.Nodes[2], env)
	if val := v.Inspect(); val != "5" {
		t.Fatalf("got=%s expect=%s\n", val, "5")
	}
}

func TestBlockStatement(t *testing.T) {
	input := `
	a = 0
	{
		a = 5
		b = 9
		b 戻す
	}
	a 戻す
	`
	head := lexer.Tokenize(input)
	program, errors := parser.Parse(head)
	if len(errors) > 0 {
		t.Fatalf("Error\n")
	}

	env := object.NewEnvironment()
	Eval(program.Nodes[0], env)
	v1 := Eval(program.Nodes[1], env)
	v2 := Eval(program.Nodes[2], env)

	if val := v1.Inspect(); val != "9" {
		t.Fatalf("got=%s expect=%s\n", val, "9")
	}
	if val := v2.Inspect(); val != "5" {
		t.Fatalf("got=%s expect=%s\n", val, "5")
	}
}

func TestFuncCall(t *testing.T) {
	input := `
	関数 abc(a, b, c) {
		a + b - c 戻す
	}
	関数 「あ、い」足す ｛
		あ＋い 戻す
	｝

	c = 90
	
	b = abc(10, 足す(2,3), c)
	b 戻す
	`
	head := lexer.Tokenize(input)
	program, errors := parser.Parse(head)
	if len(errors) > 0 {
		for _, err := range errors {
			fmt.Println(err.Message())
		}
	}

	env := object.NewEnvironment()
	Eval(program.Nodes[0], env)
	Eval(program.Nodes[1], env)
	Eval(program.Nodes[2], env)
	Eval(program.Nodes[3], env)
	v := Eval(program.Nodes[4], env)

	if val := v.Inspect(); val != "-75" {
		t.Fatalf("got=%s expect=%s\n", val, "-75")
	}
}

func TestRowComment(t *testing.T) {
	input := `
	// こんにちは = 100
	こんにちは = 800
	／／ こんにちは ＝ こんにちは ＋ 100
	こんにちは
	`
	head := lexer.Tokenize(input)
	program, errors := parser.Parse(head)
	if len(errors) > 0 {
		t.Fatalf("Error\n")
	}

	env := object.NewEnvironment()
	Eval(program.Nodes[0], env)
	v := Eval(program.Nodes[1], env)
	
	if val := v.Inspect(); val != "800" {
		t.Fatalf("got=%s expect%s\n", val, "800")
	}
}

func TestBlockComment(t *testing.T) {
	input := `
	こんにちは = 800
	/*
	こんにちは = 100
	こんにちは = こんにちは + 100
	＊／
	こんにちは
	`
	head := lexer.Tokenize(input)
	program, errors := parser.Parse(head)
	if len(errors) > 0 {
		t.Fatalf("Error\n")
	}

	env := object.NewEnvironment()
	Eval(program.Nodes[0], env)
	v := Eval(program.Nodes[1], env)
	
	if val := v.Inspect(); val != "800" {
		t.Fatalf("got=%s expect%s\n", val, "800")
	}
}

func TestExtendAssign(t *testing.T) {
	input := `
	a = 10
	a += 1
	a -= 2
	a *= 3
	a /= 27
	a
	`
	head := lexer.Tokenize(input)
	program, errors := parser.Parse(head)
	if len(errors) > 0 {
		t.Fatalf("Error\n")
	}

	env := object.NewEnvironment()
	Eval(program.Nodes[0], env)
	Eval(program.Nodes[1], env)
	Eval(program.Nodes[2], env)
	Eval(program.Nodes[3], env)
	Eval(program.Nodes[4], env)
	v := Eval(program.Nodes[5], env)

	if val := v.Inspect(); val != "1" {
		t.Fatalf("got=%s expect=%s\n", val, "1")
	}
}

func TestArray(t *testing.T) {
	input := `
	a = {1, 2, 3, 4, 5}
	a[0]
	a[1]
	a[2]
	a[3]
	a[4]
	`
	head := lexer.Tokenize(input)
	program, errors := parser.Parse(head)
	if len(errors) > 0 {
		t.Fatal("Error\n")
	}

	env := object.NewEnvironment()
	Eval(program.Nodes[0], env)
	v1 := Eval(program.Nodes[1], env)
	v2 := Eval(program.Nodes[2], env)
	v3 := Eval(program.Nodes[3], env)
	v4 := Eval(program.Nodes[4], env)
	v5 := Eval(program.Nodes[5], env)

	if val := v1.Inspect(); val != "1" {
		t.Fatalf("got=%s expect=1\n", val)
	}

	if val := v2.Inspect(); val != "2" {
		t.Fatalf("got=%s expect=2\n", val)
	}

	if val := v3.Inspect(); val != "3" {
		t.Fatalf("got=%s expect=3\n", val)
	}

	if val := v4.Inspect(); val != "4" {
		t.Fatalf("got=%s expect=4\n", val)
	}

	if val := v5.Inspect(); val != "5" {
		t.Fatalf("got=%s expect=5\n", val)
	}
}

func TestMultidimensionalArray(t *testing.T) {
	input := `
	a = {{1, 2}, {3, 4}}
	a[0][0]
	a[0][1]
	a[1][0]
	a[1][1]
	`
	head := lexer.Tokenize(input)
	program, errors := parser.Parse(head)
	if len(errors) > 0 {
		t.Fatal("Error\n")
	}

	env := object.NewEnvironment()
	Eval(program.Nodes[0], env)
	v1 := Eval(program.Nodes[1], env)
	v2 := Eval(program.Nodes[2], env)
	v3 := Eval(program.Nodes[3], env)
	v4 := Eval(program.Nodes[4], env)

	if val := v1.Inspect(); val != "1" {
		t.Fatalf("got=%s expect=1\n", val)
	}

	if val := v2.Inspect(); val != "2" {
		t.Fatalf("got=%s expect=2\n", val)
	}

	if val := v3.Inspect(); val != "3" {
		t.Fatalf("got=%s expect=3\n", val)
	}

	if val := v4.Inspect(); val != "4" {
		t.Fatalf("got=%s expect=4\n", val)
	}
}

func TestLogicalOperators(t *testing.T) {
	input := `
	a = 真　かつ　真
	b = 真　かつ　偽
	c =	偽　かつ　偽
	d = 真　または　真
	e = 真　または　偽
	f = 偽　または　偽
	g = !真
	h = !偽

	a
	b
	c
	d
	e
	f
	g
	h
	`
	expect := []bool{true, false, false, true, true, false, false, true}
	head := lexer.Tokenize(input)
	program, errors := parser.Parse(head)
	if len(errors) > 0 {
		t.Fatal("Error\n")
	}

	env := object.NewEnvironment()
	for i := 0; i < 8; i++ {
		Eval(program.Nodes[i], env)
	}

	for i := 8; i < len(program.Nodes); i++ {
		v := Eval(program.Nodes[i], env)
		val, ok := v.(*object.Boolean)
		if !ok {
			t.Fatalf("val is not *object.Boolean. got=%T", v)
		}
		if val.Value != expect[i - 8] {
			t.Fatalf("test%d : got=%t expect=%t", i - 8, val.Value, expect[i-8])
		}
	}
}

func TestGenList(t *testing.T) {
	input := `
	a = 1〜6
	a[0]
	a[1]
	a[2]
	a[3]
	a[4]
	`
	head := lexer.Tokenize(input)
	program, errors := parser.Parse(head)
	if len(errors) > 0 {
		t.Fatal("Error\n")
	}

	env := object.NewEnvironment()
	Eval(program.Nodes[0], env)
	v1 := Eval(program.Nodes[1], env)
	v2 := Eval(program.Nodes[2], env)
	v3 := Eval(program.Nodes[3], env)
	v4 := Eval(program.Nodes[4], env)
	v5 := Eval(program.Nodes[5], env)

	if val := v1.Inspect(); val != "1" {
		t.Fatalf("got=%s expect=1\n", val)
	}

	if val := v2.Inspect(); val != "2" {
		t.Fatalf("got=%s expect=2\n", val)
	}

	if val := v3.Inspect(); val != "3" {
		t.Fatalf("got=%s expect=3\n", val)
	}

	if val := v4.Inspect(); val != "4" {
		t.Fatalf("got=%s expect=4\n", val)
	}

	if val := v5.Inspect(); val != "5" {
		t.Fatalf("got=%s expect=5\n", val)
	}
}