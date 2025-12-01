package evaluator

import (
	"fmt"

	"github.com/zawlinnnaing/monkey-language-in-golang/object"
)

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
	"rest": {
		Fn: restBuiltIn,
	},
	"push": {
		Fn: pushBuiltIn,
	},
	"print": {
		Fn: printBuiltIn,
	},
}

var lenBuiltIn object.BuiltInFunction = func(args ...object.Object) object.Object {
	if err := validateArgsLen(1, args...); err != nil {
		return err
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
	err := validateArgsLen(1, args...)
	if err != nil {
		return err
	}
	err = validateArrayArgs("first", args...)
	if err != nil {
		return err
	}
	arrayObj := args[0].(*object.Array)
	if len(arrayObj.Elements) == 0 {
		return NULL
	}
	return arrayObj.Elements[0]
}

var lastBuiltIn object.BuiltInFunction = func(args ...object.Object) object.Object {
	err := validateArgsLen(1, args...)
	if err != nil {
		return err
	}
	err = validateArrayArgs("last", args...)
	if err != nil {
		return err
	}
	arrayObj := args[0].(*object.Array)
	if len(arrayObj.Elements) == 0 {
		return NULL
	}
	return arrayObj.Elements[len(arrayObj.Elements)-1]
}

var restBuiltIn object.BuiltInFunction = func(args ...object.Object) object.Object {
	err := validateArgsLen(1, args...)
	if err != nil {
		return err
	}
	err = validateArrayArgs("rest", args...)
	if err != nil {
		return err
	}
	arr := args[0].(*object.Array)
	if len(arr.Elements) == 0 {
		return NULL
	}
	newArray := &object.Array{}
	copy(arr.Elements[1:len(arr.Elements)], newArray.Elements)
	return newArray
}

var pushBuiltIn object.BuiltInFunction = func(args ...object.Object) object.Object {
	err := validateArgsLen(2, args...)
	if err != nil {
		return err
	}
	err = validateArrayArgs("push", args...)
	if err != nil {
		return err
	}
	array := args[0].(*object.Array)
	newArray := &object.Array{}
	arrLen := len(array.Elements)
	newArray.Elements = make([]object.Object, arrLen+1)
	copy(array.Elements, newArray.Elements)
	newArray.Elements[arrLen] = args[1]
	return newArray
}

var printBuiltIn object.BuiltInFunction = func(args ...object.Object) object.Object {
	for _, arg := range args {
		fmt.Println(arg.Inspect())
	}
	return NULL
}

func validateArrayArgs(fnName string, args ...object.Object) object.Object {
	if args[0].Type() != object.ARRAY_OBJ {
		return object.NewError("argument to `%s` must be ARRAY, received %s", fnName, args[0].Type())
	}
	return nil
}

func validateArgsLen(expectedLen int, args ...object.Object) object.Object {
	if len(args) != expectedLen {
		return object.NewError("wrong number of arguments: received %d, expected %d", len(args), expectedLen)
	}
	return nil
}
