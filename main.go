package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"log/slog"
	"metasource/metasource/config"
	"metasource/metasource/driver"
	"metasource/metasource/routes"
	"net/http"
	"os"
	"time"
)

func main() {
	var expt error
	var lglvtext, location, port *string
	var database, dispense *flag.FlagSet

	lglvtext = flag.String("loglevel", "info", "Set the application loglevel")
	location = flag.String("location", config.DBFOLDER, "Set the database location")
	flag.Parse()

	config.DBFOLDER = *location

	database = flag.NewFlagSet("database", flag.ExitOnError)
	dispense = flag.NewFlagSet("dispense", flag.ExitOnError)
	port = dispense.String("port", "8080", "Port to run the server on")
	config.SetLogger(lglvtext)

	if flag.NArg() < 1 {
		slog.Log(context.Background(), slog.LevelError, "Invalid subcommand")
		slog.Log(context.Background(), slog.LevelInfo, "Expected either 'database' or 'dispense' subcommand")
		os.Exit(1)
	}

	switch flag.Arg(0) {
	case "database":
		expt = database.Parse(os.Args[2:])
		if expt != nil {
			slog.Log(context.Background(), slog.LevelError, expt.Error())
			os.Exit(1)
		}
		expt = driver.Database()
		if expt != nil {
			slog.Log(context.Background(), slog.LevelError, expt.Error())
			os.Exit(1)
		}
		os.Exit(0)
	case "dispense":
		expt = dispense.Parse(os.Args[2:])
		if expt != nil {
			slog.Log(context.Background(), slog.LevelError, expt.Error())
			os.Exit(1)
		}

		var router *chi.Mux
		var server *http.Server

		router = chi.NewRouter()
		router.Use(middleware.Logger)
		router.Use(middleware.Recoverer)
		router.Use(middleware.StripSlashes)
		router.Use(cors.Handler(cors.Options{
			AllowedOrigins: []string{"*"},
			AllowedMethods: []string{"GET"},
			AllowedHeaders: []string{"*"},
		}))

		router.Get("/", routes.RetrieveHome)
		router.Get("/assets/*", routes.RetrieveStatic)
		router.Get("/branches", routes.RetrieveBranches)
		router.Get("/{vers}/changelog/{name}", routes.RetrieveOthr)
		router.Get("/{vers}/pkg/{name}", routes.RetrievePrmy)
		router.Get("/{vers}/files/{name}", routes.RetrieveFileList)
		router.Get("/{vers}/srcpkg/{name}", routes.RetrieveSrce)
		router.Get("/{vers}/{rela}/{name}", routes.RetrieveRelation)
		server = &http.Server{Addr: ":" + *port, Handler: router, ReadHeaderTimeout: 10 * time.Second}

		expt = server.ListenAndServe()
		if expt != nil {
			slog.Log(context.Background(), slog.LevelError, fmt.Sprintf("Error occurred. %s.", expt.Error()))
			os.Exit(1)
		}
	default:
		slog.Log(context.Background(), slog.LevelError, "Invalid subcommand")
		slog.Log(context.Background(), slog.LevelInfo, "Expected either 'database' or 'dispense' subcommand")
		os.Exit(1)
	}
}
