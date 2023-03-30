package token

import (
	"regexp"
	"unicode"
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
	for l.ch == '\n' || l.ch == '\t' || l.ch == ' ' || l.ch == '　' {
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

func (l *Lexer) readString() string {
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

func Tokenize(input string) *Token {
	l := newLexer(input)

	head := &Token{}
	cur := head

	for l.ch != 0 {
		l.skipSpecialChar()

		switch l.ch {
		case '+', '＋':
			cur = newToken(PLUS, cur, string(l.ch))
		case '-', 'ー':
			cur = newToken(MINUS, cur, string(l.ch))
		case '*', '＊', '×':
			cur = newToken(ASTERISK, cur, string(l.ch))
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
			} else {
				cur = newToken(SLASH, cur, string(l.ch))
			}
		case '÷':
			cur = newToken(SLASH, cur, string(l.ch))
		case '^', '＾':
			cur = newToken(CALET, cur, string(l.ch))
		case '(', '（', '「':
			cur = newToken(LPAREN, cur, string(l.ch))
		case ')', '）', '」':
			cur = newToken(RPAREN, cur, string(l.ch))
		case '{', '｛':
			cur = newToken(LBRACE, cur, string(l.ch))
		case '}', '｝':
			cur = newToken(RBRACE, cur, string(l.ch))
		case '<', '＜':
			if ch := l.peekChar(); ch == '=' || ch == '＝' {
				cur = newToken(GE, cur, string([]rune{l.ch, ch}))
				l.readChar()
			} else {
				cur = newToken(GT, cur, string(l.ch))
			}
		case '>', '＞':
			if ch := l.peekChar(); ch == '=' || ch == '＝' {
				cur = newToken(LE, cur, string([]rune{l.ch, ch}))
				l.readChar()
			} else {
				cur = newToken(LT, cur, string(l.ch))
			}
		case '=', '＝':
			if ch := l.peekChar(); ch == '=' || ch == '＝' {
				cur = newToken(EQ, cur, string([]rune{l.ch, ch}))
				l.readChar()
			} else {
				cur = newToken(ASSIGN, cur, string(l.ch))
			}
		case '!', '！':
			if ch := l.peekChar(); ch == '=' || ch == '＝' {
				cur = newToken(NOT_EQ, cur, string([]rune{l.ch, ch}))
				l.readChar()
			} else {
				cur = newToken(ILLEGAL, cur, string([]rune{l.ch, l.peekChar()}))
			}
		case ',', '、', '，':
			cur = newToken(COMMA, cur, string(l.ch))
		case 0:
			cur = newToken(EOF, cur, "")
		default:
			if isNum(l.ch) {
				cur = newIntegerToken(cur, l.readNum())
				continue
			} else if isJapanese(l.ch) || isAlphabet(l.ch) {
				str := l.readString()
				kind := lookUpIdent(str)
				cur = newToken(kind, cur, str)
				continue
			} else {
				cur = newToken(ILLEGAL, cur, string(l.ch))
			}
		}
		l.readChar()
	}

	cur = newToken(EOF, cur, "")
	return head.Next
}
