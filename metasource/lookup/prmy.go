package lookup

import (
	"database/sql"
	"errors"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"metasource/metasource/config"
	"metasource/metasource/models/home"
	"os"
	"regexp"
)

func RetrievePrmy(vers *string, name *string, actn *string, srcn *string) (home.PackUnit, string, error) {
	var base *sql.DB
	var rows *sql.Rows
	var stmt *sql.Stmt
	var expt error
	var item, path, sqlq string
	var exst bool
	var pkls []home.PackUnit
	var rslt, pkit home.PackUnit
	var list []string
	var expr *regexp.Regexp

	list = []string{"updates-testing", "updates", "testing", ""}

	fmt.Println(config.DBFOLDER)
	fmt.Println("Hello")

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

	if *actn != "" {
		sqlq = fmt.Sprintf(config.OBTAIN_PACKAGE_BY, *actn)

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
	} else if *srcn != "" {
		sqlq = fmt.Sprintf(config.OBTAIN_PACKAGE_BY_SOURCE)

		stmt, expt = base.Prepare(sqlq)
		if expt != nil {
			return rslt, item, expt
		}
		defer stmt.Close()

		rows, expt = stmt.Query(*srcn + "-%")
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
			break
		}

		fmt.Println(len(pkls))

		for _, pkit = range pkls {
			if pkit.Source == *srcn {
				rslt = pkit
				break
			}
		}

		expr = regexp.MustCompile(fmt.Sprintf("%s-[0-9]", regexp.QuoteMeta(*srcn)))
		for _, pkit = range pkls {
			if expr.MatchString(pkit.Source) {
				rslt = pkit
				break
			}
		}

		expt = rows.Err()
		if expt != nil {
			return rslt, item, expt
		}
	} else {
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
	}

	return rslt, item, nil
}
