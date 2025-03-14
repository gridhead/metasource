package models

import "encoding/xml"

type Other struct {
	XMLName  xml.Name    `xml:"otherdata"`
	XMLNS    string      `xml:"xmlns,attr"`
	Packages uint64      `xml:"packages,attr"`
	List     []UnitOther `xml:"package"`
}

type UnitOther struct {
	UnitBase
	Changelog []Changelog `xml:"changelog"`
}

type Changelog struct {
	XMLName xml.Name `xml:"changelog"`
	Author  string   `xml:"author,attr"`
	Date    uint64   `xml:"date,attr"`
	Data    string   `xml:",chardata"`
}

type EntriesBase struct {
	Entries []Entry `xml:"http://linux.duke.edu/metadata/rpm entry"`
}
