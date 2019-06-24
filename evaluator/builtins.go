package evaluator

import (
	"fmt"
	"strings"

	"github.com/d2verb/monkey/object"
)

var builtins = map[string]*object.Builtin{
	"len": &object.Builtin{
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return newError("wrong number of arguments. got=%d, want=1", len(args))
			}

			switch arg := args[0].(type) {
			case *object.Array:
				return &object.Integer{Value: int64(len(arg.Elements))}
			case *object.String:
				return &object.Integer{Value: int64(len(arg.Value))}
			default:
				return newError("argument to `len` not supported, got %s", args[0].Type())
			}
		},
	},

	"print": &object.Builtin{
		Fn: func(args ...object.Object) object.Object {
			outputs := []string{}
			for _, arg := range args {
				outputs = append(outputs, arg.Inspect())
			}
			fmt.Print(strings.Join(outputs, ""))
			return NULL
		},
	},

	"println": &object.Builtin{
		Fn: func(args ...object.Object) object.Object {
			outputs := []string{}
			for _, arg := range args {
				outputs = append(outputs, arg.Inspect())
			}
			fmt.Println(strings.Join(outputs, " "))
			return NULL
		},
	},

	"push": &object.Builtin{
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 2 {
				return newError("wrong number of arguments. got=%d, want=2",
					len(args))
			}
			if args[0].Type() != object.ARRAY_OBJ {
				return newError("argument to `first` must be ARRAY, got %s",
					args[0].Type())
			}
			arr := args[0].(*object.Array)
			arr.Elements = append(arr.Elements, args[1])
			return NULL
		},
	},

	"copy": &object.Builtin{
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return newError("wrong number of arguments. got=%d, want=1",
					len(args))
			}
			switch t := args[0].Type(); t {
			case object.ARRAY_OBJ:
				arr := args[0].(*object.Array)
				newElements := make([]object.Object, len(arr.Elements))
				copy(newElements, arr.Elements)
				return &object.Array{Elements: newElements}
			case object.HASH_OBJ:
				arr := args[0].(*object.Hash)
				newPairs := make(map[object.HashKey]object.HashPair)
				for key, value := range arr.Pairs {
					newPairs[key] = value
				}
				return &object.Hash{Pairs: newPairs}
			default:
				return newError("copy operation not supported: %s", t)
			}
		},
	},
}
