package driver

import (
	"fmt"
	"log/slog"
	"os"
)

func InitPath(loca string) error {
	var expt error

	_, expt = os.Stat(loca)
	if os.IsNotExist(expt) {
		expt = os.MkdirAll(loca, 0755)
		if expt != nil {
			return expt
		}
		expt = os.MkdirAll(fmt.Sprintf("%s/sxml", loca), 0755)
		if expt != nil {
			return expt
		}
		expt = os.MkdirAll(fmt.Sprintf("%s/comp", loca), 0755)
		if expt != nil {
			return expt
		}
	}

	slog.Log(nil, slog.LevelDebug, "Directories initialized")
	return nil
}

func KillTemp(loca string) error {
	var expt error

	expt = MoveGeneratedFiles(loca)
	if expt != nil {
		return expt
	}

	expt = os.RemoveAll(fmt.Sprintf("%s", loca))
	if expt != nil {
		return expt
	}

	slog.Log(nil, slog.LevelDebug, "Directories removed")
	return nil
}
