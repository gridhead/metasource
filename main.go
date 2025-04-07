package main

import (
	"flag"
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"log/slog"
	"metasource/metasource/config"
	"metasource/metasource/driver"
	"metasource/metasource/routes"
	"net/http"
	"os"
)

func main() {
	var expt error
	var lglvtext, location *string
	var database, dispense *flag.FlagSet

	lglvtext = flag.String("loglevel", "info", "Set the application loglevel")
	location = flag.String("location", config.DBFOLDER, "Set the database location")
	flag.Parse()

	config.DBFOLDER = *location

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

		var router *chi.Mux
		var server *http.Server

		router = chi.NewRouter()
		router.Use(middleware.Logger)
		router.Use(middleware.Recoverer)

		router.Get("/branches", routes.RetrieveBranches)
		router.Get("/{vers}/changelog/{name}", routes.RetrieveOther)
		router.Get("/{vers}/pkg/{name}", routes.RetrievePrimary)
		router.Get("/{vers}/files/{name}", routes.RetrieveFileList)
		router.Get("/{vers}/srcpkg/{name}", routes.RetrieveSrce)
		router.Get("/{vers}/{rela}/{name}", routes.RetrieveRelation)

		server = &http.Server{Addr: ":8080", Handler: router}

		expt = server.ListenAndServe()
		if expt != nil {
			slog.Log(nil, slog.LevelError, fmt.Sprintf("Error occurred. %s.", expt.Error()))
		}
	default:
		slog.Log(nil, slog.LevelError, "Invalid subcommand")
		slog.Log(nil, slog.LevelInfo, "Expected either 'database' or 'dispense' subcommand")
		os.Exit(1)
	}
}
