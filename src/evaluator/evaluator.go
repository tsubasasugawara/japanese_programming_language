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

func evalIntegerExpression(nodeKind ast.NodeKind, left object.Object, right object.Object) object.Object {
	if left.Type() != object.INTEGER || right.Type() != object.INTEGER {
		return newError("数値が必要です。")
	}

	lval := left.(*object.Integer).Value
	rval := right.(*object.Integer).Value

	switch nodeKind {
	case ast.ADD:
		return &object.Integer{Value: lval + rval}
	case ast.SUB:
		return &object.Integer{Value: lval - rval}
	case ast.MUL:
		return &object.Integer{Value: lval * rval}
	case ast.DIV:
		return &object.Integer{Value: lval / rval}
	case ast.EXPONENT:
		return &object.Integer{Value: int(math.Pow(float64(lval), float64(rval)))}
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

func evalIfStatement(node *ast.Node, env *object.Environment) object.Object {
	condition := Eval(node.Condition, env)
	if isError(condition) {
		return condition
	}

	if isTruthly(condition) {
		Eval(node.Then, env)
	} else if node.Else != nil {
		Eval(node.Else, env)
	}
	return NULL
}

func evalForStatement(node *ast.Node, env *object.Environment) object.Object {
	for {
		condition := Eval(node.Condition, env)
		if isError(condition) {
			return condition
		}

		if isTruthly(condition) {
			Eval(node.Then, env)
		} else {
			return NULL
		}
	}
}

func evalBlock(node *ast.Node, env *object.Environment) object.Object {
	var res object.Object
	blockEnv := object.NewEnclosedEnvironment(env)

	for _, stmt:= range node.Stmts {
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

func genFuncObj(node *ast.Node, env *object.Environment) object.Object {
	funcObj := &object.Function{}

	for _, v := range node.Params {
		funcObj.Params = append(funcObj.Params, v)
	}

	funcObj.Body = node.Body
	return funcObj
}

func evalCallFunc(node *ast.Node, env *object.Environment) object.Object {
	builtin, ok := builtins[node.Ident]
	if ok {
		params := []object.Object{}
		for _, v := range node.Params {
			p := Eval(v, env)
			if isError(p) {
				return p
			}
			params = append(params, p)
		}
		return builtin.Fn(params...)
	}

	obj, ok := env.Get(node.Ident)
	if !ok || obj.Type() != object.FUNCTION {
		return newError("関数が宣言されていません。")
	}
	if len(obj.(*object.Function).Params) != len(node.Params) {
		return newError("引数の個数が正しくありません。")
	}

	callEnv := object.NewEnclosedEnvironment(env)
	for i, v := range obj.(*object.Function).Params {
		p := Eval(node.Params[i], env)
		if isError(p) {
			return p
		}
		callEnv.SetCurrentEnv(v.Ident, p)
	}

	res := Eval(obj.(*object.Function).Body, callEnv)
	if res != nil {
		return res.(*object.ReturnValue).Value
	}
	return NULL
}

func evalExtendAssign(node *ast.Node, env *object.Environment, opeType ast.NodeKind) object.Object {
	if node.Lhs == nil {
		return newError("変数が宣言されていません")
	}

	lhs , ok := env.Get(node.Lhs.Ident)
	if !ok {
		return newError("変数が宣言されていません")
	}
	rhs := Eval(node.Rhs, env)

	val := evalIntegerExpression(opeType, lhs, rhs)
	if isError(val) {
		return val
	}
	env.Set(node.Lhs.Ident, val)
	return NULL
}

func Eval(node *ast.Node, env *object.Environment) object.Object {
	switch node.NodeKind {
	case ast.ASSIGN:
		val := Eval(node.Rhs, env)
		env.Set(node.Lhs.Ident, val)
		return NULL
	case ast.IDENT:
		object, ok := env.Get(node.Ident)
		if !ok {
			return newError("変数が宣言されていません")
		}
		return object
	case ast.INTEGER:
		return &object.Integer{Value: node.Num}
	case ast.RETURN:
		val := Eval(node.Lhs, env)
		if isError(val) {
			return val
		}
		return &object.ReturnValue{Value: val}
	case ast.IF:
		return evalIfStatement(node, env)
	case ast.FOR:
		return evalForStatement(node, env)
	case ast.BLOCK:
		return evalBlock(node, env)
	case ast.FUNC:
		env.Set(node.Ident, genFuncObj(node, env))
		return NULL
	case ast.CALL:
		return evalCallFunc(node, env)
	case ast.PA:
		return evalExtendAssign(node, env, ast.ADD)
	case ast.MA:
		return evalExtendAssign(node, env, ast.SUB)
	case ast.AA:
		return evalExtendAssign(node, env, ast.MUL)
	case ast.SA:
		return evalExtendAssign(node, env, ast.DIV)
	}

	lhs := Eval(node.Lhs, env)
	rhs := Eval(node.Rhs, env)

	return evalIntegerExpression(node.NodeKind, lhs, rhs)
}
