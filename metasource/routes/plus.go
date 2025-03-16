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
	var name string
	var task string
	var list []string
	var rslt_plus *[]sxml.UnitPrimary
	var wait sync.WaitGroup
	var rslt *[]dict.UnitPrimary
	var expt error
	var utbs dict.UnitBase
	var jobs []string
	var okay bool
	var tkut string

	rslt = &[]dict.UnitPrimary{}

	jobs = []string{"suggests", "enhances", "requires", "provides", "obsoletes", "conflicts", "recommends", "supplements"}

	list = strings.Split(r.URL.Path, "/")

	if len(list) != 4 {
		w.WriteHeader(http.StatusBadRequest)
		_, expt = fmt.Fprintf(w, fmt.Sprintf("%d: %s", http.StatusBadRequest, http.StatusText(http.StatusBadRequest)))
		if expt != nil {
			w.WriteHeader(http.StatusInternalServerError)
			_, expt = fmt.Fprintf(w, fmt.Sprintf("%d: %s", http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError)))
			return
		}
		return
	}

	task = list[2]
	for _, tkut = range jobs {
		if task == tkut {
			okay = true
			break
		}
	}

	if !okay {
		w.WriteHeader(http.StatusBadRequest)
		_, expt = fmt.Fprintf(w, fmt.Sprintf("%d: %s", http.StatusBadRequest, http.StatusText(http.StatusBadRequest)))
		if expt != nil {
			w.WriteHeader(http.StatusInternalServerError)
			_, expt = fmt.Fprintf(w, fmt.Sprintf("%d: %s", http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError)))
			return
		}
		return
	}

	name = list[3]
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
	go runner.ImportPlus(&wait, &rslt_plus, &name, &task)
	wait.Wait()

	if rslt_plus == nil {
		w.WriteHeader(http.StatusNotFound)
		_, expt = fmt.Fprintf(w, fmt.Sprintf("%d: %s", http.StatusNotFound, http.StatusText(http.StatusNotFound)))
		if expt != nil {
			w.WriteHeader(http.StatusInternalServerError)
			_, expt = fmt.Fprintf(w, fmt.Sprintf("%d: %s", http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError)))
			return
		}
		return
	}

	for _, unit := range *rslt_plus {
		var data dict.UnitPrimary

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

		*rslt = append(*rslt, data)
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
