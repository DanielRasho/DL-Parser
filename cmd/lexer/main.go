package main

import (
	"fmt"
	"os"
)

func main() {

	if len(os.Args) > 2 {
		fmt.Println("Usage: task lex:run -- <input file>")
		os.Exit(1)
	}

	lexer, err := NewLexer(os.Args[1])
	if err != nil {
		fmt.Println(err.Error())
	}

	for i := 0; i < 50; i++ {
		token, err := lexer.GetNextToken()
		if err != nil {
			fmt.Println(err.Error())
			os.Exit(1)
		}
		fmt.Print(token.String() + "\n")
	}
}
