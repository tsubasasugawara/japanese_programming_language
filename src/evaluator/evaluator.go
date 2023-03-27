package evaluator

import (
	"fmt"

	"jpl/ast"
	"jpl/object"
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

func Eval(node *ast.Node, env *object.Environment) object.Object {
	switch node.NodeKind {
	case ast.ASSIGN:
		val := Eval(node.Rhs, env)
		env.Set(node.Lhs.Ident, val)
		return &object.Null{}
	case ast.IDENT:
		object, ok := env.Get(node.Ident)
		if !ok {
			return newError("identifier not found")
		}
		return object
	case ast.NUMBER:
		return &object.Integer{Value: node.Num}
	case ast.RETURN:
		return Eval(node.Lhs, env)
	}

	lhs := Eval(node.Lhs, env)
	rhs := Eval(node.Rhs, env)

	return evalIntegerExpression(node.NodeKind, lhs, rhs)
}
