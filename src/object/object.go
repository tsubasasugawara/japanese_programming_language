package object

import (
	"fmt"
)

type ObjectType string

const (
	INTEGER ObjectType = "INTEGER"
	ERROR = "ERROR"
	BOOLEAN = "BOOLEAN"
	NULL = "NULL"
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
