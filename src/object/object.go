package object

import (
	"fmt"
	"strings"

	"jpl/ast"
)

type ObjectType string

const (
	INTEGER ObjectType = "INTEGER"
	ERROR = "ERROR"
	BOOLEAN = "BOOLEAN"
	NULL = "NULL"
	RETURN_VALUE = "RETURN_VALUE"
	FUNCTION = "FUNCTION"
	BUILTIN = "BUILTIN"
)

type BuiltinFunction func(args ...Object) Object

type Object interface {
	Type() ObjectType
	Inspect() string
}

type Error struct {
	Message string
}
func (e *Error) Type() ObjectType {
	return ERROR
}
func (e *Error) Inspect() string {
	return fmt.Sprintf("Error:%s", e.Message)
}

type Integer struct {
	Value int64
}
func (i *Integer) Type() ObjectType {
	return INTEGER
}
func (i *Integer) Inspect() string {
	return fmt.Sprintf("%d", i.Value)
}

type Boolean struct {
	Value bool
}
func (b *Boolean) Type() ObjectType {
	return BOOLEAN
}
func (b *Boolean) Inspect() string {
	return fmt.Sprintf("%t", b.Value)
}

type Null struct {}
func (n *Null) Type() ObjectType {
	return NULL
}
func (n *Null) Inspect() string {
	return fmt.Sprintf("null")
}

type ReturnValue struct {
	Value Object
}
func (r *ReturnValue) Type() ObjectType {
	return RETURN_VALUE
}
func (r *ReturnValue) Inspect() string {
	return r.Value.Inspect()
}

type Function struct {
	Params []*ast.Ident
	Body *ast.BlockStmt
}
func (f *Function) Type() ObjectType {
	return FUNCTION
}
func (f *Function) Inspect() string {
	params := []string{}
	for _, v := range f.Params {
		params = append(params, v.Name)
	}

	return fmt.Sprintf("関数(%s)\n", strings.Join(params, ","))
}

type Builtin struct {
	Fn BuiltinFunction
}
func (b *Builtin) Type() string {
	return BUILTIN
}
func (b * Builtin) Inspect() string {
	return "buitin function"
}