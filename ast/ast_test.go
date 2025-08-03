package ast

import (
	"testing"

	"github.com/zawlinnnaing/monkey-language-in-golang/token"
)

func TestString(t *testing.T) {
	program := &Program{
		Statements: []Statement{
			&LetStatement{
				Token: token.Token{Type: token.LET, Literal: "let"},
				Name: &Identifier{
					Token: token.Token{Type: token.IDENT, Literal: "myVar"},
					Value: "myVar",
				},
				Value: &Identifier{
					Token: token.Token{Type: token.IDENT, Literal: "anotherVar"},
					Value: "anotherVar",
				},
			},
		},
	}
	expectedString := "let myVar = anotherVar;"
	if program.String() != expectedString {
		t.Errorf("Expected %v, received %v", expectedString, program.String())
	}
}
