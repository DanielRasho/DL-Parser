package automata

import (
	parser "github.com/DanielRasho/Parser/internal/Parser"
)

func NewAutomata(df *parser.ParserDefinition) *Automata {

	// runtime.Breakpoint()
	productionsDictionary := extendGrammar(df)
	// runtime.Breakpoint()

	getRootNode(productionsDictionary)

	return nil
}

// Insert Root production S -> E
// Which points to the final production
func extendGrammar(df *parser.ParserDefinition) []*parser.ParserProduction {

	productions := make([]*parser.ParserProduction, 0, len(df.Productions)+1)

	// Inserting root production "S"
	productions = append(productions,
		&parser.ParserProduction{
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

func getRootNode(productions []*parser.ParserProduction) *metaNode {

	visited := make(map[parser.ParserSymbol]struct{})
	visited[productions[0].Head] = struct{}{}

	rootId := make(metaNodeId)
	rootId[0] = struct{}{}

	rootBody := []metaProduction{
		{id: 0, isRoot: true, completed: false, index: 0}}

	// Enqueue the fist element of the extended list of productions
	// This will always be the extra production "S"
	queue := []*parser.ParserSymbol{&productions[0].Body[0]}

	for len(queue) > 0 {

		symbol := queue[0]
		queue = queue[1:]
		if _, ok := visited[*symbol]; ok {
			continue
		}
		visited[*symbol] = struct{}{}

		// Search for productions with head == Symbol
		for _, p := range productions {
			if p.Head != *symbol {
				continue
			}

			// Add id to the root ID
			rootId[p.Id] = struct{}{}

			// Add productions found to rootBody
			rootBody = append(rootBody, metaProduction{
				id:        p.Id,
				isRoot:    false,
				completed: false,
				index:     0})

			// If first element of the production's body is NON TERMINAL
			// Add it to the queue
			if p.Body[0].Id == parser.NON_TERMINAL_ID {
				queue = append(queue, &p.Body[0])
			}
		}
	}

	root := metaNode{
		id:          rootId,
		name:        0,
		productions: rootBody,
		completed:   false,
		isFinal:     false,
		Transitions: make(map[Symbol]metaNodeId),
	}

	root.print()

	return &root
}

func getFinalAutomata(metaAutomata *metaAutomata) *Automata {
	return nil
}

func (n *metaNode) getSymbolsToEvaluate() []*parser.ParserSymbol {

	return nil
}

func (n *metaNode) evaluate(symbol *parser.ParserSymbol) *metaNode {

	return nil
}
