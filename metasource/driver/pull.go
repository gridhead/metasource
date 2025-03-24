package driver

import (
	"errors"
	"fmt"
	"io"
	"log/slog"
	"metasource/metasource/config"
	"metasource/metasource/models/home"
	"net/http"
	"os"
	"time"
)

func DownloadRepositories(unit *home.FileUnit, name *string, stab int64) (string, error) {
	if stab >= config.ATTEMPTS {
		return "", errors.New("most attempts failed")
	}

	var expt error
	var urlx, loca string
	var file *os.File
	var oper *http.Client
	var rqst *http.Request
	var resp *http.Response

	urlx = fmt.Sprintf("%s", unit.Path)
	loca = fmt.Sprintf("%s/comp/%s", config.DBFOLDER, *name)

	file, expt = os.Create(loca)
	if expt != nil {
		stab += 1
		slog.Log(nil, slog.LevelDebug, fmt.Sprintf("Stab failed due to %s", expt.Error()))
		return DownloadRepositories(unit, name, stab)
	}
	defer file.Close()

	oper = &http.Client{Timeout: time.Second * 60}
	rqst, expt = http.NewRequest("GET", urlx, nil)
	if expt != nil {
		stab += 1
		slog.Log(nil, slog.LevelDebug, fmt.Sprintf("Stab failed due to %s", expt.Error()))
		return DownloadRepositories(unit, name, stab)
	}

	resp, expt = oper.Do(rqst)
	if expt != nil {
		stab += 1
		slog.Log(nil, slog.LevelDebug, fmt.Sprintf("Stab failed due to %s", expt.Error()))
		return DownloadRepositories(unit, name, stab)
	}
	defer resp.Body.Close()

	_, expt = io.Copy(file, resp.Body)
	if expt != nil {
		stab += 1
		slog.Log(nil, slog.LevelDebug, fmt.Sprintf("Stab failed due to %s", expt.Error()))
		return DownloadRepositories(unit, name, stab)
	}

	return loca, nil
}
