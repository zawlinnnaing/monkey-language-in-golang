package evaluator

import "github.com/zawlinnnaing/monkey-language-in-golang/object"

var builtInEnvironment = map[string]*object.BuiltIn{
	"len": {
		Fn: lenBuiltIn,
	},
	"first": {
		Fn: firstBuiltIn,
	},
	"last": {
		Fn: lastBuiltIn,
	},
}

var lenBuiltIn object.BuiltInFunction = func(args ...object.Object) object.Object {
	if len(args) != 1 {
		return object.NewError("wrong number of arguments: received %d, expected %d", len(args), 1)
	}
	switch arg := args[0].(type) {
	case *object.String:
		{
			return &object.Integer{Value: int64(len(arg.Value))}
		}
	case *object.Array:
		{
			return &object.Integer{Value: int64(len(arg.Elements))}
		}
	default:
		return object.NewError("argument to `len` not supported, received %s", arg.Type())
	}
}

var firstBuiltIn object.BuiltInFunction = func(args ...object.Object) object.Object {
	if len(args) != 1 {
		return object.NewError("wrong number of arguments: received %d, expected %d", len(args), 1)
	}
	if args[0].Type() != object.ARRAY_OBJ {
		return object.NewError("argument to `first` must be ARRAY, received %s", args[0].Type())
	}
	arrayObj := args[0].(*object.Array)
	if len(arrayObj.Elements) == 0 {
		return NULL
	}
	return arrayObj.Elements[0]
}

var lastBuiltIn object.BuiltInFunction = func(args ...object.Object) object.Object {
	if len(args) != 1 {
		return object.NewError("wrong number of arguments: received %d, expected %d", len(args), 1)
	}
	if args[0].Type() != object.ARRAY_OBJ {
		return object.NewError("argument to `last` must be ARRAY, received %s", args[0].Type())
	}
	arrayObj := args[0].(*object.Array)
	if len(arrayObj.Elements) == 0 {
		return NULL
	}
	return arrayObj.Elements[len(arrayObj.Elements)-1]
}
