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

func TestReturnStatement(t *testing.T) {
	input := `return 5; return 10; return 993322;`
	lexer := lexer.New(input)
	parser := New(lexer)
	program := parser.ParseProgram()
	checkParseErrors(t, parser)

	if len(program.Statements) != 3 {
		t.Fatalf("Expected return statements to be %d, received: %d", 3, len(program.Statements))
	}

	for _, statement := range program.Statements {
		returnStatement, ok := statement.(*ast.ReturnStatement)
		if !ok {
			t.Errorf("Expected return statement, got: %T", returnStatement)
			continue
		}
		if returnStatement.TokenLiteral() != "return" {
			t.Errorf("Expected token literal to be 'return', received: %s", returnStatement.TokenLiteral())
		}
	}
}

func TestIdentifierExpression(t *testing.T) {
	input := "foobar"

	lexer := lexer.New(input)
	parser := New(lexer)
	program := parser.ParseProgram()
	checkParseErrors(t, parser)
	if len(program.Statements) != 1 {
		t.Fatalf("Expected program statements to be 1, received %d", len(program.Statements))
	}
	statement, ok := (program.Statements[0]).(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("Expected expression statement, received %v", statement)
	}
	ident, ok := statement.Expression.(*ast.Identifier)
	if !ok {
		t.Fatalf("Expected intend statement, received %v", ident)
	}
	if ident.Token.Literal != input {
		t.Fatalf("Expected token literal to be %v, received %v", ident.Token.Literal, input)
	}

}

func TestIntegerLiteralExpression(t *testing.T) {
	input := "5;"
	lexer := lexer.New(input)
	parser := New(lexer)
	program := parser.ParseProgram()
	if len(program.Statements) != 1 {
		t.Fatalf("Expected program statements to be 1, received %d", len(program.Statements))
	}
	statement, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("Expected statement to be ExpressionStatement, received %T", program.Statements[0])
	}

	literal, ok := statement.Expression.(*ast.IntegerLiteral)
	if !ok {
		t.Fatalf("Expected expression to be IntegerLiteral, received %T", statement.Expression)
	}
	if literal.Value != 5 {
		t.Fatalf("Expected value to be %d, received %d", 5, literal.Value)
	}
	if literal.TokenLiteral() != "5" {
		t.Fatalf("Expected token literal to be %s, received %s", "5", literal.TokenLiteral())
	}

}
