package types

type OauthApp struct {
	ID      string `json:"id"`
	Enabled bool   `json:"enabled"`
	Created int64  `json:"created"`
	Name    string `json:"name"`
	Domain  string `json:"domain"`
	Secret  string `json:"secret"`
}

func (m *OauthApp) Columns(action string) []string {
	return []string{
		"id",
		"enabled",
		"created",
		"name",
		"domain",
		"secret",
	}
}

func (m *OauthApp) Table() string {
	return "oauth_app"
}

func (m *OauthApp) Scan(action string, r func(dest ...any) error) error {
	return r(
		&m.ID,
		&m.Enabled,
		&m.Created,
		&m.Name,
		&m.Domain,
		&m.Secret,
	)
}

func (m *OauthApp) Values(action string) []any {
	return []any{
		m.ID,
		m.Enabled,
		m.Created,
		m.Name,
		m.Domain,
		m.Secret,
	}
}

func (m *OauthApp) GetID() string {
	return m.Name
}

func (m *OauthApp) GetSecret() string {
	return m.Secret
}

func (m *OauthApp) GetDomain() string {
	return m.Domain
}

func (m *OauthApp) IsPublic() bool {
	return true
}

func (m *OauthApp) GetUserID() string {
	return ""
}
