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
	Keep bool
}

type Checksum struct {
	Data string
	Type string
}
