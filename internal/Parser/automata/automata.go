package automata

import (
	"fmt"

	parser "github.com/DanielRasho/Parser/internal/Parser"
)

const ROOT_PRODUCTION_INDEX int = 0

var ROOT_PRODUCTION = metaProductionId{originalId: ROOT_PRODUCTION_INDEX, index: 0}

func NewAutomata(df *parser.ParserDefinition, showLogs bool) *Automata {

	productionsDictionary := extendGrammar(df)

	if showLogs {
		for i := range productionsDictionary {
			fmt.Println(productionsDictionary[i].String())
		}
	}

	// Build the Root node, with required productions
	root := getRootNode(productionsDictionary)

	if showLogs {
		root.print()
	}

	queue := []*metaNode{root}
	nodes := []*metaNode{root}

	// Build remaining roots progresivelly
	for len(queue) > 0 {
		currentNode := queue[0]
		queue = queue[1:]

		// Get which symbols must be evaluated based on the productions
		// of the current node
		toCheck := currentNode.getSymbolsToEvaluate(productionsDictionary)

		// Evalute each symbol, and produce a new child node
		for _, symbol := range toCheck {
			newNode := currentNode.evaluate(symbol, productionsDictionary, len(nodes))
			// If child node doest not currently exist add it to the automata
			nodeExist := checkIdExist(nodes, newNode)
			if nodeExist != nil {
				currentNode.transitions[*symbol] = nodeExist
				continue
			}
			// if not add, it to the automata
			nodes = append(nodes, newNode)
			// add it to the queue for future evaluation
			queue = append(queue, newNode)
			// Add a transition state from the current node to the child node
			currentNode.transitions[*symbol] = newNode
			if showLogs {
				fmt.Println(symbol.Value)
				fmt.Println("========================")
				newNode.print()
			}
		}
	}

	if showLogs {
		fmt.Println("TOTAL NODES:")
		fmt.Println(len(nodes))
	}

	// BUILD FINAL AUTOMATA
	states := make([]*State, 0, len(nodes))

	for i := range nodes {
		node := nodes[i]
		productions := make([]parser.ParserProduction, 0, len(node.metaProds))
		transitions := make(map[parser.ParserSymbol]*State, len(node.transitions))

		for _, p := range node.metaProds {
			productions = append(productions, *productionsDictionary[p.getDictIndex()])
		}

		newState := State{
			Id:          node.name,
			Productions: productions,
			Transitions: transitions,
			IsFinal:     node.isFinal,
			IsAccepted:  node.completed}
		states = append(states, &newState)
	}

	for i := range nodes {
		node := nodes[i]
		state := states[i]

		for key, value := range node.transitions {
			state.Transitions[key] = states[value.name]
		}
	}

	automata := Automata{
		StartState: states[0],
		States:     states,
	}

	dot := GenerateDOT_SLR0(&automata)

	GenerateImage(dot, "./diagrams/SLR0_Automata.png")

	return &automata
}

func checkIdExist(nodes []*metaNode, newNode *metaNode) *metaNode {
	for _, node := range nodes {
		if areEqual := setsEqual(node.id, newNode.id); areEqual {
			return node
		}
	}
	return nil
}

// Insert Root production S -> E
// Which points to the final production
func extendGrammar(df *parser.ParserDefinition) []*parser.ParserProduction {

	productions := make([]*parser.ParserProduction, 0, len(df.Productions)+1)

	// Inserting root production "S"
	productions = append(productions,
		&parser.ParserProduction{
			Id: ROOT_PRODUCTION_INDEX,
			Head: parser.ParserSymbol{
				Id: -1,
				// Empty space is a value a production would not receive
				// under normal circumstances, therefore secure to use to
				// inserte a new production
				Value:      "",
				IsTerminal: false,
			},
			Body: []parser.ParserSymbol{df.Productions[0].Head},
		})

	for i := range df.Productions {
		productions = append(productions, &df.Productions[i])
	}

	return productions
}

// Builds the root node for the SLR automata
func getRootNode(productions []*parser.ParserProduction) *metaNode {
	rootId := make(metaNodeId)
	rootId[metaProductionId{originalId: 0, index: 0}] = struct{}{}

	rootBody := []metaProduction{
		{id: metaProductionId{originalId: 0, index: 0},
			isRoot:    true,
			completed: false,
			length:    1,
		}}

	_, closure := getProductionClosure(productions, productions[0].Body[0])
	for _, p := range closure {
		if _, ok := rootId[p.id]; ok {
			continue
		}
		rootId[p.id] = struct{}{}
		rootBody = append(rootBody, p)
	}

	root := metaNode{
		id:          rootId,
		name:        0,
		metaProds:   rootBody,
		completed:   false,
		isFinal:     false,
		transitions: make(map[parser.ParserSymbol]*metaNode),
	}

	return &root
}

func getProductionClosure(dictionary []*parser.ParserProduction,
	target parser.ParserSymbol) (metaNodeId, []metaProduction) {

	visited := make(map[parser.ParserSymbol]struct{})

	nodeId := make(metaNodeId)

	closure := []metaProduction{}

	// Enqueue the fist element of the extended list of productions
	// This will always be the extra production "S"
	queue := []*parser.ParserSymbol{&target}

	for len(queue) > 0 {

		symbol := queue[0]
		queue = queue[1:]
		if _, ok := visited[*symbol]; ok {
			continue
		}
		visited[*symbol] = struct{}{}

		// Search for productions with head == Symbol
		for _, p := range dictionary {
			if p.Head != *symbol {
				continue
			}

			// Build new productoin Id
			newId := metaProductionId{originalId: p.Id, index: 0}
			// Add id to the root ID
			nodeId[newId] = struct{}{}

			// print(len(p.Body))

			// Add productions found to rootBody
			closure = append(closure, metaProduction{
				id:        newId,
				isRoot:    false,
				completed: false,
				length:    len(p.Body),
			})

			// If first element of the production's body is NON TERMINAL
			// Add it to the queue
			if p.Body[0].Id == parser.NON_TERMINAL_ID {
				queue = append(queue, &p.Body[0])
			}
		}
	}

	return nodeId, closure
}

func (n *metaNode) getSymbolsToEvaluate(dictionary []*parser.ParserProduction) []*parser.ParserSymbol {

	toCheck := make([]*parser.ParserSymbol, 0, len(n.metaProds))

	inserted := make(map[parser.ParserSymbol]struct{})

	for _, p := range n.metaProds {
		if p.completed {
			continue
		}

		// Getting the target production for dictionary
		// And then current target symbol
		targetSymbol := dictionary[p.id.originalId].Body[p.getIndex()]
		// If production already added ignore it
		if _, ok := inserted[targetSymbol]; ok {
			continue
		}
		// Add symbol to the list
		toCheck = append(toCheck, &targetSymbol)
		inserted[targetSymbol] = struct{}{}
	}

	return toCheck
}

func (n *metaNode) evaluate(symbol *parser.ParserSymbol, dictionary []*parser.ParserProduction, nextName int) *metaNode {

	nodeId := make(metaNodeId)
	nodeBody := make([]metaProduction, 0)
	isCompleted := false
	isFinal := false
	symbolsToClosure := make([]parser.ParserSymbol, 0)

	// runtime.Breakpoint()
	// Loop over the productions of the current node
	for _, metaProd := range n.metaProds {

		// If its is completed ignore it
		if metaProd.completed {
			continue
		}

		production := dictionary[metaProd.getDictIndex()]
		targetSymbol := production.Body[metaProd.getIndex()]

		// Select all productions were current_symbol == symbol
		if targetSymbol != *symbol {
			continue
		}

		newIndex := metaProd.getIndex() + 1

		// Move its point
		if newIndex <= metaProd.length {
			// If index its already at the end, the scanned is completed
			if metaProd.getIndex() == metaProd.length-1 {
				metaProd.completed = true
				isCompleted = true
				metaProd.id.index++
			} else {
				metaProd.id.index++
				symbolsToClosure = append(symbolsToClosure, production.Body[metaProd.getIndex()])
			}
		}

		// Check if its final
		if metaProd.id == ROOT_PRODUCTION {
			isFinal = true
		}

		nodeId[metaProd.id] = struct{}{}
		nodeBody = append(nodeBody, metaProd)
		nodeId[metaProd.id] = struct{}{}
	}

	for _, s := range symbolsToClosure {
		_, closure := getProductionClosure(dictionary, s)
		for _, p := range closure {
			if _, ok := nodeId[p.id]; ok {
				continue
			}
			nodeId[p.id] = struct{}{}
			nodeBody = append(nodeBody, p)
		}
	}

	childNode := metaNode{
		id:          nodeId,
		name:        nextName,
		metaProds:   nodeBody,
		completed:   isCompleted,
		isFinal:     isFinal,
		transitions: make(map[parser.ParserSymbol]*metaNode),
	}

	return &childNode
}
