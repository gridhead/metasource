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

func ImportFileList(wait *sync.WaitGroup, rslt_filelist **models.FileList) {
	defer wait.Done()

	var rslt models.FileList
	var expt error
	var file *os.File
	var data *bufio.Reader
	var deco *xml.Decoder

	file, expt = os.Open("/home/fedohide-origin/projects/metasource/rawhide/4182e96bacb8bb0ccdcc9d446977416a2a18b49a4aed13d6c550be45a1bf061e-filelists.xml")
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

	*rslt_filelist = &rslt
}
