// package should be specified after file generation
package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"
)

// =====================
//	  HEADER
// =====================
// Contains the exact same content defined on the Yaaalex file
// Tokens IDs should be defined here.

    // The entire contents of this section will be copied to the beginning of the generated Lexer.go file

    //  ------ TOKENS ID -----

    // Define the token types that the lexer will recognize

    const (

        IF = iota

        ELSE

        WHILE

        RETURN 

        ASSIGN

        PLUS

        MINUS

        MULT

        DIV

        LPAREN

        RPAREN

        LBRACE

        RBRACE

        ID

        NUMBER

        WS

        FUNC

    )



// =====================
//	  Lexer
// =====================

const NO_LEXEME = -1 // Flag constant that is used when no lexeme is recognized nor 
const SKIP_LEXEME = -2 // Flag when an action require the lexer to IGNORE the current lexeme

// PatternNotFound represents an error when a pattern is not found in a file
type PatternNotFound struct {
	Line    int
	Column  int
	Pattern string
}

// Error implements the error interface for PatternNotFound
func (e *PatternNotFound) Error() string {
	return fmt.Sprintf("error line %d column %d \n\tpattern not found. current pattern not recognized by the language: %s",
		e.Line,
		e.Column,
		e.Pattern)
}

type Symbol = string

// Definition of a Lexer
type Lexer struct {
	file         *os.File        // File to read from
	reader       *bufio.Reader   // Reader to get the symbols from file
	automata     dfa             // Automata for lexeme recognition
	symbolBuffer strings.Builder // Buffer to store the symbols of the current lexeme
	bytesRead    int             // Number of bytes the lexer has read
}

// Represents a piece of information withing the file
type Token struct {
	Value   Symbol // Actual string read by the lexer
	TokenID int    // Token Id (defined by the user above)
	Offset  int    // No of bytes from the start of the file to the current lexeme
}

// Converts the string to a human readable version
func (t *Token) String() string {
	return fmt.Sprintf("{ID: %d, OFFSET: %d ,VALUE: %s}", t.TokenID, t.Offset, t.Value)
}

// Creates a new Lexer that reads from a given path. Return error if cant open file.
func NewLexer(filePath string) (*Lexer, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	return &Lexer{
		file:         file,
		reader:       bufio.NewReader(file),
		automata:     *createDFA(),
		symbolBuffer: strings.Builder{}}, nil
}

// Close, closes the file that was being read by the Lexer.
func (l *Lexer) Close() {
	l.file.Close()
}

// GetNextToken return the next larger token that can find within the file
// starting from the last position it was left.
func (l *Lexer) GetNextToken() (Token, error) {

	// For every new lexeme we start an initial configurations
	lastTokenID := NO_LEXEME
	currentState := l.automata.startState
	lexemeBytesSize := 0 // Lenght of current lexeme in bytes.

	for {
		// fmt.Println(currentState.id)
		// 1. First check if in the current state there are any possible actions
		if actions := currentState.actions; len(currentState.actions) > 0 {
			newTokenID := actions[0]() // Get action with higher priority
			if newTokenID == SKIP_LEXEME {
				currentState = l.automata.startState
				l.bytesRead += lexemeBytesSize
				lexemeBytesSize = 0
				l.symbolBuffer.Reset()
				continue
			} else {
				// If TokenID returned, update lastToken read.
				lastTokenID = newTokenID
			}
		}
		// 2. Read the next rune 
		r, size, err := l.reader.ReadRune()
		if err != nil {
			// return the last recognized lexeme
			if lastTokenID != NO_LEXEME {
				break
			}
			// If no lexeme hast been recognized after endint the file, the file has invalid lexemes.
			return Token{}, err
		}

		nextState, ok := currentState.transitions[string(r)]

		// 3. Check if exist another state to jump to
		if !ok && lastTokenID == NO_LEXEME {
			l.symbolBuffer.WriteRune(r)
			line, columns, _ := l.getLineAndColumn(l.bytesRead)
			return Token{}, &PatternNotFound{Line: line, Column: columns, Pattern: l.symbolBuffer.String()}
		} else if !ok {
			l.reader.UnreadRune()
			break
		}

		// 4. update state
		l.symbolBuffer.WriteRune(r)
		lexemeBytesSize += size
		currentState = nextState
	}

	// 5. Build recognized token
	offset := l.bytesRead
	token := Token{
		TokenID: lastTokenID,
		Value:   l.symbolBuffer.String(),
		Offset:  offset,
	}
	l.symbolBuffer.Reset()
	l.bytesRead += lexemeBytesSize

	return token, nil
}

// getLineAndColumn takes an open file and an offset (in bytes),
// and returns the line and column where that byte is located.
func (l *Lexer) getLineAndColumn(offset int) (line, column int, err error) {

	// Reset file position to the beginning (because the lexer reader moved the file cursor previously)
	_, err = l.file.Seek(0, io.SeekStart)
	if err != nil {
		return 0, 0, err
	}

	// Create a buffered reader from the open file
	reader := bufio.NewReader(l.file)

	var currentByte int = 0
	line = 1
	column = 1

	// Read byte-by-byte
	for {
		// Read one byte at a time
		byteRead, err := reader.ReadByte()
		if err != nil && err.Error() != "EOF" {
			return 0, 0, err
		}

		// If we've read the required byte offset, stop and return the position
		if currentByte == offset {
			return line, column, nil
		}

		// Increment byte offset
		currentByte++

		// If the byte is a newline, increment line and reset column
		if byteRead == '\n' {
			line++
			column = 1
		} else {
			column++
		}

		// If we've reached the end of the file, break
		if err != nil {
			break
		}
	}

	return 0, 0, fmt.Errorf("Offset exceeds the number of bytes in the file")
}

// =====================
//	  DFA
// =====================

type dfa struct {
	startState *state
	states     []*state
}

type state struct {
	id          string
	actions     []action          // Sorted by highest too lower priority ( 0 has the hightes priority )
	transitions map[Symbol]*state // {"a": STATE1, "b": STATE2, "NUMBER": STATEFINAL}
	isFinal     bool
}

// Representes a user defined action that should happen
// when a pattern is recognized. The function should return an int, that represents a 
// tokenID. Its shape should be look something like : 
// 
// 	func () int {
// 		tokenID := SKIP_LEXEM
//		<user defined code>
//		return tokenID
//  }
//
type action func() int

// createDFA constructs the DFA that recognizes the user language.
func createDFA() *dfa {
	state16 := &state{id: "16" , transitions: make(map[Symbol]*state), isFinal: false}
state17 := &state{id: "17" , 
actions: []action{ 
 func() int { return IF 
return SKIP_LEXEME } , 
}, transitions: make(map[Symbol]*state), isFinal: false}
state19 := &state{id: "19" , transitions: make(map[Symbol]*state), isFinal: false}
state20 := &state{id: "20" , transitions: make(map[Symbol]*state), isFinal: false}
state5 := &state{id: "5" , 
actions: []action{ 
 func() int { return RPAREN
return SKIP_LEXEME } , 
}, transitions: make(map[Symbol]*state), isFinal: false}
state6 := &state{id: "6" , 
actions: []action{ 
 func() int { return ID 
return SKIP_LEXEME } , 
}, transitions: make(map[Symbol]*state), isFinal: false}
state24 := &state{id: "24" , transitions: make(map[Symbol]*state), isFinal: false}
state28 := &state{id: "28" , transitions: make(map[Symbol]*state), isFinal: false}
state30 := &state{id: "30" , transitions: make(map[Symbol]*state), isFinal: false}
state31 := &state{id: "31" , 
actions: []action{ 
 func() int { return WHILE 
return SKIP_LEXEME } , 
}, transitions: make(map[Symbol]*state), isFinal: false}
state25 := &state{id: "25" , transitions: make(map[Symbol]*state), isFinal: false}
state4 := &state{id: "4" , 
actions: []action{ 
 func() int { return MULT 
return SKIP_LEXEME } , 
}, transitions: make(map[Symbol]*state), isFinal: false}
state18 := &state{id: "18" , transitions: make(map[Symbol]*state), isFinal: true}
state29 := &state{id: "29" , 
actions: []action{ 
 func() int { return ELSE 
return SKIP_LEXEME } , 
}, transitions: make(map[Symbol]*state), isFinal: false}
state3 := &state{id: "3" , 
actions: []action{ 
 func() int { return MINUS 
return SKIP_LEXEME } , 
}, transitions: make(map[Symbol]*state), isFinal: false}
state11 := &state{id: "11" , 
actions: []action{ 
 func() int { return DIV 
return SKIP_LEXEME } , 
}, transitions: make(map[Symbol]*state), isFinal: false}
state15 := &state{id: "15" , 
actions: []action{ 
 func() int {
return SKIP_LEXEME } , 
}, transitions: make(map[Symbol]*state), isFinal: false}
state22 := &state{id: "22" , transitions: make(map[Symbol]*state), isFinal: false}
state23 := &state{id: "23" , transitions: make(map[Symbol]*state), isFinal: false}
state32 := &state{id: "32" , transitions: make(map[Symbol]*state), isFinal: false}
state0 := &state{id: "0" , 
actions: []action{ 
}, transitions: make(map[Symbol]*state), isFinal: false}
state7 := &state{id: "7" , 
actions: []action{ 
 func() int { return LPAREN
return SKIP_LEXEME } , 
}, transitions: make(map[Symbol]*state), isFinal: false}
state8 := &state{id: "8" , 
actions: []action{ 
 func() int { return NUMBER 
return SKIP_LEXEME } , 
}, transitions: make(map[Symbol]*state), isFinal: false}
state21 := &state{id: "21" , transitions: make(map[Symbol]*state), isFinal: false}
state2 := &state{id: "2" , transitions: make(map[Symbol]*state), isFinal: false}
state12 := &state{id: "12" , transitions: make(map[Symbol]*state), isFinal: false}
state14 := &state{id: "14" , 
actions: []action{ 
 func() int { return PLUS 
return SKIP_LEXEME } , 
}, transitions: make(map[Symbol]*state), isFinal: false}
state26 := &state{id: "26" , transitions: make(map[Symbol]*state), isFinal: false}
state27 := &state{id: "27" , 
actions: []action{ 
 func() int { return FUNC 
return SKIP_LEXEME } , 
}, transitions: make(map[Symbol]*state), isFinal: false}
state33 := &state{id: "33" , 
actions: []action{ 
 func() int { return RETURN 
return SKIP_LEXEME } , 
}, transitions: make(map[Symbol]*state), isFinal: false}
state13 := &state{id: "13" , transitions: make(map[Symbol]*state), isFinal: false}
state1 := &state{id: "1" , transitions: make(map[Symbol]*state), isFinal: false}
state9 := &state{id: "9" , transitions: make(map[Symbol]*state), isFinal: false}
state10 := &state{id: "10" , 
actions: []action{ 
 func() int { return ASSIGN 
return SKIP_LEXEME } , 
}, transitions: make(map[Symbol]*state), isFinal: false}

state16.transitions["20"] = state2
state16.transitions["c"] = state2
state16.transitions["h"] = state2
state16.transitions["18"] = state2
state16.transitions["n"] = state2
state16.transitions["B"] = state2
state16.transitions["15"] = state2
state16.transitions["b"] = state2
state16.transitions["+"] = state2
state16.transitions["	"] = state2
state16.transitions["0"] = state2
state16.transitions["="] = state2
state16.transitions["A"] = state2
state16.transitions["e"] = state22
state16.transitions["2"] = state2
state16.transitions["16"] = state2
state16.transitions["t"] = state2
state16.transitions["r"] = state2
state16.transitions["14"] = state2
state16.transitions["l"] = state2
state16.transitions["*"] = state2
state16.transitions[")"] = state2
state16.transitions["23"] = state2
state16.transitions["s"] = state2
state16.transitions["1"] = state2
state16.transitions["19"] = state2
state16.transitions["21"] = state2
state16.transitions["("] = state2
state16.transitions["u"] = state2
state16.transitions["13"] = state2
state16.transitions["f"] = state2
state16.transitions["10"] = state2
state16.transitions["12"] = state2
state16.transitions["w"] = state2
state16.transitions["a"] = state2
state16.transitions["i"] = state2
state16.transitions["24"] = state2
state16.transitions["17"] = state2
state16.transitions["22"] = state2
state16.transitions["/"] = state2
state16.transitions[" "] = state2
state16.transitions["11"] = state2
state16.transitions["\n"] = state2
state16.transitions["-"] = state2
state17.transitions["n"] = state2
state17.transitions["t"] = state2
state17.transitions[" "] = state2
state17.transitions["-"] = state2
state17.transitions["h"] = state2
state17.transitions["17"] = state2
state17.transitions["+"] = state2
state17.transitions["r"] = state2
state17.transitions["	"] = state2
state17.transitions["24"] = state2
state17.transitions["b"] = state2
state17.transitions["23"] = state2
state17.transitions["1"] = state2
state17.transitions["w"] = state2
state17.transitions["21"] = state2
state17.transitions["11"] = state2
state17.transitions["22"] = state2
state17.transitions["/"] = state2
state17.transitions["a"] = state2
state17.transitions["B"] = state2
state17.transitions["\n"] = state2
state17.transitions["15"] = state2
state17.transitions["A"] = state2
state17.transitions[")"] = state2
state17.transitions["i"] = state2
state17.transitions["l"] = state2
state17.transitions["("] = state2
state17.transitions["u"] = state2
state17.transitions["13"] = state2
state17.transitions["f"] = state2
state17.transitions["19"] = state2
state17.transitions["18"] = state2
state17.transitions["="] = state2
state17.transitions["10"] = state18
state17.transitions["e"] = state2
state17.transitions["14"] = state2
state17.transitions["0"] = state2
state17.transitions["*"] = state2
state17.transitions["c"] = state2
state17.transitions["16"] = state2
state17.transitions["2"] = state2
state17.transitions["20"] = state2
state17.transitions["s"] = state2
state17.transitions["12"] = state2
state19.transitions["14"] = state2
state19.transitions["s"] = state2
state19.transitions["/"] = state2
state19.transitions["10"] = state2
state19.transitions["w"] = state2
state19.transitions[" "] = state2
state19.transitions["11"] = state2
state19.transitions["15"] = state2
state19.transitions["17"] = state2
state19.transitions["12"] = state2
state19.transitions["a"] = state2
state19.transitions["t"] = state2
state19.transitions["21"] = state2
state19.transitions["20"] = state2
state19.transitions["18"] = state2
state19.transitions["23"] = state2
state19.transitions["*"] = state2
state19.transitions["22"] = state2
state19.transitions["13"] = state2
state19.transitions["19"] = state2
state19.transitions["n"] = state23
state19.transitions["B"] = state2
state19.transitions["i"] = state2
state19.transitions["l"] = state2
state19.transitions["24"] = state2
state19.transitions["h"] = state2
state19.transitions["1"] = state2
state19.transitions["A"] = state2
state19.transitions["u"] = state2
state19.transitions["+"] = state2
state19.transitions["	"] = state2
state19.transitions["-"] = state2
state19.transitions["f"] = state2
state19.transitions["="] = state2
state19.transitions["r"] = state2
state19.transitions["2"] = state2
state19.transitions["\n"] = state2
state19.transitions["("] = state2
state19.transitions["c"] = state2
state19.transitions[")"] = state2
state19.transitions["0"] = state2
state19.transitions["b"] = state2
state19.transitions["16"] = state2
state19.transitions["e"] = state2
state20.transitions[")"] = state2
state20.transitions["i"] = state24
state20.transitions["*"] = state2
state20.transitions["("] = state2
state20.transitions["24"] = state2
state20.transitions["t"] = state2
state20.transitions["+"] = state2
state20.transitions["11"] = state2
state20.transitions["0"] = state2
state20.transitions["u"] = state2
state20.transitions["A"] = state2
state20.transitions["w"] = state2
state20.transitions["n"] = state2
state20.transitions["18"] = state2
state20.transitions["f"] = state2
state20.transitions["="] = state2
state20.transitions["a"] = state2
state20.transitions["2"] = state2
state20.transitions["-"] = state2
state20.transitions["b"] = state2
state20.transitions["1"] = state2
state20.transitions["12"] = state2
state20.transitions["r"] = state2
state20.transitions["14"] = state2
state20.transitions["l"] = state2
state20.transitions["17"] = state2
state20.transitions["23"] = state2
state20.transitions[" "] = state2
state20.transitions["B"] = state2
state20.transitions["h"] = state2
state20.transitions["13"] = state2
state20.transitions["s"] = state2
state20.transitions["16"] = state2
state20.transitions["15"] = state2
state20.transitions["20"] = state2
state20.transitions["22"] = state2
state20.transitions["/"] = state2
state20.transitions["19"] = state2
state20.transitions["e"] = state2
state20.transitions["21"] = state2
state20.transitions["\n"] = state2
state20.transitions["c"] = state2
state20.transitions["10"] = state2
state20.transitions["	"] = state2
state5.transitions["2"] = state2
state5.transitions["-"] = state2
state5.transitions["*"] = state2
state5.transitions["24"] = state2
state5.transitions["0"] = state2
state5.transitions["21"] = state2
state5.transitions["r"] = state2
state5.transitions["B"] = state2
state5.transitions["("] = state2
state5.transitions["u"] = state2
state5.transitions["16"] = state2
state5.transitions["n"] = state2
state5.transitions["i"] = state2
state5.transitions["f"] = state2
state5.transitions["w"] = state2
state5.transitions[" "] = state2
state5.transitions["14"] = state2
state5.transitions["	"] = state2
state5.transitions["20"] = state18
state5.transitions["h"] = state2
state5.transitions["11"] = state2
state5.transitions["1"] = state2
state5.transitions["10"] = state2
state5.transitions["l"] = state2
state5.transitions["17"] = state2
state5.transitions[")"] = state2
state5.transitions["22"] = state2
state5.transitions["b"] = state2
state5.transitions["13"] = state2
state5.transitions["15"] = state2
state5.transitions["18"] = state2
state5.transitions["="] = state2
state5.transitions["19"] = state2
state5.transitions["t"] = state2
state5.transitions["+"] = state2
state5.transitions["\n"] = state2
state5.transitions["23"] = state2
state5.transitions["s"] = state2
state5.transitions["12"] = state2
state5.transitions["c"] = state2
state5.transitions["/"] = state2
state5.transitions["A"] = state2
state5.transitions["a"] = state2
state5.transitions["e"] = state2
state6.transitions["22"] = state18
state6.transitions["u"] = state2
state6.transitions["b"] = state6
state6.transitions["14"] = state2
state6.transitions["B"] = state6
state6.transitions["\n"] = state2
state6.transitions["l"] = state2
state6.transitions[")"] = state2
state6.transitions["e"] = state2
state6.transitions["i"] = state2
state6.transitions["("] = state2
state6.transitions["20"] = state2
state6.transitions["c"] = state6
state6.transitions["0"] = state6
state6.transitions["1"] = state6
state6.transitions["A"] = state6
state6.transitions["2"] = state6
state6.transitions["*"] = state2
state6.transitions["13"] = state2
state6.transitions["	"] = state2
state6.transitions["17"] = state2
state6.transitions["18"] = state2
state6.transitions["/"] = state2
state6.transitions["19"] = state2
state6.transitions["a"] = state6
state6.transitions["n"] = state2
state6.transitions["11"] = state2
state6.transitions["15"] = state2
state6.transitions["24"] = state2
state6.transitions["16"] = state2
state6.transitions["10"] = state2
state6.transitions["+"] = state2
state6.transitions["21"] = state2
state6.transitions[" "] = state2
state6.transitions["-"] = state2
state6.transitions["h"] = state2
state6.transitions["23"] = state2
state6.transitions["f"] = state2
state6.transitions["t"] = state2
state6.transitions["s"] = state2
state6.transitions["="] = state2
state6.transitions["12"] = state2
state6.transitions["w"] = state2
state6.transitions["r"] = state2
state24.transitions[")"] = state2
state24.transitions["l"] = state28
state24.transitions["24"] = state2
state24.transitions["16"] = state2
state24.transitions["A"] = state2
state24.transitions["f"] = state2
state24.transitions["s"] = state2
state24.transitions["="] = state2
state24.transitions["t"] = state2
state24.transitions["11"] = state2
state24.transitions["-"] = state2
state24.transitions["17"] = state2
state24.transitions["18"] = state2
state24.transitions["23"] = state2
state24.transitions["a"] = state2
state24.transitions["	"] = state2
state24.transitions["*"] = state2
state24.transitions["h"] = state2
state24.transitions["0"] = state2
state24.transitions["("] = state2
state24.transitions["12"] = state2
state24.transitions["n"] = state2
state24.transitions[" "] = state2
state24.transitions["r"] = state2
state24.transitions["14"] = state2
state24.transitions["\n"] = state2
state24.transitions["i"] = state2
state24.transitions["15"] = state2
state24.transitions["u"] = state2
state24.transitions["/"] = state2
state24.transitions["+"] = state2
state24.transitions["21"] = state2
state24.transitions["B"] = state2
state24.transitions["2"] = state2
state24.transitions["20"] = state2
state24.transitions["c"] = state2
state24.transitions["22"] = state2
state24.transitions["b"] = state2
state24.transitions["13"] = state2
state24.transitions["10"] = state2
state24.transitions["w"] = state2
state24.transitions["e"] = state2
state24.transitions["1"] = state2
state24.transitions["19"] = state2
state28.transitions["13"] = state2
state28.transitions["f"] = state2
state28.transitions["	"] = state2
state28.transitions["24"] = state2
state28.transitions["w"] = state2
state28.transitions[" "] = state2
state28.transitions["14"] = state2
state28.transitions["B"] = state2
state28.transitions["c"] = state2
state28.transitions["h"] = state2
state28.transitions["18"] = state2
state28.transitions["1"] = state2
state28.transitions["A"] = state2
state28.transitions["10"] = state2
state28.transitions["12"] = state2
state28.transitions["s"] = state2
state28.transitions["e"] = state31
state28.transitions["15"] = state2
state28.transitions["20"] = state2
state28.transitions["t"] = state2
state28.transitions["21"] = state2
state28.transitions["\n"] = state2
state28.transitions["l"] = state2
state28.transitions["-"] = state2
state28.transitions["("] = state2
state28.transitions["0"] = state2
state28.transitions["="] = state2
state28.transitions["11"] = state2
state28.transitions["2"] = state2
state28.transitions[")"] = state2
state28.transitions["b"] = state2
state28.transitions["/"] = state2
state28.transitions["19"] = state2
state28.transitions["a"] = state2
state28.transitions["+"] = state2
state28.transitions["22"] = state2
state28.transitions["u"] = state2
state28.transitions["16"] = state2
state28.transitions["n"] = state2
state28.transitions["r"] = state2
state28.transitions["i"] = state2
state28.transitions["*"] = state2
state28.transitions["17"] = state2
state28.transitions["23"] = state2
state30.transitions["\n"] = state2
state30.transitions["/"] = state2
state30.transitions["l"] = state2
state30.transitions["*"] = state2
state30.transitions["17"] = state2
state30.transitions["18"] = state2
state30.transitions["b"] = state2
state30.transitions["s"] = state2
state30.transitions["-"] = state2
state30.transitions["("] = state2
state30.transitions["12"] = state2
state30.transitions["a"] = state2
state30.transitions["15"] = state2
state30.transitions["20"] = state2
state30.transitions["24"] = state2
state30.transitions["="] = state2
state30.transitions["A"] = state2
state30.transitions[" "] = state2
state30.transitions["14"] = state2
state30.transitions["B"] = state2
state30.transitions["c"] = state2
state30.transitions["22"] = state2
state30.transitions["11"] = state2
state30.transitions["1"] = state2
state30.transitions["h"] = state2
state30.transitions["23"] = state2
state30.transitions["19"] = state2
state30.transitions["e"] = state2
state30.transitions["t"] = state2
state30.transitions["+"] = state2
state30.transitions["21"] = state2
state30.transitions[")"] = state2
state30.transitions["0"] = state2
state30.transitions["f"] = state2
state30.transitions["10"] = state2
state30.transitions["r"] = state32
state30.transitions["2"] = state2
state30.transitions["	"] = state2
state30.transitions["16"] = state2
state30.transitions["n"] = state2
state30.transitions["i"] = state2
state30.transitions["u"] = state2
state30.transitions["13"] = state2
state30.transitions["w"] = state2
state31.transitions["*"] = state2
state31.transitions["b"] = state2
state31.transitions["A"] = state2
state31.transitions["w"] = state2
state31.transitions["e"] = state2
state31.transitions["l"] = state2
state31.transitions["24"] = state2
state31.transitions["u"] = state2
state31.transitions["="] = state2
state31.transitions["1"] = state2
state31.transitions["10"] = state2
state31.transitions["19"] = state2
state31.transitions["2"] = state2
state31.transitions["-"] = state2
state31.transitions["("] = state2
state31.transitions["n"] = state2
state31.transitions["+"] = state2
state31.transitions["h"] = state2
state31.transitions["23"] = state2
state31.transitions["s"] = state2
state31.transitions["11"] = state2
state31.transitions["14"] = state2
state31.transitions["20"] = state2
state31.transitions["17"] = state2
state31.transitions["22"] = state2
state31.transitions["13"] = state18
state31.transitions["\n"] = state2
state31.transitions["18"] = state2
state31.transitions["t"] = state2
state31.transitions["i"] = state2
state31.transitions["c"] = state2
state31.transitions[")"] = state2
state31.transitions["0"] = state2
state31.transitions["16"] = state2
state31.transitions["/"] = state2
state31.transitions["21"] = state2
state31.transitions["r"] = state2
state31.transitions["15"] = state2
state31.transitions["f"] = state2
state31.transitions["12"] = state2
state31.transitions["a"] = state2
state31.transitions[" "] = state2
state31.transitions["B"] = state2
state31.transitions["	"] = state2
state25.transitions["*"] = state2
state25.transitions["b"] = state2
state25.transitions["23"] = state2
state25.transitions["A"] = state2
state25.transitions["t"] = state2
state25.transitions["-"] = state2
state25.transitions["e"] = state29
state25.transitions["20"] = state2
state25.transitions["24"] = state2
state25.transitions["1"] = state2
state25.transitions["+"] = state2
state25.transitions["15"] = state2
state25.transitions["c"] = state2
state25.transitions["0"] = state2
state25.transitions["s"] = state2
state25.transitions["19"] = state2
state25.transitions["2"] = state2
state25.transitions["h"] = state2
state25.transitions["u"] = state2
state25.transitions["12"] = state2
state25.transitions["a"] = state2
state25.transitions["11"] = state2
state25.transitions["r"] = state2
state25.transitions["14"] = state2
state25.transitions["17"] = state2
state25.transitions["13"] = state2
state25.transitions["w"] = state2
state25.transitions["n"] = state2
state25.transitions["i"] = state2
state25.transitions["18"] = state2
state25.transitions[")"] = state2
state25.transitions["22"] = state2
state25.transitions["f"] = state2
state25.transitions["B"] = state2
state25.transitions["	"] = state2
state25.transitions["\n"] = state2
state25.transitions["l"] = state2
state25.transitions["("] = state2
state25.transitions["="] = state2
state25.transitions["16"] = state2
state25.transitions["/"] = state2
state25.transitions["10"] = state2
state25.transitions["21"] = state2
state25.transitions[" "] = state2
state4.transitions["i"] = state2
state4.transitions["="] = state2
state4.transitions["10"] = state2
state4.transitions["e"] = state2
state4.transitions[" "] = state2
state4.transitions["13"] = state2
state4.transitions["16"] = state2
state4.transitions["n"] = state2
state4.transitions["t"] = state2
state4.transitions["21"] = state2
state4.transitions["11"] = state2
state4.transitions["\n"] = state2
state4.transitions["l"] = state2
state4.transitions["15"] = state2
state4.transitions["20"] = state2
state4.transitions["18"] = state18
state4.transitions[")"] = state2
state4.transitions["2"] = state2
state4.transitions["	"] = state2
state4.transitions["("] = state2
state4.transitions["17"] = state2
state4.transitions["B"] = state2
state4.transitions["*"] = state2
state4.transitions["h"] = state2
state4.transitions["23"] = state2
state4.transitions["14"] = state2
state4.transitions["f"] = state2
state4.transitions["1"] = state2
state4.transitions["/"] = state2
state4.transitions["A"] = state2
state4.transitions["12"] = state2
state4.transitions["w"] = state2
state4.transitions["+"] = state2
state4.transitions["-"] = state2
state4.transitions["24"] = state2
state4.transitions["c"] = state2
state4.transitions["0"] = state2
state4.transitions["22"] = state2
state4.transitions["b"] = state2
state4.transitions["s"] = state2
state4.transitions["19"] = state2
state4.transitions["u"] = state2
state4.transitions["a"] = state2
state4.transitions["r"] = state2
state18.transitions["l"] = state2
state18.transitions["15"] = state2
state18.transitions["16"] = state2
state18.transitions["A"] = state2
state18.transitions["n"] = state2
state18.transitions["t"] = state2
state18.transitions["+"] = state2
state18.transitions["	"] = state2
state18.transitions["17"] = state2
state18.transitions["13"] = state2
state18.transitions["12"] = state2
state18.transitions["w"] = state2
state18.transitions["21"] = state2
state18.transitions["i"] = state2
state18.transitions["b"] = state2
state18.transitions["B"] = state2
state18.transitions["-"] = state2
state18.transitions["("] = state2
state18.transitions["20"] = state2
state18.transitions["0"] = state2
state18.transitions["f"] = state2
state18.transitions["s"] = state2
state18.transitions["/"] = state2
state18.transitions["19"] = state2
state18.transitions["h"] = state2
state18.transitions[")"] = state2
state18.transitions["22"] = state2
state18.transitions["10"] = state2
state18.transitions["a"] = state2
state18.transitions["11"] = state2
state18.transitions["2"] = state2
state18.transitions["\n"] = state2
state18.transitions["23"] = state2
state18.transitions["="] = state2
state18.transitions["1"] = state2
state18.transitions["*"] = state2
state18.transitions["24"] = state2
state18.transitions["c"] = state2
state18.transitions["18"] = state2
state18.transitions["u"] = state2
state18.transitions["e"] = state2
state18.transitions[" "] = state2
state18.transitions["r"] = state2
state18.transitions["14"] = state2
state29.transitions["r"] = state2
state29.transitions["2"] = state2
state29.transitions[")"] = state2
state29.transitions["u"] = state2
state29.transitions["b"] = state2
state29.transitions["n"] = state2
state29.transitions["+"] = state2
state29.transitions["B"] = state2
state29.transitions["*"] = state2
state29.transitions["("] = state2
state29.transitions["c"] = state2
state29.transitions["23"] = state2
state29.transitions["13"] = state2
state29.transitions["18"] = state2
state29.transitions["22"] = state2
state29.transitions["16"] = state2
state29.transitions["A"] = state2
state29.transitions["w"] = state2
state29.transitions["11"] = state2
state29.transitions["	"] = state2
state29.transitions["15"] = state2
state29.transitions["20"] = state2
state29.transitions["f"] = state2
state29.transitions["a"] = state2
state29.transitions["\n"] = state2
state29.transitions["t"] = state2
state29.transitions["/"] = state2
state29.transitions["12"] = state18
state29.transitions["l"] = state2
state29.transitions["-"] = state2
state29.transitions["h"] = state2
state29.transitions["19"] = state2
state29.transitions["e"] = state2
state29.transitions[" "] = state2
state29.transitions["i"] = state2
state29.transitions["24"] = state2
state29.transitions["17"] = state2
state29.transitions["0"] = state2
state29.transitions["="] = state2
state29.transitions["1"] = state2
state29.transitions["10"] = state2
state29.transitions["14"] = state2
state29.transitions["s"] = state2
state29.transitions["21"] = state2
state3.transitions["16"] = state2
state3.transitions["1"] = state2
state3.transitions["/"] = state2
state3.transitions["12"] = state2
state3.transitions["e"] = state2
state3.transitions[" "] = state2
state3.transitions["24"] = state2
state3.transitions["n"] = state2
state3.transitions["+"] = state2
state3.transitions["("] = state2
state3.transitions["23"] = state2
state3.transitions["t"] = state2
state3.transitions["\n"] = state2
state3.transitions["l"] = state2
state3.transitions["0"] = state2
state3.transitions["b"] = state2
state3.transitions["10"] = state2
state3.transitions["a"] = state2
state3.transitions["r"] = state2
state3.transitions["	"] = state2
state3.transitions["c"] = state2
state3.transitions["22"] = state2
state3.transitions["u"] = state2
state3.transitions["13"] = state2
state3.transitions["A"] = state2
state3.transitions["11"] = state2
state3.transitions["14"] = state2
state3.transitions["*"] = state2
state3.transitions["17"] = state18
state3.transitions["f"] = state2
state3.transitions["="] = state2
state3.transitions["i"] = state2
state3.transitions["20"] = state2
state3.transitions[")"] = state2
state3.transitions["s"] = state2
state3.transitions["w"] = state2
state3.transitions["2"] = state2
state3.transitions["15"] = state2
state3.transitions["-"] = state2
state3.transitions["19"] = state2
state3.transitions["21"] = state2
state3.transitions["B"] = state2
state3.transitions["h"] = state2
state3.transitions["18"] = state2
state11.transitions["18"] = state2
state11.transitions["t"] = state2
state11.transitions["2"] = state2
state11.transitions["\n"] = state2
state11.transitions["/"] = state2
state11.transitions["l"] = state2
state11.transitions["*"] = state2
state11.transitions["b"] = state2
state11.transitions["13"] = state2
state11.transitions["19"] = state18
state11.transitions["B"] = state2
state11.transitions["20"] = state2
state11.transitions["24"] = state2
state11.transitions["23"] = state2
state11.transitions["a"] = state2
state11.transitions["e"] = state2
state11.transitions[" "] = state2
state11.transitions["r"] = state2
state11.transitions["15"] = state2
state11.transitions["c"] = state2
state11.transitions["u"] = state2
state11.transitions["16"] = state2
state11.transitions["+"] = state2
state11.transitions["21"] = state2
state11.transitions["14"] = state2
state11.transitions["	"] = state2
state11.transitions["17"] = state2
state11.transitions["n"] = state2
state11.transitions["-"] = state2
state11.transitions["("] = state2
state11.transitions["h"] = state2
state11.transitions["22"] = state2
state11.transitions["1"] = state2
state11.transitions["w"] = state2
state11.transitions["A"] = state2
state11.transitions["12"] = state2
state11.transitions[")"] = state2
state11.transitions["0"] = state2
state11.transitions["f"] = state2
state11.transitions["="] = state2
state11.transitions["i"] = state2
state11.transitions["s"] = state2
state11.transitions["10"] = state2
state11.transitions["11"] = state2
state15.transitions["B"] = state2
state15.transitions["	"] = state15
state15.transitions["i"] = state2
state15.transitions["("] = state2
state15.transitions["13"] = state2
state15.transitions["A"] = state2
state15.transitions["a"] = state2
state15.transitions["2"] = state2
state15.transitions[")"] = state2
state15.transitions["16"] = state2
state15.transitions["18"] = state2
state15.transitions["-"] = state2
state15.transitions["24"] = state18
state15.transitions["22"] = state2
state15.transitions["1"] = state2
state15.transitions["t"] = state2
state15.transitions["11"] = state2
state15.transitions["r"] = state2
state15.transitions["20"] = state2
state15.transitions["h"] = state2
state15.transitions["12"] = state2
state15.transitions["+"] = state2
state15.transitions["\n"] = state15
state15.transitions["15"] = state2
state15.transitions["*"] = state2
state15.transitions["c"] = state2
state15.transitions["14"] = state2
state15.transitions["0"] = state2
state15.transitions["b"] = state2
state15.transitions["f"] = state2
state15.transitions["s"] = state2
state15.transitions["="] = state2
state15.transitions["w"] = state2
state15.transitions["n"] = state2
state15.transitions[" "] = state15
state15.transitions["17"] = state2
state15.transitions["/"] = state2
state15.transitions["10"] = state2
state15.transitions["19"] = state2
state15.transitions["e"] = state2
state15.transitions["l"] = state2
state15.transitions["u"] = state2
state15.transitions["23"] = state2
state15.transitions["21"] = state2
state22.transitions["18"] = state2
state22.transitions["w"] = state2
state22.transitions["14"] = state2
state22.transitions["l"] = state2
state22.transitions["17"] = state2
state22.transitions["a"] = state2
state22.transitions["t"] = state26
state22.transitions["\n"] = state2
state22.transitions["24"] = state2
state22.transitions["b"] = state2
state22.transitions["s"] = state2
state22.transitions["B"] = state2
state22.transitions["i"] = state2
state22.transitions["13"] = state2
state22.transitions["1"] = state2
state22.transitions["12"] = state2
state22.transitions["15"] = state2
state22.transitions["-"] = state2
state22.transitions["*"] = state2
state22.transitions["16"] = state2
state22.transitions["A"] = state2
state22.transitions[" "] = state2
state22.transitions["11"] = state2
state22.transitions["("] = state2
state22.transitions["c"] = state2
state22.transitions["h"] = state2
state22.transitions[")"] = state2
state22.transitions["10"] = state2
state22.transitions["19"] = state2
state22.transitions["+"] = state2
state22.transitions["	"] = state2
state22.transitions["0"] = state2
state22.transitions["22"] = state2
state22.transitions["u"] = state2
state22.transitions["="] = state2
state22.transitions["/"] = state2
state22.transitions["n"] = state2
state22.transitions["r"] = state2
state22.transitions["23"] = state2
state22.transitions["f"] = state2
state22.transitions["e"] = state2
state22.transitions["21"] = state2
state22.transitions["2"] = state2
state22.transitions["20"] = state2
state23.transitions["c"] = state27
state23.transitions["0"] = state2
state23.transitions["16"] = state2
state23.transitions["1"] = state2
state23.transitions["i"] = state2
state23.transitions["("] = state2
state23.transitions["="] = state2
state23.transitions["14"] = state2
state23.transitions["B"] = state2
state23.transitions["l"] = state2
state23.transitions["f"] = state2
state23.transitions["w"] = state2
state23.transitions["r"] = state2
state23.transitions["\n"] = state2
state23.transitions["20"] = state2
state23.transitions["b"] = state2
state23.transitions["23"] = state2
state23.transitions["13"] = state2
state23.transitions["s"] = state2
state23.transitions["n"] = state2
state23.transitions[" "] = state2
state23.transitions["	"] = state2
state23.transitions["18"] = state2
state23.transitions["/"] = state2
state23.transitions["t"] = state2
state23.transitions["2"] = state2
state23.transitions["-"] = state2
state23.transitions[")"] = state2
state23.transitions["u"] = state2
state23.transitions["e"] = state2
state23.transitions["+"] = state2
state23.transitions["17"] = state2
state23.transitions["22"] = state2
state23.transitions["12"] = state2
state23.transitions["19"] = state2
state23.transitions["11"] = state2
state23.transitions["15"] = state2
state23.transitions["h"] = state2
state23.transitions["a"] = state2
state23.transitions["21"] = state2
state23.transitions["A"] = state2
state23.transitions["10"] = state2
state23.transitions["*"] = state2
state23.transitions["24"] = state2
state32.transitions["11"] = state2
state32.transitions["2"] = state2
state32.transitions["10"] = state2
state32.transitions["0"] = state2
state32.transitions["l"] = state2
state32.transitions["18"] = state2
state32.transitions["="] = state2
state32.transitions["16"] = state2
state32.transitions["1"] = state2
state32.transitions[" "] = state2
state32.transitions[")"] = state2
state32.transitions["a"] = state2
state32.transitions["e"] = state2
state32.transitions["+"] = state2
state32.transitions["	"] = state2
state32.transitions["13"] = state2
state32.transitions["A"] = state2
state32.transitions["n"] = state33
state32.transitions["12"] = state2
state32.transitions["s"] = state2
state32.transitions["r"] = state2
state32.transitions["14"] = state2
state32.transitions["i"] = state2
state32.transitions["b"] = state2
state32.transitions["t"] = state2
state32.transitions["B"] = state2
state32.transitions["23"] = state2
state32.transitions["19"] = state2
state32.transitions["21"] = state2
state32.transitions["\n"] = state2
state32.transitions["-"] = state2
state32.transitions["*"] = state2
state32.transitions["20"] = state2
state32.transitions["24"] = state2
state32.transitions["c"] = state2
state32.transitions["22"] = state2
state32.transitions["/"] = state2
state32.transitions["w"] = state2
state32.transitions["15"] = state2
state32.transitions["("] = state2
state32.transitions["h"] = state2
state32.transitions["17"] = state2
state32.transitions["u"] = state2
state32.transitions["f"] = state2
state0.transitions["B"] = state6
state0.transitions["l"] = state2
state0.transitions["*"] = state4
state0.transitions["u"] = state2
state0.transitions["1"] = state8
state0.transitions["a"] = state6
state0.transitions["20"] = state2
state0.transitions["17"] = state2
state0.transitions["11"] = state2
state0.transitions["r"] = state16
state0.transitions["-"] = state3
state0.transitions["24"] = state2
state0.transitions["f"] = state9
state0.transitions["10"] = state2
state0.transitions["e"] = state13
state0.transitions["14"] = state2
state0.transitions["("] = state5
state0.transitions["0"] = state8
state0.transitions["22"] = state2
state0.transitions["12"] = state2
state0.transitions["2"] = state8
state0.transitions["\n"] = state15
state0.transitions["b"] = state6
state0.transitions["23"] = state2
state0.transitions["16"] = state2
state0.transitions["w"] = state12
state0.transitions["21"] = state2
state0.transitions[" "] = state15
state0.transitions["	"] = state15
state0.transitions["i"] = state1
state0.transitions[")"] = state7
state0.transitions["13"] = state2
state0.transitions["="] = state10
state0.transitions["A"] = state6
state0.transitions["19"] = state2
state0.transitions["15"] = state2
state0.transitions["c"] = state6
state0.transitions["18"] = state2
state0.transitions["s"] = state2
state0.transitions["/"] = state11
state0.transitions["t"] = state2
state0.transitions["+"] = state14
state0.transitions["h"] = state2
state0.transitions["n"] = state2
state7.transitions["c"] = state2
state7.transitions["f"] = state2
state7.transitions["19"] = state2
state7.transitions["+"] = state2
state7.transitions["	"] = state2
state7.transitions["*"] = state2
state7.transitions["u"] = state2
state7.transitions["10"] = state2
state7.transitions[" "] = state2
state7.transitions["14"] = state2
state7.transitions["20"] = state2
state7.transitions["17"] = state2
state7.transitions["23"] = state2
state7.transitions["13"] = state2
state7.transitions["1"] = state2
state7.transitions["12"] = state2
state7.transitions["w"] = state2
state7.transitions["e"] = state2
state7.transitions["h"] = state2
state7.transitions[")"] = state2
state7.transitions["21"] = state18
state7.transitions["l"] = state2
state7.transitions["b"] = state2
state7.transitions["i"] = state2
state7.transitions["18"] = state2
state7.transitions["="] = state2
state7.transitions["16"] = state2
state7.transitions["A"] = state2
state7.transitions["a"] = state2
state7.transitions["r"] = state2
state7.transitions["B"] = state2
state7.transitions["15"] = state2
state7.transitions["22"] = state2
state7.transitions["2"] = state2
state7.transitions["\n"] = state2
state7.transitions["-"] = state2
state7.transitions["("] = state2
state7.transitions["0"] = state2
state7.transitions["s"] = state2
state7.transitions["/"] = state2
state7.transitions["n"] = state2
state7.transitions["t"] = state2
state7.transitions["11"] = state2
state7.transitions["24"] = state2
state8.transitions["18"] = state2
state8.transitions["22"] = state2
state8.transitions["u"] = state2
state8.transitions["b"] = state2
state8.transitions["t"] = state2
state8.transitions["r"] = state2
state8.transitions["\n"] = state2
state8.transitions["i"] = state2
state8.transitions["*"] = state2
state8.transitions["("] = state2
state8.transitions["20"] = state2
state8.transitions["24"] = state2
state8.transitions["1"] = state8
state8.transitions["n"] = state2
state8.transitions["e"] = state2
state8.transitions["11"] = state2
state8.transitions["2"] = state8
state8.transitions["l"] = state2
state8.transitions["15"] = state2
state8.transitions["h"] = state2
state8.transitions["23"] = state18
state8.transitions["19"] = state2
state8.transitions["a"] = state2
state8.transitions["+"] = state2
state8.transitions[" "] = state2
state8.transitions["10"] = state2
state8.transitions["12"] = state2
state8.transitions["/"] = state2
state8.transitions["c"] = state2
state8.transitions["17"] = state2
state8.transitions["13"] = state2
state8.transitions["f"] = state2
state8.transitions["w"] = state2
state8.transitions[")"] = state2
state8.transitions["16"] = state2
state8.transitions["14"] = state2
state8.transitions["B"] = state2
state8.transitions["	"] = state2
state8.transitions["0"] = state8
state8.transitions["21"] = state2
state8.transitions["A"] = state2
state8.transitions["-"] = state2
state8.transitions["s"] = state2
state8.transitions["="] = state2
state21.transitions["i"] = state2
state21.transitions["l"] = state2
state21.transitions["b"] = state2
state21.transitions["="] = state2
state21.transitions["e"] = state2
state21.transitions["B"] = state2
state21.transitions["12"] = state2
state21.transitions["h"] = state2
state21.transitions["0"] = state2
state21.transitions["13"] = state2
state21.transitions["r"] = state2
state21.transitions["15"] = state2
state21.transitions["/"] = state2
state21.transitions["A"] = state2
state21.transitions["19"] = state2
state21.transitions["n"] = state2
state21.transitions["t"] = state2
state21.transitions["2"] = state2
state21.transitions["17"] = state2
state21.transitions["1"] = state2
state21.transitions["a"] = state2
state21.transitions["11"] = state2
state21.transitions["	"] = state2
state21.transitions["*"] = state2
state21.transitions["20"] = state2
state21.transitions["10"] = state2
state21.transitions["+"] = state2
state21.transitions[" "] = state2
state21.transitions["16"] = state2
state21.transitions["24"] = state2
state21.transitions["22"] = state2
state21.transitions["u"] = state2
state21.transitions["23"] = state2
state21.transitions["f"] = state2
state21.transitions["21"] = state2
state21.transitions["c"] = state2
state21.transitions[")"] = state2
state21.transitions["w"] = state2
state21.transitions["14"] = state2
state21.transitions["\n"] = state2
state21.transitions["-"] = state2
state21.transitions["("] = state2
state21.transitions["18"] = state2
state21.transitions["s"] = state25
state2.transitions["20"] = state2
state2.transitions[")"] = state2
state2.transitions["12"] = state2
state2.transitions["+"] = state2
state2.transitions["24"] = state2
state2.transitions["c"] = state2
state2.transitions["1"] = state2
state2.transitions["i"] = state2
state2.transitions["h"] = state2
state2.transitions["17"] = state2
state2.transitions["0"] = state2
state2.transitions["b"] = state2
state2.transitions["/"] = state2
state2.transitions["10"] = state2
state2.transitions["a"] = state2
state2.transitions["l"] = state2
state2.transitions["*"] = state2
state2.transitions["23"] = state2
state2.transitions["s"] = state2
state2.transitions["16"] = state2
state2.transitions["21"] = state2
state2.transitions[" "] = state2
state2.transitions["14"] = state2
state2.transitions["18"] = state2
state2.transitions["n"] = state2
state2.transitions["	"] = state2
state2.transitions["-"] = state2
state2.transitions["13"] = state2
state2.transitions["="] = state2
state2.transitions["t"] = state2
state2.transitions["B"] = state2
state2.transitions["2"] = state2
state2.transitions["\n"] = state2
state2.transitions["("] = state2
state2.transitions["u"] = state2
state2.transitions["19"] = state2
state2.transitions["15"] = state2
state2.transitions["22"] = state2
state2.transitions["f"] = state2
state2.transitions["A"] = state2
state2.transitions["w"] = state2
state2.transitions["e"] = state2
state2.transitions["11"] = state2
state2.transitions["r"] = state2
state12.transitions["l"] = state2
state12.transitions["15"] = state2
state12.transitions["24"] = state2
state12.transitions["/"] = state2
state12.transitions["n"] = state2
state12.transitions["i"] = state2
state12.transitions["c"] = state2
state12.transitions["h"] = state20
state12.transitions["f"] = state2
state12.transitions["="] = state2
state12.transitions["10"] = state2
state12.transitions["19"] = state2
state12.transitions["2"] = state2
state12.transitions["-"] = state2
state12.transitions["20"] = state2
state12.transitions["u"] = state2
state12.transitions["w"] = state2
state12.transitions["+"] = state2
state12.transitions["	"] = state2
state12.transitions["b"] = state2
state12.transitions["13"] = state2
state12.transitions["0"] = state2
state12.transitions["a"] = state2
state12.transitions["e"] = state2
state12.transitions["t"] = state2
state12.transitions["14"] = state2
state12.transitions["*"] = state2
state12.transitions["("] = state2
state12.transitions["17"] = state2
state12.transitions[")"] = state2
state12.transitions["s"] = state2
state12.transitions["21"] = state2
state12.transitions[" "] = state2
state12.transitions["r"] = state2
state12.transitions["22"] = state2
state12.transitions["23"] = state2
state12.transitions["16"] = state2
state12.transitions["12"] = state2
state12.transitions["11"] = state2
state12.transitions["B"] = state2
state12.transitions["\n"] = state2
state12.transitions["18"] = state2
state12.transitions["1"] = state2
state12.transitions["A"] = state2
state14.transitions["13"] = state2
state14.transitions["17"] = state2
state14.transitions["1"] = state2
state14.transitions["A"] = state2
state14.transitions["12"] = state2
state14.transitions["w"] = state2
state14.transitions[" "] = state2
state14.transitions["r"] = state2
state14.transitions["\n"] = state2
state14.transitions["22"] = state2
state14.transitions["19"] = state2
state14.transitions["("] = state2
state14.transitions["c"] = state2
state14.transitions["s"] = state2
state14.transitions["11"] = state2
state14.transitions["14"] = state2
state14.transitions["	"] = state2
state14.transitions["i"] = state2
state14.transitions["*"] = state2
state14.transitions["20"] = state2
state14.transitions["="] = state2
state14.transitions["/"] = state2
state14.transitions["10"] = state2
state14.transitions["21"] = state2
state14.transitions["24"] = state2
state14.transitions["a"] = state2
state14.transitions["e"] = state2
state14.transitions["l"] = state2
state14.transitions["-"] = state2
state14.transitions[")"] = state2
state14.transitions["u"] = state2
state14.transitions["b"] = state2
state14.transitions["n"] = state2
state14.transitions["B"] = state2
state14.transitions["h"] = state2
state14.transitions["f"] = state2
state14.transitions["16"] = state18
state14.transitions["t"] = state2
state14.transitions["+"] = state2
state14.transitions["2"] = state2
state14.transitions["15"] = state2
state14.transitions["18"] = state2
state14.transitions["0"] = state2
state14.transitions["23"] = state2
state26.transitions["2"] = state2
state26.transitions["18"] = state2
state26.transitions["15"] = state2
state26.transitions["10"] = state2
state26.transitions["r"] = state2
state26.transitions["("] = state2
state26.transitions["0"] = state2
state26.transitions["n"] = state2
state26.transitions["e"] = state2
state26.transitions["+"] = state2
state26.transitions["\n"] = state2
state26.transitions["17"] = state2
state26.transitions["*"] = state2
state26.transitions["24"] = state2
state26.transitions["h"] = state2
state26.transitions["12"] = state2
state26.transitions["B"] = state2
state26.transitions["	"] = state2
state26.transitions["13"] = state2
state26.transitions["s"] = state2
state26.transitions["="] = state2
state26.transitions["i"] = state2
state26.transitions["c"] = state2
state26.transitions["A"] = state2
state26.transitions["t"] = state2
state26.transitions[" "] = state2
state26.transitions["11"] = state2
state26.transitions["14"] = state2
state26.transitions["-"] = state2
state26.transitions[")"] = state2
state26.transitions["22"] = state2
state26.transitions["u"] = state30
state26.transitions["b"] = state2
state26.transitions["w"] = state2
state26.transitions["19"] = state2
state26.transitions["l"] = state2
state26.transitions["20"] = state2
state26.transitions["23"] = state2
state26.transitions["/"] = state2
state26.transitions["21"] = state2
state26.transitions["f"] = state2
state26.transitions["16"] = state2
state26.transitions["1"] = state2
state26.transitions["a"] = state2
state27.transitions["*"] = state2
state27.transitions["u"] = state2
state27.transitions["="] = state2
state27.transitions["12"] = state2
state27.transitions["11"] = state18
state27.transitions["l"] = state2
state27.transitions["13"] = state2
state27.transitions[" "] = state2
state27.transitions["h"] = state2
state27.transitions["0"] = state2
state27.transitions["10"] = state2
state27.transitions["t"] = state2
state27.transitions["21"] = state2
state27.transitions["r"] = state2
state27.transitions["2"] = state2
state27.transitions["	"] = state2
state27.transitions["-"] = state2
state27.transitions["20"] = state2
state27.transitions["24"] = state2
state27.transitions["f"] = state2
state27.transitions["1"] = state2
state27.transitions["/"] = state2
state27.transitions["a"] = state2
state27.transitions["e"] = state2
state27.transitions["22"] = state2
state27.transitions["14"] = state2
state27.transitions["\n"] = state2
state27.transitions["c"] = state2
state27.transitions["17"] = state2
state27.transitions[")"] = state2
state27.transitions["b"] = state2
state27.transitions["s"] = state2
state27.transitions["16"] = state2
state27.transitions["19"] = state2
state27.transitions["B"] = state2
state27.transitions["15"] = state2
state27.transitions["18"] = state2
state27.transitions["23"] = state2
state27.transitions["A"] = state2
state27.transitions["w"] = state2
state27.transitions["n"] = state2
state27.transitions["i"] = state2
state27.transitions["("] = state2
state27.transitions["+"] = state2
state33.transitions["\n"] = state2
state33.transitions[")"] = state2
state33.transitions["f"] = state2
state33.transitions["/"] = state2
state33.transitions["19"] = state2
state33.transitions["21"] = state2
state33.transitions["i"] = state2
state33.transitions["15"] = state2
state33.transitions["18"] = state2
state33.transitions["+"] = state2
state33.transitions["11"] = state2
state33.transitions["r"] = state2
state33.transitions["l"] = state2
state33.transitions["17"] = state2
state33.transitions["0"] = state2
state33.transitions["1"] = state2
state33.transitions["12"] = state2
state33.transitions["w"] = state2
state33.transitions["B"] = state2
state33.transitions["20"] = state2
state33.transitions["c"] = state2
state33.transitions["t"] = state2
state33.transitions["2"] = state2
state33.transitions["24"] = state2
state33.transitions["10"] = state2
state33.transitions["e"] = state2
state33.transitions["	"] = state2
state33.transitions["-"] = state2
state33.transitions["22"] = state2
state33.transitions["13"] = state2
state33.transitions["s"] = state2
state33.transitions["A"] = state2
state33.transitions["a"] = state2
state33.transitions["n"] = state2
state33.transitions["*"] = state2
state33.transitions["("] = state2
state33.transitions["h"] = state2
state33.transitions["u"] = state2
state33.transitions["16"] = state2
state33.transitions["b"] = state2
state33.transitions["23"] = state2
state33.transitions["="] = state2
state33.transitions[" "] = state2
state33.transitions["14"] = state18
state13.transitions["*"] = state2
state13.transitions["20"] = state2
state13.transitions["h"] = state2
state13.transitions["18"] = state2
state13.transitions["1"] = state2
state13.transitions["12"] = state2
state13.transitions["21"] = state2
state13.transitions["B"] = state2
state13.transitions["17"] = state2
state13.transitions[")"] = state2
state13.transitions["s"] = state2
state13.transitions["10"] = state2
state13.transitions["n"] = state2
state13.transitions["r"] = state2
state13.transitions["2"] = state2
state13.transitions["	"] = state2
state13.transitions["24"] = state2
state13.transitions["22"] = state2
state13.transitions["A"] = state2
state13.transitions["a"] = state2
state13.transitions["+"] = state2
state13.transitions["-"] = state2
state13.transitions["11"] = state2
state13.transitions["l"] = state21
state13.transitions["u"] = state2
state13.transitions["f"] = state2
state13.transitions["16"] = state2
state13.transitions["/"] = state2
state13.transitions["\n"] = state2
state13.transitions["i"] = state2
state13.transitions["15"] = state2
state13.transitions["("] = state2
state13.transitions["c"] = state2
state13.transitions["13"] = state2
state13.transitions["19"] = state2
state13.transitions["e"] = state2
state13.transitions["t"] = state2
state13.transitions["0"] = state2
state13.transitions["b"] = state2
state13.transitions["23"] = state2
state13.transitions["="] = state2
state13.transitions["w"] = state2
state13.transitions[" "] = state2
state13.transitions["14"] = state2
state1.transitions["16"] = state2
state1.transitions["c"] = state2
state1.transitions["17"] = state2
state1.transitions["23"] = state2
state1.transitions["n"] = state2
state1.transitions["14"] = state2
state1.transitions["2"] = state2
state1.transitions["\n"] = state2
state1.transitions["*"] = state2
state1.transitions["24"] = state2
state1.transitions["b"] = state2
state1.transitions["w"] = state2
state1.transitions["+"] = state2
state1.transitions[" "] = state2
state1.transitions["11"] = state2
state1.transitions["20"] = state2
state1.transitions["18"] = state2
state1.transitions["22"] = state2
state1.transitions["u"] = state2
state1.transitions["13"] = state2
state1.transitions["s"] = state2
state1.transitions["e"] = state2
state1.transitions["l"] = state2
state1.transitions["("] = state2
state1.transitions["t"] = state2
state1.transitions["1"] = state2
state1.transitions[")"] = state2
state1.transitions["0"] = state2
state1.transitions["i"] = state2
state1.transitions["15"] = state2
state1.transitions["f"] = state17
state1.transitions["12"] = state2
state1.transitions["r"] = state2
state1.transitions["B"] = state2
state1.transitions["-"] = state2
state1.transitions["h"] = state2
state1.transitions["="] = state2
state1.transitions["A"] = state2
state1.transitions["19"] = state2
state1.transitions["a"] = state2
state1.transitions["/"] = state2
state1.transitions["10"] = state2
state1.transitions["21"] = state2
state1.transitions["	"] = state2
state9.transitions["23"] = state2
state9.transitions["13"] = state2
state9.transitions["12"] = state2
state9.transitions[" "] = state2
state9.transitions["r"] = state2
state9.transitions["	"] = state2
state9.transitions[")"] = state2
state9.transitions["h"] = state2
state9.transitions["n"] = state2
state9.transitions["e"] = state2
state9.transitions["*"] = state2
state9.transitions["0"] = state2
state9.transitions["u"] = state19
state9.transitions["="] = state2
state9.transitions["2"] = state2
state9.transitions["\n"] = state2
state9.transitions["15"] = state2
state9.transitions["("] = state2
state9.transitions["c"] = state2
state9.transitions["/"] = state2
state9.transitions["w"] = state2
state9.transitions["t"] = state2
state9.transitions["l"] = state2
state9.transitions["-"] = state2
state9.transitions["f"] = state2
state9.transitions["s"] = state2
state9.transitions["B"] = state2
state9.transitions["i"] = state2
state9.transitions["24"] = state2
state9.transitions["1"] = state2
state9.transitions["10"] = state2
state9.transitions["a"] = state2
state9.transitions["+"] = state2
state9.transitions["11"] = state2
state9.transitions["14"] = state2
state9.transitions["18"] = state2
state9.transitions["22"] = state2
state9.transitions["b"] = state2
state9.transitions["16"] = state2
state9.transitions["A"] = state2
state9.transitions["19"] = state2
state9.transitions["21"] = state2
state9.transitions["20"] = state2
state9.transitions["17"] = state2
state10.transitions["\n"] = state2
state10.transitions["f"] = state2
state10.transitions["1"] = state2
state10.transitions["/"] = state2
state10.transitions["A"] = state2
state10.transitions["l"] = state2
state10.transitions["-"] = state2
state10.transitions["b"] = state2
state10.transitions[" "] = state2
state10.transitions["r"] = state2
state10.transitions["23"] = state2
state10.transitions["21"] = state2
state10.transitions["20"] = state2
state10.transitions["0"] = state2
state10.transitions["13"] = state2
state10.transitions["s"] = state2
state10.transitions["="] = state2
state10.transitions["16"] = state2
state10.transitions["10"] = state2
state10.transitions["2"] = state2
state10.transitions["18"] = state2
state10.transitions["19"] = state2
state10.transitions["a"] = state2
state10.transitions["11"] = state2
state10.transitions["	"] = state2
state10.transitions["i"] = state2
state10.transitions[")"] = state2
state10.transitions["22"] = state2
state10.transitions["u"] = state2
state10.transitions["12"] = state2
state10.transitions["w"] = state2
state10.transitions["t"] = state2
state10.transitions["+"] = state2
state10.transitions["*"] = state2
state10.transitions["("] = state2
state10.transitions["24"] = state2
state10.transitions["c"] = state2
state10.transitions["h"] = state2
state10.transitions["17"] = state2
state10.transitions["e"] = state2
state10.transitions["14"] = state2
state10.transitions["15"] = state18
state10.transitions["n"] = state2
state10.transitions["B"] = state2

return &dfa{ 
startState: state0,
states: []*state{ state0, state1, state2, state3, state4, state5, state6, state7, state8, state9, state10, state11, state12, state13, state14, state15, state16, state17, state18, state19, state20, state21, state22, state23, state24, state25, state26, state27, state28, state29, state30, state31, state32, state33, }, 
}
}

// =====================
//	Footer
// =====================
// Contains the exact same content defined on the Yaaalex file


    // The entire contents of this section will be copied to the beginning of the generated Lexer.go file

    //  ------ TOKENS ID -----

    // Define the token types that the lexer will recognize

    //This is a footer


