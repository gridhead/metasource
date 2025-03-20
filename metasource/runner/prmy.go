package runner

import (
	"encoding/xml"
	"fmt"
	"github.com/antchfx/xmlquery"
	"log"
	"metasource/metasource/models/sxml"
	"os"
	"sync"
)

func ImportPrimary(wait *sync.WaitGroup, rslt_primary **sxml.UnitPrimary, name *string) {
	defer wait.Done()

	var file *os.File
	var data []byte

	filePath := "/home/fedohide-origin/projects/metasource/rawhide/10beaa5fb8bb9b8710f4608ea9bf84aff2fb68e5efc7e82bf12b421867ad3d8f-primary.xml"

	// Open the XML file
	file, err := os.Open(filePath)
	if err != nil {
		log.Fatalf("File could not be loaded: %s", err)
	}
	defer file.Close()

	// STREAM
	//
	//doc, err := xmlquery.CreateStreamParser(file, "/metadata/package", fmt.Sprintf("/metadata/package[name='%s']", *name))
	//
	//
	//if err != nil {
	//	log.Fatalf("File could not be parsed: %s", err)
	//}
	//
	//for {
	//	n, err := doc.Read()
	//	if err == io.EOF {
	//		break
	//	}
	//	if err != nil {
	//		log.Fatalf("File could not be parsed: %s", err)
	//	}
	//
	//	fmt.Println(n.OutputXML(true))
	//
	//	data = []byte(n.OutputXML(true))
	//	_ = xml.Unmarshal(data, rslt_primary)
	//
	//	break
	//}

	doc, err := xmlquery.Parse(file)

	// Find all <package> nodes
	//nameNode := xmlquery.FindOne(doc, fmt.Sprintf("/metadata/package[name='%s']", *name))
	nameNode := xmlquery.FindOne(doc, fmt.Sprintf("//package[name='%s']", *name))

	//DECODING TECHNIQUE #1
	var item sxml.UnitPrimary
	data = []byte(nameNode.OutputXML(true))
	_ = xml.Unmarshal(data, &item)
	*rslt_primary = &item

	//for _, pkg := range packages {
	//	nameNode := xmlquery.FindOne(pkg, "name")
	//	if nameNode != nil && nameNode.InnerText() == targetName {
	//		fmt.Printf("Found package: %s\n", nameNode.InnerText())
	//
	//		// Extract other details if needed
	//		archNode := xmlquery.FindOne(pkg, "arch")
	//		if archNode != nil {
	//			fmt.Printf("Architecture: %s\n", archNode.InnerText())
	//		}
	//
	//		data = []byte(nameNode.OutputXML(true))
	//		_ = xml.Unmarshal(data, rslt_primary)
	//
	//		// If you only need the first match, exit the loop
	//		return
	//	}
	//}

}
