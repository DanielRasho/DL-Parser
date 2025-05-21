package transitiontable

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
