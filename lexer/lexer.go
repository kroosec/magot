package lexer

import "magot/token"

type Lexer struct {
	input     string
	index     int  // current char index
	readIndex int  // current read index (after current char)
	ch        byte // char being examined
}

func New(input string) *Lexer {
	l := &Lexer{input: input}
	l.readChar()
	return l
}

func (l *Lexer) readChar() {
	if l.readIndex >= len(l.input) {
		l.ch = 0
	} else {
		l.ch = l.input[l.readIndex]
	}
	l.index = l.readIndex
	l.readIndex++
}

func (l *Lexer) peekChar() byte {
	if l.readIndex >= len(l.input) {
		return 0
	}
	return l.input[l.readIndex]
}

func (l *Lexer) skipWhitespace() {
	for l.ch == ' ' || l.ch == '\t' || l.ch == '\n' || l.ch == '\r' {
		l.readChar()
	}
}

func newToken(tokType token.TokenType, ch byte) token.Token {
	return token.Token{Type: tokType, Literal: string(ch)}
}

func (l *Lexer) NextToken() token.Token {
	var tok token.Token

	l.skipWhitespace()
	switch l.ch {
	case '=':
		if l.peekChar() == '=' {
			l.readChar()
			tok.Type = token.EQ
			tok.Literal = "=="
		} else {
			tok = newToken(token.ASSIGN, l.ch)
		}
	case '+':
		tok = newToken(token.PLUS, l.ch)
	case '-':
		tok = newToken(token.MINUS, l.ch)
	case '!':
		if l.peekChar() == '=' {
			l.readChar()
			tok.Type = token.NEQ
			tok.Literal = "!="
		} else {
			tok = newToken(token.BANG, l.ch)
		}
	case '*':
		tok = newToken(token.MUL, l.ch)
	case '/':
		tok = newToken(token.DIV, l.ch)
	case '<':
		tok = newToken(token.LT, l.ch)
	case '>':
		tok = newToken(token.GT, l.ch)
	case ',':
		tok = newToken(token.COMMA, l.ch)
	case ';':
		tok = newToken(token.SEMICOLON, l.ch)
	case '(':
		tok = newToken(token.LPAREN, l.ch)
	case ')':
		tok = newToken(token.RPAREN, l.ch)
	case '{':
		tok = newToken(token.LBRACE, l.ch)
	case '}':
		tok = newToken(token.RBRACE, l.ch)
	case '[':
		tok = newToken(token.LBRACKET, l.ch)
	case ']':
		tok = newToken(token.RBRACKET, l.ch)
	case '"':
		tok.Type = token.STRING
		tok.Literal = l.readString()
	case 0:
		tok.Type = token.EOF
		tok.Literal = ""
	default:
		if isIdentifierLetter(l.ch) {
			tok.Literal = l.readIdentifier()
			tok.Type = token.LookupTokenType(tok.Literal)
			return tok
		} else if isDigit(l.ch) {
			tok.Literal = l.readNumber()
			tok.Type = token.INT
			return tok
		} else {
			tok = newToken(token.ILLEGAL, l.ch)
		}
	}
	l.readChar()
	return tok
}

func (l *Lexer) readString() string {
	index := l.index + 1
	for {
		l.readChar()
		if l.ch == '"' || l.ch == 0 {
			break
		}
	}
	return l.input[index:l.index]
}

func (l *Lexer) readNumber() string {
	index := l.index
	for isDigit(l.ch) {
		l.readChar()
	}
	return l.input[index:l.index]
}

func isDigit(ch byte) bool {
	return ch >= '0' && ch <= '9'
}

func isIdentifierLetter(ch byte) bool {
	return (ch >= 'A' && ch <= 'Z') || (ch >= 'a' && ch <= 'z') || ch == '_'
}

func (l *Lexer) readIdentifier() string {
	index := l.index
	for isIdentifierLetter(l.ch) {
		l.readChar()
	}
	return l.input[index:l.index]
}
