package parser

import (
	"fmt"
	"testing"

	"jpl/ast"
	"jpl/lexer"
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
	case bool:
		return testBoolean(t, expr, bool(v))
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

func testBoolean(t *testing.T, boolExpr ast.Expr, value bool) bool {
	boolean, ok := boolExpr.(*ast.Boolean)
	if ! ok {
		t.Errorf("boolExpr is not *ast.Boolean. got=%T", boolExpr)
		return false
	}

	if boolean.Value != value {
		t.Errorf("boolean.Value is not %t. got=%t", value, boolean.Value)
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
		head := lexer.Tokenize(v.input)
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
		head := lexer.Tokenize(v.input)
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
		head := lexer.Tokenize(v.input)
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
		head := lexer.Tokenize(v.input)
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
		head := lexer.Tokenize(v.input)
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
	head:= lexer.Tokenize(input)
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
	head := lexer.Tokenize(input)
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
	head := lexer.Tokenize(input)
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
	head := lexer.Tokenize(input)
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
	head := lexer.Tokenize(input)
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
		head := lexer.Tokenize(v.input)
		program, _ := Parse(head)

		for _, n := range program.Nodes {
			node := n.(*ast.ExprStmt).Expr
			if !testInfixExpr(t, node, v.lhs, v.operator, v.rhs) {
				return
			}
		}
	}
}

func TestIndexExpr(t *testing.T) {
	tests := []struct {
		input string
		identifier string
		indexList []int64
	}{
		{"a[10]", "a", []int64{10}},
		{"a[109]", "a", []int64{109}},
		{"a[109][45][32]", "a", []int64{109, 45, 32}},
	}

	for _, v := range tests {
		head := lexer.Tokenize(v.input)
		program, _ := Parse(head)

		for _, n := range program.Nodes {
			node, ok := n.(*ast.ExprStmt).Expr.(*ast.IndexExpr)
			if !ok {
				t.Fatalf("node is not *ast.IndexExpr. got=%T", n.(*ast.ExprStmt).Expr)
			}

			if node.Ident == nil {
				t.Fatalf("node.Ident is nil")
			}

			if node.Ident.Name != v.identifier {
				t.Fatalf("ident.Name is not %s. got=%s", v.identifier, node.Ident.Name)
			}

			indexList := node.IndexList
			if len(indexList) != len(v.indexList) {
				t.Fatalf("indexList is not %d indexes. got=%d", len(v.indexList), len(indexList))
			}

			for i, ele := range indexList {
				index, ok := ele.(*ast.Integer)
				if !ok {
					t.Fatalf("index is not *ast.Integer. got=%T", ele)
				}
				if !testLiteral(t, index, v.indexList[i]) {
					return
				}
			}
		}
	}
}

func TestInfixInIndexExpr(t *testing.T) {
	type IndexList struct {
		left int64
		operator ast.OperatorKind
		right int64
	}

	tests := []struct {
		input string
		identifier string
		indexList []IndexList
	}{
		{
			"a[10+9]",
			"a",
			[]IndexList{
				IndexList{left: 10, operator: ast.ADD, right: 9},
			},
		},
		{
			"a[10+9][9*8]",
			"a",
			[]IndexList{
				IndexList{left: 10, operator: ast.ADD, right: 9},
				IndexList{left: 9, operator: ast.MUL, right: 8},
			},
		},
	}

	for _, v := range tests {
		head := lexer.Tokenize(v.input)
		program, _ := Parse(head)

		for _, n := range program.Nodes {
			node, ok := n.(*ast.ExprStmt).Expr.(*ast.IndexExpr)
			if !ok {
				t.Fatalf("node is not *ast.IndexExpr. got=%T", n.(*ast.ExprStmt).Expr)
			}

			if node.Ident == nil {
				t.Fatalf("node.Ident is nil")
			}

			if node.Ident.Name != v.identifier {
				t.Fatalf("ident.Name is not %s. got=%s", v.identifier, node.Ident.Name)
			}

			indexList := node.IndexList
			if len(indexList) != len(v.indexList) {
				t.Fatalf("indexList is not %d indexes. got=%d", len(v.indexList), len(indexList))
			}

			for i, ele := range indexList {
				index, ok := ele.(*ast.InfixExpr)
				if !ok {
					t.Fatalf("index is not *ast.InfixExpr. got=%T", ele)
				}
				if !testInfixExpr(t, index, v.indexList[i].left, v.indexList[i].operator, v.indexList[i].right) {
					return
				}
			}
		}
	}
}

func TestListElements(t *testing.T) {
	tests := []struct {
		input string
		elements []int64
	}{
		{"a = {1, 2, 3, 4}", []int64{1, 2, 3, 4}},
	}

	for _, v := range tests {
		head := lexer.Tokenize(v.input)
		program, _ := Parse(head)

		for _, n := range program.Nodes {
			node, ok := n.(*ast.ExprStmt).Expr.(*ast.InfixExpr)
			if !ok {
				t.Fatalf("node is not *ast.InfixExpr. got=%T", n.(*ast.ExprStmt).Expr)
			}

			elements, ok := node.Right.(*ast.ArrayExpr)
			if !ok {
				t.Fatalf("elements is not *ast.ArrayExpr. got=%T", node.Right)
			}

			if len(v.elements) != len(elements.Elements) {
				t.Fatalf("elements.Elements is not %d expressions. got=%d", len(v.elements), len(elements.Elements))
			}

			for i, ele := range elements.Elements {
				if !testInteger(t, ele, v.elements[i]) {
					return
				}
			}
		}
	}
}

func TestBoolean(t *testing.T) {
	tests := []struct {
		input string
		expect bool
	} {
		{"真", true},
		{"偽", false},
	}

	for _, v := range tests {
		head := lexer.Tokenize(v.input)
		program, _ := Parse(head)

		for _, n := range program.Nodes {
			node, ok := n.(*ast.ExprStmt).Expr.(*ast.Boolean)
			if !ok {
				t.Fatalf("node is not *ast.Boolean. got=%T", n.(*ast.ExprStmt).Expr)
			}

			if node.Value != v.expect {
				t.Fatalf("node.Value is not %t. got=%t", v.expect, node.Value)
			}
		}
	}
}

func TestLogicalOperators(t *testing.T) {
	tests := []struct {
		input string
		left interface{}
		operator ast.OperatorKind
		right interface{}
	} {
		{"真　かつ　真", true, ast.AND, true},
		{"真　または　偽", true, ast.OR, false},
	}

	for _, v := range tests {
		head := lexer.Tokenize(v.input)
		program, _ := Parse(head)

		for _, n := range program.Nodes {
			node, ok := n.(*ast.ExprStmt).Expr.(*ast.InfixExpr)
			if !ok {
				t.Fatalf("node is not *ast.InfixExpr. got=%T", n.(*ast.ExprStmt).Expr)
			}

			if !testInfixExpr(t, node, v.left, v.operator, v.right) {
				return
			}
		}
	}
}

func TestNotOperator(t *testing.T) {
	input := "!真"
	head := lexer.Tokenize(input)
	program, _ := Parse(head)
	node, ok := program.Nodes[0].(*ast.ExprStmt).Expr.(*ast.PrefixExpr)
	if !ok {
		t.Fatalf("node is not *ast.PrefixExpr. got=%T", program.Nodes[0].(*ast.ExprStmt).Expr)
	}

	if node.Operator != ast.NOT {
		t.Fatalf("node.Operator is not %d. got=%d", ast.NOT, node.Operator)
	}

	right, ok := node.Right.(*ast.Boolean)
	if !ok {
		t.Fatalf("right is not *ast.Boolean. got=%T", node.Right)
	}	

	if !right.Value {
		t.Fatalf("right.Value is not true. got=%t", right.Value)
	}
}