package test

import (
	"fmt"
	"metasource/metasource/config"
	"metasource/metasource/driver"
	"metasource/metasource/lookup"
	"testing"
)

func TestReadSrce_Failure_AbsentDB(t *testing.T) {
	original := config.DBFOLDER
	config.DBFOLDER = fmt.Sprintf("/var/tmp/test-%s", driver.GenerateIdentity(&config.RANDOM_LENGTH))
	defer func() { config.DBFOLDER = original }()

	vers, name := "rawhide", "systemd"
	_, _, expt := lookup.ReadSrce(&vers, &name)
	if expt == nil {
		t.Errorf("Received nil, Expected %s", "database file does not exist")
	}
}

func TestReadSrce_Failure_FaultyDB(t *testing.T) {
	Path_Init_Faulty(t)

	vers, name := "rawhide", "systemd"
	_, _, expt := lookup.ReadSrce(&vers, &name)
	if expt == nil {
		t.Errorf("Received nil, Expected %s", "file is not a database")
	}
}

func TestReadSrce_Failure_FaultyDriver(t *testing.T) {
	Path_Init(t, "")

	original := config.DBDRIVER
	config.DBDRIVER = "very-good-database-driver"
	defer func() { config.DBDRIVER = original }()

	vers, name := "rawhide", "systemd"
	_, _, expt := lookup.ReadSrce(&vers, &name)
	if expt == nil {
		t.Errorf("Received nil, Expected %s", "sql: unknown driver \"very-good-database-driver\" (forgotten import?)")
	}
}

func TestReadSrce_Failure_FaultyPlea(t *testing.T) {
	Path_Init(t, "")

	original := config.OBTAIN_PACKAGE_BY_SOURCE
	config.OBTAIN_PACKAGE_BY_SOURCE = "SELECT name, arch FROM packages WHERE rpm_sourcerpm LIKE ? ORDER BY epoch DESC, version DESC, release DESC"
	defer func() { config.OBTAIN_PACKAGE_BY_SOURCE = original }()

	vers, name := "rawhide", "systemd"
	_, _, expt := lookup.ReadSrce(&vers, &name)
	if expt == nil {
		t.Errorf("Received nil, Expected %s", "sql: expected 2 destination arguments in Scan, not 11")
	}
}

func TestReadSrce_Failure_FaultyName(t *testing.T) {
	Path_Init(t, "")

	vers, name := "rawhide", ""
	_, _, expt := lookup.ReadSrce(&vers, &name)
	if expt == nil {
		t.Errorf("Received nil, Expected %s", "no result found")
	}
}
