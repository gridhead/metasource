package dict

type Primary struct {
	Version
	Repo        string     `json:"repo"`
	Arch        string     `json:"arch"`
	Summary     string     `json:"summary"`
	Description string     `json:"description"`
	Basename    string     `json:"basename"`
	URL         string     `json:"url"`
	Conflicts   []UnitBase `json:"conflicts"`
	Obsoletes   []UnitBase `json:"obsoletes"`
	Provides    []UnitBase `json:"provides"`
	Requires    []UnitBase `json:"requires"`
	Enhances    []UnitBase `json:"enhances"`
	Recommends  []UnitBase `json:"recommends"`
	Suggests    []UnitBase `json:"suggests"`
	Supplements []UnitBase `json:"supplements"`
	CoPackages  []string   `json:"co-packages"`
}
