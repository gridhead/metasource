package home

type LinkUnit struct {
	Name string
	Link string
}

type FileUnit struct {
	Name string
	Path string
	Type string
	Hash Checksum
}

type Checksum struct {
	Data string
	Type string
}
