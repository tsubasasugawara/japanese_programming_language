package evaluator

import (
	"fmt"
	"math"

	"jpl/ast"
	"jpl/object"
)

var (
	NULL = &object.Null{}
)

func newError(format string, a ...interface{}) *object.Error {
	return &object.Error{Message: fmt.Sprintf(format, a...)}
}

func isError(obj object.Object) bool {
	if obj != nil {
		return obj.Type() == object.ERROR
	}
	return false
}

func isTruthly(obj object.Object) bool {
	switch obj.Type() {
	case object.BOOLEAN:
		return obj.(*object.Boolean).Value
	case object.INTEGER:
		return obj.(*object.Integer).Value != 0
	case object.NULL:
		return false
	default:
		return true
	}
}

func evalInfixExpr(node ast.Node, env *object.Environment) object.Object {
	expr := node.(*ast.InfixExpr)

	switch expr.Operator {
	case ast.PA, ast.MA, ast.SA, ast.AA:
		return evalExtendAssign(expr, env, expr.Operator)
	case ast.ASSIGN:
		val := Eval(expr.Right, env)
		left, ok := expr.Left.(*ast.Ident)
		if !ok {
			return newError("変数が必要です。")
		}
		env.Set(left.Name, val)
		return NULL
	}

	left := Eval(expr.Left, env)
	right := Eval(expr.Right, env)
	switch right.(type) {
	case *object.Integer:
		return evalIntegerExpression(expr.Operator, left, right)
	}

	return newError("対応していない型が検出されました")
}

func evalExtendAssign(node ast.Node, env *object.Environment, opeType ast.OperatorKind) object.Object {
	stmt := node.(*ast.InfixExpr)

	if stmt.Left == nil {
		return newError("変数が宣言されていません")
	}

	left, ok := stmt.Left.(*ast.Ident)
	if !ok {
		return newError("代入演算子の左辺には、変数が必要です。")
	}

	lhs , ok := env.Get(left.Name)
	if !ok {
		return newError("変数が宣言されていません")
	}

	rhs := Eval(stmt.Right, env)

	var val object.Object
	switch rhs.Type() {
	case object.INTEGER:
		val = evalIntegerExpression(opeType, lhs, rhs)
		if isError(val) {
			return val
		}
	default:
		return newError("対応していない型が検出されました。")
	}

	env.Set(left.Name, val)
	return NULL
}

func evalIntegerExpression(opeKind ast.OperatorKind, left object.Object, right object.Object) object.Object {
	if left.Type() != object.INTEGER || right.Type() != object.INTEGER {
		return newError("数値が必要です。")
	}

	lval := left.(*object.Integer).Value
	rval := right.(*object.Integer).Value

	switch opeKind {
	case ast.ADD, ast.PA:
		return &object.Integer{Value: lval + rval}
	case ast.SUB, ast.MA:
		return &object.Integer{Value: lval - rval}
	case ast.MUL, ast.AA:
		return &object.Integer{Value: lval * rval}
	case ast.DIV, ast.SA:
		return &object.Integer{Value: lval / rval}
	case ast.EXPONENT:
		return &object.Integer{Value: int64(math.Pow(float64(lval), float64(rval)))}
	case ast.MODULUS:
		return &object.Integer{Value: lval % rval}
	case ast.EQ:
		return &object.Boolean{Value: lval == rval}
	case ast.NOT_EQ:
		return &object.Boolean{Value: lval != rval}
	case ast.GT:
		return &object.Boolean{Value: lval < rval}
	case ast.GE:
		return &object.Boolean{Value: lval <= rval}
	default:
		return newError("対応していない演算子です")
	}
}

func evalPrefixExpr(node ast.Node, env *object.Environment) object.Object {
	expr := node.(*ast.PrefixExpr)

	switch expr.Right.(type) {
	case *ast.Integer:
		return evalIntegerExpression(expr.Operator, &object.Integer{Value: 0}, Eval(expr.Right, env))
	}

	return newError("対応していない型が検出されました。")
}

func evalIfStatement(node ast.Node, env *object.Environment) object.Object {
	stmt := node.(*ast.IfStmt)
	condition := Eval(stmt.Condition, env)
	if isError(condition) {
		return condition
	}

	if isTruthly(condition) {
		Eval(stmt.Body, env)
	} else if stmt.Else != nil {
		Eval(stmt.Else, env)
	}
	return NULL
}

func evalForStatement(node ast.Node, env *object.Environment) object.Object {
	stmt := node.(*ast.ForStmt)
	for {
		condition := Eval(stmt.Condition, env)
		if isError(condition) {
			return condition
		}

		if isTruthly(condition) {
			Eval(stmt.Body, env)
		} else {
			return NULL
		}
	}
}

func evalBlock(node ast.Node, env *object.Environment) object.Object {
	var res object.Object
	blockEnv := object.NewEnclosedEnvironment(env)

	block := node.(*ast.BlockStmt)
	for _, stmt:= range block.List {
		res = Eval(stmt, blockEnv)

		if res == nil {
			continue
		}

		if rt := res.Type(); rt == object.RETURN_VALUE || rt == object.ERROR {
			return res
		}
	}

	if res == nil {
		res = NULL
	}

	return res
}

func genFuncObj(node ast.Node, env *object.Environment) object.Object {
	funcObj := &object.Function{}

	stmt := node.(*ast.FuncStmt)
	for _, v := range stmt.Params {
		funcObj.Params = append(funcObj.Params, v)
	}

	funcObj.Body = stmt.Body
	return funcObj
}

func evalCallFunc(node ast.Node, env *object.Environment) object.Object {
	callExpr := node.(*ast.CallExpr)

	builtin, ok := builtins[callExpr.Name]
	if ok {
		params := []object.Object{}
		for _, v := range callExpr.Params {
			p := Eval(v, env)
			if isError(p) {
				return p
			}
			params = append(params, p)
		}
		return builtin.Fn(params...)
	}

	obj, ok := env.Get(callExpr.Name)
	if !ok || obj.Type() != object.FUNCTION {
		return newError("関数が宣言されていません。")
	}
	if len(obj.(*object.Function).Params) != len(callExpr.Params) {
		return newError("引数の個数が正しくありません。")
	}

	callEnv := object.NewEnclosedEnvironment(env)
	for i, v := range obj.(*object.Function).Params {
		p := Eval(callExpr.Params[i], env)
		if isError(p) {
			return p
		}
		callEnv.SetCurrentEnv(v.Name, p)
	}

	res := Eval(obj.(*object.Function).Body, callEnv)
	if res != nil {
		return res.(*object.ReturnValue).Value
	}
	return NULL
}

func Eval(node ast.Node, env *object.Environment) object.Object {
	switch node.(type) {
	case *ast.InfixExpr:
		return evalInfixExpr(node, env)
	case *ast.PrefixExpr:
		return evalPrefixExpr(node, env)
	case *ast.Ident:
		object, ok := env.Get(node.(*ast.Ident).Name)
		if !ok {
			return newError("変数が宣言されていません")
		}
		return object
	case *ast.Integer:
		return &object.Integer{Value: node.(*ast.Integer).Value}
	case *ast.CallExpr:
		return evalCallFunc(node, env)
	case *ast.ExprStmt:
		return Eval(node.(*ast.ExprStmt).Expr, env)
	case *ast.ReturnStmt:
		val := Eval(node.(*ast.ReturnStmt).Value, env)
		if isError(val) {
			return val
		}
		return &object.ReturnValue{Value: val}
	case *ast.IfStmt:
		return evalIfStatement(node, env)
	case *ast.ForStmt:
		return evalForStatement(node, env)
	case *ast.BlockStmt:
		return evalBlock(node, env)
	case *ast.FuncStmt:
		env.Set(node.(*ast.FuncStmt).Name, genFuncObj(node, env))
		return NULL
	}

	return newError("エラー")
}
