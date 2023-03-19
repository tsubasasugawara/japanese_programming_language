package token

type TokenKind int

const (
	NUMBER TokenKind = iota

	IDENT //識別子

	PLUS     //+,＋
	MINUS    // -, ー
	SLASH    // /,／,÷
	ASTERISK // *,＊,×
	ASSIGN // =

	GT // <, ＜
	LT // >, ＞
	GE // <=, ＜＝
	LE // >=, ＞＝

	EQ // ==, ＝＝
	NOT_EQ // !=, ！＝

	LPAREN // (,（
	RPAREN // ),）

	EOF
	ILLEGAL
)

/*
var keywords = map[string]TokenKind{
}
*/

type Token struct {
	Kind    TokenKind
	Next    *Token
	Literal string
}

func newToken(kind TokenKind, cur *Token, literal string) *Token {
	token := &Token{Kind: kind, Literal: literal}
	cur.Next = token
	return token
}

func newNumberToken(cur *Token, literal string) *Token {
	token := newToken(NUMBER, cur, literal)
	return token
}

/*
func lookUpIdent(key string) TokenKind {
	if tok, ok := keywords[key]; ok {
		return tok
	}
	return IDENT
}
*/
