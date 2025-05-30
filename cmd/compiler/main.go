package main

import (
	"fmt"
	"io"
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
	defer lexer.Close()

	for {
		token, err := lexer.GetNextToken()
		if err != nil {
			if err == io.EOF {
				break
			}
			fmt.Println(err.Error())
			os.Exit(1)
		}
		fmt.Print(token.String() + "\n")
	}

}
