package token

import (
	"testing"
)

func tokenKindToString(kind TokenKind) string {
	var res string
	switch kind {
	case PLUS:
		res = "PLUS"
	case MINUS:
		res = "MINUS"
	case ASTERISK:
		res = "ASTERISK"
	case SLASH:
		res = "SLASH"
	case NUMBER:
		res = "NUMBER"
	case LPAREN:
		res = "LPAREN"
	case RPAREN:
		res = "RPAREN"
	case ASSIGN:
		res = "ASSIGN"
	case EQ:
		res = "EQ"
	case NOT_EQ:
		res = "NOT_EQ"
	case GT:
		res = "GT"
	case GE:
		res = "GE"
	case LT:
		res = "LT"
	case LE:
		res = "LE"
	case EOF:
		res = "EOF"
	default:
		res = "ILLEGAL"
	}
	return res
}

func TestToken(t *testing.T) {
	input := `
	+ ＋  - ー * ＊ × / ／ ÷ 　21 02356 ０９ １２０ ()（）「」 == ＝＝ != ！＝ < ＜ <= ＜＝ > ＞ >= ＞＝
	`

	tests := []struct {
		expectedTokenKind TokenKind
		expectedLiteral   string
	}{
		{PLUS, "+"},
		{PLUS, "＋"},
		{MINUS, "-"},
		{MINUS, "ー"},
		{ASTERISK, "*"},
		{ASTERISK, "＊"},
		{ASTERISK, "×"},
		{SLASH, "/"},
		{SLASH, "／"},
		{SLASH, "÷"},
		{NUMBER, "21"},
		{NUMBER, "02356"},
		{NUMBER, "０９"},
		{NUMBER, "１２０"},
		{LPAREN, "("},
		{RPAREN, ")"},
		{LPAREN, "（"},
		{RPAREN, "）"},
		{LPAREN, "「"},
		{RPAREN, "」"},
		{EQ, "=="},
		{EQ, "＝＝"},
		{NOT_EQ, "!="},
		{NOT_EQ, "！＝"},
		{GT, "<"},
		{GT, "＜"},
		{GE, "<="},
		{GE, "＜＝"},
		{LT, ">"},
		{LT, "＞"},
		{LE, ">="},
		{LE, "＞＝"},
		{EOF, ""},
	}

	token := Tokenize(input)
	for i, v := range tests {
		if token.Kind != v.expectedTokenKind {
			t.Fatalf("test%d : got=%s expected=%s\n", i, tokenKindToString(token.Kind), tokenKindToString(v.expectedTokenKind))
		}

		if token.Literal != v.expectedLiteral {
			t.Fatalf("test%d : got=\"%s\" expected=\"%s\"\n", i, token.Literal, v.expectedLiteral)
		}

		if token.Next == nil {
			break
		}
		token = token.Next
	}
}
