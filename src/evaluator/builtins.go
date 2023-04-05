package evaluator

import (
	"fmt"

	"jpl/object"
)

var builtins = map[string]*object.Builtin {
	"表示": &object.Builtin {
		Fn: func(args ...object.Object) object.Object {
			for _, arg := range args {
				fmt.Println(arg.Inspect())
			}

			return NULL
		},
	},
	"配列": &object.Builtin {
		Fn: func(args ...object.Object) object.Object {
			var defVal object.Object

			if len(args) == 1 {
				defVal = NULL
			} else if len(args) == 2 {
				defVal = args[1]
			} else {
				return newError("引数の個数が正しくありません。")
			}

			length, ok := args[0].(*object.Integer)
			if !ok {
				return newError("数値が必要です")
			}

			elements := make([]object.Object, length.Value)
			for i, _ := range elements {
				elements[i] = defVal
			}
			return &object.Array{Elements: elements}
		},
	},
	"追加": &object.Builtin {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 2 {
				return newError("引数の個数が正しくありません。")
			}

			array, ok := args[0].(*object.Array)
			if !ok {
				return newError("配列を指定してください。")
			}

			array.Elements = append(array.Elements, args[1])
			return NULL
		},
	},
	"削除": &object.Builtin {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 2 {
				return newError("引数の個数が正しくありません。")
			}

			array, ok := args[0].(*object.Array)
			if !ok {
				return newError("配列を指定してください。")
			}

			index, ok := args[1].(*object.Integer)
			if !ok {
				return newError("数値が必要です。")
			}

			length := len(array.Elements)
			if length < 0 || int64(length) <= index.Value {
				return newError("範囲外です。")
			} else if index.Value == 0 {
				array.Elements = array.Elements[1:]
			} else if index.Value == int64(length - 1) {
				array.Elements = array.Elements[:index.Value]
			} else {
				array.Elements = append(array.Elements[:index.Value], array.Elements[index.Value + 1:]...)
			}

			return 	NULL
		},
	},
}