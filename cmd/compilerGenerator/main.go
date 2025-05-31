package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"

	lex "github.com/DanielRasho/Parser/internal/Lexer/Generator"
	parser "github.com/DanielRasho/Parser/internal/Parser/Generator"
)

func main() {
	// Define the flags
	yalexFile := flag.String("l", "", "Yapar file")
	yaparFile := flag.String("p", "", "Parser file path")
	outputFlag := flag.String("d", "", "Output file path")
	verbose := flag.Bool("verbose", true, "Render automata diagrams")

	// Parse the command line flags
	flag.Parse()

	// Check if both flags are provided, if not print usage
	if *yalexFile == "" || *yaparFile == "" || *outputFlag == "" {
		fmt.Println("Usage: task compiler:generate -- -l <yalex-file> -p <yapar-file> -d <output-dir> -t <template-file>")
		os.Exit(1)
	}

	// Print the values of the flags (just as an example)
	fmt.Printf("Yalex file: %s\n", *yalexFile)
	fmt.Printf("Yapar file: %s\n", *yaparFile)
	fmt.Printf("Output folder: %s\n", *outputFlag)
	fmt.Printf("Verbose: %t\n", *verbose)

	lexerFile := filepath.Join(*outputFlag, "lexer.go")
	parserFile := filepath.Join(*outputFlag, "parser.go")

	// CODE FOR GENERATING LEXER ...
	err := lex.Compile(*yalexFile, lexerFile, *verbose, *verbose)
	if err != nil {
		fmt.Println(err)
	}

	// CODE FOR GENERATING PARSER ...
	err = parser.Compile(*yaparFile, "./template/ParserTemplate.go", parserFile, *verbose)
	if err != nil {
		fmt.Println(err)
	}
}
