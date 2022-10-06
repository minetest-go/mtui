package types

type Feature struct {
	Name    string `json:"name"`
	Enabled bool   `json:"enabled"`
}

func (m *Feature) Columns(action string) []string {
	return []string{"name", "enabled"}
}

func (m *Feature) Table() string {
	return "feature"
}

func (m *Feature) Scan(action string, r func(dest ...any) error) error {
	return r(&m.Name, &m.Enabled)
}

func (m *Feature) Values(action string) []any {
	return []any{m.Name, m.Enabled}
}
