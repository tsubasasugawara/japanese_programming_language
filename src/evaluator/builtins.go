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
				return newError("引数の個数が合っていません。")
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
}