package main

import (
	"fmt"
	"log/slog"
	"metasource/metasource/models"
	"metasource/metasource/runner"
	"sync"
)

func main() {
	var rslt_metadata *models.Primary
	var rslt_other *models.Other
	var expt error
	var wait sync.WaitGroup

	wait.Add(1)
	rslt_metadata, expt = runner.ImportPrimary(&wait)
	if expt != nil {
		slog.Log(nil, slog.LevelError, "File could not be loaded")
		return
	}
	wait.Wait()

	wait.Add(1)
	rslt_other, expt = runner.ImportOther(&wait)
	if expt != nil {
		slog.Log(nil, slog.LevelError, "File could not be loaded")
		return
	}
	wait.Wait()

	for _, unit := range rslt_metadata.List {
		if unit.Name == "obserware" {
			fmt.Println(unit)
		}
	}

	for _, unit := range rslt_other.List {
		if unit.Name == "obserware" {
			fmt.Println(unit)
		}
	}

}
