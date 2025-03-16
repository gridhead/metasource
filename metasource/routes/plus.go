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

func RetrievePlus(w http.ResponseWriter, r *http.Request) {
	var name, task, tkut string
	var jobs, list []string
	var rslt_plus *[]sxml.UnitPrimary
	var wait sync.WaitGroup
	var rslt []dict.UnitPrimary
	var expt error
	var okay bool

	rslt = []dict.UnitPrimary{}
	jobs = []string{"suggests", "enhances", "requires", "provides", "obsoletes", "conflicts", "recommends", "supplements"}
	list = strings.Split(r.URL.Path, "/")

	if len(list) != 4 {
		http.Error(w, fmt.Sprintf("%d: %s", http.StatusBadRequest, http.StatusText(http.StatusBadRequest)), http.StatusBadRequest)
		slog.Log(nil, slog.LevelError, fmt.Sprintf("[%s] <%s> %d - Malformed request", r.Method, r.RequestURI, http.StatusBadRequest))
		return
	}

	task = strings.TrimSpace(list[2])
	for _, tkut = range jobs {
		if task == tkut {
			okay = true
			break
		}
	}

	if !okay {
		http.Error(w, fmt.Sprintf("%d: %s", http.StatusBadRequest, http.StatusText(http.StatusBadRequest)), http.StatusBadRequest)
		slog.Log(nil, slog.LevelError, fmt.Sprintf("[%s] <%s> %d - Invalid operation", r.Method, r.RequestURI, http.StatusBadRequest))
		return
	}

	name = strings.TrimSpace(list[3])
	if name == "" {
		http.Error(w, fmt.Sprintf("%d: %s", http.StatusNotFound, http.StatusText(http.StatusNotFound)), http.StatusNotFound)
		slog.Log(nil, slog.LevelError, fmt.Sprintf("[%s] <%s> %d - Absent package", r.Method, r.RequestURI, http.StatusNotFound))
		return
	}

	wait.Add(1)
	go runner.ImportPlus(&wait, &rslt_plus, &name, &task)
	wait.Wait()

	if rslt_plus == nil {
		http.Error(w, fmt.Sprintf("%d: %s", http.StatusNotFound, http.StatusText(http.StatusNotFound)), http.StatusNotFound)
		slog.Log(nil, slog.LevelError, fmt.Sprintf("[%s] <%s> %d - Absent result", r.Method, r.RequestURI, http.StatusNotFound))
		return
	}

	for _, unit := range *rslt_plus {
		var data dict.UnitPrimary
		var utbs dict.UnitBase

		data = dict.UnitPrimary{
			Repo:        "release",
			Arch:        unit.Arch,
			Epoch:       unit.Version.Epoch,
			Version:     unit.Version.Ver,
			Release:     unit.Version.Rel,
			Summary:     unit.Summary,
			Description: unit.Description,
			Basename:    unit.Name,
			URL:         unit.URL,
			CoPackages:  []string{unit.Name},
		}

		data.Supplements = []dict.UnitBase{}
		for _, item := range unit.Format.Supplements.Entries {
			utbs = dict.UnitBase{
				Version: item.Ver,
				Epoch:   item.Epoch,
				Release: item.Rel,
				Name:    item.Name,
				Flags:   item.Flags,
			}
			data.Supplements = append(data.Supplements, utbs)
		}

		data.Recommends = []dict.UnitBase{}
		for _, item := range unit.Format.Recommends.Entries {
			utbs = dict.UnitBase{
				Version: item.Ver,
				Epoch:   item.Epoch,
				Release: item.Rel,
				Name:    item.Name,
				Flags:   item.Flags,
			}
			data.Recommends = append(data.Recommends, utbs)
		}

		data.Conflicts = []dict.UnitBase{}
		for _, item := range unit.Format.Conflicts.Entries {
			utbs = dict.UnitBase{
				Version: item.Ver,
				Epoch:   item.Epoch,
				Release: item.Rel,
				Name:    item.Name,
				Flags:   item.Flags,
			}
			data.Conflicts = append(data.Conflicts, utbs)
		}

		data.Obsoletes = []dict.UnitBase{}
		for _, item := range unit.Format.Obsoletes.Entries {
			utbs = dict.UnitBase{
				Version: item.Ver,
				Epoch:   item.Epoch,
				Release: item.Rel,
				Name:    item.Name,
				Flags:   item.Flags,
			}
			data.Obsoletes = append(data.Obsoletes, utbs)
		}

		data.Provides = []dict.UnitBase{}
		for _, item := range unit.Format.Provides.Entries {
			utbs = dict.UnitBase{
				Version: item.Ver,
				Epoch:   item.Epoch,
				Release: item.Rel,
				Name:    item.Name,
				Flags:   item.Flags,
			}
			data.Provides = append(data.Provides, utbs)
		}

		data.Requires = []dict.UnitBase{}
		for _, item := range unit.Format.Requires.Entries {
			utbs = dict.UnitBase{
				Version: item.Ver,
				Epoch:   item.Epoch,
				Release: item.Rel,
				Name:    item.Name,
				Flags:   item.Flags,
			}
			data.Requires = append(data.Requires, utbs)
		}

		data.Enhances = []dict.UnitBase{}
		for _, item := range unit.Format.Enhances.Entries {
			utbs = dict.UnitBase{
				Version: item.Ver,
				Epoch:   item.Epoch,
				Release: item.Rel,
				Name:    item.Name,
				Flags:   item.Flags,
			}
			data.Enhances = append(data.Enhances, utbs)
		}

		data.Suggests = []dict.UnitBase{}
		for _, item := range unit.Format.Suggests.Entries {
			utbs = dict.UnitBase{
				Version: item.Ver,
				Epoch:   item.Epoch,
				Release: item.Rel,
				Name:    item.Name,
				Flags:   item.Flags,
			}
			data.Suggests = append(data.Suggests, utbs)
		}

		rslt = append(rslt, data)
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
