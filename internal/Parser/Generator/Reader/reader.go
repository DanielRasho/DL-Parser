package reader

import (
	"strings"

	io "github.com/DanielRasho/Parser/internal/IO"
	Parser "github.com/DanielRasho/Parser/internal/Parser"
)

// TODO: Lo que tienes que implementar andre
func Parse(filePath string) (*Parser.ParserDefinition, error) {

	var line string
	var token []string
	var Tokens []Parser.ParserSymbol
	var Productions *Parser.ParserProduction = new(Parser.ParserProduction)
	var arrProductions []Parser.ParserProduction
	var nonterminals []Parser.ParserSymbol
	ignoredTokens := make(map[int]Parser.ParserSymbol)
	var is_product = false
	var head string
	var err error

	tokensReaded := 0

	filereader, _ := io.ReadFile(filePath)
	for filereader.NextLine(&line) {

		//Starts parsing tokens
		if strings.Contains(line, "%token") && !is_product {
			token = strings.Split(line, "%token")
			token = strings.Split(token[1], " ")
			for i := 1; i < len(token); i++ {
				token[i] = strings.TrimSpace(token[i])
				Tokens = append(Tokens, Parser.ParserSymbol{Id: tokensReaded, Value: token[i], IsTerminal: true})
				tokensReaded++
			}
		} else if strings.Contains(line, "IGNORE") && !is_product {
			token = strings.Split(line, "IGNORE")
			token = strings.Split(token[1], " ")
			for i := 1; i < len(token); i++ {
				token[i] = strings.TrimSpace(token[i])
				ignoredTokens[tokensReaded] = Parser.ParserSymbol{Id: tokensReaded, Value: token[i], IsTerminal: true}
				tokensReaded++
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
							Productions = &Parser.ParserProduction{Head: samehead, Id: (len(arrProductions) + 1)}

							// fmt.Println("array", Productions.Body)

							line = strings.Split(line, "|")[1]
							line = strings.TrimSpace(line)
						}

						token = strings.Split(line, " ")
						nonTerminalIndexCounter := -1

						for i := range len(token) {
							index_val := findIndex(Tokens, token[i])
							if index_val == -1 {

								index_valnon := findIndex(nonterminals, token[i])
								if index_valnon == -1 {
									nonterminals = append(nonterminals, Parser.ParserSymbol{Id: -1, Value: token[i]})
									index_valnon = findIndex(nonterminals, token[i])
									nonTerminalIndexCounter--
								}
								Productions.Body = append(Productions.Body, nonterminals[index_valnon])
								Productions.Id = len(arrProductions) + 1
							} else {
								Productions.Id = len(arrProductions) + 1
								Productions.Body = append(Productions.Body, Tokens[index_val])
							}

						}
					} else {
						arrProductions = append(arrProductions, *Productions)
						// fmt.Println("array2", arrProductions)
						head = ""
					}

				}

			}

			//Lee la primera produccion y luego de eso si tiene un valor entonces se agrega al parser y las producciones en otro lado
			if strings.Contains(line, ":") {
				head = strings.Split(line, ":")[0]
				index_val := findIndex(Tokens, head)
				if index_val == -1 {

					Productions = new(Parser.ParserProduction)
					Productions.Head = Parser.ParserSymbol{Id: -1, Value: head}
				}
			}

		}
		//Inicializamos el lugar donde empieza a leer las producciones
		if strings.Contains(line, "%%") {
			is_product = true
		}

	}
	for i := range len(Tokens) {
		if Tokens[i].Id == -1 {
			nonterminals = append(nonterminals, Tokens[i])
		}
	}
	for i := 0; i < len(arrProductions); i++ {
		if !AlreadyinList(nonterminals, arrProductions[i].Head) {
			nonterminals = append(nonterminals, arrProductions[i].Head)

		}

	}

	return &Parser.ParserDefinition{
		NonTerminals:  nonterminals,
		Terminals:     Tokens,
		Productions:   arrProductions,
		IgnoredSymbol: ignoredTokens,
	}, err
}

// FINDS THE INDEX VALUE OF THE ARRAY OF PARSE SYMBOLS
func findIndex(tokens []Parser.ParserSymbol, target string) int {
	for i, token := range tokens {
		if token.Value == target {
			return i
		}
	}
	return -1 // not found
}

// FINDS THE INDEX VALUE OF THE ARRAY OF PARSE SYMBOLS
func AlreadyinList(nonterminals []Parser.ParserSymbol, target Parser.ParserSymbol) bool {

	for i := 0; i < len(nonterminals); i++ {
		if nonterminals[i].Id == target.Id && nonterminals[i].IsTerminal == target.IsTerminal && nonterminals[i].Value == target.Value {
			return true
		}
	}
	return false

}
