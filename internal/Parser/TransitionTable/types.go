package transitiontable

// Transition Table
// NOTE: ADD THE $ as a new terminal that works as sentinel
type TransitionTbl = map[int]TransitionTblRow
type TransitionTblRow = map[string]Movement

// Goto table
type GotoTbl = map[int]GotoTblRow
type GotoTblRow = map[string]Movement

// Movements
type Movement struct {
	MovementType int
	NextRow      string
}

type MovementType = int

const (
	SHIFT MovementType = iota
	REDUCE
	GOTO
	ACCEPT
)
