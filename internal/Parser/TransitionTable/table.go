package transitiontable

import (
	"fmt"

	parser "github.com/DanielRasho/Parser/internal/Parser"
)

func GetFirst(def *parser.ParserDefinition) map[string][]parser.ParserSymbol {

	firstSet := make(map[string]parser.SymbolSet, len(def.NonTerminals))

	// Initialize symbol set for each non terminal.
	for _, nonTerminal := range def.NonTerminals {
		head := nonTerminal.Value
		if _, ok := firstSet[head]; !ok {
			firstSet[head] = make(map[parser.ParserSymbol]struct{})
		}
	}

	// For each NON-Terminal
	changed := true

	for changed {
		changed = false
		for _, prod := range def.Productions {
			head := prod.Head.Value
			firstSymbol := prod.Body[0]
			// For terminal symbols
			if firstSymbol.Id != parser.NON_TERMINAL_ID {
				if _, exists := firstSet[head][firstSymbol]; !exists {
					firstSet[head][firstSymbol] = struct{}{}
					changed = true
				}
			}
			// FOR NON Terminal symbols.
			for terminal := range firstSet[firstSymbol.Value] {
				if _, exists := firstSet[head][terminal]; !exists {
					firstSet[head][terminal] = struct{}{}
					changed = true
				}
			}
		}
	}

	// 🖨️ Print the FIRST sets
	fmt.Println("=== FIRST Sets ===")
	for nt, set := range firstSet {
		fmt.Printf("FIRST(%s) = { ", nt)
		for sym := range set {
			fmt.Printf("%s ", sym.Value)
		}
		fmt.Println("}")
	}

	return nil
}
func GetFollow(def *parser.ParserDefinition,
	firsts map[string][]parser.ParserSymbol) map[parser.ParserSymbol]int {
	return nil
}

func NewTable(*parser.ParserDefinition) {

}
