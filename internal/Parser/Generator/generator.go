package generator

import (
	"fmt"
	"strconv"
	"strings"

	parser "github.com/DanielRasho/Parser/internal/Parser"
	reader "github.com/DanielRasho/Parser/internal/Parser/Generator/Reader"
	generator "github.com/DanielRasho/Parser/internal/Parser/Generator/Writer"
	table "github.com/DanielRasho/Parser/internal/Parser/TransitionTable"
	"github.com/DanielRasho/Parser/internal/Parser/automata"
	"github.com/golang-collections/collections/queue"
	"github.com/golang-collections/collections/stack"
)

// Given a file to read and a output path, writes a parser definition to the desired path.
func Compile(filePathparser, filepathtemplate, outputPath string, showLogs bool) error {

	// Parse Yalex file definition
	parserDef, err := reader.Parse(filePathparser)
	if err != nil {
		return err
	}

	// runtime.Breakpoint()
	first := table.GetFirst(parserDef)
	follow := table.GetFollow(parserDef, first)

	auto := automata.NewAutomata(parserDef, showLogs)

	transitable, gotable, _ := table.NewTable(auto, first, follow, *parserDef)

	err = generator.WriteParserFile(filepathtemplate, outputPath, parserDef, transitable, gotable)
	if err != nil {
		return err
	}

	return nil
}

// Funcion de referencia
func ParseInput(transit table.TransitionTbl, parserdef parser.ParserDefinition, gotable table.GotoTbl, token []Token) *[]Token {

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
		if !table.CheckNonTerminal(estackval, parserdef) && table.CheckTerminal(queval, parserdef) {
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
		if table.CheckNonTerminal(estackval, parserdef) {

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
		if !table.CheckNonTerminal(queval, parserdef) && !table.CheckTerminal(queval, parserdef) {
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
