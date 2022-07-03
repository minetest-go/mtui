package modmanager

type ModType string
type SourceType string

const (
	ModTypeGame         ModType    = "game"
	ModTypeRegular      ModType    = "mod"
	ModTypeTexturepack  ModType    = "textures"
	ModTypeWorldmods    ModType    = "worldmods"
	SourceTypeContentDB SourceType = "cdb"
	SourceTypeGit       SourceType = "git"
)

type Mod struct {
	Name       string     `json:"name"`
	ModType    ModType    `json:"mod_type"`
	SourceType SourceType `json:"source_type"`
	URL        string     `json:"url"`
	Version    string     `json:"version"`
}

type ModStatus struct {
	CurrentVersion string `json:"current_version"`
	LatestVersion  string `json:"latest_version"`
}
