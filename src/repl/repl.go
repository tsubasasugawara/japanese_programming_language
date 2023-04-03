package repl

import (
	"bufio"
	"fmt"
	"io"

	"jpl/lexer"
	"jpl/parser"
	"jpl/evaluator"
	"jpl/object"
	"jpl/ast"
)

const STANDARD_PROMPT = ">> "
const MULTI_LINE_PROMPT = "..."

func printParserErrors(errors []ast.Error) {
	for _, err := range errors {
		fmt.Println(err.Message())
	}
}

func Start(in io.Reader, out io.Writer) {
	scanner := bufio.NewScanner(in)
	env := object.NewEnvironment()

	input := ""
	var prompt = STANDARD_PROMPT

	for {
		fmt.Print(prompt)
		scanned := scanner.Scan()
		if !scanned {
			return
		}

		input = input + scanner.Text() + "\n"

		head := lexer.Tokenize(input)
		program, errors := parser.Parse(head)
		if len(errors) > 0 {
			errMessages := ""
			stmtNotCompleted := false //括弧が閉じているかどうか
			for _, err := range errors {
				e, ok := err.(*ast.SyntaxError)
				if ok && (e.Category() == ast.MISSING_RBRACE || e.Category() == ast.MISSING_RPAREN || e.Category() == ast.MISSING_R_SQUARE_BRACE) {
					prompt = MULTI_LINE_PROMPT
					stmtNotCompleted = true
				}
				errMessages += err.Message() + "\n"
			}

			// 括弧が閉じていなければ続ける
			if stmtNotCompleted {
				continue
			} else {
				fmt.Print(errMessages)
				input = ""
				prompt = STANDARD_PROMPT
				continue
			}
		}

		for _, v := range program.Nodes {
			o := evaluator.Eval(v, env)
			if o.Type() != object.NULL {
				fmt.Println(o.Inspect())
			}
		}

		input = ""
		prompt = STANDARD_PROMPT
	}
}
