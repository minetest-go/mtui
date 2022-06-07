package types

type ModType int
type SourceType int

const (
	ModTypeGame         ModType    = 0
	ModTypeRegular      ModType    = 1
	ModTypeTexturepack  ModType    = 2
	ModTypeWorldmods    ModType    = 3 // all mods in a submodule repo
	SourceTypeContentDB SourceType = 0
	SourceTypeGit       SourceType = 1
)

type Mod struct {
	ID         string     `json:"id"`
	Name       string     `json:"name"`
	ModType    ModType    `json:"mod_type"`
	SourceType SourceType `json:"source_type"`
	URL        string     `json:"url"`
	Version    string     `json:"version"`
	AutoUpdate bool       `json:"auto_update"`
}
