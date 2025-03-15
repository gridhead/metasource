package json

type Files struct {
	Repo
	Files []File `json:"files"`
}

type File struct {
	DirName   string `json:"dirname"`
	FileNames string `json:"filenames"`
	FileTypes string `json:"filetypes"`
}
