package main

// import (
// 	"fmt"
// 	"io"
// 	"os"
// 	"testing"
// )

// func Test_check2(t *testing.T) {
// 	if len(os.Args) > 2 {
// 		fmt.Println("Usage: task lex:run -- <input file>")
// 		os.Exit(1)
// 	}

// 	lexer, err := NewLexer(os.Args[1])
// 	if err != nil {
// 		fmt.Println(err.Error())
// 		os.Exit(1)
// 	}
// 	defer lexer.Close()

// 	var allLines [][]Token  // Slice to store tokens per line
// 	var currentLine []Token // Current line's tokens

// 	for {
// 		token, err := lexer.GetNextToken()
// 		if err != nil {
// 			if err == io.EOF {
// 				if len(currentLine) > 0 {
// 					allLines = append(allLines, currentLine)
// 				}
// 				break
// 			}
// 			fmt.Println(err.Error())
// 			os.Exit(1)
// 		}

// 		currentLine = append(currentLine, token)

// 		if token.Type == TokenNewline { // Adjust this to match your actual newline token type
// 			allLines = append(allLines, currentLine)
// 			currentLine = nil
// 		}
// 	}

// 	// Print result
// 	for i, line := range allLines {
// 		fmt.Printf("Line %d:\n", i+1)
// 		for _, token := range line {
// 			fmt.Println("  " + token.String())
// 		}
// 	}
// }
