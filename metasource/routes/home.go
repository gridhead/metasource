package routes

import (
	_ "embed"
	"fmt"
	"github.com/dustin/go-humanize"
	"html/template"
	"metasource/metasource/config"
	"metasource/metasource/lookup"
	"metasource/metasource/models/page"
	"net/http"
	"os"
	"strings"
)

//go:embed shapes/home.html
var homeHTML []byte

func LastModified(list []string) (map[string]page.Vary, error) {
	datelist := map[string]page.Vary{}

	for _, name := range list {
		name_primary := fmt.Sprintf("metasource-%s-primary.sqlite", name)
		location_primary := fmt.Sprintf("%s/%s", config.DBFOLDER, name_primary)
		name_changelog := fmt.Sprintf("metasource-%s-other.sqlite", name)
		location_changelog := fmt.Sprintf("%s/%s", config.DBFOLDER, name_changelog)
		name_filelists := fmt.Sprintf("metasource-%s-filelists.sqlite", name)
		location_filelists := fmt.Sprintf("%s/%s", config.DBFOLDER, name_filelists)
		info_primary, expt := os.Stat(location_primary)
		if expt != nil {
			return datelist, expt
		}
		info_changelog, expt := os.Stat(location_changelog)
		if expt != nil {
			return datelist, expt
		}
		info_filelists, expt := os.Stat(location_filelists)
		if expt != nil {
			return datelist, expt
		}
		datelist[name] = page.Vary{
			Primary: page.Date{
				When: info_primary.ModTime().Format("2006-01-02 15:04:05 MST"),
				Past: humanize.Time(info_primary.ModTime()),
			},
			Changelog: page.Date{
				When: info_changelog.ModTime().Format("2006-01-02 15:04:05 MST"),
				Past: humanize.Time(info_changelog.ModTime()),
			},
			Filelists: page.Date{
				When: info_filelists.ModTime().Format("2006-01-02 15:04:05 MST"),
				Past: humanize.Time(info_filelists.ModTime()),
			},
			Safe: strings.Replace(name, ".", "_", -1),
		}
	}
	return datelist, nil
}

func RetrieveHome(w http.ResponseWriter, r *http.Request) {
	var expt error

	tmpl, expt := template.New("home").Parse(string(homeHTML))
	if expt != nil {
		http.Error(w, fmt.Sprintf("%d: %s", http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError)), http.StatusInternalServerError)
		return
	}

	list, expt := lookup.ReadBranches()
	if expt != nil {
		http.Error(w, fmt.Sprintf("%d: %s", http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError)), http.StatusInternalServerError)
		return
	}

	dict, expt := LastModified(list)
	if expt != nil {
		fmt.Println("Hello %s", expt)
		http.Error(w, fmt.Sprintf("%d: %s", http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError)), http.StatusInternalServerError)
		return
	}

	park := []page.Card{
		{
			Iden: "package",
			Head: "Package",
			Desc: "You can retrieve the information about a specific package on a specific branch by querying its package name.",
			Path: "pkg",
		},
		{
			Iden: "source",
			Head: "Source",
			Desc: "You can retrieve the information about a specific package on a specific branch by querying its source name.",
			Path: "srcpkg",
		},
		{
			Iden: "filelist",
			Head: "Filelist",
			Desc: "You can retrieve the list of files present in a specific package on a specific branch by querying its package name.",
			Path: "files",
		},
		{
			Iden: "changelog",
			Head: "Changelog",
			Desc: "You can retrieve the changelog of a specific package on a specific branch by querying its package name.",
			Path: "changelog",
		},
		{
			Iden: "requires",
			Head: "Requires",
			Desc: "You can retrieve the list of packages requiring a specific package on a specific branch by querying its package name.",
			Path: "requires",
		},
		{
			Iden: "provides",
			Head: "Provides",
			Desc: "You can retrieve the list of packages providing a specific package on a specific branch by querying its package name.",
			Path: "provides",
		},
		{
			Iden: "obsoletes",
			Head: "Obsoletes",
			Desc: "You can retrieve the list of packages obsoleting a specific package on a specific branch by querying its package name.",
			Path: "obsoletes",
		},
		{
			Iden: "conflicts",
			Head: "Conflicts",
			Desc: "You can retrieve the list of packages conflicting a specific package on a specific branch by querying its package name.",
			Path: "conflicts",
		},
		{
			Iden: "enhances",
			Head: "Enhances",
			Desc: "You can retrieve the list of packages enhancing a specific package on a specific branch by querying its package name.",
			Path: "enhances",
		},
		{
			Iden: "recommends",
			Head: "Recommends",
			Desc: "You can retrieve the list of packages recommending a specific package on a specific branch by querying its package name.",
			Path: "recommends",
		},
		{
			Iden: "suggests",
			Head: "Suggests",
			Desc: "You can retrieve the list of packages suggesting a specific package on a specific branch by querying its package name.",
			Path: "suggests",
		},
		{
			Iden: "supplements",
			Head: "Supplements",
			Desc: "You can retrieve the list of packages providing a specific package on a specific branch by querying its package name.",
			Path: "supplements",
		},
	}

	data := page.Page{
		Name: "MetaSource",
		Vers: "v0.1.0",
		Host: "metasource.gridhead.net",
		Conn: "https",
		Dict: dict,
		Park: park,
	}

	w.Header().Set("Content-Type", "text/html")
	expt = tmpl.Execute(w, data)
	if expt != nil {
		http.Error(w, fmt.Sprintf("%d: %s", http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError)), http.StatusInternalServerError)
		return
	}
}
