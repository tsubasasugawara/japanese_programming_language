package token

type TokenKind int

const (
	NUMBER TokenKind = iota

	IDENT

	PLUS
	MINUS
	SLASH
	ASTERISK

	LPAREN
	RPAREN

	EOF
	ILLEGAL
)

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
