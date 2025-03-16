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
	go runner.ImportFileList(&wait, &rslt_filelist, &name)
	wait.Wait()

	if rslt_filelist == nil {
		http.Error(w, fmt.Sprintf("%d: %s", http.StatusNotFound, http.StatusText(http.StatusNotFound)), http.StatusNotFound)
		slog.Log(nil, slog.LevelError, fmt.Sprintf("[%s] <%s> %d - Absent result", r.Method, r.RequestURI, http.StatusNotFound))
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
	expt = json.NewEncoder(w).Encode(rslt)
	if expt != nil {
		http.Error(w, fmt.Sprintf("%d: %s", http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError)), http.StatusInternalServerError)
		slog.Log(nil, slog.LevelError, fmt.Sprintf("[%s] <%s> %d - Marshalling failed", r.Method, r.RequestURI, http.StatusInternalServerError))
		return
	}
	slog.Log(nil, slog.LevelInfo, fmt.Sprintf("[%s] <%s> %d - Result dispatched", r.Method, r.RequestURI, http.StatusOK))
}
