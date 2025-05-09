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
	"strings"
	"time"
)

func DownloadRepositories(unit *home.FileUnit, vers *string, stab int64, cast *int, loca *string) error {
	if stab >= config.ATTEMPTS {
		unit.Keep = false
		return errors.New("most attempts failed")
	}

	var expt error
	var urlx, path string
	var head, name string
	var file *os.File
	var oper *http.Client
	var rqst *http.Request
	var resp *http.Response

	head = strings.Split(unit.Name, ".")[0]
	name = strings.Replace(unit.Name, head, fmt.Sprintf(config.FILENAME, *vers, unit.Type), -1)
	urlx = fmt.Sprintf("%s", unit.Path)
	path = fmt.Sprintf("%s/comp/%s", *loca, name)

	file, expt = os.Create(path)
	if expt != nil {
		stab += 1
		slog.Log(nil, slog.LevelDebug, fmt.Sprintf("[%s] Stab failed for %s due to %s", *vers, unit.Name, expt.Error()))
		return DownloadRepositories(unit, vers, stab, cast, loca)
	}
	defer file.Close()

	oper = &http.Client{Timeout: time.Second * 60}
	rqst, expt = http.NewRequest("GET", urlx, nil)
	if expt != nil {
		stab += 1
		slog.Log(nil, slog.LevelDebug, fmt.Sprintf("[%s] Stab failed for %s due to %s", *vers, unit.Name, expt.Error()))
		return DownloadRepositories(unit, vers, stab, cast, loca)
	}

	resp, expt = oper.Do(rqst)
	if expt != nil {
		stab += 1
		slog.Log(nil, slog.LevelDebug, fmt.Sprintf("[%s] Stab failed for %s due to %s", *vers, unit.Name, expt.Error()))
		return DownloadRepositories(unit, vers, stab, cast, loca)
	}
	defer resp.Body.Close()

	_, expt = io.Copy(file, resp.Body)
	if expt != nil {
		stab += 1
		slog.Log(nil, slog.LevelDebug, fmt.Sprintf("[%s] Stab failed for %s due to %s", *vers, unit.Name, expt.Error()))
		return DownloadRepositories(unit, vers, stab, cast, loca)
	}

	unit.Keep = true
	unit.Name = name
	unit.Path = path
	*cast++
	slog.Log(nil, slog.LevelDebug, fmt.Sprintf("[%s] Stab complete for %s", *vers, unit.Name))
	return nil
}
