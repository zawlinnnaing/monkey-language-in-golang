package evaluator

import (
	"github.com/zawlinnnaing/monkey-language-in-golang/ast"
	"github.com/zawlinnnaing/monkey-language-in-golang/object"
)

var (
	TRUE  = &object.Boolean{Value: true}
	FALSE = &object.Boolean{Value: false}
	NULL  = &object.Null{}
)

func Eval(node ast.Node) object.Object {
	switch n := node.(type) {
	case *ast.Program:
		return evalProgram(n)
	case *ast.ExpressionStatement:
		return Eval(n.Expression)
	case *ast.IntegerLiteral:
		return &object.Integer{Value: n.Value}
	case *ast.BooleanLiteral:
		return evalBooleanLiteral(n)
	case *ast.PrefixExpression:
		right := Eval(n.Right)
		return evalPrefixExpression(n.Operator, right)
	case *ast.InfixExpression:
		left := Eval(n.Left)
		right := Eval(n.Right)
		return evalInfixExpression(n.Operator, left, right)
	case *ast.IfExpression:
		return evalIfExpression(n)
	case *ast.BlockStatement:
		return evalBlockStatement(n)
	case *ast.ReturnStatement:
		val := Eval(n.ReturnValue)
		return &object.ReturnValue{Value: val}
	}
	return nil
}

func evalIfExpression(node *ast.IfExpression) object.Object {
	condition := Eval(node.Condition)
	if isTruthy(condition) {
		return Eval(node.Consequence)
	} else if node.Alternative != nil {
		return Eval(node.Alternative)
	}
	return NULL
}

func evalInfixExpression(operator string, left, right object.Object) object.Object {
	switch {
	case left.Type() == object.INTEGER_OBJ && right.Type() == object.INTEGER_OBJ:
		return evalIntegerInfixExpression(operator, left, right)
	case operator == "==":
		return object.NewBoolean(left == right)
	case operator == "!=":
		return object.NewBoolean(left != right)
	default:
		// Temporarily return NULL
		return NULL
	}

}

func evalIntegerInfixExpression(operator string, left, right object.Object) object.Object {
	leftVal := left.(*object.Integer).Value
	rightVal := right.(*object.Integer).Value
	switch operator {
	case "+":
		return &object.Integer{Value: leftVal + rightVal}
	case "-":
		return &object.Integer{Value: leftVal - rightVal}
	case "*":
		return &object.Integer{Value: leftVal * rightVal}
	case "/":
		return &object.Integer{Value: leftVal / rightVal}
	case "<":
		if leftVal < rightVal {
			return TRUE
		}
		return FALSE
	case ">":
		if leftVal > rightVal {
			return TRUE
		}
		return FALSE
	case "==":
		if leftVal == rightVal {
			return TRUE
		}
		return FALSE
	case "!=":
		if leftVal != rightVal {
			return TRUE
		}
		return FALSE
	case ">=":
		if leftVal >= rightVal {
			return TRUE
		}
		return FALSE
	case "<=":
		if leftVal <= rightVal {
			return TRUE
		}
		return FALSE
	default:
		// Temporarily return NULL for unsupported operators
		return NULL
	}
}

func evalBooleanLiteral(node *ast.BooleanLiteral) object.Object {
	if node.Value {
		return TRUE
	}
	return FALSE
}

func evalProgram(program *ast.Program) object.Object {
	var result object.Object
	for _, statement := range program.Statements {
		result = Eval(statement)
		if returnObj, ok := result.(*object.ReturnValue); ok {
			return returnObj.Value
		}
	}
	return result
}

func evalBlockStatement(block *ast.BlockStatement) object.Object {
	var result object.Object
	for _, statement := range block.Statements {
		result = Eval(statement)
		if result != nil && result.Type() == object.RETURN_VALUE_OBJ {
			return result
		}
	}
	return result
}

func evalPrefixExpression(operator string, right object.Object) object.Object {
	switch operator {
	case "!":
		return evalBangOperatorExpression(right)
	case "-":
		return evalMinusPrefixOperatorExpression(right)
	default:
		// Temporarily return NULL for unsupported prefix operators
		return NULL
	}
}

func evalMinusPrefixOperatorExpression(right object.Object) object.Object {
	integer, ok := right.(*object.Integer)
	if !ok {
		// Temporarily return NULL for non-integer negation
		return NULL
	}
	return &object.Integer{Value: -integer.Value}
}

func evalBangOperatorExpression(right object.Object) object.Object {
	integer, ok := right.(*object.Integer)
	if ok {
		if integer.Value == 0 {
			return TRUE
		}
		return FALSE
	}
	switch right {
	case TRUE:
		return FALSE
	case FALSE:
		return TRUE
	case NULL:
		return TRUE
	default:
		return FALSE
	}
}

func isTruthy(obj object.Object) bool {
	switch obj {
	case TRUE:
		return true
	case NULL, FALSE:
		return false
	default:
		return true
	}
}
