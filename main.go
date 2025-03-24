package main

import (
	"flag"
	"fmt"
	"log/slog"
	"metasource/metasource/config"
	"metasource/metasource/driver"
	"os"
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

	var expt error
	var lglvtext, location *string
	var database, dispense *flag.FlagSet

	lglvtext = flag.String("loglevel", "info", "Set the application loglevel")
	location = flag.String("location", config.DBFOLDER, "Set the database location")
	flag.Parse()

	database = flag.NewFlagSet("database", flag.ExitOnError)
	dispense = flag.NewFlagSet("dispense", flag.ExitOnError)
	config.SetLogger(lglvtext)

	if flag.NArg() < 1 {
		slog.Log(nil, slog.LevelError, "Invalid subcommand")
		slog.Log(nil, slog.LevelInfo, "Expected either 'database' or 'dispense' subcommand")
		os.Exit(1)
	}

	switch flag.Arg(0) {
	case "database":
		expt = database.Parse(os.Args[2:])
		if expt != nil {
			slog.Log(nil, slog.LevelError, fmt.Sprintf("%s", expt.Error()))
			os.Exit(1)
		}
		expt = driver.Database(*location)
		if expt != nil {
			slog.Log(nil, slog.LevelError, fmt.Sprintf("%s", expt.Error()))
			os.Exit(1)
		}
		os.Exit(0)
	case "dispense":
		expt = dispense.Parse(os.Args[2:])
		if expt != nil {
			slog.Log(nil, slog.LevelError, fmt.Sprintf("%s", expt.Error()))
			os.Exit(1)
		}
		slog.Log(nil, slog.LevelError, "Dispense is not implemented yet")
		os.Exit(1)
	default:
		slog.Log(nil, slog.LevelError, "Invalid subcommand")
		slog.Log(nil, slog.LevelInfo, "Expected either 'database' or 'dispense' subcommand")
		os.Exit(1)
	}
}
