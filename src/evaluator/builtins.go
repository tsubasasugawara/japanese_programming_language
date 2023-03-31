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
}