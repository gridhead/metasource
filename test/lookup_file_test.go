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

var vers_lookup_file string = "rawhide"
var pack_lookup_file home.PackUnit = home.PackUnit{
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
var repo_lookup_file string = "release"

func TestReadFile_Failure_AbsentDB(t *testing.T) {
	original := config.DBFOLDER
	config.DBFOLDER = fmt.Sprintf("/var/tmp/test-%s", driver.GenerateIdentity(&config.RANDOM_LENGTH))
	defer func() { config.DBFOLDER = original }()

	_, expt := lookup.ReadFile(&vers_lookup_file, &pack_lookup_file, &repo_lookup_file)
	if expt == nil {
		t.Errorf("Received nil, Expected %s", "database file does not exist")
	}
}

func TestReadFile_Failure_FaultyDB(t *testing.T) {
	Path_Init_Faulty(t)

	_, expt := lookup.ReadFile(&vers_lookup_file, &pack_lookup_file, &repo_lookup_file)
	if expt == nil {
		t.Errorf("Received nil, Expected %s", "file is not a database")
	}
}

func TestReadFile_Failure_FaultyDriver(t *testing.T) {
	Path_Init(t, "")

	original := config.DBDRIVER
	config.DBDRIVER = "very-good-database-driver"
	defer func() { config.DBDRIVER = original }()

	_, expt := lookup.ReadFile(&vers_lookup_file, &pack_lookup_file, &repo_lookup_file)
	if expt == nil {
		t.Errorf("Received nil, Expected %s", "sql: unknown driver \"very-good-database-driver\" (forgotten import?)")
	}
}

func TestReadFile_Failure_FaultyPlea(t *testing.T) {
	Path_Init(t, "")

	original := config.OBTAIN_FILEUNIT
	config.OBTAIN_FILEUNIT = "SELECT f.pkgKey FROM filelist f JOIN packages p ON p.pkgId = ? WHERE f.pkgKey = p.pkgKey ORDER BY f.filenames"
	defer func() { config.OBTAIN_FILEUNIT = original }()

	_, expt := lookup.ReadFile(&vers_lookup_file, &pack_lookup_file, &repo_lookup_file)
	if expt == nil {
		t.Errorf("Received nil, Expected %s", "sql: expected 1 destination arguments in Scan, not 4")
	}
}
