package main

type Lexer struct {
	input        string
	position     int  // current position in input (points to current char)
	readPosition int  // current reading position in input (after current char)
	ch           byte // current char under examination
}

func NewLexer(input string) *Lexer {
	l := &Lexer{input: input}
	l.readChar()
	return l
}

func (l *Lexer) NextToken() Token {
	var token Token

	l.skipWhitespace()

	switch l.ch {
	case '=':
		if l.peakChar() == '=' {
			ch := l.ch
			l.readChar()
			token = Token{Type: EQ, Literal: string(ch) + string(l.ch)}
		} else {
			token = newToken(ASSIGN, l.ch)
		}
	case '+':
		token = newToken(PLUS, l.ch)
	case '-':
		token = newToken(MINUS, l.ch)
	case '!':
		if l.peakChar() == '=' {
			ch := l.ch
			l.readChar()
			token = Token{Type: NOT_EQ, Literal: string(ch) + string(l.ch)}
		} else {
			token = newToken(BANG, l.ch)
		}
	case '/':
		token = newToken(SLASH, l.ch)
	case '*':
		token = newToken(ASTERISK, l.ch)
	case '<':
		token = newToken(LT, l.ch)
	case '>':
		token = newToken(GT, l.ch)
	case ';':
		token = newToken(SEMICOLON, l.ch)
	case '(':
		token = newToken(LPAREN, l.ch)
	case ')':
		token = newToken(RPAREN, l.ch)
	case ',':
		token = newToken(COMMA, l.ch)
	case '{':
		token = newToken(LBRACE, l.ch)
	case '}':
		token = newToken(RBRACE, l.ch)
	case '"':
		token.Type = STRING
		token.Literal = l.readString()
	case 0:
		token.Literal = ""
		token.Type = EOF
	default:
		if isLetter(l.ch) {
			token.Literal = l.readIdentifier()
			token.Type = LookupIdent(token.Literal)
			return token
		} else if isDigit(l.ch) {
			token.Type = INT
			token.Literal = l.readNumber()
			return token
		} else {
			token = newToken(ILLEGAL, l.ch)
		}
	}

	l.readChar()

	return token
}

func (l *Lexer) readNumber() string {
	position := l.position
	for isDigit(l.ch) {
		l.readChar()
	}
	return l.input[position:l.position]
}

func (l *Lexer) skipWhitespace() {
	for l.ch == ' ' || l.ch == '\t' || l.ch == '\n' || l.ch == '\r' {
		l.readChar()
	}
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

func (l *Lexer) peakChar() byte {
	if l.readPosition >= len(l.input) {
		return 0
	} else {
		return l.input[l.readPosition]
	}
}

func (l *Lexer) readIdentifier() string {
	position := l.position
	for isLetter(l.ch) {
		l.readChar()
	}
	return l.input[position:l.position]
}

func isLetter(ch byte) bool {
	return 'a' <= ch && ch <= 'z' || 'A' <= ch && ch <= 'Z' || ch == '_'
}

func isDigit(ch byte) bool {
	return '0' <= ch && ch <= '9'
}

func LookupIdent(ident string) TokenType {
	if token, ok := keywords[ident]; ok {
		return token
	}
	return IDENT
}

func (l *Lexer) readString() string {
	position := l.position + 1
	for {
		l.readChar()
		if l.ch == '"' || l.ch == 0 {
			break
		}
	}
	return l.input[position:l.position]
}
