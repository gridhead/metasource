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

func RetrievePrimary(w http.ResponseWriter, r *http.Request) {
	var name string
	var list []string
	var rslt_primary *sxml.UnitPrimary
	var wait sync.WaitGroup
	var rslt dict.UnitPrimary
	var expt error
	var utbs dict.UnitBase

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
	go runner.ImportPrimary(&wait, &rslt_primary, &name)
	wait.Wait()

	if rslt_primary == nil {
		w.WriteHeader(http.StatusNotFound)
		_, expt = fmt.Fprintf(w, fmt.Sprintf("%d: %s", http.StatusNotFound, http.StatusText(http.StatusNotFound)))
		if expt != nil {
			w.WriteHeader(http.StatusInternalServerError)
			_, expt = fmt.Fprintf(w, fmt.Sprintf("%d: %s", http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError)))
			return
		}
		return
	}

	rslt = dict.UnitPrimary{
		Repo:        "release",
		Arch:        rslt_primary.Arch,
		Epoch:       rslt_primary.Version.Epoch,
		Version:     rslt_primary.Version.Ver,
		Release:     rslt_primary.Version.Rel,
		Summary:     rslt_primary.Summary,
		Description: rslt_primary.Description,
		Basename:    rslt_primary.Name,
		URL:         rslt_primary.URL,
		CoPackages:  []string{rslt_primary.Name},
	}

	for _, item := range rslt_primary.Format.Supplements.Entries {
		utbs = dict.UnitBase{
			Version: item.Ver,
			Epoch:   item.Epoch,
			Release: item.Rel,
			Name:    item.Name,
			Flags:   item.Flags,
		}
		rslt.Supplements = append(rslt.Supplements, utbs)
	}

	for _, item := range rslt_primary.Format.Recommends.Entries {
		utbs = dict.UnitBase{
			Version: item.Ver,
			Epoch:   item.Epoch,
			Release: item.Rel,
			Name:    item.Name,
			Flags:   item.Flags,
		}
		rslt.Recommends = append(rslt.Recommends, utbs)
	}

	for _, item := range rslt_primary.Format.Conflicts.Entries {
		utbs = dict.UnitBase{
			Version: item.Ver,
			Epoch:   item.Epoch,
			Release: item.Rel,
			Name:    item.Name,
			Flags:   item.Flags,
		}
		rslt.Conflicts = append(rslt.Conflicts, utbs)
	}

	for _, item := range rslt_primary.Format.Obsoletes.Entries {
		utbs = dict.UnitBase{
			Version: item.Ver,
			Epoch:   item.Epoch,
			Release: item.Rel,
			Name:    item.Name,
			Flags:   item.Flags,
		}
		rslt.Obsoletes = append(rslt.Obsoletes, utbs)
	}

	for _, item := range rslt_primary.Format.Provides.Entries {
		utbs = dict.UnitBase{
			Version: item.Ver,
			Epoch:   item.Epoch,
			Release: item.Rel,
			Name:    item.Name,
			Flags:   item.Flags,
		}
		rslt.Provides = append(rslt.Provides, utbs)
	}

	for _, item := range rslt_primary.Format.Requires.Entries {
		utbs = dict.UnitBase{
			Version: item.Ver,
			Epoch:   item.Epoch,
			Release: item.Rel,
			Name:    item.Name,
			Flags:   item.Flags,
		}
		rslt.Requires = append(rslt.Requires, utbs)
	}

	for _, item := range rslt_primary.Format.Enhances.Entries {
		utbs = dict.UnitBase{
			Version: item.Ver,
			Epoch:   item.Epoch,
			Release: item.Rel,
			Name:    item.Name,
			Flags:   item.Flags,
		}
		rslt.Enhances = append(rslt.Enhances, utbs)
	}

	for _, item := range rslt_primary.Format.Suggests.Entries {
		utbs = dict.UnitBase{
			Version: item.Ver,
			Epoch:   item.Epoch,
			Release: item.Rel,
			Name:    item.Name,
			Flags:   item.Flags,
		}
		rslt.Suggests = append(rslt.Suggests, utbs)
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
