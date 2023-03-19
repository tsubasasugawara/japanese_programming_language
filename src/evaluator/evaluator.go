package evaluator

import (
	"fmt"

	"jpl/ast"
	"jpl/object"
)

func newError(format string, a ...interface{}) *object.Error {
	return &object.Error{Message: fmt.Sprintf(format, a...)}
}

func evalIntegerExpression(nodeKind ast.NodeKind, left object.Object, right object.Object) object.Object {
	leftVal := left.(*object.Integer).Value
	rightVal := right.(*object.Integer).Value

	switch nodeKind {
	case ast.ADD:
		return &object.Integer{Value: leftVal + rightVal}
	case ast.SUB:
		return &object.Integer{Value: leftVal - rightVal}
	case ast.MUL:
		return &object.Integer{Value: leftVal * rightVal}
	case ast.DIV:
		return &object.Integer{Value: leftVal / rightVal}
	default:
		return newError("unknown operator")
	}
}

func Eval(node *ast.Node) object.Object {
	switch node.NodeKind {
	case ast.NUMBER:
		return &object.Integer{Value: node.Num}
	}

	lhs := Eval(node.Lhs)
	rhs := Eval(node.Rhs)

	return evalIntegerExpression(node.NodeKind, lhs, rhs)
}
