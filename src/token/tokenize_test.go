package token

import (
	"testing"
)

func TestToken(t *testing.T) {
	input := `
	+  足 - 引 * 掛 / 割 　21 02356 ０９ １２０ ()（）
	`

	tests := []struct {
		expectedTokenKind TokenKind
		expectedLiteral   string
	}{
		{PLUS, "+"},
		{PLUS, "足"},
		{MINUS, "-"},
		{MINUS, "引"},
		{ASTERISK, "*"},
		{ASTERISK, "掛"},
		{SLASH, "/"},
		{SLASH, "割"},
		{NUMBER, "21"},
		{NUMBER, "02356"},
		{NUMBER, "０９"},
		{NUMBER, "１２０"},
		{LPAREN, "("},
		{RPAREN, ")"},
		{LPAREN, "（"},
		{RPAREN, "）"},
		{EOF, ""},
	}

	token := Tokenize(input)
	for i, v := range tests {
		if token.Kind != v.expectedTokenKind {
			t.Fatalf("test%d : got=%d expected=%d\n", i, token.Kind, v.expectedTokenKind)
		}

		if token.Literal != v.expectedLiteral {
			t.Fatalf("test%d : got=%s expected=%s\n", i, token.Literal, v.expectedLiteral)
		}

		if token.Next == nil {
			break
		}
		token = token.Next
	}
}
