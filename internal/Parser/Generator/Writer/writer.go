package writer

import (
	"fmt"
	"os"
	"text/template"

	reader "github.com/DanielRasho/Parser/internal/Parser/Generator/Reader"
	table "github.com/DanielRasho/Parser/internal/Parser/TransitionTable"
	"github.com/DanielRasho/Parser/internal/Parser/automata"
)

// Writes a parser.go file in the desired location.
// Possible errors:
//   - file paths invalids/not found
//   - invalid parsing table.
//
// REMINDER!!!!! DONT LOAD THE ENTIRE FILE ON A STRING, use buffers instead.
func WriteParserFile(templateFilePath string, outputFilePath string) error {
	parserdef, _ := reader.Parse("../../../../examples/productions2.y")
	first := table.GetFirst(parserdef)
	follow := table.GetFollow(parserdef, first)
	automa := automata.NewAutomata(parserdef, false)

	transitionTbl, gotoTbl, _ := table.NewTable(automa, first, follow, *parserdef)

	// Load and parse the template
	fmt.Println("PRINTING")
	tmpl, err := template.ParseFiles(templateFilePath)
	if err != nil {
		return fmt.Errorf("failed to parse template: %w", err)
	}

	// Create the data context
	data := templateLexwrite{
		ParserDefinition: *parserdef,
		TransitTable:     *transitionTbl,
		Gotable:          *gotoTbl,
	}

	// Open output file
	outFile, err := os.Create(outputFilePath)
	if err != nil {
		return fmt.Errorf("failed to create output file: %w", err)
	}
	defer outFile.Close()

	// Execute the template
	err = tmpl.Execute(outFile, data)
	if err != nil {
		return fmt.Errorf("failed to execute template: %w", err)
	}

	return nil
}
