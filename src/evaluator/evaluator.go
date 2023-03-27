package evaluator

import (
	"fmt"

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
	case ast.EQ:
		return &object.Boolean{Value: lval == rval}
	case ast.NOT_EQ:
		return &object.Boolean{Value: lval != rval}
	case ast.GT:
		return &object.Boolean{Value: lval < rval}
	case ast.GE:
		return &object.Boolean{Value: lval <= rval}
	default:
		return newError("unknown operator")
	}
}

func evalIfStatement(node *ast.Node, env *object.Environment) object.Object {
	condition := Eval(node.Condition, env)
	if isError(condition) {
		return condition
	}

	if isTruthly(condition) {
		return Eval(node.Then, env)
	} else if node.Else != nil {
		return Eval(node.Else, env)
	}
	return NULL
}

func evalForStatement(node *ast.Node, env *object.Environment) object.Object {
	var fnode object.Object

	for {
		condition := Eval(node.Condition, env)
		if isError(condition) {
			return condition
		}

		if isTruthly(condition) {
			fnode = Eval(node.Then, env)
		} else {
			return fnode
		}
	}
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
			return newError("identifier not found")
		}
		return object
	case ast.NUMBER:
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
	}

	lhs := Eval(node.Lhs, env)
	rhs := Eval(node.Rhs, env)

	return evalIntegerExpression(node.NodeKind, lhs, rhs)
}
