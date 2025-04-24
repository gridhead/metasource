package driver

import (
	"crypto/rand"
	"fmt"
	"log/slog"
	"os"
	"strings"
)

func GenerateIdentity(length *int64) string {
	var randBytes []byte

	randBytes = make([]byte, *length/2)
	_, _ = rand.Read(randBytes)

	return fmt.Sprintf("%x", randBytes)
}

func InitPath(vers *string, loca *string) error {
	var expt error

	_, expt = os.Stat(*loca)
	if os.IsNotExist(expt) {
		expt = os.MkdirAll(fmt.Sprintf("%s/sxml", *loca), 0755)
		if expt != nil {
			return expt
		}
		expt = os.MkdirAll(fmt.Sprintf("%s/sxml", *loca), 0755)
		if expt != nil {
			return expt
		}
		expt = os.MkdirAll(fmt.Sprintf("%s/comp", *loca), 0755)
		if expt != nil {
			return expt
		}
	}

	slog.Log(nil, slog.LevelDebug, fmt.Sprintf("[%s] Directories initialized", *vers))
	return nil
}

func KillTemp(vers *string, loca *string) error {
	var expt error

	expt = TransferResult(vers, loca)
	if expt != nil {
		return expt
	}

	expt = os.RemoveAll(*loca)
	if expt != nil {
		return expt
	}

	slog.Log(nil, slog.LevelDebug, fmt.Sprintf("[%s] Directories removed", *vers))
	return nil
}

func TransferResult(vers *string, loca *string) error {
	var expt error
	var files []os.DirEntry
	var oldPath, newPath string

	files, expt = os.ReadDir(*loca)
	if expt != nil {
		return expt
	}

	for _, file := range files {
		if strings.HasPrefix(file.Name(), "metasource-") && strings.HasSuffix(file.Name(), ".sqlite") {
			oldPath = fmt.Sprintf("%s/%s", *loca, file.Name())
			newPath = fmt.Sprintf("%s/../%s", *loca, file.Name())
			expt = os.Rename(oldPath, newPath)
			if expt != nil {
				return expt
			}
		}
	}

	slog.Log(nil, slog.LevelDebug, fmt.Sprintf("[%s] Results transferred", *vers))
	return nil
}
