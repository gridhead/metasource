package runner

import (
	"bufio"
	"encoding/xml"
	"fmt"
	"log/slog"
	"metasource/metasource/models"
	"os"
	"sync"
)

func ImportOther(wait *sync.WaitGroup) (*models.Other, error) {
	defer wait.Done()

	var rslt models.Other
	var expt error
	var file *os.File
	var data *bufio.Reader
	var deco *xml.Decoder

	file, expt = os.Open("/home/fedohide-origin/projects/metasource/rawhide/e3ac902af73897fe77cbc4df42d1c87d72ff4d69c9e792224d0ff3857f630e92-other.xml")
	if expt != nil {
		slog.Log(nil, slog.LevelError, fmt.Sprintf("File could not be loaded. %s.", expt.Error()))
		return &rslt, expt
	}

	data = bufio.NewReader(file)
	deco = xml.NewDecoder(data)
	expt = deco.Decode(&rslt)
	if expt != nil {
		slog.Log(nil, slog.LevelError, fmt.Sprintf("File could not be parsed. %s.", expt.Error()))
		return &rslt, expt
	}

	expt = file.Close()
	if expt != nil {
		return &rslt, expt
	}

	return &rslt, expt
}
