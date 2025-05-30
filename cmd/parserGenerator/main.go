package main

import (
	"flag"
	"fmt"
	"os"

	parser "github.com/DanielRasho/Parser/internal/Parser/Generator"
)

func main() {
	// Define the flags
	fileFlag := flag.String("f", "", "Yalex file path")
	outputFlag := flag.String("o", "", "Output file path")
	template := flag.String("t", "", "Template Parser")
	diagramFlag := flag.Bool("diagram", true, "Render automata diagrams")

	// Parse the command line flags
	flag.Parse()

	// Check if both flags are provided, if not print usage
	if *fileFlag == "" || *outputFlag == "" || *template == "" {
		fmt.Println("Usage: task lex:generate -- -f <input-file> -o <output-file> -t <template-parser>")
		os.Exit(1)
	}

	// Print the values of the flags (just as an example)
	fmt.Printf("Input file: %s\n", *fileFlag)
	fmt.Printf("Output file: %s\n", *outputFlag)
	fmt.Printf("Render diagramas: %t\n", *diagramFlag)

	// CODE FOR GENERATING LPARSER ...
	err := parser.Compile(*fileFlag, *template, *outputFlag, true)
	if err != nil {
		fmt.Println(err)
	}
}
