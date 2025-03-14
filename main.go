package main

import (
	"fmt"
	"log/slog"
	"metasource/metasource/models"
	"metasource/metasource/runner"
)

func main() {
	var meta *models.Metadata
	var expt error

	meta, expt = runner.Import()
	if expt != nil {
		slog.Log(nil, slog.LevelError, "File could not be loaded")
		return
	}

	for _, unit := range meta.List {
		if unit.Name == "obserware" {
			fmt.Println(unit)
		}
	}

	// fmt.Println(meta.List[0])
}
