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
)

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
	Value int
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
	Params []*ast.Node
	Body *ast.Node
	Env *Environment
}
func (f *Function) Type() ObjectType {
	return FUNCTION
}
func (f *Function) Inspect() string {
	params := []string{}
	for _, v := range f.Params {
		params = append(params, v.Ident)
	}

	return fmt.Sprintf("関数(%s)\n", strings.Join(params, ","))
}