package writer

import (
	parserdef "github.com/DanielRasho/Parser/internal/Parser"
	table "github.com/DanielRasho/Parser/internal/Parser/TransitionTable"
)

// This module is in charge of writing the final Lexer.go file, based on a template file

// Definition of variable fields withing a template
type templateLexwrite struct {
	Gotable          table.GotoTbl
	TransitTable     table.TransitionTbl
	ParserDefinition parserdef.ParserDefinition
}
