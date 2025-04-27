package parser

// Its a programatically representation of a yapar file.
type ParserDefinition struct {
	NonTerminals []ParserSymbol
	Tokens       []ParserSymbol
	Productions  []ParserProduction
}

// Represents a single production declaration
//
//	{Head : "A", Body: ["A",+"A"]}
type ParserProduction struct {
	Head ParserSymbol
	// List of symbols that comprehend a production
	Body []ParserSymbol
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
	// Set to -1 if symbol is NON-TERMINAL
	Id int
	// The string value itself of the symbol
	Value string
}

const NON_TERMINAL_ID = -1

type SymbolSet = map[ParserSymbol]struct{}
