package transitiontable

import (
	"fmt"
	"testing"

	reader "github.com/DanielRasho/Parser/internal/Parser/Generator/Reader"
	"github.com/DanielRasho/Parser/internal/Parser/automata"
)

func Test_check1(t *testing.T) {

	parserdef, _ := reader.Parse("../../../examples/productions2.y")
	first := GetFirst(parserdef)
	follow := GetFollow(parserdef, first)
	var automa = automata.NewAutomata(parserdef)

	TransitionTabl, gotable, _ := NewTable(automa, first, follow, *parserdef)

	fmt.Println("TABLA DE TRANSICION")
	fmt.Println(TransitionTabl)
	fmt.Println("GOTO")
	fmt.Println(gotable)
}
