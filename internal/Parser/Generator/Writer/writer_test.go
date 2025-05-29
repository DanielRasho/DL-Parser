package writer

import (
	"testing"

	reader "github.com/DanielRasho/Parser/internal/Parser/Generator/Reader"
	table "github.com/DanielRasho/Parser/internal/Parser/TransitionTable"
	automata "github.com/DanielRasho/Parser/internal/Parser/automata"
)

func Test_check2(t *testing.T) {

	parserDef, err := reader.Parse("../../../../examples/productions2.y")
	if err != nil {
	}

	first := table.GetFirst(parserDef)
	follow := table.GetFollow(parserDef, first)
	automa := automata.NewAutomata(parserDef, false)

	transitionTbl, gotoTbl, _ := table.NewTable(automa, first, follow, *parserDef)

	WriteParserFile("../../../../template/ParserTemplate.go", "../../../../cmd/compiler/parser.go", parserDef, transitionTbl, gotoTbl)

}
