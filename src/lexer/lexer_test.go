package lexer

import (
	"testing"

	"jpl/token"
)

func TestToken(t *testing.T) {
	input := `
	+ ＋  - ー * ＊ × / ／ ÷ 　
	21 02356 ０９ １２０
	()（） 「」
	== ＝＝ != ！＝ < ＜ <= ＜＝ > ＞ >= ＞＝
	あ 日 a z ア A Z こんにちは 世界 戻す もし それ以外 ならば 繰り返す
	{} ｛｝
	関数
	, 、 ^ ＾ % ％
	+= ＋＝ -= ー＝ *= ＊＝ ×＝ /= ／＝ ÷＝
	[]
	真 偽 
	&& ＆＆ かつ 
	|| ｜｜ または 
	!　！　ではない 
	~ 〜
	それぞれ繰り返す
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
		{token.TRUE, "真"},
		{token.FALSE, "偽"},
		{token.AND, "&&"},
		{token.AND, "＆＆"},
		{token.AND, "かつ"},
		{token.OR, "||"},
		{token.OR, "｜｜"},
		{token.OR, "または"},
		{token.NOT, "!"},
		{token.NOT, "！"},
		{token.NOT, "ではない"},
		{token.RANGE, "~"},
		{token.RANGE, "〜"},
		{token.FOREACH, "それぞれ繰り返す"},
		{token.EOF, ""},
	}

	head := Tokenize(input)
	for i, v := range tests {
		if head.Kind != v.expectedTokenKind {
			t.Fatalf("test%d : got=%d expected=%d\n", i, head.Kind, v.expectedTokenKind)
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
