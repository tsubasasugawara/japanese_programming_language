package object

import (
	"fmt"
)

type ObjectType string

const (
	INTEGER = "INTEGER"
	ERROR = "ERROR"
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
