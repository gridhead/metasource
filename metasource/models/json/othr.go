package json

type Other struct {
	Repo
	Changelogs []Changelog `json:"changelogs"`
}

type Changelog struct {
	Author    string `json:"author"`
	Changelog string `json:"changelog"`
	Date      uint64 `json:"date"`
}
