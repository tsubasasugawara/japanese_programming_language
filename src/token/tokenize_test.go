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
	case CALET:
		res = "CALET"
	case INTEGER:
		res = "INTEGER"
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
	case IDENT:
		res = "IDENT"
	case COMMA:
		res = "COMMA"
	default:
		res = "ILLEGAL"
	}
	return res
}

func TestToken(t *testing.T) {
	input := `
	+ ＋  - ー * ＊ × / ／ ÷ 　21 02356 ０９ １２０ ()（）「」 == ＝＝ != ！＝ < ＜ <= ＜＝ > ＞ >= ＞＝ あ 日 a z ア A Z こんにちは 世界 戻す もし それ以外 ならば 繰り返す {} ｛｝ 関数 , 、 ^ ＾ % ％
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
		{INTEGER, "21"},
		{INTEGER, "02356"},
		{INTEGER, "０９"},
		{INTEGER, "１２０"},
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
		{IDENT, "あ"},
		{IDENT, "日"},
		{IDENT, "a"},
		{IDENT, "z"},
		{IDENT, "ア"},
		{IDENT, "A"},
		{IDENT, "Z"},
		{IDENT, "こんにちは"},
		{IDENT, "世界"},
		{RETURN, "戻す"},
		{IF, "もし"},
		{ELSE, "それ以外"},
		{THEN, "ならば"},
		{FOR, "繰り返す"},
		{LBRACE, "{"},
		{RBRACE, "}"},
		{LBRACE, "｛"},
		{RBRACE, "｝"},
		{FUNC, "関数"},
		{COMMA, ","},
		{COMMA, "、"},
		{CALET, "^"},
		{CALET, "＾"},
		{PARCENT, "%"},
		{PARCENT, "％"},
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
