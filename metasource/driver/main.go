package driver

import (
	"fmt"
	"log/slog"
	"metasource/metasource/models/home"
	"os"
	"path/filepath"
)

func InitPath(loca string) error {
	var expt error

	_, expt = os.Stat(loca)
	if os.IsNotExist(expt) {
		expt = os.MkdirAll(loca, 0755)
		if expt != nil {
			return expt
		}
		expt = os.MkdirAll(filepath.Join(loca, "sxml"), 0755)
		if expt != nil {
			return expt
		}
		expt = os.MkdirAll(filepath.Join(loca, "comp"), 0755)
		if expt != nil {
			return expt
		}
	} else {
		expt = os.RemoveAll(loca)
		if expt != nil {
			return expt
		}
		return InitPath(loca)
	}

	slog.Log(nil, slog.LevelWarn, "Directories initialized")
	return nil
}

func KillTemp(loca string) error {
	var expt error

	expt = os.RemoveAll(filepath.Join(loca, "sxml"))
	if expt != nil {
		return expt
	}
	expt = os.RemoveAll(filepath.Join(loca, "comp"))
	if expt != nil {
		return expt
	}

	slog.Log(nil, slog.LevelWarn, "Directories removed")
	return nil
}

func Database(loca string) error {
	var expt error
	var list []home.LinkUnit
	var item home.LinkUnit

	list, expt = PopulateRepositories()
	if expt != nil {
		return expt
	}

	expt = InitPath(loca)
	if expt != nil {
		return expt
	}

	for _, item = range list {
		expt = HandleRepositories(&item)
		if expt != nil {
			slog.Log(nil, slog.LevelWarn, fmt.Sprintf("[%s] Repository handling failed due to %s", item.Name, expt.Error()))
		}
	}

	expt = KillTemp(loca)
	if expt != nil {
		return expt
	}

	return nil
}
