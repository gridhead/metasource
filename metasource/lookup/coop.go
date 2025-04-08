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

func RetrieveCoop(vers *string, pack *home.PackUnit, repo *string) ([]string, error) {
	var base *sql.DB
	var rows *sql.Rows
	var stmt *sql.Stmt
	var expt error
	var path, sqlq string
	var cpit sql.NullString
	var rslt []string

	rslt = []string{}

	switch *repo {
	case "updates-testing", "updates", "testing":
		path = filepath.Join(config.DBFOLDER, fmt.Sprintf("metasource-%s-%s-primary.sqlite", *vers, *repo))
	default:
		path = filepath.Join(config.DBFOLDER, fmt.Sprintf("metasource-%s-primary.sqlite", *vers))
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

	sqlq = fmt.Sprintf(config.OBTAIN_CO_PACKAGE)

	stmt, expt = base.Prepare(sqlq)
	if expt != nil {
		return rslt, expt
	}
	defer stmt.Close()

	rows, expt = stmt.Query(pack.Source.String)
	if expt != nil {
		return rslt, expt
	}
	defer rows.Close()

	for rows.Next() {
		expt = rows.Scan(&cpit)
		if expt != nil {
			return rslt, expt
		}
		if cpit.Valid {
			rslt = append(rslt, cpit.String)
		}
	}

	expt = rows.Err()
	if expt != nil {
		return rslt, expt
	}

	return rslt, expt
}
