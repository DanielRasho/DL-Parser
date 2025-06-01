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

	parser, err := NewParser(os.Args[1])
	if err != nil {
		fmt.Println(err.Error())
	}

	slicetokens := []Token{}
	for {
		token, err := lexer.GetNextToken()
		if err != nil {
			if err == io.EOF {
				break
			}
			fmt.Println(err.Error())
			os.Exit(1)
		}

		_, ok := parser.parsedefinition.IgnoredSymbols[token.TokenID]

		if !ok {
			slicetokens = append(slicetokens, token)
			fmt.Print(token.String() + "\n")
		}

		if token.Value == "\n" {
			slicetokens = *parser.ParseInput(slicetokens, parser.parsedefinition.Terminals, *parser.parsedefinition)
		}

	}

	if len(slicetokens) >= 0 {
		slicetokens = *parser.ParseInput(slicetokens, parser.parsedefinition.Terminals, *parser.parsedefinition)
	}

	if len(slicetokens) == 0 {
		fmt.Println("ALL LINES ARE ACCEPTED")
	} else {
		fmt.Printf("\nERROR PARSING FROM %d, to %d", slicetokens[0].Offset, slicetokens[len(slicetokens)-1].Offset)
		fmt.Println(slicetokens)
	}

}
