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

// pseudo-automata. Its an intermediate representation for the SLR automata.
type metaAutomata struct {
	StartState *metaNode
	States     []*metaNode
}

// ID for metaNode, represents a set of production's id
type metaNodeId = map[metaProductionId]struct{}

type metaProductionId struct {
	originalId int
	index      int
}

// Representation of automata node, using for intermediate computation steps
// for its ligh
type metaNode struct {
	// The ID is denoted by a set of ID of productions.
	// To check if to nodes are different just check that its sets are different
	// rather that comparing each production symbol by symbol.
	id metaNodeId
	// Will become the final automata name
	name      int
	metaProds []metaProduction
	// if the node contains root production completed scanned.
	completed bool
	// if the noded contains the root production
	isFinal     bool
	transitions map[Symbol]metaNodeId
}

// Representation of pseudo production. Rather than contain a production
// body or head, it store the id of real production
type metaProduction struct {
	// Id of a real production
	id metaProductionId
	// If production is root (special production added during computation)
	isRoot bool
	// If this production has already scanned completely
	completed bool
	// Producton's body length
	length int
}

func (p *metaProduction) getIndex() int {
	return p.id.index
}
func (p *metaProduction) getDictIndex() int {
	return p.id.originalId
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

// PrettyPrint prints the metaNode in a readable format
func (m *metaNode) print() {
	fmt.Println("MetaNode:")
	fmt.Printf("  ID: %v\n", m.formatNodeId())
	fmt.Printf("  Name: %d\n", m.name)
	fmt.Printf("  Completed: %t\n", m.completed)
	fmt.Printf("  Is Final: %t\n", m.isFinal)

	fmt.Println("  Productions:")
	for _, p := range m.metaProds {
		fmt.Printf("    - ID: %d, IsRoot: %t, Completed: %t, Index: (%d/%d)\n", p.id, p.isRoot, p.completed, p.id.index, p.length)
	}

	fmt.Println("  Transitions:")
	for symbol, nodes := range m.transitions {
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

func prettyPrintSymbols(symbols []*parser.ParserSymbol) string {
	var sb strings.Builder
	sb.WriteString("Parser Symbols:\n")
	for _, sym := range symbols {
		if sym.IsTerminal {
			sb.WriteString(fmt.Sprintf("  [T] Id: %d, Value: %s\n", sym.Id, sym.Value))
		} else {
			sb.WriteString(fmt.Sprintf("  [NT] Id: %d, Value: %s\n", sym.Id, sym.Value))
		}
	}
	return sb.String()
}

func prettyPrintMetaProductions(productions []metaProduction) string {
	var sb strings.Builder
	sb.WriteString("Meta Productions:\n")
	for _, prod := range productions {
		sb.WriteString(fmt.Sprintf(
			"  [ID: %d] Root: %t, Completed: %t, Index: %d/%d\n",
			prod.id, prod.isRoot, prod.completed, prod.id.index, prod.length,
		))
	}
	return sb.String()
}
