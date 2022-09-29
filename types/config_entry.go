package types

type ConfigEntry struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

const (
	ConfigJWTKey = "jwt_key"
	ConfigApiKey = "api_key"
)

func (m *ConfigEntry) Columns(action string) []string {
	return []string{"key", "value"}
}

func (m *ConfigEntry) Table() string {
	return "config"
}

func (m *ConfigEntry) Scan(action string, r func(dest ...any) error) error {
	return r(&m.Key, &m.Value)
}

func (m *ConfigEntry) Values(action string) []any {
	return []any{m.Key, m.Value}
}
