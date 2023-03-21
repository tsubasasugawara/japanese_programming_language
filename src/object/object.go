package object

import (
	"fmt"
)

type ObjectType string

const (
	INTEGER ObjectType = "INTEGER"
	ERROR = "ERROR"
	BOOLEAN = "BOOLEAN"
)

type Object interface {
	Type() ObjectType
	Inspect() string
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

type Error struct {
	Message string
}
func (e *Error) Type() ObjectType {
	return ERROR
}
func (e *Error) Inspect() string {
	return fmt.Sprintf("Error:%s", e.Message)
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
