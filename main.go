package main

import (
	"fmt"
	"log/slog"
	"metasource/metasource/routes"
	"net/http"
)

func main() {
	//var name string
	//var rslt_primary *sxml.UnitPrimary
	//var rslt_other *sxml.UnitOther
	//var rslt_filelist *sxml.UnitFileList
	//var wait sync.WaitGroup
	//
	//name = "obserware"
	//
	//wait.Add(3)
	//go runner.ImportPrimary(&wait, &rslt_primary, &name)
	//go runner.ImportOther(&wait, &rslt_other, &name)
	//go runner.ImportFileList(&wait, &rslt_filelist, &name)
	//wait.Wait()
	//
	//if rslt_primary != nil {
	//	fmt.Println(rslt_primary)
	//}
	//
	//if rslt_other != nil {
	//	fmt.Println(rslt_other)
	//}
	//
	//if rslt_filelist != nil {
	//	fmt.Println(rslt_filelist)
	//}

	var server *http.Server
	var expt error

	server = &http.Server{Addr: ":8080", Handler: http.DefaultServeMux}

	http.HandleFunc("/rawhide/changelog/", routes.RetrieveChangelog)

	expt = server.ListenAndServe()
	if expt != nil {
		slog.Log(nil, slog.LevelError, fmt.Sprintf("Error occurred. %s.", expt.Error()))
	}

}
