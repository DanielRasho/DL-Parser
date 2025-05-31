package transitiontable

import (
	"fmt"
	"os"
	"sort"
	"strconv"

	"github.com/olekukonko/tablewriter"
)

// Transition Table
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

func PrintMovementTable(title string, tbl map[string]map[string]Movement) {
	// Step 1: Collect all unique column names
	columnSet := make(map[string]struct{})
	for _, row := range tbl {
		for col := range row {
			columnSet[col] = struct{}{}
		}
	}

	// Convert map to sorted slice
	var columns []string
	for col := range columnSet {
		columns = append(columns, col)
	}
	sort.Strings(columns)

	// Step 2: Prepare sorted row keys
	var rowKeys []string
	for row := range tbl {
		rowKeys = append(rowKeys, row)
	}
	sort.Strings(rowKeys)

	// Step 3: Setup tablewriter
	table := tablewriter.NewWriter(os.Stdout)
	header := append([]string{"State"}, columns...)
	table.Header(header)

	// Step 4: Fill rows
	for _, rowKey := range rowKeys {
		row := []string{rowKey}
		for _, col := range columns {
			if move, ok := tbl[rowKey][col]; ok {
				row = append(row, movementToString(move))
			} else {
				row = append(row, "")
			}
		}
		table.Append(row)
	}

	// Step 5: Print title and table
	fmt.Println("=== " + title + " ===")
	table.Render()
}

// Helper to convert Movement to readable string
func movementToString(m Movement) string {
	switch m.MovementType {
	case SHIFT:
		return "s" + strconv.Itoa(m.NextRow)
	case REDUCE:
		return "r" + strconv.Itoa(m.NextRow)
	case GOTO:
		return strconv.Itoa(m.NextRow)
	case ACCEPT:
		return "acc"
	default:
		return "?"
	}
}
