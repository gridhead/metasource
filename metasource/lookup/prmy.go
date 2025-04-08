package lookup

import (
	"database/sql"
	"errors"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"metasource/metasource/config"
	"metasource/metasource/models/home"
	"os"
	"path/filepath"
)

func RetrievePrmy(vers *string, name *string) (home.PackUnit, string, error) {
	var base *sql.DB
	var rows *sql.Rows
	var stmt *sql.Stmt
	var expt error
	var item, path, sqlq string
	var exst bool
	var rslt home.PackUnit
	var list []string

	list = []string{"updates-testing", "updates", "testing", ""}

	for _, item = range list {
		switch item {
		case "updates-testing", "updates", "testing":
			path = filepath.Join(config.DBFOLDER, fmt.Sprintf("metasource-%s-%s-primary.sqlite", *vers, item))
		default:
			path = filepath.Join(config.DBFOLDER, fmt.Sprintf("metasource-%s-primary.sqlite", *vers))
		}
		_, expt = os.Stat(path)
		if os.IsNotExist(expt) {
			continue
		}
		exst = true
	}

	if !exst {
		return rslt, item, errors.New("database file does not exist")
	}

	base, expt = sql.Open("sqlite3", path)
	if expt != nil {
		return rslt, item, expt
	}
	defer base.Close()

	sqlq = fmt.Sprintf(config.OBTAIN_PACKAGE)

	stmt, expt = base.Prepare(sqlq)
	if expt != nil {
		return rslt, item, expt
	}
	defer stmt.Close()

	rows, expt = stmt.Query(*name)
	if expt != nil {
		return rslt, item, expt
	}
	defer rows.Close()

	for rows.Next() {
		expt = rows.Scan(&rslt.Key, &rslt.Id, &rslt.Name, &rslt.Source, &rslt.Epoch, &rslt.Version, &rslt.Release, &rslt.Arch, &rslt.Summary, &rslt.Desc, &rslt.Link)
		if expt != nil {
			return rslt, item, expt
		}
		break
	}

	expt = rows.Err()
	if expt != nil {
		return rslt, item, expt
	}

	return rslt, item, nil
}
