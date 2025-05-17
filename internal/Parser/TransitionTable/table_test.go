package transitiontable

import (
	"fmt"
	"testing"

	reader "github.com/DanielRasho/Parser/internal/Parser/Generator/Reader"
)

func Test_check1(t *testing.T) {

	el, _ := reader.Parse("../../../../examples/productions2.y")
	first := GetFirst(el)
	GetFollow(el, first)

	// automata := &automata.Automata{
	// 	States: []*automata.State{
	// 		{Id: "I1", Transitions: map[Symbol]*State{},IsAccepted: true ,IsFinal: true ,Productions: []parser.ParserProduction{  parser.ParserProduction{ Head: parser.ParserSymbol{Id: -1, Value: "S'", IsTerminal: true}, Body:[]parser.ParserSymbol{ parser.ParserSymbol{ Id: -1, Value: "E", IsTerminal: false} }}}
	// 	}
	// }

	fmt.Println(el.Terminals)
	fmt.Println(el.NonTerminals)
	fmt.Println(el.Productions)

}
