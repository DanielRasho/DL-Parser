package reader

import (
	"errors"
	"fmt"
	"log"
	"strings"

	io "github.com/DanielRasho/Parser/internal/IO"
)

// TODO: Lo que tienes que implementar andre
func Parse(filePath string) (*ParserDefinition, error) {

	var line string
	var token []string
	var Tokens []ParserSymbol
	var Productions *ParserProduction = new(ParserProduction)
	var arrProductions []ParserProduction
	var is_product = false
	var head string

	filereader, _ := io.ReadFile(filePath)
	for filereader.NextLine(&line) {

		//Starts parsing tokens
		if strings.Contains(line, "%token") && !is_product {
			token = strings.Split(line, "%token")
			token = strings.Split(token[1], " ")
			for i := 1; i < len(token); i++ {
				token[i] = strings.TrimSpace(token[i])
				Tokens = append(Tokens, ParserSymbol{Id: len(Tokens), Value: token[i]})
			}

		}

		//Una vez terminado de leer los tokens terminales empezamos a leer las producciones
		if is_product {

			//Aqui es donde se lee las producciones y termina cuando no hay mas producciones a leer si tiene el simbolo ; entonces se termina
			if head != "" {
				line = strings.TrimSpace(line)
				if line != "" {

					if !strings.Contains(line, ";") {
						if strings.Contains(line, "|") {

							arrProductions = append(arrProductions, *Productions)
							samehead := Productions.Head
							Productions = &ParserProduction{Head: samehead}

							// fmt.Println("array", Productions.Body)

							line = strings.Split(line, "|")[1]
							line = strings.TrimSpace(line)
						}

						token = strings.Split(line, " ")

						for i := range len(token) {
							index_val := findIndex(Tokens, token[i])
							Productions.Body = append(Productions.Body, Tokens[index_val])
						}
					} else {
						arrProductions = append(arrProductions, *Productions)
						// fmt.Println("array2", arrProductions)
						Productions = new(ParserProduction)
						head = ""
					}

				}

			}

			//Lee la primera produccion y luego de eso si tiene un valor entonces se agrega al parser y las producciones en otro lado
			if strings.Contains(line, ":") {
				head = strings.Split(line, ":")[0]
				index_val := findIndex(Tokens, head)
				if index_val != -1 {
					Tokens[index_val].Id = -1
					Productions.Head = Tokens[index_val]

				} else {
					err := errors.New("NO TOKEN FOUND EXITING PROGRAM")
					log.Fatal(err)
				}
			}

		}
		//Inicializamos el lugar donde empieza a leer las producciones
		if strings.Contains(line, "%%") {
			is_product = true
		}

	}

	fmt.Println(Tokens)
	fmt.Println(" PRODUCCIONES")
	fmt.Println(arrProductions)

	return nil, nil
}

// FINDS THE INDEX VALUE OF THE ARRAY OF PARSE SYMBOLS
func findIndex(tokens []ParserSymbol, target string) int {
	for i, token := range tokens {
		if token.Value == target {
			return i
		}
	}
	return -1 // not found
}
