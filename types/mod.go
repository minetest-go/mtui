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
	ID            string     `json:"id" gorm:"primarykey;column:id"`
	Name          string     `json:"name" gorm:"column:name"`
	Author        string     `json:"author" gorm:"column:author"`
	ModType       ModType    `json:"mod_type" gorm:"column:mod_type"`
	SourceType    SourceType `json:"source_type" gorm:"column:source_type"`
	URL           string     `json:"url" gorm:"column:url"`
	Branch        string     `json:"branch" gorm:"column:branch"`
	Version       string     `json:"version" gorm:"column:version"`
	LatestVersion string     `json:"latest_version" gorm:"column:latest_version"`
	AutoUpdate    bool       `json:"auto_update" gorm:"column:auto_update"`
}

func (m *Mod) TableName() string {
	return "mod"
}
