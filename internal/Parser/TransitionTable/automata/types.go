package automata

import parser "github.com/DanielRasho/Parser/internal/Parser"

type Automata struct {
	StartState *State
	States     []*State
}

type Symbol = int

type State struct {
	Id          string
	Productions []parser.ParserProduction // Sorted by highest too lower priority ( 0 has the hightes priority )
	Transitions map[Symbol]*State         // {"a": STATE1, "b": STATE2, "NUMBER": STATEFINAL}
	IsFinal     bool                      // es el final final el que tiene el signo de dolar
	IsAccepted  bool
}
