package transitiontable

import (
	"fmt"

	parser "github.com/DanielRasho/Parser/internal/Parser"
	"github.com/DanielRasho/Parser/internal/Parser/TransitionTable/automata"
)

func NewTransitionTable(a *automata.Automata) (*TransitionTbl, error) {

	return nil, nil
}

func GetFirst(def *parser.ParserDefinition) map[string]parser.SymbolSet {

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

	fmt.Println("=== FIRST Sets ===")
	for nt, set := range firstSet {
		fmt.Printf("FIRST(%s) = { ", nt)
		for sym := range set {
			fmt.Printf("%s ", sym.Value)
		}
		fmt.Println("}")
	}

	return firstSet
}

func GetFollow(def *parser.ParserDefinition,
	firstSet map[string]parser.SymbolSet) map[string]parser.SymbolSet {

	followSet := make(map[string]parser.SymbolSet, len(def.NonTerminals))

	// Initialize symbol set for each non terminal.
	for _, nonTerminal := range def.NonTerminals {
		head := nonTerminal.Value
		if _, ok := followSet[head]; !ok {
			followSet[head] = make(map[parser.ParserSymbol]struct{})
		}
	}

	fmt.Printf("%v\n", def.Productions[0].Head.Value)
	initialValue := parser.ParserSymbol{Id: 0, Value: "$"}

	// Add initial symbol
	followSet[def.Productions[0].Head.Value][initialValue] = struct{}{}

	changed := true

	// runtime.Breakpoint()

	for changed {
		changed = false
		for _, prod := range def.Productions {
			for i := 0; i < len(prod.Body); i++ {
				symbol := prod.Body[i]
				// If is terminal ignore it
				if symbol.Id != parser.NON_TERMINAL_ID {
					continue
				}

				// If form A -> a B b
				if i+1 < len(prod.Body) {
					target := prod.Body[i+1]
					if target.Id == parser.NON_TERMINAL_ID {
						// FOR NON Terminal symbols.
						for terminal := range firstSet[target.Value] {
							if _, exists := followSet[symbol.Value][terminal]; !exists {
								followSet[symbol.Value][terminal] = struct{}{}
								changed = true
							}
						}
					} else {
						if _, exist := followSet[symbol.Value][target]; !exist {
							followSet[symbol.Value][target] = struct{}{}
							changed = true
						}
					}
					// If form A -> a B
				} else {
					target := prod.Head
					for terminal := range followSet[target.Value] {
						if _, exists := followSet[symbol.Value][terminal]; !exists {
							followSet[symbol.Value][terminal] = struct{}{}
							changed = true
						}
					}
				}

			}
		}
	}

	fmt.Println("=== FOLLOW Sets ===")
	for nt, set := range followSet {
		fmt.Printf("FOLLOW(%s) = { ", nt)
		for sym := range set {
			fmt.Printf("%s ", sym.Value)
		}
		fmt.Println("}")
	}

	return followSet
}

func NewTable(*parser.ParserDefinition) {

}
