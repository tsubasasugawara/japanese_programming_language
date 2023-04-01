package parser

import (
	"fmt"
	"testing"

	"jpl/ast"
	"jpl/token"
)

func testInfixExpr(
	t *testing.T,
	node ast.Expr,
	left interface{},
	operator ast.OperatorKind,
	right interface{},
) bool {
	ie, ok := node.(*ast.InfixExpr)
	if !ok {
		t.Errorf("This is not *ast.InfixExpr. got=%T", node)
	}

	if !testLiteral(t, ie.Left, left) {
		return false
	}

	if ie.Operator != operator {
		t.Errorf("node.Operator is not %d. got=%d", operator, ie.Operator)
	}

	if !testLiteral(t, ie.Right, right) {
		return false
	}

	return true
}

func testLiteral(t *testing.T, expr ast.Expr, expected interface{}) bool {
	switch v := expected.(type) {
	case int64:
		return testInteger(t, expr, int64(v))
	case string:
		return testIdentifier(t, expr, string(v))
	}
	t.Errorf("type of expr not handled. got=%T", expr)
	return false
}

func testInteger(t *testing.T, intExpr ast.Expr, value int64) bool {
	integ, ok := intExpr.(*ast.Integer)
	if !ok {
		t.Errorf("This is not *ast.Integer. got=%T", intExpr)
		return false
	}

	if integ.Value != value {
		t.Errorf("integ.Value not %d. got=%d", value, integ.Value)
		return false
	}

	return true
}

func testIdentifier(t *testing.T, identExpr ast.Expr, value string) bool {
	ident, ok := identExpr.(*ast.Ident)
	if !ok {
		t.Errorf("This is not *ast.Ident. got=%T", identExpr)
		return false
	}

	if ident.Name != value {
		t.Errorf("ident.Name not %s. got=%s", value, ident.Name)
		return false
	}

	return true
}

func TestOperator(t *testing.T) {
	tests := []struct {
		input    string
		lhs      int64
		operator ast.OperatorKind
		rhs      int64
	}{
		{ "6５+6", 65, ast.ADD, 6,},
		{ "５ ＋ 6", 5, ast.ADD, 6,},
		{ "5５ - 5", 55, ast.SUB, 5,},
		{"６ー6", 6, ast.SUB, 6,},
		{"５*5", 5, ast.MUL, 5},
		{"５＊5", 5, ast.MUL, 5},
		{"2５×5", 25, ast.MUL, 5},
		{"５/5", 5, ast.DIV, 5},
		{"５／5", 5, ast.DIV, 5},
		{"５÷45", 5, ast.DIV, 45},
		{"5^10", 5, ast.EXPONENT, 10},
		{"５＾１０", 5, ast.EXPONENT, 10},
		{"5%3", 5, ast.MODULUS, 3},
		{"５％３", 5, ast.MODULUS, 3},
	}

	for _, v := range tests {
		head := token.Tokenize(v.input)
		program, _ := Parse(head)

		for _, n := range program.Nodes {
			node := n.(*ast.ExprStmt).Expr
			if !testInfixExpr(t, node, v.lhs, v.operator, v.rhs) {
				return
			}
		}
	}
}

func TestUnaryOperator(t *testing.T) {
	tests := []struct {
		input string
		rhs int64
	} {
		{"+5", 5},
		{"-5", 5},
	}

	for i, v := range tests {
		head := token.Tokenize(v.input)
		program, _ := Parse(head)

		for _, n := range program.Nodes {
			node := n.(*ast.ExprStmt).Expr
			if node.(*ast.PrefixExpr).Right == nil || node.(*ast.PrefixExpr).Right.(*ast.Integer).Value != v.rhs {
				t.Fatalf("test%d(lhs) : got=%d expect=%d\n", i, node.(*ast.PrefixExpr).Right.(*ast.Integer).Value, v.rhs)
			}
		}
	}
}

func TestComparisonOperators(t *testing.T) {
	tests := []struct {
		input string
		lhs int64
		operator ast.OperatorKind
		rhs int64
	} {
		{"5 < 9", 5, ast.GT, 9},
		{"５＜９", 5, ast.GT, 9},
		{"5 <= 9", 5, ast.GE, 9},
		{"５＜＝９", 5, ast.GE, 9},
		{"9>5", 5, ast.GT, 9},
		{"９＞５", 5, ast.GT, 9},
		{"9>=5", 5, ast.GE, 9},
		{"9＞＝5", 5, ast.GE, 9},
		{"5==5", 5, ast.EQ, 5},
		{"５＝＝５", 5, ast.EQ, 5},
		{"5!=9", 5, ast.NOT_EQ, 9},
		{"５！＝９", 5, ast.NOT_EQ, 9},
	}

	for _, v := range tests {
		head := token.Tokenize(v.input)
		program, _ := Parse(head)

		for _, n := range program.Nodes {
			node := n.(*ast.ExprStmt).Expr
			if !testInfixExpr(t, node, v.lhs, v.operator, v.rhs) {
				return
			}
		}
	}
}

func TestIdentifier(t *testing.T) {
	tests := []struct {
		input string
		lhs string
		operator ast.OperatorKind
		rhs int64
	} {
		{"こ=5", "こ", ast.ASSIGN, 5},
		{"a＝10", "a", ast.ASSIGN, 10},
	}

	for _, v := range tests {
		head := token.Tokenize(v.input)
		program, _ := Parse(head)

		for _, n := range program.Nodes {
			node := n.(*ast.ExprStmt).Expr
			if !testInfixExpr(t, node, v.lhs, v.operator, v.rhs) {
				return
			}
		}
	}
}

func TestReturnStatement(t *testing.T) {
	tests := []struct {
		input string
		lhs int64
		operator ast.OperatorKind
		rhs int64
	} {
		{"5+5 戻す", 5, ast.ADD, 5},
	}

	for _, v := range tests {
		head := token.Tokenize(v.input)
		program, _ := Parse(head)

		for _, n := range program.Nodes {
			node, ok := n.(*ast.ReturnStmt)
			if !ok {
				t.Fatalf("This is not *ast.ReturnStmt. got=%T", n)
				return
			}

			if !testInfixExpr(t, node.Value, v.lhs, v.operator, v.rhs) {
				return
			}
		}
	}
}

func TestIfStatement(t *testing.T) {
	input := "もし 5 == 5 ならば {10+10}"
	head:= token.Tokenize(input)
	program, errors := Parse(head)
	if len(errors) > 0 {
		for _, err := range errors {
			fmt.Println(err.Message())
		}
		t.Fatalf("error")
	}

	node, ok := program.Nodes[0].(*ast.IfStmt)
	if !ok {
		t.Fatalf("This is not *ast.IfStmt. got=%T", program.Nodes[0])
	}

	if !testInfixExpr(t, node.Condition, int64(5), ast.EQ, int64(5)) {
		return
	}

	if node.Body == nil {
		t.Errorf("node.Body is nil.")
	}

	if len(node.Body.List) != 1 {
		t.Errorf("Body is not 1 statement. got=%d\n", len(node.Body.List))
	}

	con, ok := node.Body.List[0].(*ast.ExprStmt)
	if !ok {
		t.Fatalf("This is not *ast.ExprStmt. got=%T", node.Body.List[0])
	}

	if !testInfixExpr(t, con.Expr, int64(10), ast.ADD, int64(10)) {
		return
	}

	if node.Else != nil {
		t.Errorf("node.Else is not nil. got=%+v", node.Else)
	}
}

func TestIfElseStatement(t *testing.T) {
	input := "もし 5 == 5 ならば {10+10} それ以外 {15-5}"
	head := token.Tokenize(input)
	program, _ := Parse(head)

	node, ok := program.Nodes[0].(*ast.IfStmt)
	if !ok {
		t.Fatalf("This is not *ast.IfStmt. got=%T", program.Nodes[0])
	}

	if !testInfixExpr(t, node.Condition, int64(5), ast.EQ, int64(5)) {
		return
	}

	if node.Body == nil {
		t.Errorf("node.Body is nil.")
	}

	if len(node.Body.List) != 1 {
		t.Errorf("Body is not 1 statement. got=%d\n", len(node.Body.List))
	}

	con, ok := node.Body.List[0].(*ast.ExprStmt)
	if !ok {
		t.Fatalf("This is not *ast.ExprStmt. got=%T", node.Body.List[0])
	}

	if !testInfixExpr(t, con.Expr, int64(10), ast.ADD, int64(10)) {
		return
	}

	if node.Else == nil {
		t.Errorf("node.Else is not nil. got=%+v", node.Else)
	}

	if len(node.Else.List) != 1 {
		t.Errorf("Else is not 1 statement. got=%d\n", len(node.Else.List))
	}

	alter, ok := node.Else.List[0].(*ast.ExprStmt)
	if !ok {
		t.Fatalf("This is not *ast.ExprStmt. got=%T", node.Else.List[0])
	}

	if !testInfixExpr(t, alter.Expr, int64(15), ast.SUB, int64(5)) {
		return
	}
}

func TestForStatement(t *testing.T) {
	input := `
	a = 1
	a < 5 ならば 繰り返す {a = 1}
	`
	head := token.Tokenize(input)
	program, _ := Parse(head)

	node, ok := program.Nodes[1].(*ast.ForStmt)
	if !ok {
		t.Fatalf("This is not *ast.ForStmt. got=%T", program.Nodes[1])
	}

	if !testInfixExpr(t, node.Condition, "a", ast.GT, int64(5)) {
		return
	}

	if node.Body == nil {
		t.Fatalf("node.Body is nil.")
	}

	if len(node.Body.List) != 1 {
		t.Errorf("node.Body.List is not 1 statement. got=%d", len(node.Body.List))
	}

	body, ok := node.Body.List[0].(*ast.ExprStmt)
	if !ok {
		t.Fatalf("This is not *ast.ExprStmt. got=%T", node.Body.List[0])
	}

	if !testInfixExpr(t, body.Expr, "a", ast.ASSIGN, int64(1)) {
		return
	}
}

func TestFuncDeclation(t *testing.T) {
	input := `
	関数 足し算(a, b) {
		a + b
	}`
	head := token.Tokenize(input)
	program, _ := Parse(head)

	node, ok := program.Nodes[0].(*ast.FuncStmt)
	if !ok {
		t.Fatalf("This is not *ast.FuncStmt. got=%T", program.Nodes[0])
	}

	if node.Name != "足し算" {
		t.Fatalf("This name is not 足し算. got=%s", node.Name)
	}

	if len(node.Params) != 2 {
		t.Fatalf("node.Params is not 2 identifiers. got=%d", len(node.Params))
	}

	if node.Params[0].Name != "a" {
		t.Fatalf("first arg : got=%s expect=%s\n", node.Params[0].Name, "a")
	}

	if node.Params[1].Name != "b" {
		t.Fatalf("second arg : got=%s expect=%s\n", node.Params[1].Name, "b")
	}

	if node.Body == nil {
		t.Fatalf("node.Body is nil.")
	}

	if len(node.Body.List) != 1 {
		t.Errorf("node.Body.List is not 1 statement. got=%d", len(node.Body.List))
	}

	body, ok := node.Body.List[0].(*ast.ExprStmt)
	if !ok {
		t.Fatalf("This is not *ast.ExprStmt. got=%T", node.Body.List[0])
	}

	if !testInfixExpr(t, body.Expr, "a", ast.ADD, "b") {
		return
	}
}

func TestFuncCall(t *testing.T) {
	input := `
	こんにちは(世界)
	`
	head := token.Tokenize(input)
	program, _ := Parse(head)

	node, ok := program.Nodes[0].(*ast.ExprStmt).Expr.(*ast.CallExpr)
	if !ok {
		t.Fatalf("This is not *ast.CallExpr. got=%T", program.Nodes[0])
	}

	if node.Name != "こんにちは" {
		t.Fatalf("node.Name is not こんにちは. got=%s", node.Name)
	}

	if len(node.Params) != 1 {
		t.Fatalf("node.Params is not 1 expression. got=%d", len(node.Params))
	}

	param, ok := node.Params[0].(*ast.Ident)
	if !ok {
		t.Fatalf("This is not *ast.Ident. got=%T", node.Params[0])
	}

	if param.Name != "世界" {
		t.Fatalf("param.Name is not 世界. got=%s", param.Name)
	}
}

func TestExtendAssign(t *testing.T) {
	tests := []struct {
		input string
		lhs string
		operator ast.OperatorKind
		rhs int64
	} {
	{"a = 10", "a", ast.ASSIGN, 10},
	{"a += 1", "a", ast.PA, 1},
	{"a -= 2", "a", ast.MA, 2},
	{"a *= 3", "a", ast.AA, 3},
	{"a /= 4", "a", ast.SA, 4},
	}
	
	for _, v := range tests {
		head := token.Tokenize(v.input)
		program, _ := Parse(head)

		for _, n := range program.Nodes {
			node := n.(*ast.ExprStmt).Expr
			if !testInfixExpr(t, node, v.lhs, v.operator, v.rhs) {
				return
			}
		}
	}
}