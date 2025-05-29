<<<<<<< HEAD
package compiler

import (
	"fmt"
	"io"
=======
package main

import (
	"fmt"
>>>>>>> 9ca2b225d9d60e8d8330e554f31cb9e5b77fb218
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

<<<<<<< HEAD
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

=======
>>>>>>> 9ca2b225d9d60e8d8330e554f31cb9e5b77fb218
}
