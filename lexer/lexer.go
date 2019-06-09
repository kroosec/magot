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
		tok = newToken(token.ASSIGN, l.ch)
	case '+':
		tok = newToken(token.PLUS, l.ch)
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
	case '0':
		tok = newToken(token.EOF, l.ch)
	default:
		if isIdentifierLetter(l.ch) {
			tok.Literal = readIdentifier(l)
			tok.Type = token.LookupTokenType(tok.Literal)
			return tok
		} else {
			tok.Type = token.ILLEGAL
		}
	}
	l.readChar()
	return tok
}

func isIdentifierLetter(ch byte) bool {
	if (ch >= 'A' && ch <= 'Z') || (ch >= 'a' && ch <= 'z') || ch == '_' {
		return true
	}
	return false
}

func readIdentifier(l *Lexer) string {
	index := l.index
	for isIdentifierLetter(l.ch) {
		l.readChar()
	}
	return l.input[index:l.index]
}
