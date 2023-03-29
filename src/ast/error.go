package ast

import (
	"fmt"
)

type SyntaxErrorCategory string
const (
	MISSING_RPAREN SyntaxErrorCategory	= "括弧が閉じていません"
	MISSING_RBRACE 						= "括弧が閉じていません"
	MISSING_FUNCTION_NAME 				= "関数名が指定されていません"

	UNEXPECTED_TOKEN					= "予期しないトークンがあります"
	ILLEGAL_CHARACTER					= "対応していない文字が含まれています"
)

type Error interface {
	Message() string
}

type SyntaxError struct {
	category SyntaxErrorCategory
	msg string
}
func (s *SyntaxError) Type() SyntaxErrorCategory {
	return "シンタックスエラー"
}
func (s *SyntaxError) Category() SyntaxErrorCategory {
	return s.category
}
func (s *SyntaxError) Message() string {
	return s.msg
}
func NewSyntaxError(category SyntaxErrorCategory, format string, args ...interface{}) *SyntaxError {
	s := &SyntaxError{
		category : category,
		msg : fmt.Sprintf(format, args...),
	}
	return s
}