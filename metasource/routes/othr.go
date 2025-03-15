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

	if len(list) > 3 {
		name = list[3]
	} else {
		w.WriteHeader(http.StatusNotFound)
		_, expt = fmt.Fprintf(w, fmt.Sprintf("%d: %s", http.StatusNotFound, http.StatusText(http.StatusNotFound)))
		if expt != nil {
			w.WriteHeader(http.StatusInternalServerError)
			_, expt = fmt.Fprintf(w, fmt.Sprintf("%d: %s", http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError)))
			return
		}
		return
	}

	if name == "" {
		w.WriteHeader(http.StatusNotFound)
		_, expt = fmt.Fprintf(w, fmt.Sprintf("%d: %s", http.StatusNotFound, http.StatusText(http.StatusNotFound)))
		if expt != nil {
			w.WriteHeader(http.StatusInternalServerError)
			_, expt = fmt.Fprintf(w, fmt.Sprintf("%d: %s", http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError)))
			return
		}
		return
	}

	wait.Add(1)
	go runner.ImportOther(&wait, &rslt_other, &name)
	wait.Wait()

	if rslt_other == nil {
		w.WriteHeader(http.StatusNotFound)
		_, expt = fmt.Fprintf(w, fmt.Sprintf("%d: %s", http.StatusNotFound, http.StatusText(http.StatusNotFound)))
		if expt != nil {
			w.WriteHeader(http.StatusInternalServerError)
			_, expt = fmt.Fprintf(w, fmt.Sprintf("%d: %s", http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError)))
			return
		}
		return
	}

	rslt.Repo = "release"
	for _, item := range rslt_other.Changelog {
		var unit dict.Changelog
		unit = dict.Changelog{Author: item.Author, Changelog: item.Data, Date: item.Date}
		rslt.Changelogs = append(rslt.Changelogs, unit)
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)

	expt = json.NewEncoder(w).Encode(rslt)
	if expt != nil {
		w.WriteHeader(http.StatusInternalServerError)
		slog.Log(nil, slog.LevelError, fmt.Sprintf("JSON could not be marshalled. %s", expt.Error()))
		return
	}
}
