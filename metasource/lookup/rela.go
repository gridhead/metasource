package lookup

import (
	"database/sql"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"metasource/metasource/config"
	"metasource/metasource/models/home"
	"os"
)

func ReadRelation(vers *string, pack *home.PackUnit, repo *string, relation *string) ([]home.PackUnit, error) {
	var base *sql.DB
	var rows *sql.Rows
	var stmt *sql.Stmt
	var expt error
	var item, path, sqlq string
	var list []string
	var rslt []home.PackUnit
	var pkit home.PackUnit
	var okay bool

	list = []string{
		"supplements",
		"recommends",
		"conflicts",
		"obsoletes",
		"provides",
		"requires",
		"enhances",
		"suggests",
	}

	switch *repo {
	case "updates-testing", "updates", "testing":
		path = fmt.Sprintf("%s/%s", config.DBFOLDER, fmt.Sprintf("metasource-%s-%s-primary.sqlite", *vers, *repo))
	default:
		path = fmt.Sprintf("%s/%s", config.DBFOLDER, fmt.Sprintf("metasource-%s-primary.sqlite", *vers))
	}
	_, expt = os.Stat(path)
	if os.IsNotExist(expt) {
		return rslt, expt
	}

	base, expt = sql.Open("sqlite3", path)
	if expt != nil {
		return rslt, expt
	}
	defer base.Close()

	for _, item = range list {
		if item == *relation {
			okay = true
			break
		}
	}

	if !okay {
		return rslt, fmt.Errorf("mistaken relationship")
	}

	sqlq = fmt.Sprintf(config.OBTAIN_PACKAGE_BY, *relation)

	stmt, expt = base.Prepare(sqlq)
	if expt != nil {
		return rslt, expt
	}
	defer stmt.Close()

	rows, expt = stmt.Query(pack.Name)
	if expt != nil {
		return rslt, expt
	}
	defer rows.Close()

	rslt = []home.PackUnit{}

	for rows.Next() {
		expt = rows.Scan(&pkit.Key, &pkit.Id, &pkit.Name, &pkit.Source, &pkit.Epoch, &pkit.Version, &pkit.Release, &pkit.Arch, &pkit.Summary, &pkit.Desc, &pkit.Link)
		if expt != nil {
			return rslt, expt
		}
		rslt = append(rslt, pkit)
	}

	expt = rows.Err()
	if expt != nil {
		return rslt, expt
	}

	return rslt, expt
}
