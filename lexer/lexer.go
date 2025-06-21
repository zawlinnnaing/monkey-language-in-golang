package lexer

import (
	"github.com/zawlinnnaing/monkey-language-in-golang/token"
)

type Lexer struct {
	input        string
	position     int
	readPosition int
	ch           byte
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

func (l *Lexer) NextToken() token.Token {
	var tok token.Token

	l.skipsWhitespace()

	switch l.ch {
	case '=':
		tok = *token.New(token.ASSIGN, string(l.ch))
	case '+':
		tok = *token.New(token.PLUS, string(l.ch))
	case ',':
		tok = *token.New(token.COMMA, string(l.ch))
	case ';':
		tok = *token.New(token.SEMICOLON, string(l.ch))
	case '(':
		tok = *token.New(token.LPAREN, string(l.ch))
	case ')':
		tok = *token.New(token.RPAREN, string(l.ch))
	case '{':
		tok = *token.New(token.LBRACE, string(l.ch))
	case '}':
		tok = *token.New(token.RBRACE, string(l.ch))
	case 0:
		tok = *token.New(token.EOF, "")
	default:
		if isLetter(l.ch) {
			tok.Literal = l.readIdentifier()
			tok.Type = token.LookupIdentifier(tok.Literal)
			return tok
		}
		if isDigit(l.ch) {
			tok.Literal = l.readDigit()
			tok.Type = token.INT
			return tok
		}
		tok = *token.New(token.ILLEGAL, string(l.ch))

	}

	l.readChar()

	return tok
}

func (l *Lexer) readIdentifier() string {
	position := l.position
	for isLetter(l.ch) {
		l.readChar()
	}
	return l.input[position:l.position]
}

func (l *Lexer) readDigit() string {
	position := l.position
	for isDigit(l.ch) {
		l.readChar()
	}
	return l.input[position:l.position]
}

func (l *Lexer) skipsWhitespace() {
	for l.ch == ' ' || l.ch == '\t' || l.ch == '\n' || l.ch == '\r' {
		l.readChar()
	}
}

func isDigit(ch byte) bool {
	return '0' <= ch && ch <= '9'
}

func isLetter(ch byte) bool {
	return ('a' <= ch && ch <= 'z') || ('A' <= ch && ch <= 'Z') || ch == '_'
}

func New(input string) *Lexer {
	l := &Lexer{
		input: input,
	}
	l.readChar()
	return l
}
