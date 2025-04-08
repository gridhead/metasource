package lookup

import (
	"database/sql"
	"errors"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"metasource/metasource/config"
	"metasource/metasource/models/home"
	"os"
)

func ReadSrce(vers *string, name *string) (home.PackUnit, string, error) {
	var base *sql.DB
	var rows *sql.Rows
	var stmt *sql.Stmt
	var expt error
	var item, path, sqlq string
	var exst bool
	var pkls []home.PackUnit
	var rslt, pkit home.PackUnit
	var list []string

	list = []string{"updates-testing", "updates", "testing", ""}

	for _, item = range list {
		switch item {
		case "updates-testing", "updates", "testing":
			path = fmt.Sprintf("%s/%s", config.DBFOLDER, fmt.Sprintf("metasource-%s-%s-primary.sqlite", *vers, item))
		default:
			path = fmt.Sprintf("%s/%s", config.DBFOLDER, fmt.Sprintf("metasource-%s-primary.sqlite", *vers))
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

	sqlq = fmt.Sprintf(config.OBTAIN_PACKAGE_BY_SOURCE)

	stmt, expt = base.Prepare(sqlq)
	if expt != nil {
		return rslt, item, expt
	}
	defer stmt.Close()

	rows, expt = stmt.Query(*name + "-%")
	if expt != nil {
		return rslt, item, expt
	}
	defer rows.Close()

	for rows.Next() {
		var pack home.PackUnit
		expt = rows.Scan(&pack.Key, &pack.Id, &pack.Name, &pack.Source, &pack.Epoch, &pack.Version, &pack.Release, &pack.Arch, &pack.Summary, &pack.Desc, &pack.Link)
		if expt != nil {
			return rslt, item, expt
		}
		pkls = append(pkls, pack)
	}

	for _, pkit = range pkls {
		if pkit.Name.String == *name {
			rslt = pkit
			break
		}
	}

	if !rslt.Id.Valid {
		return rslt, item, errors.New("no result found")
	}

	return rslt, item, expt
}
