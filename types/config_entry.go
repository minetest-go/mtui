package types

type ConfigEntry struct {
	Key   ConfigKey `json:"key" gorm:"primarykey;column:key"`
	Value string    `json:"value" gorm:"column:value"`
}

type ConfigKey string

const (
	ConfigJWTKey             ConfigKey = "jwt_key"
	ConfigApiKey             ConfigKey = "api_key"
	ConfigThemeKey           ConfigKey = "theme"
	ConfigLogStreamTimestamp ConfigKey = "logstream_timestamp"
)

func (m *ConfigEntry) TableName() string {
	return "config"
}
