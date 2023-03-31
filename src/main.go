package main

import (
	"fmt"
	"os"
	"os/user"

	"jpl/repl"
	"jpl/reader"
)

func main() {
	if len(os.Args) == 1 {
		user, err := user.Current()
		if err != nil {
			panic(err)
		}
		
		fmt.Printf("Hello %s!\n", user.Username)
		repl.Start(os.Stdin, os.Stdout)
	} else {
		path := os.Args[1]
		reader.Read(path)
	}
}
