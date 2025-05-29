package transitiontable

import (
	"fmt"
	"testing"

	reader "github.com/DanielRasho/Parser/internal/Parser/Generator/Reader"
	automata "github.com/DanielRasho/Parser/internal/Parser/automata"
)

func Test_check1(t *testing.T) {

	parserdef, _ := reader.Parse("../../../examples/productions2.y")
	first := GetFirst(parserdef)
	follow := GetFollow(parserdef, first)
	var automa = automata.NewAutomata(parserdef, false)

	TransitionTabl, gotable, _ := NewTable(automa, first, follow, *parserdef)

	var debug = true

	if debug {

		fmt.Println("TABLA DE TRANSICION")
		fmt.Println(TransitionTabl)
		fmt.Println("GOTO")
		fmt.Println(gotable)

		for state, transitions := range *TransitionTabl {
			fmt.Printf("State %s:\n", state)
			for symbol, action := range transitions {
				fmt.Printf("  %s => {%d %d}\n", symbol, action.MovementType, action.NextRow)
			}
		}

	}

	tokens := []Token{Token{Value: "int", TokenID: 2, Offset: 0}, Token{Value: "+", TokenID: 2, Offset: 0}, Token{Value: "int", TokenID: 2, Offset: 0}}

	ParseInput(*TransitionTabl, *parserdef, *gotable, tokens)

	tokens = []Token{Token{Value: "int", TokenID: 2, Offset: 0}, Token{Value: "int", TokenID: 2, Offset: 0}}

	ParseInput(*TransitionTabl, *parserdef, *gotable, tokens)

}

func Test_check2(t *testing.T) {

	parserdef, _ := reader.Parse("../../../examples/productions.y")
	first := GetFirst(parserdef)
	follow := GetFollow(parserdef, first)
	var automa = automata.NewAutomata(parserdef, false)

	TransitionTabl, gotable, _ := NewTable(automa, first, follow, *parserdef)

	fmt.Println("TABLA DE TRANSICION")
	fmt.Println(TransitionTabl)
	fmt.Println("GOTO")
	fmt.Println(gotable)

	for state, transitions := range *TransitionTabl {
		fmt.Printf("State %s:\n", state)
		for symbol, action := range transitions {
			fmt.Printf("  %s => {%d %d}\n", symbol, action.MovementType, action.NextRow)
		}
	}

}

type Symbol = string

type Token struct {
	Value   Symbol // Actual string read by the lexer
	TokenID int    // Token Id (defined by the user above)
	Offset  int    // No of bytes from the start of the file to the current lexeme
}
