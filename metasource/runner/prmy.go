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

func ImportPrimary(wait *sync.WaitGroup, rslt_primary **models.Primary) {
	defer wait.Done()

	var rslt models.Primary
	var expt error
	var file *os.File
	var data *bufio.Reader
	var deco *xml.Decoder

	file, expt = os.Open("/home/fedohide-origin/projects/metasource/rawhide/10beaa5fb8bb9b8710f4608ea9bf84aff2fb68e5efc7e82bf12b421867ad3d8f-primary.xml")
	if expt != nil {
		slog.Log(nil, slog.LevelError, fmt.Sprintf("File could not be loaded. %s.", expt.Error()))
	}

	data = bufio.NewReader(file)
	deco = xml.NewDecoder(data)
	expt = deco.Decode(&rslt)
	if expt != nil {
		slog.Log(nil, slog.LevelError, fmt.Sprintf("File could not be parsed. %s.", expt.Error()))
	}

	expt = file.Close()
	if expt != nil {
		slog.Log(nil, slog.LevelError, fmt.Sprintf("File could not be closed. %s.", expt.Error()))
	}

	*rslt_primary = &rslt
}
