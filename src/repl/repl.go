package repl

import (
	"bufio"
	"fmt"
	"io"

	"jpl/token"
	"jpl/parser"
)

const PROMPT = ">> "

func printParserErrors(errors []string) {
	for _, err := range errors {
		fmt.Println(err)
	}
}

func Start(in io.Reader, out io.Writer) {
	scanner := bufio.NewScanner(in)

	for {
		fmt.Print(PROMPT)
		scanned := scanner.Scan()
		if !scanned {
			return
		}

		line := scanner.Text()

		token := token.Tokenize(line)

		program, errors := parser.Parse(token)
		if len(errors) > 0 {
			printParserErrors(errors)
			continue
		}

		for {
			fmt.Printf("%+v\n", token)

			if token.Next == nil {
				break
			} else {
				token = token.Next
			}
		}

		fmt.Println(program)
	}
}
