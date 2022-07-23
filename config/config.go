package config

type Config struct {
	JWTKey       string `json:"jwt_key"`
	APIKey       string `json:"api_key"`
	CookieDomain string `json:"cookie_domain"`
	CookieSecure bool   `json:"cookie_secure"`
	CookiePath   string `json:"cookie_path"`
}
