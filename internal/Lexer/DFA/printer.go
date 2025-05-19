package dfa

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"

	"github.com/olekukonko/tablewriter"
)

func intSliceToString(slice []int) string {
	strs := make([]string, len(slice))
	for i, v := range slice {
		strs[i] = strconv.Itoa(v)
	}
	return strings.Join(strs, ", ")
}

func printPositionTable(table map[int]positionTableRow) {
	header := []string{"Key", "Token", "Nullable", "IsFinal", "FirstPos", "LastPos", "FollowPos", "Actions"}

	data := make([][]string, 0, len(table))
	for key, row := range table {
		actions := fmt.Sprintf("{%d, %s}", row.action.Priority, row.action.Code)
		data = append(data, []string{
			fmt.Sprintf("%d", key),
			row.token,
			fmt.Sprintf("%t", row.nullable),
			fmt.Sprintf("%t", row.isFinal),
			intSliceToString(row.firstPos),
			intSliceToString(row.lastPos),
			intSliceToString(row.followPos),
			actions,
		})
	}

	tableWriter := tablewriter.NewWriter(os.Stdout)
	tableWriter.Header(header)
	tableWriter.Bulk(data)
	tableWriter.Render()
}

func printStateSetTable(states []*nodeSet, transitionTokens []string) {
	// Define the header
	header := []string{"ID", "Value", "isFinal", "Actions"}
	header = append(header, transitionTokens...)

	// Prepare the data for tablewriter
	data := make([][]string, 0, len(states))
	for _, state := range states {
		// Convert value slice to a comma-separated string
		valueStr := strings.Trim(strings.Join(strings.Fields(fmt.Sprint(state.value)), ","), "[]")
		actions := strings.Trim(strings.Join(strings.Fields(fmt.Sprint(state.actions)), ","), "[]")

		// Initialize the row with ID, Value, isFinal, and Actions
		row := []string{
			fmt.Sprintf("%d", state.id),
			valueStr,
			fmt.Sprintf("%t", state.isFinal),
			actions,
		}

		// Fill in the transitions for each token
		for _, token := range transitionTokens {
			if nextState, exists := state.transitions[token]; exists {
				row = append(row, fmt.Sprintf("%d", nextState.id))
			} else {
				row = append(row, "-")
			}
		}

		// Append the row to the data slice
		data = append(data, row)
	}

	// Create and render the table
	tableWriter := tablewriter.NewWriter(os.Stdout)
	tableWriter.Header(header)
	tableWriter.Bulk(data)
	tableWriter.Render()
}

func PrintDFA(dfa *DFA) {
	fmt.Println("DFA Representation:")
	fmt.Println("===================")
	fmt.Printf("Start State: %s\n\n", dfa.StartState.Id)

	for _, state := range dfa.States {
		fmt.Printf("State: %s\n", state.Id)
		if state.IsFinal {
			fmt.Println("  [Final State]")
		}
		if len(state.Actions) > 0 {
			fmt.Println("  Actions:")
			for _, action := range state.Actions {
				fmt.Printf("    - Code: %s (Priority: %d)\n", action.Code, action.Priority)
			}
		}
		if len(state.Transitions) > 0 {
			fmt.Println("  Transitions:")
			for symbol, target := range state.Transitions {
				fmt.Printf("    - %s -> %s\n", symbol, target.Id)
			}
		}
		fmt.Println(strings.Repeat("-", 25))
	}
}

// GenerateDOTFromRoot creates a DOT graph from a root Node and saves it as an image
func RenderAST(root node, outputPath string) error {
	// Generate the DOT representation
	dot := GenerateDOT_AST(root)

	// Print the DOT representation (for debugging purposes)
	// fmt.Println(dot)

	// Generate the image from the DOT representation
	return GenerateImage(dot, outputPath)
}

func RenderDFA(dfa *DFA, filename string) error {
	DOT := GenerateDOT_DFA(dfa)
	err := GenerateImage(DOT, filename)
	return err
}

// GenerateDOT_AST generates the DOT representation of the AST
func GenerateDOT_AST(root node) string {
	var buf bytes.Buffer
	buf.WriteString("digraph AST {\n")

	var addNode func(node, string) string
	nodeCount := 0

	addNode = func(n node, parentID string) string {
		nodeID := fmt.Sprintf("node%d", nodeCount)
		nodeCount++
		nodeLabel := strings.ReplaceAll(n.Value, "\"", "\\\"")

		buf.WriteString(fmt.Sprintf("  %s [label=\"%s\"];\n", nodeID, nodeLabel))

		if parentID != "" {
			buf.WriteString(fmt.Sprintf("  %s -> %s;\n", parentID, nodeID))
		}

		if n.IsOperator {
			for _, operand := range n.Children {
				addNode(operand, nodeID)
			}
		}

		return nodeID
	}

	addNode(root, "")
	buf.WriteString("}\n")
	return buf.String()
}

// GenerateDOT generates a DOT representation of a DFA as a string.
func GenerateDOT_DFA(dfa *DFA) string {
	var sb strings.Builder

	// Write the Graphviz dot header
	sb.WriteString("digraph DFA {\n")
	sb.WriteString("    rankdir=LR;\n") // Left to right orientation

	// Check if the DFA has any states
	if len(dfa.States) == 0 {
		panic("DFA has no states defined.")
	}

	// Define the nodes (states)
	for _, state := range dfa.States {
		shape := "circle"
		if state.IsFinal {
			shape = "doublecircle"
		}
		sb.WriteString(fmt.Sprintf("    \"%s\" [shape=%s];\n", state.Id, shape))

		// Define the transitions

		for symbol, toState := range state.Transitions {
			sb.WriteString(fmt.Sprintf("    \"%s\" -> \"%s\" [label=\"%s\"];\n",
				state.Id, toState.Id, symbol))
		}

	}

	// Define the start state
	sb.WriteString(fmt.Sprintf("    \"\" [shape=plaintext,label=\"\"];\n"))
	sb.WriteString(fmt.Sprintf("    \"\" -> \"%s\";\n", dfa.StartState.Id))

	sb.WriteString("}\n")

	return sb.String()
}

// getShape returns the shape for the state node based on whether it's a final state.
func getShape(isFinal bool) string {
	if isFinal {
		return "doublecircle"
	}
	return "circle"
}

// GenerateImage generates an image from the DOT representation using Graphviz
func GenerateImage(dot string, outputPath string) error {
	cmd := exec.Command("dot", "-Tpng", "-o", outputPath)
	cmd.Stdin = strings.NewReader(dot)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}
