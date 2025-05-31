package parser

import (
	"fmt"
	"strings"
)

// Its a programatically representation of a yapar file.
type ParserDefinition struct {
	NonTerminals  []ParserSymbol
	Terminals     []ParserSymbol
	Productions   []ParserProduction
	IgnoredSymbol map[int]ParserSymbol
}

// Represents a single production declaration
//
//	{Head : "A", Body: ["A",+"A"]}
type ParserProduction struct {
	// Given by the order of definition in the yapar file, starting from 1
	Id   int
	Head ParserSymbol
	// List of symbols that comprehend a production
	Body []ParserSymbol
}

func (p *ParserProduction) String() string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d: %s → ", p.Id, p.Head.Value))

	for i, symbol := range p.Body {
		if i > 0 {
			sb.WriteString(" ")
		}
		sb.WriteString(symbol.Value)
	}
	return sb.String()
}

// Smallest information unit, the parser can read. Which is basically a symbol
// which can a terminal or non terminal.
// Example TERMINAL:
//
//	{Id: 2, Value: "NUMBER"}
//
// Example TERMINAL:
//
//	{Id: -1, Value: "A"}
type ParserSymbol struct {
	// If token is terminal, the id comes from the order of declaration in
	// the Yapar file.
	// Start from 1
	Id int
	// The string value itself of the symbol
	Value string

	IsTerminal bool
}

const NON_TERMINAL_ID = -1

// Used for first-follow computations
type SymbolSet = map[ParserSymbol]struct{}
