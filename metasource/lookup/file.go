package lookup

import (
	"database/sql"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"metasource/metasource/config"
	"metasource/metasource/models/home"
	"os"
	"path/filepath"
)

func RetrieveFile(vers *string, pack *home.PackUnit, repo *string) (home.FilelistRslt, error) {
	var base *sql.DB
	var rows *sql.Rows
	var stmt *sql.Stmt
	var expt error
	var path, sqlq string
	var flit home.FilelistUnit
	var rslt home.FilelistRslt

	rslt = home.FilelistRslt{List: []home.FilelistUnit{}}

	switch *repo {
	case "updates-testing", "updates", "testing":
		path = filepath.Join(config.DBFOLDER, fmt.Sprintf("metasource-%s-%s-filelists.sqlite", *vers, *repo))
	default:
		path = filepath.Join(config.DBFOLDER, fmt.Sprintf("metasource-%s-filelists.sqlite", *vers))
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

	sqlq = fmt.Sprintf(config.OBTAIN_FILEUNIT)

	stmt, expt = base.Prepare(sqlq)
	if expt != nil {
		return rslt, expt
	}
	defer stmt.Close()

	rows, expt = stmt.Query(pack.Id)
	if expt != nil {
		return rslt, expt
	}
	defer rows.Close()

	for rows.Next() {
		expt = rows.Scan(&flit.Key, &flit.Directory, &flit.Name, &flit.Type)
		if expt != nil {
			return rslt, expt
		}
		rslt.List = append(rslt.List, flit)
	}

	expt = rows.Err()
	if expt != nil {
		return rslt, expt
	}

	return rslt, expt
}
