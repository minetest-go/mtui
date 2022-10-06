package types

type Config struct {
	JWTKey       string
	APIKey       string
	CookieDomain string
	CookieSecure bool
	CookiePath   string
}
