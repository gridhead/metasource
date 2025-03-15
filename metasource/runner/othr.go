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

func ImportOther(wait *sync.WaitGroup, rslt_other **sxml.UnitOther, name *string) {
	defer wait.Done()

	var expt error
	var file *os.File
	var data *bufio.Reader
	var deco *xml.Decoder
	var chip xml.Token

	file, expt = os.Open("/home/fedohide-origin/projects/metasource/rawhide/e3ac902af73897fe77cbc4df42d1c87d72ff4d69c9e792224d0ff3857f630e92-other.xml")
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
				var item sxml.UnitOther
				expt = deco.DecodeElement(&item, &se)
				if expt != nil {
					slog.Log(nil, slog.LevelError, fmt.Sprintf("File could not be solved. %s.", expt.Error()))
					return
				}
				if item.Name == *name {
					if rslt_other != nil {
						*rslt_other = &item
					}
					return
				}
			}
		}
	}
}
