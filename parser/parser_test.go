package parser

import (
	"testing"

	"github.com/zawlinnnaing/monkey-language-in-golang/ast"
	"github.com/zawlinnnaing/monkey-language-in-golang/lexer"
)

func TestLetStatement(t *testing.T) {
	input := `let x= 5;
let y = 10;
let foobar = 838383;`

	lexer := lexer.New(input)
	parser := New(lexer)

	program := parser.ParseProgram()
	checkParseErrors(t, parser)

	if program == nil {
		t.Fatalf("ParseProgram() return nil")
	}

	if len(program.Statements) != 3 {
		t.Fatalf("program.Statements does not contain 3 statements. got=%d", len(program.Statements))
	}

	tests := []struct {
		expectedIdentifier string
	}{
		{"x"},
		{"y"},
		{"foobar"},
	}
	for i, testCase := range tests {
		if !testLetStatement(t, program.Statements[i], testCase.expectedIdentifier) {
			return
		}
	}

}

func testLetStatement(t *testing.T, statement ast.Statement, expectedIdentifier string) bool {
	if statement.TokenLiteral() != "let" {
		t.Errorf("s.TokenLiteral not 'let'. got=%q", statement.TokenLiteral())
		return false
	}
	letStatement, ok := statement.(*ast.LetStatement)
	if !ok {
		t.Errorf("statement is not a let statement. got %v", statement)
		return false
	}
	if letStatement.Name.TokenLiteral() != expectedIdentifier {
		t.Errorf("Expected token literal %v, received %v", expectedIdentifier, letStatement.Name.TokenLiteral())
		return false
	}
	if letStatement.Name.Value != expectedIdentifier {
		t.Errorf("Expected identifier %v, received %v", expectedIdentifier, letStatement.Name.Value)
		return false
	}
	return true
}

func checkParseErrors(t *testing.T, parser *Parser) {
	errors := parser.Errors()

	if len(errors) == 0 {
		return
	}

	t.Errorf("parser has %d errors", len(errors))
	for _, msg := range errors {
		t.Errorf("parser error: %v", msg)
	}
	t.FailNow()
}
