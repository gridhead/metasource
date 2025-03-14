package runner

import (
	"encoding/xml"
	"log/slog"
	"metasource/metasource/models"
)

func ImportPrimary() (*models.Metadata, error) {
	var rslt models.Metadata
	var data []byte
	var expt error
	data, expt = Load("/home/fedohide-origin/projects/metasource/rawhide/10beaa5fb8bb9b8710f4608ea9bf84aff2fb68e5efc7e82bf12b421867ad3d8f-primary.xml")
	if expt != nil {
		slog.Log(nil, slog.LevelError, "File could not be loaded")
		return &rslt, expt
	}
	expt = xml.Unmarshal(data, &rslt)
	if expt != nil {
		slog.Log(nil, slog.LevelError, "File could not be parsed")
		return &rslt, expt
	}
	return &rslt, expt
}
