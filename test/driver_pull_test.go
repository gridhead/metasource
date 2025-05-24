package test

import (
	"metasource/metasource/config"
	"metasource/metasource/driver"
	"metasource/metasource/models/home"
	"testing"
)

func TestDownload_Success(t *testing.T) {
	Path_Init_Vacant(t)

	unit := home.LinkUnit{
		Name: "rawhide",
		Link: "https://dl.fedoraproject.org/pub/fedora/linux/development/rawhide/Everything/x86_64/os/repodata/",
	}

	list, expt := driver.ReadMetadata(&unit)
	if expt != nil {
		t.Errorf("Test requires active internet connection")
	}
	if len(list) != 3 {
		t.Errorf("Received %d, Expected %d", len(list), 3)
	}

	castup, entire := 0, 3
	for indx := range list {
		expt = driver.DownloadRepositories(&list[indx], &unit.Name, 0, &castup, &config.DBFOLDER)
		if expt != nil {
			t.Errorf("Test requires active internet connection")
		}
	}

	if castup != entire {
		t.Errorf("Received %d, Expected %d", castup, entire)
	}
}
