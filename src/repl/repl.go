package repl

import (
	"bufio"
	"fmt"
	"io"
	"jpl/token"
)

const PROMPT = ">> "

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

		for {
			fmt.Printf("%+v\n", token)

			if token.Next == nil {
				break
			} else {
				token = token.Next
			}
		}
	}
}
