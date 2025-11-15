package evaluator

import "github.com/zawlinnnaing/monkey-language-in-golang/object"

var builtInEnvironment = map[string]*object.BuiltIn{
	"len": &object.BuiltIn{
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return object.NewError("wrong number of arguments: received %d, expected %d", len(args), 1)
			}
			switch arg := args[0].(type) {
			case *object.String:
				{
					return &object.Integer{Value: int64(len(arg.Value))}
				}
			default:
				return object.NewError("argument to `len` not supported, received %s", arg.Type())
			}
		},
	},
}
