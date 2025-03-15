package json

type Version struct {
	Epoch   string `json:"epoch"`
	Version string `json:"version"`
	Release string `json:"release"`
}

type Repo struct {
	Repo string `json:"repo"`
}

type UnitBase struct {
	Version
	Name  string   `json:"name"`
	Flags []string `json:"flags"`
}
