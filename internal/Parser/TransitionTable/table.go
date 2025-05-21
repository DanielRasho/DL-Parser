package transitiontable

import (
	"fmt"
	"strconv"

	parser "github.com/DanielRasho/Parser/internal/Parser"
	automata "github.com/DanielRasho/Parser/internal/Parser/automata"
)

func NewTable(a *automata.Automata, first map[string]parser.SymbolSet, follow map[string]parser.SymbolSet, Parserdefinition parser.ParserDefinition) (*TransitionTbl, *GotoTbl, error) {

	gototable := GotoTbl{}
	transit := TransitionTbl{}
	// Leemos para el go to
	for i := 0; i < len(a.States); i++ {

		value := strconv.Itoa(i)
		gototable[value] = GotoTblRow{}
		transit[value] = TransitionTblRow{}
		for e := range a.States[i].Transitions {

			//Identifica si es no terminal para agregarlo a la tabla de goto
			if CheckNonTerminal(e.Value, Parserdefinition) {

				idnumber := strconv.Itoa(a.States[i].Id)
				gototable[idnumber][e.Value] = Movement{MovementType: 2, NextRow: a.States[i].Transitions[e].Id}
				// fmt.Println(transit)
			}
			// Si es un terminal entonces solo se agrega los shift
			if !CheckNonTerminal(e.Value, Parserdefinition) {
				idnumber := strconv.Itoa(a.States[i].Id)
				transit[idnumber][e.Value] = Movement{MovementType: 0, NextRow: a.States[i].Transitions[e].Id}
			}

		}

		if a.States[i].IsAccepted {
			toshift := a.States[i].Productions[0]
			for e := 0; e < len(a.States[i].Productions); e++ {

				if len(a.States[i].Productions[e].Body) < len(toshift.Body) {
					toshift = a.States[i].Productions[e]
				}

			}

			if a.States[i].Id == 1 {
				idnumber := strconv.Itoa(a.States[i].Id)
				transit[idnumber]["$"] = Movement{MovementType: 3, NextRow: -1}

			} else {

				// fmt.Println(Parserdefinition.Productions[Getindexprodcutions(toshift, Parserdefinition)])
				fset := follow[Parserdefinition.Productions[Getindexprodcutions(toshift, Parserdefinition)].Head.Value]
				for sym := range fset {
					idnumber := strconv.Itoa(a.States[i].Id)
					value := Getindexprodcutions(toshift, Parserdefinition)
					transit[idnumber][sym.Value] = Movement{MovementType: 1, NextRow: value}
				}

			}

		}

	}

	return &transit, &gototable, nil
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

func CheckNonTerminal(id string, definition parser.ParserDefinition) bool {
	for i := 0; i < len(definition.NonTerminals); i++ {
		if definition.NonTerminals[i].Value == id {
			return true
		}
	}

	return false

}

func CheckTerminal(id string, definition parser.ParserDefinition) bool {
	for i := 0; i < len(definition.Terminals); i++ {
		if definition.Terminals[i].Value == id {
			return true
		}
	}

	return false

}

func Getindexprodcutions(prod parser.ParserProduction, parsedef parser.ParserDefinition) int {
	for i := 0; i < len(parsedef.Productions); i++ {
		if prod.Head.Value == parsedef.Productions[i].Head.Value && equalBodies(prod.Body, parsedef.Productions[i].Body) {
			return i
		}
	}
	return -1 // Not found
}

func equalBodies(a, b []parser.ParserSymbol) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i].Value != b[i].Value {
			return false
		}
	}
	return true
}
