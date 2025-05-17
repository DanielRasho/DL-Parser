package reader

import (
	"fmt"
	"testing"
)

func Test_check1(t *testing.T) {

	el, _ := Parse("../../../../examples/productions2.y")

	fmt.Println(el.Terminals)
	fmt.Println(el.NonTerminals)
	fmt.Println(el.Productions)

}

func Test_check2(t *testing.T) {

	el, _ := Parse("../../../../examples/exampleprod.y")

	fmt.Println(el.Terminals)
	fmt.Println(el.NonTerminals)
	fmt.Print("Productions\n")
	fmt.Println(el.Productions)

}

func Test_check3(t *testing.T) {

	el, _ := Parse("../../../../examples/exampleprod2.y")

	fmt.Println(el.Terminals)
	fmt.Println(el.NonTerminals)
	fmt.Print("Productions\n")
	fmt.Println(el.Productions)

}
