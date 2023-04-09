package lexer

import (
	"regexp"
	"unicode"

	"jpl/token"
)

type Lexer struct {
	input        []rune
	position     int
	readPosition int
	ch           rune
}

func newLexer(input string) *Lexer {
	l := &Lexer{input: []rune(input)}
	l.readChar()
	return l
}

func (l *Lexer) peekChar() rune {
	if l.readPosition >= len(l.input) {
		return 0
	} else {
		return l.input[l.readPosition]
	}
}

func (l *Lexer) skipSpecialChar() {
	for l.ch == '\n' || l.ch == '\r' || l.ch == '\t' || l.ch == ' ' || l.ch == '　' {
		l.readChar()
	}
}

func isNum(ch rune) bool {
	return ('0' <= ch && ch <= '9') || ('０' <= ch && ch <= '９')
}

func isHiragana(ch rune) bool {
	match, err := regexp.MatchString("[\u3041-\u3096]", string(ch))
	if err != nil {
		panic(err)
	}
	return match
}

func isKatakana(ch rune) bool {
	match, err := regexp.MatchString("[\u30a1-\u30fc]", string(ch))
	if err != nil {
		panic(err)
	}
	return match
}

func isKanji(ch rune) bool {
	return unicode.In(ch, unicode.Han)
}

func isJapanese(ch rune) bool {
	return isHiragana(ch) || isKatakana(ch) || isKanji(ch)
}

func isAlphabet(ch rune) bool {
	return ('a' <= ch && ch <= 'z') || ('A' <= ch && ch <= 'Z') || ('ａ' <= ch && ch <= 'ｚ') || ('Ａ' <= ch && ch <= 'Ｚ')
}

func (l *Lexer) readChar() {
	if l.readPosition >= len(l.input) {
		l.ch = 0
	} else {
		l.ch = l.input[l.readPosition]
	}
	l.position = l.readPosition
	l.readPosition += 1
}

func (l *Lexer) readNum() string {
	position := l.position
	for isNum(l.ch) {
		l.readChar()
	}
	return string(l.input[position:l.position])
}

func (l *Lexer) readIdentifier() string {
	position := l.position
	if !isAlphabet(l.ch) && !isJapanese(l.ch) {
		return ""
	}
	l.readChar()

	for isAlphabet(l.ch) || isJapanese(l.ch) || isNum(l.ch) ||
		l.ch == '_' || l.ch == '＿' {
		l.readChar()
	}
	return string(l.input[position:l.position])
}

func Tokenize(input string) *token.Token {
	l := newLexer(input)

	head := &token.Token{}
	cur := head

	for l.ch != 0 {
		l.skipSpecialChar()

		switch l.ch {
		case '+', '＋':
			if ch := l.peekChar(); ch == '=' || ch == '＝' {
				cur = token.NewToken(token.PA, cur, string([]rune{l.ch, ch}))
				l.readChar()
			} else {
				cur = token.NewToken(token.PLUS, cur, string(l.ch))
			}
		case '-', 'ー':
			if ch := l.peekChar(); ch == '=' || ch == '＝' {
				cur = token.NewToken(token.MA, cur, string([]rune{l.ch, ch}))
				l.readChar()
			} else {
				cur = token.NewToken(token.MINUS, cur, string(l.ch))
			}
		case '*', '＊', '×':
			if ch := l.peekChar(); ch == '=' || ch == '＝' {
				cur = token.NewToken(token.AA, cur, string([]rune{l.ch, ch}))
				l.readChar()
			} else {
				cur = token.NewToken(token.ASTERISK, cur, string(l.ch))
			}
		case '/', '／':
			if ch := l.peekChar(); ch == '/' || ch == '／' {
				for l.ch != '\n' {
					if l.ch == 0 {
						break
					}
					l.readChar()
				}
			} else if ch := l.peekChar(); ch == '*' || ch == '＊' {
				for !(l.ch == '*' || l.ch == '＊') || !(l.peekChar() == '/' || l.peekChar() == '／') {
					if l.ch == 0 {
						break
					}
					l.readChar()
				}
				l.readChar()
			} else if ch := l.peekChar(); ch == '=' || ch == '＝' {
				cur = token.NewToken(token.SA, cur, string([]rune{l.ch, ch}))
				l.readChar()
			} else {
				cur = token.NewToken(token.SLASH, cur, string(l.ch))
			}
		case '÷':
			if ch := l.peekChar(); ch == '=' || ch == '＝' {
				cur = token.NewToken(token.SA, cur, string([]rune{l.ch, ch}))
				l.readChar()
			} else {
				cur = token.NewToken(token.SLASH, cur, string(l.ch))
			}
		case '^', '＾':
			cur = token.NewToken(token.CALET, cur, string(l.ch))
		case '%', '％':
			cur = token.NewToken(token.PARCENT, cur, string(l.ch))
		case '(', '（', '「':
			cur = token.NewToken(token.LPAREN, cur, string(l.ch))
		case ')', '）', '」':
			cur = token.NewToken(token.RPAREN, cur, string(l.ch))
		case '{', '｛':
			cur = token.NewToken(token.LBRACE, cur, string(l.ch))
		case '}', '｝':
			cur = token.NewToken(token.RBRACE, cur, string(l.ch))
		case '[':
			cur = token.NewToken(token.L_SQUARE_BRACE, cur, string(l.ch))
		case ']':
			cur = token.NewToken(token.R_SQUARE_BRACE, cur, string(l.ch))
		case '~', '〜':
			cur = token.NewToken(token.RANGE, cur, string(l.ch))
		case '<', '＜':
			if ch := l.peekChar(); ch == '=' || ch == '＝' {
				cur = token.NewToken(token.GE, cur, string([]rune{l.ch, ch}))
				l.readChar()
			} else {
				cur = token.NewToken(token.GT, cur, string(l.ch))
			}
		case '>', '＞':
			if ch := l.peekChar(); ch == '=' || ch == '＝' {
				cur = token.NewToken(token.LE, cur, string([]rune{l.ch, ch}))
				l.readChar()
			} else {
				cur = token.NewToken(token.LT, cur, string(l.ch))
			}
		case '=', '＝':
			if ch := l.peekChar(); ch == '=' || ch == '＝' {
				cur = token.NewToken(token.EQ, cur, string([]rune{l.ch, ch}))
				l.readChar()
			} else {
				cur = token.NewToken(token.ASSIGN, cur, string(l.ch))
			}
		case '!', '！':
			if ch := l.peekChar(); ch == '=' || ch == '＝' {
				cur = token.NewToken(token.NOT_EQ, cur, string([]rune{l.ch, ch}))
				l.readChar()
			} else {
				cur = token.NewToken(token.NOT, cur, string(l.ch))
			}
		case ',', '、', '，':
			cur = token.NewToken(token.COMMA, cur, string(l.ch))
		case '&', '＆':
			if ch := l.peekChar(); ch == '&' || ch == '＆' {
				cur = token.NewToken(token.AND, cur, string([]rune{l.ch, ch}))
				l.readChar()
			} else {
				cur = token.NewToken(token.ILLEGAL, cur, string(l.ch))
			}
		case '|', '｜':
			if ch := l.peekChar(); ch == '|' || ch == '｜' {
				cur = token.NewToken(token.OR, cur, string([]rune{l.ch, ch}))
				l.readChar()
			} else {
				cur = token.NewToken(token.ILLEGAL, cur, string(l.ch))
			}
		case '"', '”':
			cur = token.NewToken(token.DOUBLE_QUOTES, cur , string(l.ch))
			l.readChar()

			position := l.position
			for l.ch != '"' && l.ch != '”' && l.ch != 0 {
				l.readChar()
			}
			cur = token.NewToken(token.STRING, cur, string(l.input[position:l.position]))

			if l.ch == '"' || l.ch == '”' {
				cur = token.NewToken(token.DOUBLE_QUOTES, cur, string(l.ch))
			}
		case 0:
			cur = token.NewToken(token.EOF, cur, "")
		default:
			if isNum(l.ch) {
				cur = token.NewIntegerToken(cur, l.readNum())
				continue
			} else if isJapanese(l.ch) || isAlphabet(l.ch) {
				str := l.readIdentifier()
				kind := token.LookUpIdent(str)
				cur = token.NewToken(kind, cur, str)
				continue
			} else {
				cur = token.NewToken(token.ILLEGAL, cur, string(l.ch))
			}
		}
		l.readChar()
	}

	cur = token.NewToken(token.EOF, cur, "")
	return head.Next
}
