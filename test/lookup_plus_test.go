package test

import (
	"database/sql"
	"fmt"
	"metasource/metasource/config"
	"metasource/metasource/driver"
	"metasource/metasource/lookup"
	"metasource/metasource/models/home"
	"testing"
)

var vers_lookup_extn string = "rawhide"
var pack_lookup_extn home.PackUnit = home.PackUnit{
	Key:     1,
	Id:      sql.NullString{Valid: true, String: "28d3c752b8f7f78aae51fb4afee36d5bfdee295df85b3fdacb4bc0357f614784"},
	Name:    sql.NullString{Valid: true, String: "systemd"},
	Source:  sql.NullString{Valid: true, String: "systemd-257.5-6.fc42.src.rpm"},
	Epoch:   sql.NullString{Valid: true, String: "0"},
	Version: sql.NullString{Valid: true, String: "257.5"},
	Release: sql.NullString{Valid: true, String: "6.f42"},
	Arch:    sql.NullString{Valid: true, String: "i686"},
	Summary: sql.NullString{Valid: true, String: "System and Service Manager"},
	Desc:    sql.NullString{Valid: true, String: ""},
	Link:    sql.NullString{Valid: true, String: "https://systemd.io"},
}
var repo_lookup_extn string = "release"

func TestReadExtn_Failure_AbsentDB(t *testing.T) {
	original := config.DBFOLDER
	config.DBFOLDER = fmt.Sprintf("/var/tmp/test-%s", driver.GenerateIdentity(&config.RANDOM_LENGTH))
	defer func() { config.DBFOLDER = original }()

	_, expt := lookup.ReadExtn(&vers_lookup_extn, &pack_lookup_extn, &repo_lookup_extn)
	if expt == nil {
		t.Errorf("Received nil, Expected %s", "database file does not exist")
	}
}

func TestReadExtn_Failure_FaultyDB(t *testing.T) {
	Path_Init_Faulty(t)

	_, expt := lookup.ReadExtn(&vers_lookup_extn, &pack_lookup_extn, &repo_lookup_extn)
	if expt == nil {
		t.Errorf("Received nil, Expected %s", "file is not a database")
	}
}

func TestReadExtn_Failure_FaultyDriver(t *testing.T) {
	Path_Init(t, "")

	original := config.DBDRIVER
	config.DBDRIVER = "very-good-database-driver"
	defer func() { config.DBDRIVER = original }()

	_, expt := lookup.ReadExtn(&vers_lookup_extn, &pack_lookup_extn, &repo_lookup_extn)
	if expt == nil {
		t.Errorf("Received nil, Expected %s", "sql: unknown driver \"very-good-database-driver\" (forgotten import?)")
	}
}

func TestReadExtn_Failure_FaultyPlea(t *testing.T) {
	Path_Init(t, "")

	original := config.OBTAIN_PACKAGE_INFO
	config.OBTAIN_PACKAGE_INFO = "SELECT rowid, name FROM %s WHERE pkgKey = ?"
	defer func() { config.OBTAIN_PACKAGE_INFO = original }()

	_, expt := lookup.ReadExtn(&vers_lookup_extn, &pack_lookup_extn, &repo_lookup_extn)
	if expt == nil {
		t.Errorf("Received nil, Expected %s", "sql: expected 1 destination arguments in Scan, not 4")
	}
}
