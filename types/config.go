package types

import (
	"os"
	"strings"
)

// env provided configuration flags
type Config struct {
	WorldDir               string
	JWTKey                 string
	APIKey                 string
	CookieDomain           string
	CookieSecure           bool
	CookiePath             string
	Servername             string
	EnabledFeatures        []string
	InstallMtuiMod         bool
	LogRetention           string
	LogStreamURL           string
	LogStreamAuthorization string
	MinetestConfig         string
	TailEngineLogfile      string
	GeoIPAPI               string
	WASMMinetestHost       string
}

func NewConfig(world_dir string) *Config {
	return &Config{
		WorldDir:               world_dir,
		CookieDomain:           os.Getenv("COOKIE_DOMAIN"),
		CookieSecure:           os.Getenv("COOKIE_SECURE") == "true",
		CookiePath:             os.Getenv("COOKIE_PATH"),
		APIKey:                 os.Getenv("API_KEY"),
		JWTKey:                 os.Getenv("JWT_KEY"),
		Servername:             os.Getenv("SERVER_NAME"),
		EnabledFeatures:        strings.Split(os.Getenv("ENABLE_FEATURES"), ","),
		InstallMtuiMod:         os.Getenv("INSTALL_MTUI_MOD") == "true",
		LogRetention:           os.Getenv("LOG_RETENTION"),
		LogStreamURL:           os.Getenv("LOG_STREAM_URL"),
		LogStreamAuthorization: os.Getenv("LOG_STREAM_AUTHORIZATION"),
		MinetestConfig:         os.Getenv("MINETEST_CONFIG"),
		TailEngineLogfile:      os.Getenv("TAIL_ENGINE_LOGFILE"),
		GeoIPAPI:               os.Getenv("GEOIP_API"),
		WASMMinetestHost:       os.Getenv("WASM_MINETEST_HOST"),
	}
}
