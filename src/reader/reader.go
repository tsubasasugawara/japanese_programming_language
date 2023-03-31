package reader

import (
	"fmt"
	"io/ioutil"

	"jpl/token"
	"jpl/parser"
	"jpl/evaluator"
	"jpl/object"
	"jpl/ast"
)

func printParserErrors(errors []ast.Error) {
	for _, err := range errors {
		fmt.Println(err.Message())
	}
}

func Read(path string) {
	env := object.NewEnvironment()
	content, err := ioutil.ReadFile(path)
	if err != nil {
		panic(err)
	}

	head := token.Tokenize(string(content))
	program, errors := parser.Parse(head)
	if len(errors) > 0 {
		printParserErrors(errors)
		return
	}
			
	for _, v := range program.Nodes {
		evaluator.Eval(v, env)
	}
}