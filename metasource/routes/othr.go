package routes

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"metasource/metasource/models/dict"
	"metasource/metasource/models/sxml"
	"metasource/metasource/runner"
	"net/http"
	"strings"
	"sync"
)

func RetrieveOther(w http.ResponseWriter, r *http.Request) {
	var name string
	var list []string
	var rslt_other *sxml.UnitOther
	var wait sync.WaitGroup
	var rslt dict.UnitOther
	var expt error

	list = strings.Split(r.URL.Path, "/")

	if len(list) != 4 {
		http.Error(w, fmt.Sprintf("%d: %s", http.StatusBadRequest, http.StatusText(http.StatusBadRequest)), http.StatusBadRequest)
		slog.Log(nil, slog.LevelError, fmt.Sprintf("[%s] <%s> %d - Malformed request", r.Method, r.RequestURI, http.StatusBadRequest))
		return
	}

	name = strings.TrimSpace(list[3])
	if name == "" {
		http.Error(w, fmt.Sprintf("%d: %s", http.StatusNotFound, http.StatusText(http.StatusNotFound)), http.StatusNotFound)
		slog.Log(nil, slog.LevelError, fmt.Sprintf("[%s] <%s> %d - Absent package", r.Method, r.RequestURI, http.StatusNotFound))
		return
	}

	wait.Add(1)
	go runner.ImportOther(&wait, &rslt_other, &name)
	wait.Wait()

	if rslt_other == nil {
		http.Error(w, fmt.Sprintf("%d: %s", http.StatusNotFound, http.StatusText(http.StatusNotFound)), http.StatusNotFound)
		slog.Log(nil, slog.LevelError, fmt.Sprintf("[%s] <%s> %d - Absent result", r.Method, r.RequestURI, http.StatusNotFound))
		return
	}

	rslt.Repo = "release"
	for _, item := range rslt_other.Changelog {
		var unit dict.Changelog
		unit = dict.Changelog{Author: item.Author, Changelog: item.Data, Date: item.Date}
		rslt.Changelogs = append(rslt.Changelogs, unit)
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	expt = json.NewEncoder(w).Encode(rslt)
	if expt != nil {
		http.Error(w, fmt.Sprintf("%d: %s", http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError)), http.StatusInternalServerError)
		slog.Log(nil, slog.LevelError, fmt.Sprintf("[%s] <%s> %d - Marshalling failed", r.Method, r.RequestURI, http.StatusInternalServerError))
		return
	}
	slog.Log(nil, slog.LevelInfo, fmt.Sprintf("[%s] <%s> %d - Result dispatched", r.Method, r.RequestURI, http.StatusOK))
}
