package types

const (
	FEATURE_MEDIASERVER     = "mediaserver"
	FEATURE_MAIL            = "mail"
	FEATURE_AREAS           = "areas"
	FEATURE_SKINSDB         = "skinsdb"
	FEATURE_LUASHELL        = "luashell"
	FEATURE_SHELL           = "shell"
	FEATURE_MODMANAGEMENT   = "modmanagement"
	FEATURE_XBAN            = "xban"
	FEATURE_MONITORING      = "monitoring"
	FEATURE_MINETEST_CONFIG = "minetest_config"
	FEATURE_OTP             = "otp"
	FEATURE_DOCKER          = "docker"
)

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
