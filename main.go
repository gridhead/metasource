package main

import (
	"fmt"
	"metasource/metasource/models/sxml"
	"metasource/metasource/runner"
	"sync"
)

func main() {
	var name string
	var rslt_primary *sxml.UnitPrimary
	var rslt_other *sxml.UnitOther
	var rslt_filelist *sxml.UnitFileList
	var wait sync.WaitGroup

	name = "obserware"

	wait.Add(3)
	go runner.ImportPrimary(&wait, &rslt_primary, &name)
	go runner.ImportOther(&wait, &rslt_other, &name)
	go runner.ImportFileList(&wait, &rslt_filelist, &name)
	wait.Wait()

	if rslt_primary != nil {
		fmt.Println(rslt_primary)
	}

	if rslt_other != nil {
		fmt.Println(rslt_other)
	}

	if rslt_filelist != nil {
		fmt.Println(rslt_filelist)
	}
}
