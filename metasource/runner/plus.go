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

func ImportPlus(wait *sync.WaitGroup, rslt_plus **[]sxml.UnitPrimary, name *string, task *string) {
	defer wait.Done()

	var expt error
	var file *os.File
	var data *bufio.Reader
	var deco *xml.Decoder
	var chip xml.Token
	var list []sxml.UnitPrimary

	file, expt = os.Open("/home/fedohide-origin/projects/metasource/rawhide/10beaa5fb8bb9b8710f4608ea9bf84aff2fb68e5efc7e82bf12b421867ad3d8f-primary.xml")
	if expt != nil {
		slog.Log(nil, slog.LevelError, fmt.Sprintf("File could not be loaded. %s.", expt.Error()))
		return
	}
	defer file.Close()

	list = []sxml.UnitPrimary{}
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

				switch *task {
				case "suggests":
					for _, unit := range item.Format.Suggests.Entries {
						if unit.Name == *name {
							list = append(list, item)
						}
					}
				case "enhances":
					for _, unit := range item.Format.Enhances.Entries {
						if unit.Name == *name {
							list = append(list, item)
						}
					}
				case "requires":
					for _, unit := range item.Format.Requires.Entries {
						if unit.Name == *name {
							list = append(list, item)
						}
					}
				case "provides":
					for _, unit := range item.Format.Provides.Entries {
						if unit.Name == *name {
							list = append(list, item)
						}
					}
				case "obsoletes":
					for _, unit := range item.Format.Obsoletes.Entries {
						if unit.Name == *name {
							list = append(list, item)
						}
					}
				case "conflicts":
					for _, unit := range item.Format.Conflicts.Entries {
						if unit.Name == *name {
							list = append(list, item)
						}
					}
				case "recommends":
					for _, unit := range item.Format.Recommends.Entries {
						if unit.Name == *name {
							list = append(list, item)
						}
					}
				case "supplements":
					for _, unit := range item.Format.Supplements.Entries {
						if unit.Name == *name {
							list = append(list, item)
						}
					}
				}
			}
		}
	}
	if rslt_plus != nil {
		*rslt_plus = &list
	}
}
