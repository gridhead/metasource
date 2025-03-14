package runner

import (
	"io"
	"log/slog"
	"os"
)

func Load(path string) ([]byte, error) {
	var expt error
	var file *os.File
	var data []byte
	file, expt = os.Open(path)
	if expt != nil {
		slog.Log(nil, slog.LevelError, "File could not be loaded")
		return data, expt
	}
	defer file.Close()
	data, expt = io.ReadAll(file)
	if expt != nil {
		slog.Log(nil, slog.LevelError, "File could not be read")
		return data, expt
	}
	return data, expt
}
