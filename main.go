package main

import (
	"fmt"
	"metasource/metasource/models"
	"metasource/metasource/runner"
	"sync"
)

func main() {
	var name string
	var rslt_primary *models.Primary
	var rslt_other *models.Other
	var rslt_filelist *models.FileList
	var wait sync.WaitGroup

	wait.Add(3)
	go runner.ImportPrimary(&wait, &rslt_primary)
	go runner.ImportOther(&wait, &rslt_other)
	go runner.ImportFileList(&wait, &rslt_filelist)
	wait.Wait()

	name = "nano"

	if rslt_primary != nil {
		for _, unit := range rslt_primary.List {
			if unit.Name == name {
				fmt.Println(unit)
			}
		}
	}

	if rslt_other != nil {
		for _, unit := range rslt_other.List {
			if unit.Name == name {
				fmt.Println(unit)
			}
		}
	}

	if rslt_filelist != nil {
		for _, unit := range rslt_filelist.List {
			if unit.Name == name {
				fmt.Println(unit)
			}
		}
	}
}
