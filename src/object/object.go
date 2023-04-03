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
	ARRAY = "ARRAY"
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

type (
	Error struct {
		Message string
	}

	Array struct {
		Elements []Object
	}

	Integer struct {
		Value int64
	}

	Boolean struct {
		Value bool
	}

	Null struct {}

	ReturnValue struct {
		Value Object
	}

	Function struct {
		Params []*ast.Ident
		Body *ast.BlockStmt
	}
	
	Builtin struct {
		Fn BuiltinFunction
	}
)

func (e *Error) Type() ObjectType { return ERROR }
func (a *Array) Type() ObjectType { return ARRAY }
func (i *Integer) Type() ObjectType { return INTEGER }
func (b *Boolean) Type() ObjectType { return BOOLEAN }
func (n *Null) Type() ObjectType { return NULL }
func (r *ReturnValue) Type() ObjectType { return RETURN_VALUE }
func (f *Function) Type() ObjectType { return FUNCTION }
func (b *Builtin) Type() string { return BUILTIN }

func (e *Error) Inspect() string { return fmt.Sprintf("Error:%s", e.Message) }
func (a *Array) Inspect() string {
	elements := []string{}
	for _, v := range a.Elements {
		elements = append(elements, v.Inspect())
	}

	return fmt.Sprintf("{" + strings.Join(elements, ",") + "}")
}
func (i *Integer) Inspect() string { return fmt.Sprintf("%d", i.Value) }
func (b *Boolean) Inspect() string { return fmt.Sprintf("%t", b.Value) }
func (n *Null) Inspect() string { return fmt.Sprintf("null") }
func (r *ReturnValue) Inspect() string { return r.Value.Inspect() }
func (f *Function) Inspect() string {
	params := []string{}
	for _, v := range f.Params {
		params = append(params, v.Name)
	}

	return fmt.Sprintf("関数(%s)\n", strings.Join(params, ","))
}
func (b * Builtin) Inspect() string { return "buitin function" }