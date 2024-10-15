package types

type OauthApp struct {
	ID      string `json:"id" gorm:"primarykey;column:id"`
	Enabled bool   `json:"enabled" gorm:"column:enabled"`
	Created int64  `json:"created" gorm:"column:created"`
	Name    string `json:"name" gorm:"column:name"`
	Domain  string `json:"domain" gorm:"column:domain"`
	Secret  string `json:"secret" gorm:"column:secret"`
}

func (m *OauthApp) TableName() string {
	return "oauth_app"
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
