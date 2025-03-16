package sxml

import "encoding/xml"

type Version struct {
	XMLName xml.Name `xml:"version"`
	Epoch   string   `xml:"epoch,attr"`
	Ver     string   `xml:"ver,attr"`
	Rel     string   `xml:"rel,attr"`
}

type UnitBase struct {
	XMLName xml.Name `xml:"package"`
	PkgID   string   `xml:"pkgid,attr"`
	Name    string   `xml:"name,attr"`
	Arch    string   `xml:"arch,attr"`
	Version Version  `xml:"version"`
}
