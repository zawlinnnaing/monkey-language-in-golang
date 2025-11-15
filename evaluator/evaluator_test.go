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
	env := object.NewEnvironment()
	return Eval(program, env)
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

func testBooleanObject(t *testing.T, obj object.Object, expected bool) bool {
	result, ok := obj.(*object.Boolean)
	if !ok {
		t.Errorf("object is not Boolean. got=%T (%+v)", obj, obj)
		return false
	}
	if result.Value != expected {
		t.Errorf("Expected: %t, received %t", expected, result.Value)
		return false
	}
	return true
}

func testNullObject(t *testing.T, obj object.Object) bool {
	if obj != NULL {
		t.Errorf("object is not NULL. got=%T (%+v)", obj, obj)
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
		{"-5", -5},
		{"-10", -10},
		{"5 + 5 + 5 + 5 - 10", 10},
		{"2 * 2 * 2 * 2 * 2", 32},
		{"-50 + 100 + -50", 0},
		{"5 * 2 + 10", 20},
		{"5 + 2 * 10", 25},
		{"20 + 2 * -10", 0},
		{"50 / 2 * 2 + 10", 60},
		{"2 * (5 + 10)", 30},
		{"3 * 3 * 3 + 10", 37},
		{"3 * (3 * 3) + 10", 37},
		{"(5 + 10 * 2 + 15 / 3) * 2 + -10", 50},
	}
	for _, tt := range tests {
		evaluated := testEval(tt.input)
		testIntegerObject(t, evaluated, tt.expected)
	}
}

func TestEvalStringExpression(t *testing.T) {
	testCases := []struct {
		input    string
		expected string
	}{
		{`"Hello World!"`, "Hello World!"},
		{`"hello" + " " + "world"`, "hello world"},
	}
	for _, testCase := range testCases {
		evaluated := testEval(testCase.input)
		str, ok := evaluated.(*object.String)
		if !ok {
			t.Fatalf("object is not String. got=%T (%+v)", evaluated, evaluated)
		}
		if str.Value != testCase.expected {
			t.Errorf("String has wrong value. got=%q, expected=%q", str.Value, testCase.expected)
		}
	}
}

func TestEvalBooleanExpression(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		{"true", true},
		{"false", false},
		{"1 < 2", true},
		{"1 > 2", false},
		{"1 < 1", false},
		{"1 > 1", false},
		{"1 == 1", true},
		{"1 != 1", false},
		{"1 == 2", false},
		{"1 != 2", true},
		{"true == true", true},
		{"false == false", true},
		{"true == false", false},
		{"true != false", true},
		{"false != true", true},
		{"(1 < 2) == true", true},
		{"(1 < 2) == false", false},
		{"(1 > 2) == true", false},
		{"(1 > 2) == false", true},
	}
	for _, tt := range tests {
		evaluated := testEval(tt.input)
		testBooleanObject(t, evaluated, tt.expected)
	}
}

func TestEvalBangOperator(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		{"!true", false},
		{"!false", true},
		{"!!true", true},
		{"!!false", false},
		{"!5", false},
		{"!0", true},
		{"!!5", true},
		{"!!0", false},
		{"!-5", false},
		{"!-0", true},
	}
	for _, testCase := range tests {
		evaluated := testEval(testCase.input)
		result := testBooleanObject(t, evaluated, testCase.expected)
		if result != true {
			t.Errorf("Error at input: %s", testCase.input)
		}
	}
}

func TestEvalIfElseExpression(t *testing.T) {
	testCases := []struct {
		input    string
		expected any
	}{
		{"if (true) { 10 }", 10},
		{"if (false) { 10 }", nil},
		{"if (1) { 10 }", 10},
		{"if (1 < 2) { 10 }", 10},
		{"if (1 > 2) { 10 }", nil},
		{"if (1 > 2) { 10 } else { 20 }", 20},
		{"if (1 < 2) { 10 } else { 20 }", 10},
	}
	for _, testCase := range testCases {
		evaluated := testEval(testCase.input)
		integer, ok := testCase.expected.(int)
		if ok {
			testIntegerObject(t, evaluated, int64(integer))
		} else {
			testNullObject(t, evaluated)
		}
	}
}

func TestReturnStatement(t *testing.T) {
	testCases := []struct {
		input    string
		expected int64
	}{
		{"return 10;", 10},
		{"return 10; 9;", 10},
		{"return 2 * 5; 9;", 10},
		{"9; return 2 * 5; 9;", 10},
		{` if (10 > 1) { if (10 > 1) { return 10; } return 1; }`, 10},
		{"let f = fn() { return 5; }; f();", 5},
		{"let f = fn() { return 5; 10; }; f();", 5},
		{"let f = fn() { 5; return 10; }; f();", 10},
	}
	for _, testCase := range testCases {
		evaluated := testEval(testCase.input)
		testIntegerObject(t, evaluated, testCase.expected)
	}
}

func TestErrorHandling(t *testing.T) {
	testCases := []struct {
		input           string
		expectedMessage string
	}{
		{
			"5 + true;",
			"type mismatch: INTEGER + BOOLEAN",
		},
		{
			"5 + true; 5;",
			"type mismatch: INTEGER + BOOLEAN",
		},
		{
			"-true",
			"unknown operator: -BOOLEAN",
		},
		{
			"true + false;",
			"unknown operator: BOOLEAN + BOOLEAN",
		},
		{
			"5; true + false; 5",
			"unknown operator: BOOLEAN + BOOLEAN",
		},
		{
			"if (10 > 1) { true + false; }",
			"unknown operator: BOOLEAN + BOOLEAN",
		},
		{
			`if (10 > 1) {
				if (10 > 1) {
					return true + false;
				}
				return 1;
			}`,
			"unknown operator: BOOLEAN + BOOLEAN",
		},
		{
			"a",
			"identifier not found: a",
		},
		{
			"let f = fn(x, y) { x + y; }; f();",
			"arguments mismatch. Defined 2, received: 0",
		},
		{
			"let f = fn(x) { x; }; f(1, 2, 3);",
			"arguments mismatch. Defined 1, received: 3",
		},
		{
			`"hello" - "world"`,
			"unknown operator: STRING - STRING",
		},
	}
	for _, testCase := range testCases {
		evaluated := testEval(testCase.input)
		errObj, ok := evaluated.(*object.Error)
		if !ok {
			t.Errorf("no error object returned. got=%T(%+v)", evaluated, evaluated)
			continue
		}
		if errObj.Message != testCase.expectedMessage {
			t.Errorf("Expected error message %s, recevied: %s", testCase.expectedMessage, errObj.Message)
		}
	}
}

func TestLetStatement(t *testing.T) {
	testCases := []struct {
		input    string
		expected any
	}{
		{"let a = 5; a;", 5},
		{"let a = 5 * 5; a;", 25},
		{"let a = 5; let b = a; b;", 5},
		{"let a = 5; let b = a; let c = a + b + 5; c;", 15},
		{"let a = 5;", NULL},
	}
	for _, testCase := range testCases {
		evaluated := testEval(testCase.input)
		switch expected := testCase.expected.(type) {
		case int:
			testIntegerObject(t, evaluated, int64(expected))
		default:
			testNullObject(t, evaluated)
		}

	}
}

func TestFunctionObject(t *testing.T) {
	input := "fn(x) { x + 2; };"
	evaluated := testEval(input)
	fn, ok := evaluated.(*object.Function)
	if !ok {
		t.Fatalf("object is not Function. got=%T (%+v)", evaluated, evaluated)
	}
	if len(fn.Parameters) != 1 {
		t.Fatalf("function has wrong parameters. Parameters=%+v", fn.Parameters)
	}
	if fn.Parameters[0].String() != "x" {
		t.Fatalf("parameter is not 'x'. got=%q", fn.Parameters[0])
	}
	expectedBody := "(x + 2)"
	if fn.Body.String() != expectedBody {
		t.Fatalf("body is not %q. got=%q", expectedBody, fn.Body.String())
	}
}

func TestFunctionApplication(t *testing.T) {
	testCases := []struct {
		input    string
		expected int64
	}{
		{"let identity = fn(x) { x; }; identity(5);", 5},
		{"let identity = fn(x) { return x; }; identity(5);", 5},
		{"let double = fn(x) { x * 2; }; double(5);", 10},
		{"let add = fn(x, y) { x + y; }; add(5, 5);", 10},
		{"let add = fn(x, y) { x + y; }; add(5 + 5, add(5, 5));", 20},
		{"fn(x) { x; }(5)", 5},
	}
	for _, testCase := range testCases {
		evaluated := testEval(testCase.input)
		testIntegerObject(t, evaluated, testCase.expected)
	}
}

func TestClosures(t *testing.T) {
	testCases := []struct {
		input    string
		expected int64
	}{
		{
			`let newAdder = fn(x) {
			  return fn(y) {x + y}
			}
			let addTwo = newAdder(2);
			addTwo(2);  
			`,
			4,
		},
		{
			`let add = fn(x, y) { x + y };
				let applyFunc = fn(f, a, b) { f(a, b) };
				applyFunc(add, 2, 2);`,
			4,
		},
	}
	for _, testCase := range testCases {
		evaluated := testEval(testCase.input)
		testIntegerObject(t, evaluated, testCase.expected)
	}
}

func TestBuiltInFunctions(t *testing.T) {
	testCases := []struct {
		input    string
		expected any
	}{
		{`len("")`, 0},
		{`len("four")`, 4},
		{`len("hello world")`, 11},
		{`len(1)`, "argument to `len` not supported, received INTEGER"},
		{`len("one", "two")`, "wrong number of arguments: received 2, expected 1"},
	}
	for _, testCase := range testCases {
		evaluated := testEval(testCase.input)
		switch expected := testCase.expected.(type) {
		case int:
			testIntegerObject(t, evaluated, int64(expected))
		case string:
			errObj, ok := evaluated.(*object.Error)
			if !ok {
				t.Errorf("object is not Error. got=%T (%+v)", evaluated, evaluated)
				continue
			}
			if errObj.Message != expected {
				t.Errorf("wrong error message. expected=%q, got=%q", expected, errObj.Message)
			}
		}
	}
}
