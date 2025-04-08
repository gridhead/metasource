package driver

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"metasource/metasource/config"
	"metasource/metasource/models/rels"
	"net/http"
	"net/url"
	"path"
	"time"
)

func ListBranches(status string) ([]string, error) {
	var burl, link string
	var prms url.Values
	var expt error
	var oper *http.Client
	var rqst *http.Request
	var resp *http.Response
	var list []string
	var care []string
	var body []byte
	var rels rels.Releases

	care = []string{"FEDORA", "FEDORA-EPEL", "FEDORA-EPEL-NEXT"}

	burl = path.Join(config.BODHIURL, "releases")
	prms = url.Values{"state": {status}}
	link = fmt.Sprintf("%s?%s", burl, prms.Encode())

	rqst, expt = http.NewRequest("GET", link, nil)
	if expt != nil {
		return list, expt
	}

	oper = &http.Client{Timeout: time.Second * 60}
	resp, expt = oper.Do(rqst)
	if expt != nil || resp.StatusCode != 200 {
		if expt != nil {
			return list, expt
		}
		if resp.StatusCode != 200 {
			return list, errors.New(fmt.Sprintf("%s", resp.Status))
		}
	}
	defer resp.Body.Close()

	body, expt = io.ReadAll(resp.Body)
	if expt != nil {
		return list, expt
	}

	expt = json.Unmarshal(body, &rels)
	if expt != nil {
		return list, expt
	}

	for _, item := range rels.List {
		for _, word := range care {
			if item.IdPrefix == word {
				list = append(list, item.Branch)
			}
		}
	}

	return list, nil
}
