package repl

import (
	"bufio"
	"fmt"
	"io"

	"jpl/token"
	"jpl/parser"
	"jpl/evaluator"
	"jpl/object"
)

const PROMPT = ">> "

func printParserErrors(errors []string) {
	for _, err := range errors {
		fmt.Println(err)
	}
}

func Start(in io.Reader, out io.Writer) {
	scanner := bufio.NewScanner(in)

	env := object.NewEnvironment()

	for {
		fmt.Print(PROMPT)
		scanned := scanner.Scan()
		if !scanned {
			return
		}

		line := scanner.Text()

		head := token.Tokenize(line)
		program, errors := parser.Parse(head)
		if len(errors) > 0 {
			printParserErrors(errors)
			continue
		}
		for _, v := range program.Nodes {
			o := evaluator.Eval(v, env)
			if o.Type() != object.NULL {
				fmt.Println(o.Inspect())
			}
		}
	}
}
