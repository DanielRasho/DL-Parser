package automata

import (
	"fmt"
	"strings"

	parser "github.com/DanielRasho/Parser/internal/Parser"
)

type Automata struct {
	StartState *State
	States     []*State
}

type Symbol = int

type State struct {
	Id          string
	Productions []parser.ParserProduction // Sorted by highest too lower priority ( 0 has the hightes priority )
	Transitions map[Symbol]*State         // {"a": STATE1, "b": STATE2, "NUMBER": STATEFINAL}
	IsFinal     bool
	IsAccepted  bool
}

// =========================
// 	INTERNAL
// =========================

type metaAutomata struct {
	StartState *metaNode
	States     []*metaNode
}

type metaNodeId = map[int]struct{}

type metaNode struct {
	id          metaNodeId
	name        int
	productions []metaProduction
	completed   bool
	isFinal     bool
	Transitions map[Symbol]metaNodeId
}

// PrettyPrint prints the metaNode in a readable format
func (m *metaNode) print() {
	fmt.Println("MetaNode:")
	fmt.Printf("  ID: %v\n", m.formatNodeId())
	fmt.Printf("  Name: %d\n", m.name)
	fmt.Printf("  Completed: %t\n", m.completed)
	fmt.Printf("  Is Final: %t\n", m.isFinal)

	fmt.Println("  Productions:")
	for _, p := range m.productions {
		fmt.Printf("    - ID: %d, IsRoot: %t, Completed: %t, Index: %d\n", p.id, p.isRoot, p.completed, p.index)
	}

	fmt.Println("  Transitions:")
	for symbol, nodes := range m.Transitions {
		fmt.Printf("    - Symbol: %d -> %v\n", symbol, formatNodeId(nodes))
	}
	fmt.Println()
}

// Helper to format metaNodeId
func (m *metaNode) formatNodeId() string {
	return formatNodeId(m.id)
}

// Helper to format metaNodeId maps
func formatNodeId(nodeId metaNodeId) string {
	var ids []string
	for k := range nodeId {
		ids = append(ids, fmt.Sprintf("%d", k))
	}
	return "{" + strings.Join(ids, ", ") + "}"
}

type metaProduction struct {
	id        int
	isRoot    bool
	completed bool
	index     int
}

func setsEqual[T comparable](set1, set2 map[T]struct{}) bool {
	if len(set1) != len(set2) {
		return false
	}
	for key := range set1 {
		if _, exists := set2[key]; !exists {
			return false
		}
	}
	return true
}
