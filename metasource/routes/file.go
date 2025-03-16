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

func RetrieveFileList(w http.ResponseWriter, r *http.Request) {
	var name string
	var list []string
	var rslt_filelist *sxml.UnitFileList
	var wait sync.WaitGroup
	var rslt dict.UnitFileList
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
	go runner.ImportFileList(&wait, &rslt_filelist, &name)
	wait.Wait()

	if rslt_filelist == nil {
		w.WriteHeader(http.StatusNotFound)
		_, expt = fmt.Fprintf(w, fmt.Sprintf("%d: %s", http.StatusNotFound, http.StatusText(http.StatusNotFound)))
		if expt != nil {
			w.WriteHeader(http.StatusInternalServerError)
			_, expt = fmt.Fprintf(w, fmt.Sprintf("%d: %s", http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError)))
			return
		}
		return
	}

	rslt = dict.UnitFileList{}
	rslt.Repo = "release"

	for _, item := range rslt_filelist.List {
		var unit dict.File
		var temp, file, loca string
		var path []string
		var indx int
		var fdtp string

		temp = item.Data
		path = strings.Split(temp, "/")
		file = path[len(path)-1]
		for indx = 0; indx < len(path)-1; indx++ {
			loca = loca + path[indx] + "/"
		}
		if item.Type == "dir" {
			fdtp = "d"
		} else {
			fdtp = "f"
		}

		unit.FileTypes = fdtp
		unit.DirName = loca
		unit.FileNames = file

		rslt.Files = append(rslt.Files, unit)
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
