package token

type TokenKind int

const (
	INTEGER TokenKind = iota

	IDENT //識別子

	PLUS     // +,＋
	MINUS    // -, ー
	SLASH    // /,／,÷
	ASTERISK // *,＊,×
	CALET	 // ^, ＾
	PARCENT // %, ％
	ASSIGN // =

	PA // +=, ＋＝ PLUS + ASSIGN
	MA // -=,ー＝ MINUS + ASSIGN
	SA // /=, ／＝ ÷＝ SLASH + ASSIGN
	AA // *=, ＊＝ ×＝ ASTERISK + ASSIGN

	GT // <, ＜
	LT // >, ＞
	GE // <=, ＜＝
	LE // >=, ＞＝

	EQ // ==, ＝＝
	NOT_EQ // !=, ！＝

	AND // &&, ＆＆, かつ
	OR // ||, ｜｜, または
	NOT // !, ！, ではない

	LPAREN // (,（, 「
	RPAREN // ), ）, 」
	LBRACE // {, ｛
	RBRACE // }, ｝
	L_SQUARE_BRACE // [
	R_SQUARE_BRACE // ],

	COMMA //, 、

	RANGE //~, 〜

	RETURN
	IF
	ELSE
	THEN
	FOR
	FUNC
	TRUE
	FALSE

	EOF
	ILLEGAL
)

var Keywords = map[string]TokenKind{
	"戻す" : RETURN,
	"もし" : IF,
	"それ以外" : ELSE,
	"ならば" : THEN,
	"繰り返す" : FOR,
	"ならば繰り返す" : FOR,
	"関数" : FUNC,
	"真" : TRUE,
	"偽" : FALSE,
	"かつ" : AND,
	"または" : OR,
	"ではない" : NOT,
}

type Token struct {
	Kind    TokenKind
	Next    *Token
	Literal string
}

func NewToken(kind TokenKind, cur *Token, literal string) *Token {
	token := &Token{Kind: kind, Literal: literal}
	cur.Next = token
	return token
}

func NewIntegerToken(cur *Token, literal string) *Token {
	token := NewToken(INTEGER, cur, literal)
	return token
}

func LookUpIdent(key string) TokenKind {
	if tok, ok := Keywords[key]; ok {
		return tok
	}
	return IDENT
}
