package generator

import (
	reader "github.com/DanielRasho/Parser/internal/Parser/Generator/Reader"
	transitiontable "github.com/DanielRasho/Parser/internal/Parser/TransitionTable"
)

// Given a file to read and a output path, writes a parser definition to the desired path.
func Compile(filePath, outputPath string, showLogs bool) error {

	// Parse Yalex file definition
	yalexDefinition, err := reader.Parse(filePath)
	if err != nil {
		return err
	}

	// runtime.Breakpoint()
	first := transitiontable.GetFirst(yalexDefinition)
	transitiontable.GetFollow(yalexDefinition, first)

	return nil
}
