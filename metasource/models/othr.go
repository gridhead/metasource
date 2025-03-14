package models

import "encoding/xml"

type Other struct {
	XMLName xml.Name    `xml:"otherdata"`
	XMLNS   string      `xml:"xmlns,attr"`
	Package uint64      `xml:"package,attr"`
	List    []UnitOther `xml:"package"`
}

type UnitOther struct {
	XMLName   xml.Name    `xml:"package"`
	PkgID     string      `xml:"pkgid,attr"`
	Name      string      `xml:"name,attr"`
	Arch      string      `xml:"arch,attr"`
	Version   Version     `xml:"version,attr"`
	Changelog []Changelog `xml:"changelog"`
}

type Changelog struct {
	XMLName xml.Name `xml:"changelog"`
	Author  string   `xml:"author,attr"`
	Date    uint64   `xml:"date,attr"`
	Data    string   `xml:",chardata"`
}
