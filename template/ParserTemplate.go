{{ define "ParserTemplate" }}
package main

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/golang-collections/collections/queue"
	"github.com/golang-collections/collections/stack"
)


// =============================
// 			TYPES
// =============================

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
	IgnoredSymbols map[int]ParserSymbol
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
	for i := range definition.NonTerminals {
		if definition.NonTerminals[i].Value == id {
			return true
		}
	}

	return false

}

func CheckTerminal(id string, definition ParserDefinition) bool {
	for i := range definition.Terminals {
		if definition.Terminals[i].Value == id {
			return true
		}
	}

	return false

}

type Parser struct {
	parsedefinition *ParserDefinition // Automata for lexeme recognition
	transitiontable *TransitionTbl    //Table to parse the input where do a shift, reduce or to accept said input
	gototable       *GotoTbl          //Table that stores which transition should go
}


func NewParser(filePath string) (*Parser, error) {
	return &Parser{
		parsedefinition: newParserdefinition(), // Automata for lexeme recognition
		transitiontable: newTransitTable(),     //Table to parse the input where do a shift, reduce or to accept said input
		gototable:       newGoToTable(),        //Table that stores which transition should go
	}, nil
}


func (p *Parser) ParseInput(token []Token, parserterminals []ParserSymbol) *[]Token {

	input := ""

	for i := 0; i < len(token); i++ {
		if token[i].TokenID <= len(parserterminals)-1 {
			input = input + " " + parserterminals[token[i].TokenID].Value
		}

	}

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
	fmt.Println(input)
	var staticCount = 0
	var lastEstackVal = ""
	var lastQueueVal = ""

	for accepted {


		//Si es terminal o algo asi osea que no sea int, (, + ) y que el segundo sea terminal, utilizamos la tabla de transition
		if !CheckNonTerminal(estackval, *p.parsedefinition) && CheckTerminal(queval, *p.parsedefinition) {
			switch (*p.transitiontable)[estackval][queval].MovementType {
			case 0:
				_, ok := (*p.transitiontable)[estackval][queval]
				if !ok {
					return &token
				} else {
					topush := strconv.Itoa((*p.transitiontable)[estackval][queval].NextRow)
					estack.Push(q.Dequeue())
					estack.Push(topush)
					estackval = estack.Peek().(string)
					queval = q.Peek().(string)
				}

			case 1:

				var reduced = true
				_, ok := (*p.transitiontable)[estackval][queval]
				if !ok {
					return &token
				} else {
					for i := 0; i < len(p.parsedefinition.Productions[(*p.transitiontable)[estackval][queval].NextRow].Body); i++ {
						for reduced {
							if (p.parsedefinition).Productions[(*p.transitiontable)[estackval][queval].NextRow].Body[i].Value == estack.Peek().(string) {
								reduced = false
								estack.Pop()
								estack.Push(p.parsedefinition.Productions[(*p.transitiontable)[estackval][queval].NextRow].Head.Value)
								estackval = estack.Peek().(string)
								queval = q.Peek().(string)
							} else {
								estack.Pop()
							}

						}
					}

				}

			}

		}
		if CheckNonTerminal(estackval, *p.parsedefinition) {

			lastval := estack.Peek().(string)
			estack.Pop()
			firstval := estack.Peek().(string)

			switch (*p.gototable)[firstval][lastval].MovementType {
			case 2:
				_, ok := (*p.gototable)[firstval][lastval]
				if !ok {
					return &token
				} else {
					estack.Push(lastval)
					topush := strconv.Itoa((*p.gototable)[firstval][lastval].NextRow)
					estack.Push(topush)
					queval = q.Peek().(string)
					estackval = estack.Peek().(string)
				}
			}

		}
		if !CheckNonTerminal(queval, *p.parsedefinition) && !CheckTerminal(queval, *p.parsedefinition) {
			switch (*p.transitiontable)[estackval][queval].MovementType {
			case 1:

				_, ok := (*p.transitiontable)[estackval][queval] //IF IT DOEsNT FIND THE VALUE FROM THE MAP
				if !ok {
					return &token
				} else {
					var reduced = true
					for i := 0; i < len(p.parsedefinition.Productions[(*p.transitiontable)[estackval][queval].NextRow].Body); i++ {
						for reduced {

							if p.parsedefinition.Productions[(*p.transitiontable)[estackval][queval].NextRow].Body[i].Value == estack.Peek().(string) {
								reduced = false
								estack.Pop()
								estack.Push(p.parsedefinition.Productions[(*p.transitiontable)[estackval][queval].NextRow].Head.Value)
								estackval = estack.Peek().(string)
								queval = q.Peek().(string)
							} else {
								estack.Pop()
							}

						}
					}
				}

			case 3:
				fmt.Println("INPUT ACCEPTED")
				accepted = false
				value := ""
				for i := 0; i < len(token); i++ {
					value = value + " " + token[i].Value
				}
				fmt.Printf("\nInput  Code Line: %s        Tokens Line: %s \n", value, input)
				for i := 0; i < estack.Len(); i++ {
					p := estack.Pop().(string)
					value = value + " " + p
				}
				return &[]Token{}

			}

		}
		staticCount++
		if lastEstackVal == estackval && lastQueueVal == queval {
			if staticCount > 3 {
				fmt.Println("Parser got stuck in an infinite loop.")
				return &token
			}
		} else {
			staticCount = 0
		}
		lastEstackVal = estackval
		lastQueueVal = queval



	}

	return nil

}


func newTransitTable() *TransitionTbl {
	return &TransitionTbl{
		{{- range $state, $row := .TransitTable }}
		"{{ $state }}": TransitionTblRow{
			{{- range $symbol, $move := $row }}
			"{{ $symbol }}": Movement{MovementType: {{ $move.MovementType }}, NextRow: {{ $move.NextRow }}},
			{{- end }}
		},
		{{- end }}
	}
}


func newGoToTable() *GotoTbl {
	return &GotoTbl{
		{{- range $state, $row := .Gotable }}
		"{{ $state }}": GotoTblRow{
			{{- range $symbol, $move := $row }}
			"{{ $symbol }}": Movement{MovementType: {{ $move.MovementType }}, NextRow: {{ $move.NextRow }}},
			{{- end }}
		},
		{{- end }}
	}
}


func newParserdefinition() *ParserDefinition {
	return &ParserDefinition{
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
		IgnoredSymbols: {{ goLiteral .ParserDefinition.IgnoredSymbol }},
	}
}
{{ end }}
