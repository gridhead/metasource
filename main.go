package main

import (
	"fmt"
	"metasource/metasource/models"
	"metasource/metasource/runner"
	"sync"
)

func main() {
	var rslt_primary *models.Primary
	var rslt_other *models.Other
	var wait sync.WaitGroup

	wait.Add(2)
	go runner.ImportPrimary(&wait, &rslt_primary)
	go runner.ImportOther(&wait, &rslt_other)
	wait.Wait()

	if rslt_primary != nil {
		for _, unit := range rslt_primary.List {
			if unit.Name == "obserware" {
				fmt.Println(unit)
			}
		}
	}

	if rslt_other != nil {
		for _, unit := range rslt_other.List {
			if unit.Name == "obserware" {
				fmt.Println(unit)
			}
		}
	}
}
