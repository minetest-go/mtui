package types

type ModType string

const (
	ModTypeMod         ModType = "mod"
	ModTypeGame        ModType = "game"
	ModTypeTexturepack ModType = "txp"
)

type SourceType string

const (
	SourceTypeCDB    SourceType = "cdb"
	SourceTypeGIT    SourceType = "git"
	SourceTypeManual SourceType = "manual"
)

type Mod struct {
	ID            string     `json:"id"`
	Name          string     `json:"name"`
	Author        string     `json:"author"`
	ModType       ModType    `json:"mod_type"`
	SourceType    SourceType `json:"source_type"`
	URL           string     `json:"url"`
	Branch        string     `json:"branch"`
	Version       string     `json:"version"`
	LatestVersion string     `json:"latest_version"`
	AutoUpdate    bool       `json:"auto_update"`
}

func (m *Mod) Columns(action string) []string {
	return []string{"id", "name", "author", "mod_type", "source_type", "url", "branch", "version", "latest_version", "auto_update"}
}

func (m *Mod) Table() string {
	return "mod"
}

func (m *Mod) Scan(action string, r func(dest ...any) error) error {
	return r(&m.ID, &m.Name, &m.Author, &m.ModType, &m.SourceType, &m.URL, &m.Branch, &m.Version, &m.LatestVersion, &m.AutoUpdate)
}

func (m *Mod) Values(action string) []any {
	return []any{m.ID, m.Name, m.Author, m.ModType, m.SourceType, m.URL, m.Branch, m.Version, m.LatestVersion, m.AutoUpdate}
}
