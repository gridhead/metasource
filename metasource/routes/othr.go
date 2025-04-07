package routes

import (
	"encoding/json"
	"fmt"
	"github.com/go-chi/chi/v5"
	"metasource/metasource/lookup"
	"metasource/metasource/models/dict"
	"metasource/metasource/models/home"
	"net/http"
)

func RetrieveOther(w http.ResponseWriter, r *http.Request) {
	var name, vers, repo string
	var rslt dict.UnitOther
	var data home.OthrRslt
	var pack home.PackUnit
	var expt error

	name = chi.URLParam(r, "name")
	vers = chi.URLParam(r, "vers")

	if name == "" || vers == "" {
		http.Error(w, fmt.Sprintf("%d: %s", http.StatusBadRequest, http.StatusText(http.StatusBadRequest)), http.StatusBadRequest)
		return
	}

	pack, repo, expt = lookup.RetrievePrmy(&vers, &name)
	if expt != nil {
		http.Error(w, fmt.Sprintf("%d: %s", http.StatusBadRequest, http.StatusText(http.StatusBadRequest)), http.StatusBadRequest)
		return
	}

	data, expt = lookup.RetrieveOthr(&vers, &pack, &repo)
	if expt != nil {
		http.Error(w, fmt.Sprintf("%d: %s", http.StatusBadRequest, http.StatusText(http.StatusBadRequest)), http.StatusBadRequest)
		return
	}

	rslt = dict.UnitOther{}
	rslt.Repo = repo
	if rslt.Repo == "" {
		rslt.Repo = "release"
	}

	for _, item := range data.List {
		var unit dict.Changelog
		unit = dict.Changelog{Author: item.Author.String, Changelog: item.Text.String, Date: uint64(item.Date)}
		rslt.Changelogs = append(rslt.Changelogs, unit)
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	expt = json.NewEncoder(w).Encode(rslt)
	if expt != nil {
		http.Error(w, fmt.Sprintf("%d: %s", http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError)), http.StatusInternalServerError)
		return
	}
}
