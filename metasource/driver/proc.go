package driver

import (
	"encoding/xml"
	"errors"
	"fmt"
	"io"
	"metasource/metasource/models/home"
	"metasource/metasource/models/sxml"
	"net/http"
	"strings"
	"time"
)

func HandleRepositories(unit *home.LinkUnit) (bool, error) {
	var mdlink string
	var expt error
	var oper *http.Client
	var rqst *http.Request
	var resp *http.Response
	var repo sxml.RepoMD
	var body []byte
	var list []home.FileUnit

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
		var hash = home.Checksum{
			Data: item.ChecksumOpen.Data,
			Type: item.ChecksumOpen.Type,
		}
		var loca = home.FileUnit{
			Name: strings.Replace(item.Location.Href, "repodata/", "", -1),
			Hash: hash,
		}
		list = append(list, loca)
	}

	fmt.Println(list)

	return true, nil
}
