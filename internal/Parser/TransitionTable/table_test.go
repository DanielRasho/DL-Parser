package transitiontable

import (
	"fmt"
	"strconv"
	"strings"
	"testing"

	parser "github.com/DanielRasho/Parser/internal/Parser"
	reader "github.com/DanielRasho/Parser/internal/Parser/Generator/Reader"
	automata "github.com/DanielRasho/Parser/internal/Parser/automata"
	queue "github.com/golang-collections/collections/queue"
	"github.com/golang-collections/collections/stack"
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

	ParseInput(*TransitionTabl, *parserdef, *gotable, "int + int")

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

func ParseInput(transit TransitionTbl, parserdef parser.ParserDefinition, gotable GotoTbl, input string) {

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

			}

		}

		if lasval == queval && onceval == estackval {
			if notaccept > 3 {
				accepted = false
				fmt.Println("NOT ACCEPTED")
			}
			notaccept++

		}

	}

}
