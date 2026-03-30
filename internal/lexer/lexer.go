package lexer

type Lexer struct {
	Input        string
	Position     int
	ReadPosition int
	Ch           byte
}

func New(input string) *Lexer {
	l := &Lexer{Input: input}
	l.ReadChar()
	return l
}

func (l *Lexer) ReadChar() {
	if l.ReadPosition >= len(l.Input) {
		l.Ch = 0
	} else {
		l.Ch = l.Input[l.ReadPosition]
	}
	l.Position = l.ReadPosition
	l.ReadPosition += 1
}

func (l *Lexer) PeekChar() byte {
	if l.ReadPosition >= len(l.Input) {
		return 0
	}
	return l.Input[l.ReadPosition]
}

func (l *Lexer) readString() string {
	position := l.Position + 1
	for {
		l.ReadChar()
		if l.Ch == '\'' || l.Ch == 0 {
			break
		}
	}
	return l.Input[position:l.Position]
}

func (l *Lexer) readIdentifier() string {
	position := l.Position
	for isLetter(l.Ch) {
		l.ReadChar()
	}
	return l.Input[position:l.Position]
}

func (l *Lexer) readNumber() string {
	position := l.Position
	for isDigit(l.Ch) {
		l.ReadChar()
	}
	return l.Input[position:l.Position]
}

func (l *Lexer) skipWhitespace() {
	for l.Ch == ' ' || l.Ch == '\t' || l.Ch == '\n' || l.Ch == '\r' {
		l.ReadChar()
	}
}

func isLetter(ch byte) bool {
	return 'a' <= ch && ch <= 'z' || 'A' <= ch && ch <= 'Z' || ch == '_'
}

func isDigit(ch byte) bool {
	return '0' <= ch && ch <= '9'
}

func (l *Lexer) NextToken() Token {
	var tok Token

	l.skipWhitespace()

	switch l.Ch {
	case '=':
		if l.PeekChar() == '=' {
			ch := l.Ch
			l.ReadChar()
			tok = Token{Type: EQ, Literal: string(ch) + string(l.Ch)}
		} else {
			tok = Token{Type: ASSIGN, Literal: string(l.Ch)}
		}
	case '>':
		if l.PeekChar() == '=' {
			ch := l.Ch
			l.ReadChar()
			tok = Token{Type: GE, Literal: string(ch) + string(l.Ch)}
		} else {
			tok = Token{Type: GT, Literal: string(l.Ch)}
		}
	case '<':
		if l.PeekChar() == '=' {
			ch := l.Ch
			l.ReadChar()
			tok = Token{Type: LE, Literal: string(ch) + string(l.Ch)}
		} else {
			tok = Token{Type: LT, Literal: string(l.Ch)}
		}
	case '!':
		if l.PeekChar() == '=' {
			ch := l.Ch
			l.ReadChar()
			tok = Token{Type: NE, Literal: string(ch) + string(l.Ch)}
		} else {
			tok = Token{Type: ILLEGAL, Literal: string(l.Ch)}
		}
	case 0:
		tok = Token{Type: EOF, Literal: ""}
	case ',':
		tok = Token{Type: COMMA, Literal: string(l.Ch)}
	case '.':
		tok = Token{Type: DOT, Literal: string(l.Ch)}
	case '*':
		tok = Token{Type: ASTERISK, Literal: string(l.Ch)}
	case '\'':
		tok.Type = STRING
		tok.Literal = l.readString()
		l.ReadChar()
		return tok
	case '(':
		tok = Token{Type: LPAREN, Literal: string(l.Ch)}
	case ')':
		tok = Token{Type: RPAREN, Literal: string(l.Ch)}
	default:
		if isLetter(l.Ch) {
			tok.Literal = l.readIdentifier()
			tok.Type = LookupIdent(tok.Literal)
			return tok
		} else if isDigit(l.Ch) {
			tok.Type = INT
			tok.Literal = l.readNumber()
			return tok
		} else {
			tok = Token{Type: ILLEGAL, Literal: string(l.Ch)}
		}
	}

	l.ReadChar()
	return tok
}
