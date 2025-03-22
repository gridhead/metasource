package driver

import (
	"encoding/xml"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"metasource/metasource/config"
	"metasource/metasource/models/home"
	"metasource/metasource/models/sxml"
	"net/http"
	"strings"
	"time"
)

func HandleRepositories(unit *home.LinkUnit) (bool, error) {
	var mdlink, name, path, head string
	var expt error
	var oper *http.Client
	var rqst *http.Request
	var resp *http.Response
	var repo sxml.RepoMD
	var hash home.Checksum
	var body []byte
	var list []home.FileUnit
	var loca home.FileUnit
	var castupDownload, entireDownload int

	entireDownload = 3
	mdlink = fmt.Sprintf("%s/repomd.xml", unit.Link)

	rqst, expt = http.NewRequest("GET", mdlink, nil)
	if expt != nil {
		return false, expt
	}

	oper = &http.Client{Timeout: time.Second * 60}
	resp, expt = oper.Do(rqst)
	if expt != nil || resp.StatusCode != 200 {
		if expt != nil {
			return false, expt
		}
		if resp.StatusCode != 200 {
			return false, errors.New(fmt.Sprintf("%s", resp.Status))
		}
	}
	defer resp.Body.Close()

	body, expt = io.ReadAll(resp.Body)
	if expt != nil {
		return false, expt
	}

	expt = xml.Unmarshal(body, &repo)
	if expt != nil {
		return false, expt
	}

	for _, item := range repo.Data {
		if item.Type != "primary" && item.Type != "filelists" && item.Type != "other" {
			continue
		}

		if !strings.Contains(item.Location.Href, ".xml") {
			continue
		}

		name = strings.Replace(item.Location.Href, "repodata/", "", -1)
		path = fmt.Sprintf("%s/%s", unit.Link, name)
		hash = home.Checksum{Data: item.ChecksumOpen.Data, Type: item.ChecksumOpen.Type}
		loca = home.FileUnit{Name: name, Path: path, Type: item.Type, Hash: hash}
		list = append(list, loca)
	}

	for _, item := range list {
		head = strings.Split(item.Name, ".")[0]
		name = strings.Replace(item.Name, head, fmt.Sprintf(config.FILENAME, unit.Name, item.Type), -1)
		path, expt = DownloadRepositories(&item, &name, 0)
		if expt != nil {
			slog.Log(nil, slog.LevelDebug, fmt.Sprintf("[%s] Download failed for %s due to %s", unit.Name, name, expt))
			castupDownload -= 1
		} else {
			slog.Log(nil, slog.LevelDebug, fmt.Sprintf("[%s] Download complete for %s", unit.Name, name))
			castupDownload += 1
		}
		item.Path = path
		item.Name = name
	}

	if castupDownload == entireDownload {
		slog.Log(nil, slog.LevelInfo, fmt.Sprintf("[%s] Metadata download complete", unit.Name))
	} else {
		slog.Log(nil, slog.LevelError, fmt.Sprintf("[%s] Metadata download failed", unit.Name))
	}

	return true, nil
}
