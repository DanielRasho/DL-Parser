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
		{{- range $state, $row := .TransitTable }}
		"{{ $state }}": TransitionTblRow{
			{{- range $symbol, $move := $row }}
			"{{ $symbol }}": Movement{MovementType: {{ $move.MovementType }}, NextRow: {{ $move.NextRow }}},
			{{- end }}
		},
		{{- end }}
	}
}

func newGoToTable() GotoTbl {
	return GotoTbl{
		{{- range $state, $row := .Gotable }}
		"{{ $state }}": GotoTblRow{
			{{- range $symbol, $move := $row }}
			"{{ $symbol }}": Movement{MovementType: {{ $move.MovementType }}, NextRow: {{ $move.NextRow }}},
			{{- end }}
		},
		{{- end }}
	}
}


func newParserdefinition() ParserDefinition {
	return ParserDefinition{
		NonTerminals: []ParserSymbol{
			{{- range .ParserDefinition.NonTerminals }}
			{Id: {{ .Id }}, Value: "{{ .Value }}", IsTerminal: {{ .IsTerminal }}},
			{{- end }}
		},
		Terminals: []ParserSymbol{
			{{- range .ParserDefinition.Terminals }}
			{Id: {{ .Id }}, Value: "{{ .Value }}", IsTerminal: {{ .IsTerminal }}},
			{{- end }}
		},
		Productions: []ParserProduction{
			{{- range .ParserDefinition.Productions }}
			{Id: {{ .Id }}, Head: ParserSymbol{Id: {{ .Head.Id }}, Value: "{{ .Head.Value }}", IsTerminal: {{ .Head.IsTerminal }}},
			 Body: []ParserSymbol{
				{{- range .Body }}
				{Id: {{ .Id }}, Value: "{{ .Value }}", IsTerminal: {{ .IsTerminal }}},
				{{- end }}
			 }},
			{{- end }}
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

func (p *Parser) Close() {
	p.file.Close()
}

func ParseInput(transit TransitionTbl, parserdef ParserDefinition, gotable GotoTbl, token []Token, tokenNames []string) *[]Token {

	input := ""

	for i := 0; i < len(token); i++ {
		input = input + " " + tokenNames[token[i].TokenID]
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
