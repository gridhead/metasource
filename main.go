package main

import (
	"fmt"
	"log/slog"
	"metasource/metasource/config"
	"metasource/metasource/driver"
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

	// ACTUAL CODE
	//var server *http.Server
	//var expt error
	//
	config.SetLogger()
	//
	//server = &http.Server{Addr: ":8080", Handler: http.DefaultServeMux}
	//
	//http.HandleFunc("/rawhide/changelog/", routes.RetrieveOther)
	//http.HandleFunc("/rawhide/pkg/", routes.RetrievePrimary)
	//http.HandleFunc("/rawhide/files/", routes.RetrieveFileList)
	//http.HandleFunc("/rawhide/", routes.RetrievePlus)
	//
	//expt = server.ListenAndServe()
	//if expt != nil {
	//	slog.Log(nil, slog.LevelError, fmt.Sprintf("Error occurred. %s.", expt.Error()))
	//}

	expt := driver.Database("/var/tmp/xyz")
	if expt != nil {
		slog.Log(nil, slog.LevelError, fmt.Sprintf("%s", expt.Error()))
	}
}
