package driver

import (
	"fmt"
	"log/slog"
	"os"
	"strings"
)

func MoveGeneratedFiles(loca string) error {
    var expt error
	var files []os.DirEntry
	var oldPath, newPath string
   
	files, expt = os.ReadDir(loca)
	if expt != nil {
		return expt
	}

	for _, file := range files {
		if strings.HasPrefix(file.Name(), "metasource-") && strings.HasSuffix(file.Name(), ".sqlite") {
            oldPath = fmt.Sprintf("%s/%s", loca, file.Name())
            newPath = fmt.Sprintf("%s/%s", loca+"/..", file.Name())
            expt = os.Rename(oldPath, newPath)
            if expt != nil {
                return expt
            }
        }
    }

    slog.Log(nil, slog.LevelDebug, fmt.Sprintf("Files moved from temporary to parent directory"))
    return nil
}