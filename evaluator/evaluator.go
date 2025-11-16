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

func Eval(node ast.Node, env *object.Environment) object.Object {
	switch n := node.(type) {
	case *ast.Program:
		return evalProgram(n, env)
	case *ast.Identifier:
		return evalIdentifier(n, env)
	case *ast.ExpressionStatement:
		return Eval(n.Expression, env)
	case *ast.IntegerLiteral:
		return &object.Integer{Value: n.Value}
	case *ast.BooleanLiteral:
		return evalBooleanLiteral(n)
	case *ast.StringLiteral:
		return &object.String{Value: n.Value}
	case *ast.LetStatement:
		return evalLetStatement(n, env)
	case *ast.PrefixExpression:
		right := Eval(n.Right, env)
		if isError(right) {
			return right
		}
		return evalPrefixExpression(n.Operator, right)
	case *ast.InfixExpression:
		left := Eval(n.Left, env)
		right := Eval(n.Right, env)
		if isError(left) {
			return left
		}
		if isError(right) {
			return right
		}
		return evalInfixExpression(n.Operator, left, right)
	case *ast.IfExpression:
		return evalIfExpression(n, env)
	case *ast.BlockStatement:
		return evalBlockStatement(n, env)
	case *ast.FunctionLiteral:
		return evalFunctionLiteral(n, env)
	case *ast.CallExpression:
		return evalCallExpression(n, env)
	case *ast.ReturnStatement:
		val := Eval(n.ReturnValue, env)
		if isError(val) {
			return val
		}
		return &object.ReturnValue{Value: val}
	case *ast.ArrayLiteral:
		elements := evalExpressions(n.Elements, env)
		if len(elements) == 1 && isError(elements[0]) {
			return elements[0]
		}
		return &object.Array{Elements: elements}
	case *ast.IndexExpression:
		return evalIndexExpression(n, env)
	}
	return nil
}

func evalIndexExpression(node *ast.IndexExpression, env *object.Environment) object.Object {
	left := Eval(node.Left, env)
	if isError(left) {
		return left
	}
	right := Eval(node.Index, env)
	if isError(right) {
		return right
	}
	switch {
	case left.Type() == object.ARRAY_OBJ && right.Type() == object.INTEGER_OBJ:
		return evalArrayIndexExpression(left, right)
	default:
		return object.NewError("index operator not supported: %s", right.Type())
	}

}

func evalArrayIndexExpression(array, index object.Object) object.Object {
	arrayObj := array.(*object.Array)
	idx := index.(*object.Integer).Value
	maxIdx := len(arrayObj.Elements) - 1
	if idx < 0 || idx > int64(maxIdx) {
		return NULL
	}
	return arrayObj.Elements[idx]
}

func evalIdentifier(node *ast.Identifier, env *object.Environment) object.Object {
	val, ok := env.Get(node.Value)
	if ok {
		return val
	}
	builtInObj, ok := builtInEnvironment[node.Value]
	if ok {
		return builtInObj
	}
	return object.NewError("identifier not found: %s", node.Value)
}

func evalLetStatement(node *ast.LetStatement, env *object.Environment) object.Object {
	val := Eval(node.Value, env)
	if isError(val) {
		return val
	}
	env.Set(node.Name.Value, val)
	return NULL
}

func evalIfExpression(node *ast.IfExpression, env *object.Environment) object.Object {
	condition := Eval(node.Condition, env)
	if isError(condition) {
		return condition
	}
	if isTruthy(condition) {
		return Eval(node.Consequence, env)
	} else if node.Alternative != nil {
		return Eval(node.Alternative, env)
	}
	return NULL
}

func evalInfixExpression(operator string, left, right object.Object) object.Object {
	switch {
	case left.Type() == object.INTEGER_OBJ && right.Type() == object.INTEGER_OBJ:
		return evalIntegerInfixExpression(operator, left, right)
	case left.Type() == object.STRING_OBJ && right.Type() == object.STRING_OBJ:
		return evalStringInfixExpression(operator, left, right)
	case operator == "==":
		return object.NewBoolean(left == right)
	case operator == "!=":
		return object.NewBoolean(left != right)
	case left.Type() != right.Type():
		return object.NewError("type mismatch: %s %s %s", left.Type(), operator, right.Type())
	default:
		return object.NewError("unknown operator: %s %s %s", left.Type(), operator, right.Type())
	}

}

func evalStringInfixExpression(operator string, left, right object.Object) object.Object {
	// TODO: add support for string comparison
	if operator != "+" {
		return object.NewError("unknown operator: %s %s %s", left.Type(), operator, right.Type())
	}
	leftVal := left.(*object.String).Value
	rightVal := right.(*object.String).Value
	return &object.String{Value: leftVal + rightVal}
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
		return object.NewError("unknown operator: %s %s %s", left.Type(), operator, right.Type())
	}
}

func evalBooleanLiteral(node *ast.BooleanLiteral) object.Object {
	if node.Value {
		return TRUE
	}
	return FALSE
}

func evalProgram(program *ast.Program, env *object.Environment) object.Object {
	var result object.Object
	for _, statement := range program.Statements {
		result = Eval(statement, env)
		switch result := result.(type) {
		case *object.ReturnValue:
			return result.Value
		case *object.Error:
			return result
		}
	}
	return result
}

func evalBlockStatement(block *ast.BlockStatement, env *object.Environment) object.Object {
	var result object.Object
	for _, statement := range block.Statements {
		result = Eval(statement, env)
		if result != nil {
			resultType := result.Type()
			if resultType == object.RETURN_VALUE_OBJ || resultType == object.ERROR_OBJ {
				return result
			}
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
		return object.NewError("unknown operator: %s%s", operator, right.Type())
	}
}

func evalMinusPrefixOperatorExpression(right object.Object) object.Object {
	integer, ok := right.(*object.Integer)
	if !ok {
		return object.NewError("unknown operator: -%s", right.Type())
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

func evalFunctionLiteral(node *ast.FunctionLiteral, env *object.Environment) *object.Function {
	return &object.Function{
		Parameters: node.Parameters,
		Body:       node.Body,
		Env:        env,
	}
}

func evalCallExpression(node *ast.CallExpression, env *object.Environment) object.Object {
	evaluated := Eval(node.Function, env)
	if isError(evaluated) {
		return evaluated
	}

	evaluatedArgs := evalExpressions(node.Arguments, env)
	if len(evaluatedArgs) == 1 && isError(evaluatedArgs[0]) {
		return evaluatedArgs[0]
	}

	switch function := evaluated.(type) {
	case *object.Function:
		{
			argErr := validateFunctionArguments(function, evaluatedArgs)
			if argErr != nil {
				return argErr
			}
			return applyFunction(function, evaluatedArgs)
		}
	case *object.BuiltIn:
		{
			return function.Fn(evaluatedArgs...)
		}
	default:
		return object.NewError("not a function: %s", evaluated.Type())
	}
}

func validateFunctionArguments(fn *object.Function, args []object.Object) *object.Error {
	// TODO: add support for optional parameters later
	if len(args) != len(fn.Parameters) {
		return object.NewError("arguments mismatch. Defined %d, received: %d", len(fn.Parameters), len(args))
	}
	return nil
}

func applyFunction(function *object.Function, args []object.Object) object.Object {
	extendedEnv := extendEnv(function, args)
	evaluated := Eval(function.Body, extendedEnv)
	return unwrappedReturnValue(evaluated)
}

func extendEnv(fn *object.Function, args []object.Object) *object.Environment {
	env := object.NewEnclosedEnvironment(fn.Env)
	for idx, param := range fn.Parameters {
		env.Set(param.Value, args[idx])
	}
	return env
}

func unwrappedReturnValue(obj object.Object) object.Object {
	returnObj, ok := obj.(*object.ReturnValue)
	if ok {
		return returnObj.Value
	}
	return obj
}

func evalExpressions(expressions []ast.Expression, env *object.Environment) []object.Object {
	var results []object.Object
	for _, expression := range expressions {
		evaluated := Eval(expression, env)
		if isError(evaluated) {
			return []object.Object{evaluated}
		}
		results = append(results, evaluated)
	}
	return results
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

func isError(obj object.Object) bool {
	_, ok := obj.(*object.Error)
	return ok
}
