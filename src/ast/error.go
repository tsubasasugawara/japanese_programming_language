package ast

import (
	"fmt"
)

const (
	// シンタックスエラー
	MISSING_RPAREN 						= "括弧が閉じていません"
	MISSING_RBRACE 						= "括弧が閉じていません"
	MISSING_R_SQUARE_BRACE				= "括弧が閉じていません"
	MISSING_FUNCTION_NAME 				= "関数名が指定されていません"
	MISSING_DOUBLE_QOUTES				="引用符(\", ”)が閉じられていません"

	UNEXPECTED_TOKEN					= "予期しないトークンがあります"
	ILLEGAL_CHARACTER					= "対応していない文字が含まれています"
)

type Error interface {
	Type() string
	Category() string
	Message() string
}

type SyntaxError struct {
	category string
	msg string
}
func (s *SyntaxError) Type() string {
	return "シンタックスエラー"
}
func (s *SyntaxError) Category() string {
	return s.category
}
func (s *SyntaxError) Message() string {
	return s.msg
}
func NewSyntaxError(category string, format string, args ...interface{}) *SyntaxError {
	s := &SyntaxError{
		category : category,
		msg : fmt.Sprintf(format, args...),
	}
	return s
}