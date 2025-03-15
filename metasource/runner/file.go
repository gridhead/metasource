package runner

import (
	"bufio"
	"encoding/xml"
	"fmt"
	"io"
	"log/slog"
	"metasource/metasource/models/sxml"
	"os"
	"sync"
)

func ImportFileList(wait *sync.WaitGroup, rslt_filelist **sxml.UnitFileList, name *string) {
	defer wait.Done()

	var expt error
	var file *os.File
	var data *bufio.Reader
	var deco *xml.Decoder
	var chip xml.Token

	file, expt = os.Open("/home/fedohide-origin/projects/metasource/rawhide/4182e96bacb8bb0ccdcc9d446977416a2a18b49a4aed13d6c550be45a1bf061e-filelists.xml")
	if expt != nil {
		slog.Log(nil, slog.LevelError, fmt.Sprintf("File could not be loaded. %s.", expt.Error()))
		return
	}
	defer file.Close()

	data = bufio.NewReader(file)
	deco = xml.NewDecoder(data)

	for {
		chip, expt = deco.Token()
		if expt != nil {
			if expt != io.EOF {
				slog.Log(nil, slog.LevelError, fmt.Sprintf("File could not be parsed. %s.", expt.Error()))
				return
			}
			break
		}

		switch se := chip.(type) {
		case xml.StartElement:
			if se.Name.Local == "package" {
				var item sxml.UnitFileList
				expt = deco.DecodeElement(&item, &se)
				if expt != nil {
					slog.Log(nil, slog.LevelError, fmt.Sprintf("File could not be solved. %s.", expt.Error()))
					return
				}
				if item.Name == *name {
					if rslt_filelist != nil {
						*rslt_filelist = &item
					}
					return
				}
			}
		}
	}
}
