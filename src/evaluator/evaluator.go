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

func accessToElementOfList(expr *ast.IndexExpr, env *object.Environment, dataToSet object.Object) object.Object {
	ident := expr.Ident
	if ident == nil {
		return newError("識別子が必要です。")
	}

	o, ok := env.Get(ident.Name)
	if !ok {
		return newError("配列が宣言されていません。")
	}

	array, ok := o.(*object.Array)
	if !ok {
		return newError("配列ではありません。")
	}

	indexList := expr.IndexList
	elements := &(array.Elements)
	for n, ele := range indexList {
		index, ok := (Eval(ele, env)).(*object.Integer)
		if !ok {
			return newError("数値が必要です。")
		}
		if int64(len(*elements)) <= index.Value || index.Value < 0 {
			return newError("範囲外です。")
		}

		_, ok = (*elements)[index.Value].(*object.Array)
		if !ok {
			if n < len(indexList) - 1 {
				return newError("範囲外です。")
			}
			break
		}

		elements = &((*elements)[index.Value].(*object.Array).Elements)
	}

	index, ok := (Eval(indexList[len(indexList) - 1], env)).(*object.Integer)
	if !ok {
		return newError("数値が必要です。")
	}

	if dataToSet != nil {
		(*elements)[index.Value] = dataToSet
		return NULL
	}

	return (*elements)[index.Value]
}

func evalInfixExpr(node ast.Node, env *object.Environment) object.Object {
	expr := node.(*ast.InfixExpr)

	switch expr.Operator {
	case ast.PA, ast.MA, ast.SA, ast.AA:
		return evalExtendAssign(expr, env, expr.Operator)
	case ast.ASSIGN:
		val := Eval(expr.Right, env)
		left := expr.Left

		switch left.(type) {
		case *ast.Ident:
			env.Set(left.(*ast.Ident).Name, val)
			return NULL
		case *ast.IndexExpr:
			accessToElementOfList(left.(*ast.IndexExpr), env, val)
			return NULL
		}

		return newError("変数が必要です。")
	}

	left := Eval(expr.Left, env)
	right := Eval(expr.Right, env)
	switch right.(type) {
	case *object.Integer:
		return evalIntegerExpression(expr.Operator, left, right)
	case *object.Boolean:
		return evalBooleanExpression(expr.Operator, left, right)
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
	case object.BOOLEAN:
		val = evalBooleanExpression(opeType, lhs, rhs)
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
		return newError("異なる型での演算は出来ません。 左オペランド:%s 右オペランド:%s", left.Type(), right.Type())
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
	case ast.RANGE:
		array := &object.Array{Elements: []object.Object{}}
		for i := lval; i < rval; i++ {
			array.Elements = append(array.Elements, &object.Integer{Value: i})
		}
		return array
	default:
		return newError("対応していない演算子です")
	}
}

func evalBooleanExpression(opeKind ast.OperatorKind, left object.Object, right object.Object) object.Object {
	if left.Type() != object.BOOLEAN || right.Type() != object.BOOLEAN {
		return newError("異なる型での演算は出来ません。 左オペランド:%s 右オペランド:%s", left.Type(), right.Type())
	}

	lval := left.(*object.Boolean).Value
	rval := right.(*object.Boolean).Value

	switch opeKind {
	case ast.EQ:
		return &object.Boolean{Value: lval == rval}
	case ast.NOT_EQ:
		return &object.Boolean{Value: lval != rval}
	case ast.AND:
		return &object.Boolean{Value: lval && rval}
	case ast.OR:
		return &object.Boolean{Value: lval || rval}
	default:
		return newError("対応していない演算子です")
	}
}

func evalPrefixExpr(node ast.Node, env *object.Environment) object.Object {
	expr := node.(*ast.PrefixExpr)

	right := Eval(expr.Right, env)
	switch right.(type) {
	case *object.Integer:
		return evalIntegerExpression(expr.Operator, &object.Integer{Value: 0}, right)
	case *object.Boolean:
		return &object.Boolean{Value: !right.(*object.Boolean).Value}
	}

	return newError("対応していない型が検出されました。")
}

func evalIdent(node ast.Node, env *object.Environment) object.Object {
	object, ok := env.Get(node.(*ast.Ident).Name)
	if !ok {
		return newError("変数が宣言されていません")
	}
	return object
}

func evalReturnStmt(node ast.Node, env *object.Environment) object.Object {
	val := Eval(node.(*ast.ReturnStmt).Value, env)
	if isError(val) {
		return val
	}
	return &object.ReturnValue{Value: val}
}

func evalIfStatement(node ast.Node, env *object.Environment) object.Object {
	stmt := node.(*ast.IfStmt)
	condition := Eval(stmt.Condition, env)
	var res object.Object = NULL
	if isError(condition) {
		return condition
	}

	if isTruthly(condition) {
		res = Eval(stmt.Body, env)
	} else if stmt.Else != nil {
		res = Eval(stmt.Else, env)
	}

	if res.Type() == object.ERROR {
		return res
	} else {
		return NULL
	}
}

func evalForStatement(node ast.Node, env *object.Environment) object.Object {
	stmt := node.(*ast.ForStmt)
	var res object.Object = NULL
	for {
		condition := Eval(stmt.Condition, env)
		if isError(condition) {
			return condition
		}

		if isTruthly(condition) {
			res = Eval(stmt.Body, env)
		} else {
			break
		}
	}

	if res.Type() == object.ERROR {
		return res
	} else {
		return NULL
	}
}

func evalForEachStatement(node ast.Node, env *object.Environment) object.Object {
	stmt := node.(*ast.ForEachStmt)
	array := Eval(stmt.Array, env)
	if array.Type() != object.ARRAY {
		return newError("配列が必要です。")
	}

	var res object.Object = NULL
	forEachEnv := object.NewEnclosedEnvironment(env)
	for i, ele := range array.(*object.Array).Elements {
		forEachEnv.SetCurrentEnv("添字", &object.Integer{Value: int64(i)})
		forEachEnv.SetCurrentEnv("要素", ele)

		res = Eval(stmt.Body, forEachEnv)
	}

	if res.Type() == object.ERROR {
		return res
	} else {
		return NULL
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

func evalArrayExpr(node ast.Node, env *object.Environment) object.Object {
	arrayExpr := node.(*ast.ArrayExpr)

	elements := []object.Object{}
	for _, v := range arrayExpr.Elements {
		elements = append(elements, Eval(v, env))
	}

	return &object.Array{Elements: elements}
}

func Eval(node ast.Node, env *object.Environment) object.Object {
	switch node.(type) {
	case *ast.InfixExpr:
		return evalInfixExpr(node, env)
	case *ast.PrefixExpr:
		return evalPrefixExpr(node, env)
	case *ast.Ident:
		return evalIdent(node, env)
	case *ast.Integer:
		return &object.Integer{Value: node.(*ast.Integer).Value}
	case *ast.Boolean:
		return &object.Boolean{Value: node.(*ast.Boolean).Value}
	case *ast.String:
		return &object.String{Value: node.(*ast.String).Value}
	case *ast.CallExpr:
		return evalCallFunc(node, env)
	case *ast.ArrayExpr:
		return evalArrayExpr(node, env)
	case *ast.IndexExpr:
		return accessToElementOfList(node.(*ast.IndexExpr), env, nil)
	case *ast.ExprStmt:
		return Eval(node.(*ast.ExprStmt).Expr, env)
	case *ast.ReturnStmt:
		return evalReturnStmt(node, env)
	case *ast.IfStmt:
		return evalIfStatement(node, env)
	case *ast.ForStmt:
		return evalForStatement(node, env)
	case *ast.ForEachStmt:
		return evalForEachStatement(node, env)
	case *ast.BlockStmt:
		return evalBlock(node, env)
	case *ast.FuncStmt:
		env.Set(node.(*ast.FuncStmt).Name, genFuncObj(node, env))
		return NULL
	}

	return newError("エラー")
}
