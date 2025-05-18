package automata

import (
	"fmt"

	parser "github.com/DanielRasho/Parser/internal/Parser"
)

const ROOT_PRODUCTION_INDEX int = 0

var ROOT_PRODUCTION = metaProductionId{originalId: ROOT_PRODUCTION_INDEX, index: 0}

func NewAutomata(df *parser.ParserDefinition) *Automata {

	// runtime.Breakpoint()
	productionsDictionary := extendGrammar(df)
	// runtime.Breakpoint()

	root := getRootNode(productionsDictionary)
	root.print()

	queue := []*metaNode{root}
	nodes := []*metaNode{root}

	// toCheck := root.getSymbolsToEvaluate(productionsDictionary)
	// newNode := root.evaluate(toCheck[0], productionsDictionary)
	// newNode.print()

	for len(queue) > 0 {
		currentNode := queue[0]
		queue = queue[1:]

		toCheck := currentNode.getSymbolsToEvaluate(productionsDictionary)

		for _, symbol := range toCheck {
			fmt.Println(symbol.Value)
			fmt.Println("========================")
			newNode := currentNode.evaluate(symbol, productionsDictionary)
			// runtime.Breakpoint()
			nodeExist := checkIdExist(nodes, newNode)
			if nodeExist != nil {
				currentNode.transitions[*symbol] = nodeExist
				continue
			}
			nodes = append(nodes, newNode)
			queue = append(queue, newNode)
			currentNode.transitions[*symbol] = newNode
			newNode.print()
		}
	}

	fmt.Println(len(nodes))

	return nil
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

func (n *metaNode) evaluate(symbol *parser.ParserSymbol, dictionary []*parser.ParserProduction) *metaNode {

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
		name:        2,
		metaProds:   nodeBody,
		completed:   isCompleted,
		isFinal:     isFinal,
		transitions: make(map[parser.ParserSymbol]*metaNode),
	}

	return &childNode
}
