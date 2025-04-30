package transitiontable

// Transition Table
// NOTE: ADD THE $ as a new terminal that works as sentinel
type TransitionTbl = map[int]TransitionTblRow
type TransitionTblRow = map[int]Movement

// Goto table
type GotoTbl = map[int]GotoTblRow
type GotoTblRow = map[int]Movement

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
