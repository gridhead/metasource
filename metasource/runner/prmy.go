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

func ImportPrimary(wait *sync.WaitGroup, rslt_primary **sxml.UnitPrimary, name *string) {
	defer wait.Done()

	var expt error
	var file *os.File
	var data *bufio.Reader
	var deco *xml.Decoder
	var chip xml.Token

	file, expt = os.Open("/home/fedohide-origin/projects/metasource/rawhide/10beaa5fb8bb9b8710f4608ea9bf84aff2fb68e5efc7e82bf12b421867ad3d8f-primary.xml")
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
				var item sxml.UnitPrimary
				expt = deco.DecodeElement(&item, &se)
				if expt != nil {
					slog.Log(nil, slog.LevelError, fmt.Sprintf("File could not be solved. %s.", expt.Error()))
					return
				}
				if item.Name == *name {
					if rslt_primary != nil {
						*rslt_primary = &item
					}
					return
				}
			}
		}
	}
}
