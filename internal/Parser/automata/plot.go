package automata

import (
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

// GenerateDOT generates a DOT representation of an SLR(0) Automata as a string.
func GenerateDOT_SLR0(automata *Automata) string {
	var sb strings.Builder

	// Write the Graphviz dot header
	sb.WriteString("digraph SLR0_Automata {\n")
	sb.WriteString("    rankdir=LR;\n") // Left to right orientation

	// Define the nodes (states)
	for _, state := range automata.States {
		shape := getShape(state.IsFinal, state.IsAccepted)
		label := strconv.Itoa(state.Id)

		// Append productions to label if any
		if len(state.Productions) > 0 {
			label += "\\n"
			for _, prod := range state.Productions {
				label += fmt.Sprintf("%s → ", prod.Head.Value)
				for i, sym := range prod.Body {
					if i > 0 {
						label += " "
					}
					label += sym.Value
				}
				label += "\\n"
			}
		}

		sb.WriteString(fmt.Sprintf("    \"%d\" [label=\"%s\", shape=%s];\n", state.Id, label, shape))
	}

	// Define the transitions
	for _, state := range automata.States {
		for symbol, toState := range state.Transitions {
			sb.WriteString(fmt.Sprintf("    \"%d\" -> \"%d\" [label=\"%s\"];\n",
				state.Id, toState.Id, symbol.Value))
		}
	}

	// Define the start state
	sb.WriteString(fmt.Sprintln("    \"\" [shape=plaintext,label=\"\"];"))
	sb.WriteString(fmt.Sprintf("    \"\" -> \"%d\";\n", automata.StartState.Id))

	sb.WriteString("}\n")

	return sb.String()
}

// getShape returns the shape for the state node based on its properties.
func getShape(isFinal bool, isAccepted bool) string {
	if isAccepted {
		return "doubleoctagon"
	}
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
