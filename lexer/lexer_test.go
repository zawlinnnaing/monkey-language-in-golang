package lexer

import (
	"testing"

	"github.com/zawlinnnaing/monkey-language-in-golang/token"
)

func TestLexer(t *testing.T) {
	input := `let five = 5; let ten = 10; let add = fn(x, y) {x + y; };let result = add(five, ten); `
	expected := []struct {
		expectedTokenType token.TokenType
		expectedLiteral   string
	}{
		{token.LET, "let"},
		{token.IDENT, "five"},
		{token.ASSIGN, "="},
		{token.INT, "5"},
		{token.SEMICOLON, ";"},
		{token.LET, "let"},
		{token.IDENT, "ten"},
		{token.ASSIGN, "="},
		{token.INT, "10"},
		{token.SEMICOLON, ";"},
		{token.LET, "let"},
		{token.IDENT, "add"},
		{token.ASSIGN, "="},
		{token.FUNCTION, "fn"},
		{token.LPAREN, "("},
		{token.IDENT, "x"},
		{token.COMMA, ","},
		{token.IDENT, "y"},
		{token.RPAREN, ")"},
		{token.LBRACE, "{"},
		{token.IDENT, "x"},
		{token.PLUS, "+"},
		{token.IDENT, "y"},
		{token.SEMICOLON, ";"},
		{token.RBRACE, "}"},
		{token.SEMICOLON, ";"},

		{token.LET, "let"},
		{token.IDENT, "result"},
		{token.ASSIGN, "="},
		{token.IDENT, "add"},
		{token.LPAREN, "("},
		{token.IDENT, "five"},
		{token.COMMA, ","},
		{token.IDENT, "ten"},
		{token.RPAREN, ")"},
		{token.SEMICOLON, ";"},
		{token.EOF, ""},
	}

	lexer := New(input)

	for i, expectedToken := range expected {
		actualToken := lexer.NextToken()
		if actualToken.Type != expectedToken.expectedTokenType {
			t.Errorf("Test[%d]: Expected token type: %s, received: %s", i, expectedToken.expectedTokenType, actualToken.Type)
		}
		if actualToken.Literal != expectedToken.expectedLiteral {
			t.Errorf("Test[%d]: Expected token literal: %s, received: %s", i, expectedToken.expectedLiteral, actualToken.Literal)
		}
	}
}
