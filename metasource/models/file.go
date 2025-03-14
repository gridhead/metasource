package models

import "encoding/xml"

type FileList struct {
	XMLName  xml.Name       `xml:"filelists"`
	XMLNS    string         `xml:"xmlns,attr"`
	Packages uint64         `xml:"packages,attr"`
	List     []UnitFileList `xml:"package"`
}

type UnitFileList struct {
	UnitBase
	List []FileName `xml:"file"`
}

type FileName struct {
	XMLName xml.Name `xml:"file"`
	Data    string   `xml:",chardata"`
}
