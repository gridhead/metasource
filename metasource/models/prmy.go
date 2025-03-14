package models

import "encoding/xml"

type Primary struct {
	XMLName  xml.Name      `xml:"metadata"`
	XMLNS    string        `xml:"xmlns,attr"`
	XMLNSRPM string        `xml:"xmlns:rpm,attr"`
	Packages uint64        `xml:"packages,attr"`
	List     []UnitPrimary `xml:"package"`
}

type UnitPrimary struct {
	XMLName     xml.Name `xml:"package"`
	Type        string   `xml:"type,attr"`
	Name        string   `xml:"name"`
	Arch        string   `xml:"arch"`
	Version     Version  `xml:"version"`
	Checksum    Checksum `xml:"checksum"`
	Summary     string   `xml:"summary"`
	Description string   `xml:"description"`
	Packager    string   `xml:"packager"`
	URL         string   `xml:"url"`
	Time        Time     `xml:"time"`
	Size        Size     `xml:"size"`
	Location    Location `xml:"location"`
	Format      Format   `xml:"format"`
}

type Checksum struct {
	XMLName xml.Name `xml:"checksum"`
	Type    string   `xml:"type,attr"`
	PkgID   string   `xml:"pkgid,attr"`
	Data    string   `xml:",chardata"`
}

type Time struct {
	XMLName xml.Name `xml:"time"`
	File    uint64   `xml:"file,attr"`
	Build   uint64   `xml:"build,attr"`
}

type Size struct {
	XMLName   xml.Name `xml:"size"`
	Package   uint64   `xml:"package,attr"`
	Installed uint64   `xml:"installed,attr"`
	Archive   uint64   `xml:"archive,attr"`
}

type Location struct {
	XMLName xml.Name `xml:"location"`
	Href    string   `xml:"href,attr"`
}

type Format struct {
	XMLName     xml.Name    `xml:"format"`
	License     string      `xml:"http://linux.duke.edu/metadata/rpm license"`
	Vendor      string      `xml:"http://linux.duke.edu/metadata/rpm vendor"`
	Group       string      `xml:"http://linux.duke.edu/metadata/rpm group"`
	BuildHost   string      `xml:"http://linux.duke.edu/metadata/rpm buildhost"`
	SourceRPM   string      `xml:"http://linux.duke.edu/metadata/rpm sourcerpm"`
	HeaderRange HeaderRange `xml:"http://linux.duke.edu/metadata/rpm header-range"`
	Provides    Provides    `xml:"http://linux.duke.edu/metadata/rpm provides"`
	Requires    Requires    `xml:"http://linux.duke.edu/metadata/rpm requires"`
	Files       []string    `xml:"files"`
}

type HeaderRange struct {
	XMLName xml.Name `xml:"http://linux.duke.edu/metadata/rpm header-range"`
	Start   uint64   `xml:"start,attr"`
	End     uint64   `xml:"end,attr"`
}

type Provides struct {
	XMLName xml.Name `xml:"http://linux.duke.edu/metadata/rpm provides"`
	EntriesBase
}

type Requires struct {
	XMLName xml.Name `xml:"http://linux.duke.edu/metadata/rpm requires"`
	EntriesBase
}

type Entry struct {
	XMLName xml.Name `xml:"http://linux.duke.edu/metadata/rpm entry"`
	Name    string   `xml:"name,attr"`
	Flags   string   `xml:"flags,attr,omitempty"`
	Epoch   string   `xml:"epoch,attr,omitempty"`
	Ver     string   `xml:"ver,attr,omitempty"`
	Rel     string   `xml:"rel,attr,omitempty"`
}
