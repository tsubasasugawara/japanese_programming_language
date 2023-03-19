package token

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

/*
func isJapanese(ch rune) bool {
	r := regexp.MustCompile("[亜-熙ぁ-んァ-ヶ]")
	return r.MatchString(string(ch))
}
*/

/*
func isAlphabet(ch rune) bool {
	return ('a' <= ch && ch <= 'z') || ('A' <= ch && ch <= 'Z') || ('ａ' <= ch && ch <= 'ｚ') || ('Ａ' <= ch && ch <= 'Ｚ')
}
*/

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
		case '/', '／', '÷':
			cur = newToken(SLASH, cur, string(l.ch))
		case '(', '（', '「':
			cur = newToken(LPAREN, cur, string(l.ch))
		case ')', '）', '」':
			cur = newToken(RPAREN, cur, string(l.ch))
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
		default:
			if isNum(l.ch) {
				cur = newNumberToken(cur, l.readNum())
				continue
			} else {
				cur = newToken(ILLEGAL, cur, "")
			}
		}
		l.readChar()
	}

	cur = newToken(EOF, cur, "")
	return head.Next
}
