package evaluator

import (
	"errors"
	"fmt"
	"strings"
	"syscall"

	"github.com/d2verb/monkey/object"
)

func parseMode(mode string) (int, error) {
	m, o := 0, 0

	if len(mode) < 1 {
		return 0, errors.New("invalid mode")
	}

	switch c := mode[0]; string(c) {
	case "r":
		m = syscall.O_RDONLY
		o = 0
	case "w":
		m = syscall.O_WRONLY
		o = syscall.O_CREAT | syscall.O_TRUNC
	case "a":
		m = syscall.O_WRONLY
		o = syscall.O_CREAT | syscall.O_APPEND
	default:
		return 0, errors.New("invalid mode")
	}

	if mode[1:] == "+" || mode[1:] == "b+" {
		m = syscall.O_RDWR
	}

	return m | o, nil
}

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
				return newError("1st argument to `push` must be ARRAY, got %s",
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

	"open": &object.Builtin{
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 2 {
				return newError("wrong number of arguments. got=%d, want=2",
					len(args))
			}

			if args[0].Type() != object.STRING_OBJ {
				return newError("1st argument to `open` must be STRING, got %s",
					args[0].Type())
			}

			if args[1].Type() != object.STRING_OBJ {
				return newError("2nd argument to `open` must be STRING, got %s",
					args[1].Type())
			}

			path := args[0].(*object.String).Value
			mode := args[1].(*object.String).Value

			flag, err := parseMode(mode)
			if err != nil {
				return newError(err.Error())
			}

			fd, err := syscall.Open(path, flag, 0666)
			if err != nil {
				return newError(err.Error())
			}

			return &object.Integer{Value: int64(fd)}
		},
	},

	"close": &object.Builtin{
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return newError("wrong number of arguments. got=%d, want=1",
					len(args))
			}

			if args[0].Type() != object.INTEGER_OBJ {
				return newError("1st argument to `close` must be INTEGER, got %s",
					args[0].Type())
			}

			fd := args[0].(*object.Integer).Value
			err := syscall.Close(int(fd))
			if err != nil {
				return newError(err.Error())
			}

			return NULL
		},
	},

	"read": &object.Builtin{
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 2 {
				return newError("wrong number of arguments. got=%d, want=2",
					len(args))
			}

			if args[0].Type() != object.INTEGER_OBJ {
				return newError("1st argument to `read` must be INTEGER, got %s",
					args[0].Type())
			}

			if args[1].Type() != object.INTEGER_OBJ {
				return newError("2nd argument to `read` must be INTEGER, got %s",
					args[0].Type())
			}

			fd := args[0].(*object.Integer).Value
			nb := args[1].(*object.Integer).Value

			p := make([]byte, nb)
			n, err := syscall.Read(int(fd), p)
			if err != nil {
				return newError(err.Error())
			}

			result := &object.Hash{Pairs: make(map[object.HashKey]object.HashPair)}

			key1 := &object.String{Value: "content"}
			result.Pairs[key1.HashKey()] = object.HashPair{Key: key1, Value: &object.String{Value: string(p)}}

			key2 := &object.String{Value: "nb"}
			result.Pairs[key2.HashKey()] = object.HashPair{Key: key2, Value: &object.Integer{Value: int64(n)}}

			return result
		},
	},
}
