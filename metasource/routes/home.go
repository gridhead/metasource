package routes

import (
	_ "embed"
	"fmt"
	"net/http"
)

//go:embed home.html
var homeHTML []byte

func RetrieveHome(w http.ResponseWriter, r *http.Request) {
	var expt error

	w.Header().Set("Content-Type", "text/html")
	_, expt = w.Write(homeHTML)
	if expt != nil {
		http.Error(w, fmt.Sprintf("%d: %s", http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError)), http.StatusInternalServerError)
		return
	}
}
