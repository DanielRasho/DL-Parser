package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/golang-collections/collections/queue"
	"github.com/golang-collections/collections/stack"
)


//WE DEFINE THE FOLLOWING STRUCTURES

// TRANSITIONTABLE

// NOTE: ADD THE $ as a new terminal that works as sentinel
type TransitionTbl = map[string]TransitionTblRow
type TransitionTblRow = map[string]Movement

// Goto table
type GotoTbl = map[string]GotoTblRow
type GotoTblRow = map[string]Movement

// Movements
type Movement struct {
	MovementType int
	NextRow      int
}

type MovementType = int

const (
	SHIFT MovementType = iota
	REDUCE
	GOTO
	ACCEPT
)

// PARSER DEFINITION
// Its a programatically representation of a yapar file.
type ParserDefinition struct {
	NonTerminals []ParserSymbol
	Terminals    []ParserSymbol
	Productions  []ParserProduction
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



func CheckNonTerminal(id string, definition ParserDefinition) bool {
	for i := 0; i < len(definition.NonTerminals); i++ {
		if definition.NonTerminals[i].Value == id {
			return true
		}
	}

	return false

}

func CheckTerminal(id string, definition ParserDefinition) bool {
	for i := 0; i < len(definition.Terminals); i++ {
		if definition.Terminals[i].Value == id {
			return true
		}
	}

	return false

}

func newTransitTable() TransitionTbl {
	return TransitionTbl{
		"0": TransitionTblRow{
			"ID": Movement{MovementType: 0, NextRow: 8},
			"LET": Movement{MovementType: 0, NextRow: 5},
			"NUMBER": Movement{MovementType: 0, NextRow: 9},
		},
		"1": TransitionTblRow{
			"$": Movement{MovementType: 3, NextRow: -1},
			"ID": Movement{MovementType: 0, NextRow: 8},
			"LET": Movement{MovementType: 0, NextRow: 5},
			"NUMBER": Movement{MovementType: 0, NextRow: 9},
		},
		"10": TransitionTblRow{
			"$": Movement{MovementType: 1, NextRow: 0},
			"ID": Movement{MovementType: 1, NextRow: 0},
			"LET": Movement{MovementType: 1, NextRow: 0},
			"NUMBER": Movement{MovementType: 1, NextRow: 0},
		},
		"11": TransitionTblRow{
			"$": Movement{MovementType: 1, NextRow: 3},
			"ID": Movement{MovementType: 1, NextRow: 3},
			"LET": Movement{MovementType: 1, NextRow: 3},
			"NUMBER": Movement{MovementType: 1, NextRow: 3},
		},
		"12": TransitionTblRow{
			"ID": Movement{MovementType: 0, NextRow: 8},
			"NUMBER": Movement{MovementType: 0, NextRow: 9},
		},
		"13": TransitionTblRow{
			"ID": Movement{MovementType: 0, NextRow: 8},
			"NUMBER": Movement{MovementType: 0, NextRow: 9},
		},
		"14": TransitionTblRow{
			"ASSIGN": Movement{MovementType: 0, NextRow: 19},
		},
		"15": TransitionTblRow{
			"ID": Movement{MovementType: 0, NextRow: 8},
			"NUMBER": Movement{MovementType: 0, NextRow: 9},
		},
		"16": TransitionTblRow{
			"ID": Movement{MovementType: 0, NextRow: 8},
			"NUMBER": Movement{MovementType: 0, NextRow: 9},
		},
		"17": TransitionTblRow{
			"DIV": Movement{MovementType: 0, NextRow: 16},
			"MINUS": Movement{MovementType: 1, NextRow: 5},
			"MULT": Movement{MovementType: 0, NextRow: 15},
			"PLUS": Movement{MovementType: 1, NextRow: 5},
			"SEMICOLON": Movement{MovementType: 1, NextRow: 5},
		},
		"18": TransitionTblRow{
			"DIV": Movement{MovementType: 0, NextRow: 16},
			"MINUS": Movement{MovementType: 1, NextRow: 6},
			"MULT": Movement{MovementType: 0, NextRow: 15},
			"PLUS": Movement{MovementType: 1, NextRow: 6},
			"SEMICOLON": Movement{MovementType: 1, NextRow: 6},
		},
		"19": TransitionTblRow{
			"ID": Movement{MovementType: 0, NextRow: 8},
			"NUMBER": Movement{MovementType: 0, NextRow: 9},
		},
		"2": TransitionTblRow{
			"$": Movement{MovementType: 1, NextRow: 1},
			"ID": Movement{MovementType: 1, NextRow: 1},
			"LET": Movement{MovementType: 1, NextRow: 1},
			"NUMBER": Movement{MovementType: 1, NextRow: 1},
		},
		"20": TransitionTblRow{
			"DIV": Movement{MovementType: 1, NextRow: 8},
			"MINUS": Movement{MovementType: 1, NextRow: 8},
			"MULT": Movement{MovementType: 1, NextRow: 8},
			"PLUS": Movement{MovementType: 1, NextRow: 8},
			"SEMICOLON": Movement{MovementType: 1, NextRow: 8},
		},
		"21": TransitionTblRow{
			"DIV": Movement{MovementType: 1, NextRow: 9},
			"MINUS": Movement{MovementType: 1, NextRow: 9},
			"MULT": Movement{MovementType: 1, NextRow: 9},
			"PLUS": Movement{MovementType: 1, NextRow: 9},
			"SEMICOLON": Movement{MovementType: 1, NextRow: 9},
		},
		"22": TransitionTblRow{
			"MINUS": Movement{MovementType: 0, NextRow: 13},
			"PLUS": Movement{MovementType: 0, NextRow: 12},
			"SEMICOLON": Movement{MovementType: 0, NextRow: 23},
		},
		"23": TransitionTblRow{
			"$": Movement{MovementType: 1, NextRow: 4},
			"ID": Movement{MovementType: 1, NextRow: 4},
			"LET": Movement{MovementType: 1, NextRow: 4},
			"NUMBER": Movement{MovementType: 1, NextRow: 4},
		},
		"3": TransitionTblRow{
			"$": Movement{MovementType: 1, NextRow: 2},
			"ID": Movement{MovementType: 1, NextRow: 2},
			"LET": Movement{MovementType: 1, NextRow: 2},
			"NUMBER": Movement{MovementType: 1, NextRow: 2},
		},
		"4": TransitionTblRow{
			"MINUS": Movement{MovementType: 0, NextRow: 13},
			"PLUS": Movement{MovementType: 0, NextRow: 12},
			"SEMICOLON": Movement{MovementType: 0, NextRow: 11},
		},
		"5": TransitionTblRow{
			"ID": Movement{MovementType: 0, NextRow: 14},
		},
		"6": TransitionTblRow{
			"DIV": Movement{MovementType: 0, NextRow: 16},
			"MINUS": Movement{MovementType: 1, NextRow: 7},
			"MULT": Movement{MovementType: 0, NextRow: 15},
			"PLUS": Movement{MovementType: 1, NextRow: 7},
			"SEMICOLON": Movement{MovementType: 1, NextRow: 7},
		},
		"7": TransitionTblRow{
			"DIV": Movement{MovementType: 1, NextRow: 10},
			"MINUS": Movement{MovementType: 1, NextRow: 10},
			"MULT": Movement{MovementType: 1, NextRow: 10},
			"PLUS": Movement{MovementType: 1, NextRow: 10},
			"SEMICOLON": Movement{MovementType: 1, NextRow: 10},
		},
		"8": TransitionTblRow{
			"DIV": Movement{MovementType: 1, NextRow: 11},
			"MINUS": Movement{MovementType: 1, NextRow: 11},
			"MULT": Movement{MovementType: 1, NextRow: 11},
			"PLUS": Movement{MovementType: 1, NextRow: 11},
			"SEMICOLON": Movement{MovementType: 1, NextRow: 11},
		},
		"9": TransitionTblRow{
			"DIV": Movement{MovementType: 1, NextRow: 12},
			"MINUS": Movement{MovementType: 1, NextRow: 12},
			"MULT": Movement{MovementType: 1, NextRow: 12},
			"PLUS": Movement{MovementType: 1, NextRow: 12},
			"SEMICOLON": Movement{MovementType: 1, NextRow: 12},
		},
	}
}

func newGoToTable() GotoTbl {
	return GotoTbl{
		"0": GotoTblRow{
			"expression": Movement{MovementType: 2, NextRow: 4},
			"factor": Movement{MovementType: 2, NextRow: 7},
			"program": Movement{MovementType: 2, NextRow: 1},
			"statement": Movement{MovementType: 2, NextRow: 2},
			"term": Movement{MovementType: 2, NextRow: 6},
			"var_decl": Movement{MovementType: 2, NextRow: 3},
		},
		"1": GotoTblRow{
			"expression": Movement{MovementType: 2, NextRow: 4},
			"factor": Movement{MovementType: 2, NextRow: 7},
			"statement": Movement{MovementType: 2, NextRow: 10},
			"term": Movement{MovementType: 2, NextRow: 6},
			"var_decl": Movement{MovementType: 2, NextRow: 3},
		},
		"10": GotoTblRow{
		},
		"11": GotoTblRow{
		},
		"12": GotoTblRow{
			"factor": Movement{MovementType: 2, NextRow: 7},
			"term": Movement{MovementType: 2, NextRow: 17},
		},
		"13": GotoTblRow{
			"factor": Movement{MovementType: 2, NextRow: 7},
			"term": Movement{MovementType: 2, NextRow: 18},
		},
		"14": GotoTblRow{
		},
		"15": GotoTblRow{
			"factor": Movement{MovementType: 2, NextRow: 20},
		},
		"16": GotoTblRow{
			"factor": Movement{MovementType: 2, NextRow: 21},
		},
		"17": GotoTblRow{
		},
		"18": GotoTblRow{
		},
		"19": GotoTblRow{
			"expression": Movement{MovementType: 2, NextRow: 22},
			"factor": Movement{MovementType: 2, NextRow: 7},
			"term": Movement{MovementType: 2, NextRow: 6},
		},
		"2": GotoTblRow{
		},
		"20": GotoTblRow{
		},
		"21": GotoTblRow{
		},
		"22": GotoTblRow{
		},
		"23": GotoTblRow{
		},
		"3": GotoTblRow{
		},
		"4": GotoTblRow{
		},
		"5": GotoTblRow{
		},
		"6": GotoTblRow{
		},
		"7": GotoTblRow{
		},
		"8": GotoTblRow{
		},
		"9": GotoTblRow{
		},
	}
}


func newParserdefinition() ParserDefinition {
	return ParserDefinition{
		NonTerminals: []ParserSymbol{
			{Id: -1, Value: "program", IsTerminal: false},
			{Id: -1, Value: "statement", IsTerminal: false},
			{Id: -1, Value: "var_decl", IsTerminal: false},
			{Id: -1, Value: "expression", IsTerminal: false},
			{Id: -1, Value: "term", IsTerminal: false},
			{Id: -1, Value: "factor", IsTerminal: false},
		},
		Terminals: []ParserSymbol{
			{Id: 0, Value: "LET", IsTerminal: true},
			{Id: 1, Value: "ASSIGN", IsTerminal: true},
			{Id: 2, Value: "PLUS", IsTerminal: true},
			{Id: 3, Value: "MINUS", IsTerminal: true},
			{Id: 4, Value: "MULT", IsTerminal: true},
			{Id: 5, Value: "DIV", IsTerminal: true},
			{Id: 6, Value: "ID", IsTerminal: true},
			{Id: 7, Value: "NUMBER", IsTerminal: true},
			{Id: 8, Value: "SEMICOLON", IsTerminal: true},
			{Id: 9, Value: "WS", IsTerminal: true},
		},
		Productions: []ParserProduction{
			{Id: 1, Head: ParserSymbol{Id: -1, Value: "program", IsTerminal: false},
			 Body: []ParserSymbol{
				{Id: -1, Value: "program", IsTerminal: false},
				{Id: -1, Value: "statement", IsTerminal: false},
			 }},
			{Id: 2, Head: ParserSymbol{Id: -1, Value: "program", IsTerminal: false},
			 Body: []ParserSymbol{
				{Id: -1, Value: "statement", IsTerminal: false},
			 }},
			{Id: 3, Head: ParserSymbol{Id: -1, Value: "statement", IsTerminal: false},
			 Body: []ParserSymbol{
				{Id: -1, Value: "var_decl", IsTerminal: false},
			 }},
			{Id: 4, Head: ParserSymbol{Id: -1, Value: "statement", IsTerminal: false},
			 Body: []ParserSymbol{
				{Id: -1, Value: "expression", IsTerminal: false},
				{Id: 8, Value: "SEMICOLON", IsTerminal: true},
			 }},
			{Id: 5, Head: ParserSymbol{Id: -1, Value: "var_decl", IsTerminal: false},
			 Body: []ParserSymbol{
				{Id: 0, Value: "LET", IsTerminal: true},
				{Id: 6, Value: "ID", IsTerminal: true},
				{Id: 1, Value: "ASSIGN", IsTerminal: true},
				{Id: -1, Value: "expression", IsTerminal: false},
				{Id: 8, Value: "SEMICOLON", IsTerminal: true},
			 }},
			{Id: 6, Head: ParserSymbol{Id: -1, Value: "expression", IsTerminal: false},
			 Body: []ParserSymbol{
				{Id: -1, Value: "expression", IsTerminal: false},
				{Id: 2, Value: "PLUS", IsTerminal: true},
				{Id: -1, Value: "term", IsTerminal: false},
			 }},
			{Id: 7, Head: ParserSymbol{Id: -1, Value: "expression", IsTerminal: false},
			 Body: []ParserSymbol{
				{Id: -1, Value: "expression", IsTerminal: false},
				{Id: 3, Value: "MINUS", IsTerminal: true},
				{Id: -1, Value: "term", IsTerminal: false},
			 }},
			{Id: 8, Head: ParserSymbol{Id: -1, Value: "expression", IsTerminal: false},
			 Body: []ParserSymbol{
				{Id: -1, Value: "term", IsTerminal: false},
			 }},
			{Id: 9, Head: ParserSymbol{Id: -1, Value: "term", IsTerminal: false},
			 Body: []ParserSymbol{
				{Id: -1, Value: "term", IsTerminal: false},
				{Id: 4, Value: "MULT", IsTerminal: true},
				{Id: -1, Value: "factor", IsTerminal: false},
			 }},
			{Id: 10, Head: ParserSymbol{Id: -1, Value: "term", IsTerminal: false},
			 Body: []ParserSymbol{
				{Id: -1, Value: "term", IsTerminal: false},
				{Id: 5, Value: "DIV", IsTerminal: true},
				{Id: -1, Value: "factor", IsTerminal: false},
			 }},
			{Id: 11, Head: ParserSymbol{Id: -1, Value: "term", IsTerminal: false},
			 Body: []ParserSymbol{
				{Id: -1, Value: "factor", IsTerminal: false},
			 }},
			{Id: 12, Head: ParserSymbol{Id: -1, Value: "factor", IsTerminal: false},
			 Body: []ParserSymbol{
				{Id: 6, Value: "ID", IsTerminal: true},
			 }},
			{Id: 13, Head: ParserSymbol{Id: -1, Value: "factor", IsTerminal: false},
			 Body: []ParserSymbol{
				{Id: 7, Value: "NUMBER", IsTerminal: true},
			 }},
		},
	}
}


type Parser struct {
	file            *os.File         // File to read from
	reader          *bufio.Reader    // Reader to get the symbols from file
	parsedefinition ParserDefinition // Automata for lexeme recognition
	transitiontable TransitionTbl    //Table to parse the input where do a shift, reduce or to accept said input
	gototable       GotoTbl          //Table that stores which transition should go
	bytesRead       int              // Number of bytes the lexer has read
}


func NewParser(filePath string) (*Parser, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	return &Parser{
		file:            file,
		reader:          bufio.NewReader(file),
		parsedefinition: newParserdefinition(), // Automata for lexeme recognition
		transitiontable: newTransitTable(),     //Table to parse the input where do a shift, reduce or to accept said input
		gototable:       newGoToTable(),        //Table that stores which transition should go
	}, nil
}


func ParseInput(transit TransitionTbl, parserdef ParserDefinition, gotable GotoTbl, token []Token) *[]Token {

	input := ""

	for i := 0; i < len(token); i++ {
		input = input + " " + token[i].Value
	}

	fmt.Println(input)
	tokens := strings.Fields(input) // ["input", "+", "input"]

	q := queue.New()

	for _, token := range tokens {
		q.Enqueue(token)
	}
	q.Enqueue("$")

	estack := stack.New()

	estack.Push("0") // Se declara el inicial
	estackval := estack.Peek().(string)

	queval := q.Peek().(string)

	// #Empezamos a parsear
	var accepted = true
	var notaccept = 0
	for accepted {

		fmt.Println(estackval, queval)

		onceval := estackval
		lasval := queval

		//Si es terminal o algo asi osea que no sea int, (, + ) y que el segundo sea terminal, utilizamos la tabla de transition
		if !CheckNonTerminal(estackval, parserdef) && CheckTerminal(queval, parserdef) {
			fmt.Println(transit[estackval][queval].MovementType)
			switch transit[estackval][queval].MovementType {
			case 0:
				fmt.Println("printing shift")
				topush := strconv.Itoa(transit[estackval][queval].NextRow)
				fmt.Println(topush)
				estack.Push(q.Dequeue())
				estack.Push(topush)
				estackval = estack.Peek().(string)
				queval = q.Peek().(string)

				fmt.Println(estackval, queval)

			case 1:
				fmt.Println("printing REDUCE")
				var reduced = true
				for i := 0; i < len(parserdef.Productions[transit[estackval][queval].NextRow].Body); i++ {
					for reduced {

						if parserdef.Productions[transit[estackval][queval].NextRow].Body[i].Value == estack.Peek().(string) {
							reduced = false
							estack.Pop()
							estack.Push(parserdef.Productions[transit[estackval][queval].NextRow].Head.Value)
							estackval = estack.Peek().(string)
							queval = q.Peek().(string)
						} else {
							estack.Pop()
						}

					}
				}
				fmt.Println(estackval, queval)

			}

		}
		if CheckNonTerminal(estackval, parserdef) {

			lastval := estack.Peek().(string)
			estack.Pop()
			firstval := estack.Peek().(string)

			switch gotable[firstval][lastval].MovementType {
			case 2:
				fmt.Println("printing GOTO")
				estack.Push(lastval)
				topush := strconv.Itoa(gotable[firstval][lastval].NextRow)
				estack.Push(topush)
				queval = q.Peek().(string)
				estackval = estack.Peek().(string)
				fmt.Println(estackval, queval)
			}

		}
		if !CheckNonTerminal(queval, parserdef) && !CheckTerminal(queval, parserdef) {
			fmt.Println("ADDING TO STACK")
			fmt.Println(estackval, queval)
			fmt.Println(transit[estackval][queval].MovementType)
			switch transit[estackval][queval].MovementType {
			case 1:
				var reduced = true
				for i := 0; i < len(parserdef.Productions[transit[estackval][queval].NextRow].Body); i++ {
					for reduced {

						if parserdef.Productions[transit[estackval][queval].NextRow].Body[i].Value == estack.Peek().(string) {
							reduced = false
							estack.Pop()
							estack.Push(parserdef.Productions[transit[estackval][queval].NextRow].Head.Value)
							estackval = estack.Peek().(string)
							queval = q.Peek().(string)
						} else {
							estack.Pop()
						}

					}
				}
				fmt.Println(estackval, queval)

			case 3:
				fmt.Println("Accepted the input")
				accepted = false
				return &[]Token{}

			}

		}

		if lasval == queval && onceval == estackval {
			if notaccept > 3 {
				accepted = false
				fmt.Println("NOT ACCEPTED")
				return &token
			}
			notaccept++

		}

	}

	return nil

}
