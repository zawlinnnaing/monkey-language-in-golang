package parser

import (
	"fmt"
	"testing"

	"github.com/zawlinnnaing/monkey-language-in-golang/ast"
	"github.com/zawlinnnaing/monkey-language-in-golang/lexer"
)

func testExpressionStatement(t *testing.T, statement ast.Statement) bool {
	_, ok := statement.(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("Expected *ast.ExpressionStatement, received %T", statement)
		return false
	}
	return true
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

func testBooleanLiteral(t *testing.T, exp ast.Expression, value bool) bool {
	booleanLiteral, ok := exp.(*ast.BooleanLiteral)
	if !ok {
		t.Errorf("Expected *ast.BooleanLiteral, received %T", exp)
		return false
	}
	if booleanLiteral.Value != value {
		t.Errorf("Expected value to be %t, received %t", value, booleanLiteral.Value)
		return false
	}
	if booleanLiteral.TokenLiteral() != fmt.Sprintf("%t", value) {
		t.Errorf("Expected token literal to be %t, received %s", value, booleanLiteral.TokenLiteral())
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
	case bool:
		return testBooleanLiteral(t, expression, v)
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
		t.Fatalf("Expected statement, received %T", program.Statements[0])
	}
	if !testIntegerLiteral(t, statement.Expression, 5) {
		return
	}
}

func TestPrefixExpression(t *testing.T) {
	testCases := []struct {
		input    string
		operator string
		value    any
	}{
		{"!5", "!", 5},
		{"-15", "-", 15},
		{"!true;", "!", true},
		{"!false", "!", false},
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

		if !testLiteralExpression(t, prefixExpr.Right, testCase.value) {
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
		{"true", "true"},
		{"false", "false"},
		{"3 > 5 == false", "((3 > 5) == false)"},
		{"3 < 5 == true", "((3 < 5) == true)"},
		{"1 + (2 + 3) + 4", "((1 + (2 + 3)) + 4)"},
		{"(5 + 5)* 2", "((5 + 5) * 2)"},
		{"2 / (5 + 5)", "(2 / (5 + 5))"},
		{"-(5 + 5)", "(-(5 + 5))"},
		{"!(true == true)", "(!(true == true))"},
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

func TestInfixExpression(t *testing.T) {
	testCases := []struct {
		input      string
		leftValue  any
		operator   string
		rightValue any
	}{
		{"5 + 5;", 5, "+", 5},
		{"5 - 5;", 5, "-", 5},
		{"5 * 5;", 5, "*", 5},
		{"5 / 5;", 5, "/", 5},
		{"5 > 5;", 5, ">", 5},
		{"5 < 5;", 5, "<", 5},
		{"5 == 5;", 5, "==", 5},
		{"5 != 5;", 5, "!=", 5},
		{"true != false", true, "!=", false},
		{"true == true", true, "==", true},
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

func TestBooleanLiteralExpression(t *testing.T) {
	tests := []struct {
		input           string
		expectedBoolean bool
	}{
		{"true;", true},
		{"false;", false},
	}

	for _, tt := range tests {
		l := lexer.New(tt.input)
		p := New(l)
		program := p.ParseProgram()
		checkParseErrors(t, p)

		if len(program.Statements) != 1 {
			t.Fatalf("Expected program statements to be 1, received %d", len(program.Statements))
		}

		statement, ok := program.Statements[0].(*ast.ExpressionStatement)
		if !ok {
			t.Fatalf("Expected statement, received %T", program.Statements[0])
		}
		testBooleanLiteral(t, statement.Expression, tt.expectedBoolean)
	}
}

func TestConditionalExpression(t *testing.T) {
	testCases := []struct {
		input             string
		expectedCondition [3]string
		consequence       string
		alternative       string
	}{
		{
			input:             "if (x < y) { x }",
			expectedCondition: [3]string{"x", "<", "y"},
			consequence:       "x",
			alternative:       "",
		},
		{
			input:             "if (x < y) { x } else { y }",
			expectedCondition: [3]string{"x", "<", "y"},
			consequence:       "x",
			alternative:       "y",
		},
	}
	for _, testCase := range testCases {
		lexer := lexer.New(testCase.input)
		parser := New(lexer)
		program := parser.ParseProgram()
		checkParseErrors(t, parser)
		if len(program.Statements) != 1 {
			t.Fatalf("Expected program statements to be %d, received %d", 1, len(program.Statements))
		}
		statement, ok := program.Statements[0].(*ast.ExpressionStatement)
		if !ok {
			t.Fatalf("Expected statement, received %T", program.Statements[0])
		}
		ifExpression, ok := statement.Expression.(*ast.IfExpression)
		if !ok {
			t.Fatalf("Expected if statement, received %T", statement.Expression)
		}
		if !testInfixExpression(t, ifExpression.Condition, testCase.expectedCondition[0], testCase.expectedCondition[1], testCase.expectedCondition[2]) {
			return
		}

		if len(ifExpression.Consequence.Statements) != 1 {
			t.Fatalf("Expected consequence statements to be %d, received %d", 1, len(ifExpression.Consequence.Statements))
		}

		if !testExpressionStatement(t, ifExpression.Consequence.Statements[0]) {
			return
		}

		consequenceStatement, _ := ifExpression.Consequence.Statements[0].(*ast.ExpressionStatement)

		if !testLiteralExpression(t, consequenceStatement.Expression, testCase.consequence) {
			return
		}

		if testCase.alternative != "" {
			if !testExpressionStatement(t, ifExpression.Alternative.Statements[0]) {
				return
			}
			alternativeStatement, _ := ifExpression.Alternative.Statements[0].(*ast.ExpressionStatement)
			if !testLiteralExpression(t, alternativeStatement.Expression, testCase.alternative) {
				return
			}
		} else {
			if ifExpression.Alternative != nil {
				t.Fatalf("Expected empty alternative, received %v", ifExpression.Alternative)
			}
		}

	}
}

func TestFunctionLiteralParsing(t *testing.T) {
	input := "fn(x, y) { x + y; }"
	lexer := lexer.New(input)
	parser := New(lexer)
	program := parser.ParseProgram()
	checkParseErrors(t, parser)
	if len(program.Statements) != 1 {
		t.Fatalf("Expected program statements to be %d, received %d", 1, len(program.Statements))
	}

	statement, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("Expected expression statement, received %T", program.Statements[0])
	}

	functionLiteral, ok := statement.Expression.(*ast.FunctionLiteral)
	if !ok {
		t.Fatalf("Expected function literal, received %T", statement.Expression)
	}
	if len(functionLiteral.Parameters) != 2 {
		t.Fatalf("Expected number of parameters %d, received %d", 2, len(functionLiteral.Parameters))
	}

	testLiteralExpression(t, functionLiteral.Parameters[0], "x")
	testLiteralExpression(t, functionLiteral.Parameters[1], "y")

	if len(functionLiteral.Body.Statements) != 1 {
		t.Fatalf("Expected body statements to be %d, received %d", len(functionLiteral.Body.Statements), 1)
	}
	testExpressionStatement(t, functionLiteral.Body.Statements[0])
	bodyStatement, _ := functionLiteral.Body.Statements[0].(*ast.ExpressionStatement)
	testInfixExpression(t, bodyStatement.Expression, "x", "+", "y")
}

func TestFunctionParametersParsing(t *testing.T) {
	testCases := []struct {
		input          string
		expectedParams []string
	}{
		{
			input:          "fn() {}",
			expectedParams: []string{},
		},
		{
			input:          "fn(x) {}",
			expectedParams: []string{"x"},
		},
		{
			input:          "fn(x, y, z) {}",
			expectedParams: []string{"x", "y", "z"},
		},
	}

	for _, testCase := range testCases {
		lexer := lexer.New(testCase.input)
		parser := New(lexer)
		program := parser.ParseProgram()
		checkParseErrors(t, parser)

		testExpressionStatement(t, program.Statements[0])
		statement := program.Statements[0].(*ast.ExpressionStatement)
		functionLiteral := statement.Expression.(*ast.FunctionLiteral)

		if len(functionLiteral.Parameters) != len(testCase.expectedParams) {
			t.Errorf("Expected %d parameters, received %d", len(testCase.expectedParams), len(functionLiteral.Parameters))
		}

		for i, param := range testCase.expectedParams {
			testLiteralExpression(t, functionLiteral.Parameters[i], param)
		}
	}
}
