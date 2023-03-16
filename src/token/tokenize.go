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

func (l *Lexer) skipSpecialChar() {
	for l.ch == '\n' || l.ch == '\t' || l.ch == ' ' || l.ch == '　' {
		l.readChar()
	}
}

func isNum(ch rune) bool {
	return ('0' <= ch && ch <= '9') || ('０' <= ch && ch <= '９')
}

func (l *Lexer) readChar() {
	if l.readPosition >= len(l.input) {
		l.ch = 0
	} else {
		l.ch = l.input[l.readPosition]
	}
	l.position = l.readPosition
	l.readPosition++
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
		case '+', '足':
			cur = newToken(PLUS, cur, string(l.ch))
		case '-', '引':
			cur = newToken(MINUS, cur, string(l.ch))
		case '*', '掛':
			cur = newToken(ASTERISK, cur, string(l.ch))
		case '/', '割':
			cur = newToken(SLASH, cur, string(l.ch))
		case '(', '（':
			cur = newToken(LPAREN, cur, string(l.ch))
		case ')', '）':
			cur = newToken(RPAREN, cur, string(l.ch))
		case 0:
			cur = newToken(EOF, cur, "")
		default:
			if isNum(l.ch) {
				cur = newNumberToken(cur, l.readNum())
			} else {
				cur = newToken(ILLEGAL, cur, "")
			}
		}
		l.readChar()
	}

	return head.Next
}
