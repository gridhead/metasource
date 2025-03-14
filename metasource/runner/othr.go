package runner

import (
	"encoding/xml"
	"log/slog"
	"metasource/metasource/models"
)

func ImportOther() (*models.Other, error) {
	var rslt models.Other
	var data []byte
	var expt error
	data, expt = Load("/home/fedohide-origin/projects/rawhide/e3ac902af73897fe77cbc4df42d1c87d72ff4d69c9e792224d0ff3857f630e92-other.xml")
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
