package lexer

import (
	"testing"

	"jpl/token"
)

func tokenKindToString(kind token.TokenKind) string {
	var res string
	switch kind {
	case token.PLUS:
		res = "PLUS"
	case token.MINUS:
		res = "MINUS"
	case token.ASTERISK:
		res = "ASTERISK"
	case token.SLASH:
		res = "SLASH"
	case token.CALET:
		res = "CALET"
	case token.PA:
		res = "PLUS+ASSIGN"
	case token.MA:
		res = "MINUS+ASSIGN"
	case token.AA:
		res = "ASTERISK+ASSIGN"
	case token.SA:
		res = "SLASH+ASSIGN"
	case token.INTEGER:
		res = "INTEGER"
	case token.LPAREN:
		res = "LPAREN"
	case token.RPAREN:
		res = "RPAREN"
	case token.LBRACE:
		res = "LBRACE"
	case token.RBRACE:
		res ="RBRACE"
	case token.L_SQUARE_BRACE:
		res = "L_SQUARE_BRACE"
	case token.R_SQUARE_BRACE:
		res = "R_SQUARE_BRACE"
	case token.ASSIGN:
		res = "ASSIGN"
	case token.EQ:
		res = "EQ"
	case token.NOT_EQ:
		res = "NOT_EQ"
	case token.GT:
		res = "GT"
	case token.GE:
		res = "GE"
	case token.LT:
		res = "LT"
	case token.LE:
		res = "LE"
	case token.EOF:
		res = "EOF"
	case token.IDENT:
		res = "IDENT"
	case token.COMMA:
		res = "COMMA"
	default:
		res = "ILLEGAL"
	}
	return res
}

func TestToken(t *testing.T) {
	input := `
	+ ＋  - ー * ＊ × / ／ ÷ 　21 02356 ０９ １２０ ()（） 「」  == ＝＝ != ！＝ < ＜ <= ＜＝ > ＞ >= ＞＝ あ 日 a z ア A Z こんにちは 世界 戻す もし それ以外 ならば 繰り返す {} ｛｝ 関数 , 、 ^ ＾ % ％ += ＋＝ -= ー＝ *= ＊＝ ×＝ /= ／＝ ÷＝ []
	`

	tests := []struct {
		expectedTokenKind token.TokenKind
		expectedLiteral   string
	}{
		{token.PLUS, "+"},
		{token.PLUS, "＋"},
		{token.MINUS, "-"},
		{token.MINUS, "ー"},
		{token.ASTERISK, "*"},
		{token.ASTERISK, "＊"},
		{token.ASTERISK, "×"},
		{token.SLASH, "/"},
		{token.SLASH, "／"},
		{token.SLASH, "÷"},
		{token.INTEGER, "21"},
		{token.INTEGER, "02356"},
		{token.INTEGER, "０９"},
		{token.INTEGER, "１２０"},
		{token.LPAREN, "("},
		{token.RPAREN, ")"},
		{token.LPAREN, "（"},
		{token.RPAREN, "）"},
		{token.LPAREN, "「"},
		{token.RPAREN, "」"},
		{token.EQ, "=="},
		{token.EQ, "＝＝"},
		{token.NOT_EQ, "!="},
		{token.NOT_EQ, "！＝"},
		{token.GT, "<"},
		{token.GT, "＜"},
		{token.GE, "<="},
		{token.GE, "＜＝"},
		{token.LT, ">"},
		{token.LT, "＞"},
		{token.LE, ">="},
		{token.LE, "＞＝"},
		{token.IDENT, "あ"},
		{token.IDENT, "日"},
		{token.IDENT, "a"},
		{token.IDENT, "z"},
		{token.IDENT, "ア"},
		{token.IDENT, "A"},
		{token.IDENT, "Z"},
		{token.IDENT, "こんにちは"},
		{token.IDENT, "世界"},
		{token.RETURN, "戻す"},
		{token.IF, "もし"},
		{token.ELSE, "それ以外"},
		{token.THEN, "ならば"},
		{token.FOR, "繰り返す"},
		{token.LBRACE, "{"},
		{token.RBRACE, "}"},
		{token.LBRACE, "｛"},
		{token.RBRACE, "｝"},
		{token.FUNC, "関数"},
		{token.COMMA, ","},
		{token.COMMA, "、"},
		{token.CALET, "^"},
		{token.CALET, "＾"},
		{token.PARCENT, "%"},
		{token.PARCENT, "％"},
		{token.PA, "+="},
		{token.PA, "＋＝"},
		{token.MA, "-="},
		{token.MA, "ー＝"},
		{token.AA, "*="},
		{token.AA, "＊＝"},
		{token.AA, "×＝"},
		{token.SA, "/="},
		{token.SA, "／＝"},
		{token.SA, "÷＝"},
		{token.L_SQUARE_BRACE, "["},
		{token.R_SQUARE_BRACE, "]"},
		{token.EOF, ""},
	}

	head := Tokenize(input)
	for i, v := range tests {
		if head.Kind != v.expectedTokenKind {
			t.Fatalf("test%d : got=%s expected=%s\n", i, tokenKindToString(head.Kind), tokenKindToString(v.expectedTokenKind))
		}

		if head.Literal != v.expectedLiteral {
			t.Fatalf("test%d : got=\"%s\" expected=\"%s\"\n", i, head.Literal, v.expectedLiteral)
		}

		if head.Next == nil {
			break
		}
		head = head.Next
	}
}
