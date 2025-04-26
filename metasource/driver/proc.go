package driver

import (
	"context"
	"encoding/xml"
	"fmt"
	"io"
	"log/slog"
	"metasource/metasource/config"
	"metasource/metasource/models/home"
	"metasource/metasource/models/sxml"
	"metasource/metasource/reader"
	"net/http"
	"strings"
	"sync"
	"time"
)

func HandleRepositories(unit *home.LinkUnit) error {
	var mdlink, name, path, loca string
	var prmyinpt, fileinpt, othrinpt string
	var prmyname, filename, othrname string
	var prmypath, filepath, othrpath string
	var expt error
	var oper *http.Client
	var rqst *http.Request
	var resp *http.Response
	var repo sxml.RepoMD
	var hash home.Checksum
	var body []byte
	var list []home.FileUnit
	var file home.FileUnit
	var castupDownload, entireDownload int
	var castupWithdraw, entireWithdraw int
	var castupChecksum, entireChecksum int
	var castupGenerate, entireGenerate int
	var castupSignalDB, entireSignalDB int
	var wait sync.WaitGroup
	var pack int64

	loca = fmt.Sprintf("%s/.%s-%s", config.DBFOLDER, unit.Name, GenerateIdentity(&config.RANDOM_LENGTH))

	entireDownload = 3
	entireWithdraw = 3
	entireChecksum = 3
	entireGenerate = 3
	entireSignalDB = 1
	mdlink = fmt.Sprintf("%s/repomd.xml", unit.Link)

	rqst, expt = http.NewRequest("GET", mdlink, nil)
	if expt != nil {
		return expt
	}

	oper = &http.Client{Timeout: time.Second * 60}
	resp, expt = oper.Do(rqst)
	if expt != nil || resp.StatusCode != 200 {
		if expt != nil {
			return expt
		}
		if resp.StatusCode != 200 {
			return fmt.Errorf("%s", resp.Status)
		}
	}
	defer resp.Body.Close()

	body, expt = io.ReadAll(resp.Body)
	if expt != nil {
		return expt
	}

	expt = xml.Unmarshal(body, &repo)
	if expt != nil {
		return expt
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
		file = home.FileUnit{Name: name, Path: path, Type: item.Type, Hash: hash, Keep: true}
		list = append(list, file)
	}

	expt = InitPath(&unit.Name, &loca)
	if expt != nil {
		return expt
	}

	for indx := range list {
		if !list[indx].Keep {
			slog.Log(context.Background(), slog.LevelDebug, fmt.Sprintf("[%s] Processing rejected as earlier midphase failed for %s", unit.Name, list[indx].Name))
			continue
		}

		expt = DownloadRepositories(&list[indx], &unit.Name, 0, &castupDownload, &loca)
		if expt != nil {
			slog.Log(context.Background(), slog.LevelDebug, fmt.Sprintf("[%s] Download failed for %s due to %s", unit.Name, list[indx].Name, expt))
		} else {
			slog.Log(context.Background(), slog.LevelDebug, fmt.Sprintf("[%s] Download complete for %s", unit.Name, list[indx].Name))
		}
	}

	if castupDownload == entireDownload {
		slog.Log(context.Background(), slog.LevelInfo, fmt.Sprintf("[%s] Metadata download complete", unit.Name))
	} else {
		slog.Log(context.Background(), slog.LevelError, fmt.Sprintf("[%s] Metadata download failed", unit.Name))
	}

	for indx := range list {
		if !list[indx].Keep {
			slog.Log(context.Background(), slog.LevelDebug, fmt.Sprintf("[%s] Processing rejected as earlier midphase failed for %s", unit.Name, list[indx].Name))
			continue
		}

		wait.Add(1)
		go WithdrawArchives(&list[indx], &unit.Name, &wait, &castupWithdraw, &loca)
	}
	wait.Wait()

	if castupWithdraw == entireWithdraw {
		slog.Log(context.Background(), slog.LevelInfo, fmt.Sprintf("[%s] Metadata extraction complete", unit.Name))
	} else {
		slog.Log(context.Background(), slog.LevelError, fmt.Sprintf("[%s] Metadata extraction failed", unit.Name))
	}

	for indx := range list {
		if !list[indx].Keep {
			slog.Log(context.Background(), slog.LevelDebug, fmt.Sprintf("[%s] Processing rejected as earlier midphase failed for %s", unit.Name, list[indx].Name))
			continue
		}

		wait.Add(1)
		go VerifyChecksum(&list[indx], &unit.Name, &wait, &castupChecksum)
	}
	wait.Wait()

	if castupChecksum == entireChecksum {
		slog.Log(context.Background(), slog.LevelInfo, fmt.Sprintf("[%s] Checksum verification complete", unit.Name))
	} else {
		slog.Log(context.Background(), slog.LevelError, fmt.Sprintf("[%s] Checksum verification failed", unit.Name))
	}

	for indx := range list {
		if !list[indx].Keep {
			slog.Log(context.Background(), slog.LevelDebug, fmt.Sprintf("[%s] Processing rejected as earlier midphase failed for %s", unit.Name, list[indx].Name))
			continue
		}

		switch list[indx].Type {
		case "primary", "filelists", "other":
			path = list[indx].Path
			list[indx].Name = strings.Replace(list[indx].Name, ".xml", ".sqlite", -1)
			list[indx].Path = fmt.Sprintf("%s/%s", loca, list[indx].Name)
			list[indx].Keep = true
			switch list[indx].Type {
			case "primary":
				prmyinpt = path
				prmyname = list[indx].Name
				prmypath = list[indx].Path
			case "filelists":
				fileinpt = path
				filename = list[indx].Name
				filepath = list[indx].Path
			case "other":
				othrinpt = path
				othrname = list[indx].Name
				othrpath = list[indx].Path
			}
		default:
			continue
		}
	}

	pack, expt = reader.MakeDatabase(&unit.Name, &castupGenerate, &prmyinpt, &fileinpt, &othrinpt, &prmyname, &filename, &othrname, &prmypath, &filepath, &othrpath)
	if expt == nil && castupGenerate == entireGenerate {
		slog.Log(context.Background(), slog.LevelInfo, fmt.Sprintf("[%s] Database generation complete with %d package(s)", unit.Name, pack))
	} else {
		slog.Log(context.Background(), slog.LevelError, fmt.Sprintf("[%s] Database generation failed due to %s", unit.Name, expt.Error()))
	}

	for indx := range list {
		if !list[indx].Keep {
			slog.Log(context.Background(), slog.LevelDebug, fmt.Sprintf("[%s] Processing rejected as earlier midphase failed for %s", unit.Name, list[indx].Name))
			continue
		}

		if list[indx].Type != "primary" {
			continue
		}

		expt = GenerateSignal(&list[indx], &castupSignalDB)
		if expt != nil {
			slog.Log(context.Background(), slog.LevelDebug, fmt.Sprintf("[%s] Indexing failed for %s due to %s", unit.Name, list[indx].Name, expt))
		} else {
			slog.Log(context.Background(), slog.LevelDebug, fmt.Sprintf("[%s] Indexing complete for %s", unit.Name, list[indx].Name))
		}
	}

	if castupSignalDB == entireSignalDB {
		slog.Log(context.Background(), slog.LevelInfo, fmt.Sprintf("[%s] Database indexing complete", unit.Name))
	} else {
		slog.Log(context.Background(), slog.LevelError, fmt.Sprintf("[%s] Database indexing failed", unit.Name))
	}

	expt = KillTemp(&unit.Name, &loca)
	if expt != nil {
		return expt
	}

	return nil
}
