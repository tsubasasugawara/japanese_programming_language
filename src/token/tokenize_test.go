package token

import (
	"testing"
)

func TestToken(t *testing.T) {
	input := `
	+ ＋  - ー * ＊ × / ／ ÷ 　21 02356 ０９ １２０ ()（）「」 = ＝ == ＝＝ != ！＝ < ＜ <= ＜＝ > ＞ >= ＞＝
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
		{ASSIGN, "="},
		{ASSIGN, "＝＝"},
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
			t.Fatalf("test%d : got=%d expected=%d\n", i, token.Kind, v.expectedTokenKind)
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
