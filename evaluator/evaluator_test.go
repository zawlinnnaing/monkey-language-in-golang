package evaluator

import (
	"testing"

	"github.com/zawlinnnaing/monkey-language-in-golang/lexer"
	"github.com/zawlinnnaing/monkey-language-in-golang/object"
	"github.com/zawlinnnaing/monkey-language-in-golang/parser"
)

func testEval(input string) object.Object {
	l := lexer.New(input)
	parser := parser.New(l)
	program := parser.ParseProgram()
	return Eval(program)
}

func testIntegerObject(t *testing.T, obj object.Object, expected int64) bool {
	result, ok := obj.(*object.Integer)
	if !ok {
		t.Errorf("object is not Integer. got=%T (%+v)", obj, obj)
		return false
	}
	if result.Value != expected {
		t.Errorf("Expected: %d, received %d", result.Value, expected)
		return false
	}
	return true
}

func TestEvalIntegerExpression(t *testing.T) {
	tests := []struct {
		input    string
		expected int64
	}{
		{"5", 5},
		{"10", 10},
	}
	for _, tt := range tests {
		evaluated := testEval(tt.input)
		testIntegerObject(t, evaluated, tt.expected)
	}
}
