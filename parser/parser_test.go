package parser

import (
	"fmt"
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

func TestPrefixExpression(t *testing.T) {
	testCases := []struct {
		input    string
		operator string
		value    int64
	}{
		{"!5", "!", 5},
		{"-15", "-", 15},
	}
	for _, testCase := range testCases {
		lexer := lexer.New(testCase.input)
		parser := New(lexer)
		program := parser.ParseProgram()
		checkParseErrors(t, parser)
		if len(program.Statements) != 1 {
			t.Fatalf("Expected program statements to be 1, received %d", len(program.Statements))
		}

		statement, ok := program.Statements[0].(*ast.ExpressionStatement)
		if !ok {
			t.Fatalf("Expected statement to be ExpressionStatement, received %T", program.Statements[0])
		}

		prefixExpr, ok := statement.Expression.(*ast.PrefixExpression)
		if !ok {
			t.Fatalf("Expected expression to be PrefixExpression, received %T", statement.Expression)
		}

		if prefixExpr.Operator != testCase.operator {
			t.Fatalf("Expected operator to be %s, received %s", testCase.operator, prefixExpr.Operator)
		}

		if !testIntegerLiteral(t, prefixExpr.Right, testCase.value) {
			return
		}
	}
}

func TestOperatorPrecedenceString(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"-a * b", "((-a) * b)"},
		{"!-a", "(!(-a))"},
		{"a + b + c", "((a + b) + c)"},
		{"a + b - c", "((a + b) - c)"},
		{"a * b * c", "((a * b) * c)"},
		{"a * b / c", "((a * b) / c)"},
		{"a + b / c", "(a + (b / c))"},
		{"a + b * c + d / e - f", "(((a + (b * c)) + (d / e)) - f)"},
		{"3 + 4; -5 * 5", "(3 + 4)((-5) * 5)"},
		{"5 > 4 == 3 < 4", "((5 > 4) == (3 < 4))"},
		{"5 < 4 != 3 > 4", "((5 < 4) != (3 > 4))"},
		{"3 + 4 * 5 == 3 * 1 + 4 * 5", "((3 + (4 * 5)) == ((3 * 1) + (4 * 5)))"},
	}

	for _, tt := range tests {
		l := lexer.New(tt.input)
		p := New(l)
		program := p.ParseProgram()
		checkParseErrors(t, p)

		actual := program.String()
		if actual != tt.expected {
			t.Errorf("Expected=%q, received=%q", tt.expected, actual)
		}
	}
}

func testIntegerLiteral(t *testing.T, exp ast.Expression, value int64) bool {
	integerLiteral, ok := exp.(*ast.IntegerLiteral)
	if !ok {
		t.Errorf("Expected *ast.IntegerLiteral, received %T", exp)
		return false
	}
	if integerLiteral.Value != value {
		t.Errorf("Expected value to be %d, received %d", value, integerLiteral.Value)
		return false
	}
	if integerLiteral.TokenLiteral() != fmt.Sprintf("%d", value) {
		t.Errorf("Expected token literal to be %d, received %s", value, integerLiteral.TokenLiteral())
		return false
	}
	return true
}

func testIdentifier(t *testing.T, exp ast.Expression, val string) bool {
	ident, ok := exp.(*ast.Identifier)
	if !ok {
		t.Errorf("expected identifier, received %T", exp)
		return false
	}
	if ident.Value != val {
		t.Errorf("expected %v, received %v", val, ident.Value)
		return false
	}
	if ident.TokenLiteral() != val {
		t.Errorf("TokenLiteral, expected %v, received %v", val, ident.TokenLiteral())
		return false
	}
	return true
}

func testLiteralExpression(t *testing.T, expression ast.Expression, expected any) bool {
	switch v := expected.(type) {
	case int:
		return testIntegerLiteral(t, expression, int64(v))
	case int64:
		return testIntegerLiteral(t, expression, v)
	case string:
		return testIdentifier(t, expression, v)
	}
	t.Errorf("Unknown expected type: %T", expected)
	return false
}

func testInfixExpression(t *testing.T, expression ast.Expression, left any, operator string, right any) bool {
	infixExpression, ok := expression.(*ast.InfixExpression)
	if !ok {
		t.Errorf("Expected infix expression, received %T", expression)
		return false
	}
	if !testLiteralExpression(t, infixExpression.Left, left) {
		return false
	}
	if infixExpression.Operator != operator {
		t.Errorf("Expected operator: %v, received %v", operator, infixExpression.Operator)
		return false
	}
	if !testLiteralExpression(t, infixExpression.Right, right) {
		return false
	}
	return true
}

func TestInfixExpression(t *testing.T) {
	testCases := []struct {
		input      string
		leftValue  int64
		operator   string
		rightValue int64
	}{
		{"5 + 5;", 5, "+", 5},
		{"5 - 5;", 5, "-", 5},
		{"5 * 5;", 5, "*", 5},
		{"5 / 5;", 5, "/", 5},
		{"5 > 5;", 5, ">", 5},
		{"5 < 5;", 5, "<", 5},
		{"5 == 5;", 5, "==", 5},
		{"5 != 5;", 5, "!=", 5},
	}

	for _, testCase := range testCases {
		lexer := lexer.New(testCase.input)
		parser := New(lexer)
		program := parser.ParseProgram()
		checkParseErrors(t, parser)

		if len(program.Statements) != 1 {
			t.Fatalf("Expected program statements to be 1, received %d", len(program.Statements))
		}

		statement, ok := program.Statements[0].(*ast.ExpressionStatement)
		if !ok {
			t.Fatalf("Expected statement to be ExpressionStatement, received %T", program.Statements[0])
		}
		if !testInfixExpression(t, statement.Expression, testCase.leftValue, testCase.operator, testCase.rightValue) {
			return
		}
	}
}
