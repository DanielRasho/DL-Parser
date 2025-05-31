package writer

import (
	"fmt"
	"os"
	"regexp"
	"text/template"

	parser "github.com/DanielRasho/Parser/internal/Parser"
	table "github.com/DanielRasho/Parser/internal/Parser/TransitionTable"
)

// Writes a parser.go file in the desired location.
// Possible errors:
//   - file paths invalids/not found
//   - invalid parsing table.
//
// REMINDER!!!!! DONT LOAD THE ENTIRE FILE ON A STRING, use buffers instead.
func WriteParserFile(templateFilePath string, outputFilePath string, parserdef *parser.ParserDefinition, transitionTbl *table.GotoTbl, gotoTbl *table.TransitionTbl) error {

	// Load and parse the template
	fmt.Println("PRINTING")
	tmpl, err := template.New("ParserTemplate").Funcs(template.FuncMap{
		"goLiteral": goLiteral,
	}).ParseFiles("./template/ParserTemplate.go")

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
	// err = tmpl.Execute(outFile, data)
	err = tmpl.ExecuteTemplate(outFile, "ParserTemplate", data)
	if err != nil {
		return fmt.Errorf("failed to execute template: %w", err)
	}

	return nil
}

func goLiteral(v any) string {
	raw := fmt.Sprintf("%#v", v)

	// Remove package prefix like "parser.ParserSymbol" -> "ParserSymbol"
	re := regexp.MustCompile(`\b\w+\.(ParserSymbol)\b`)
	return re.ReplaceAllString(raw, `$1`)
}
