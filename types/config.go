package types

import (
	"os"
	"strconv"
	"strings"
)

// env provided configuration flags
type Config struct {
	WorldDir                string
	JWTKey                  string
	APIKey                  string
	CookieDomain            string
	CookieSecure            bool
	CookiePath              string
	DefaultTheme            string
	Webdev                  bool
	Servername              string
	EnabledFeatures         []string
	InstallMtuiMod          bool
	AutoReconfigureMods     bool
	LogStreamURL            string
	LogStreamAuthorization  string
	MinetestConfig          string
	GeoIPAPI                string
	DockerMinetestConfig    string
	DockerMinetestPort      int
	WASMMinetestHost        string
	DockerHostname          string
	DockerNetwork           string
	DockerWorlddir          string
	DockerContainerPrefix   string
	DockerAutoInstallEngine bool
}

func NewConfig(world_dir string) *Config {
	port, _ := strconv.ParseInt(os.Getenv("DOCKER_MINETEST_PORT"), 10, 64)

	return &Config{
		WorldDir:                world_dir,
		CookieDomain:            os.Getenv("COOKIE_DOMAIN"),
		CookieSecure:            os.Getenv("COOKIE_SECURE") == "true",
		CookiePath:              os.Getenv("COOKIE_PATH"),
		DefaultTheme:            os.Getenv("DEFAULT_THEME"),
		APIKey:                  os.Getenv("API_KEY"),
		JWTKey:                  os.Getenv("JWT_KEY"),
		Webdev:                  os.Getenv("WEBDEV") == "true",
		Servername:              os.Getenv("SERVER_NAME"),
		EnabledFeatures:         strings.Split(os.Getenv("ENABLE_FEATURES"), ","),
		InstallMtuiMod:          os.Getenv("INSTALL_MTUI_MOD") == "true",
		AutoReconfigureMods:     os.Getenv("AUTORECONFIGURE_MODS") == "true",
		LogStreamURL:            os.Getenv("LOG_STREAM_URL"),
		LogStreamAuthorization:  os.Getenv("LOG_STREAM_AUTHORIZATION"),
		MinetestConfig:          os.Getenv("MINETEST_CONFIG"),
		GeoIPAPI:                os.Getenv("GEOIP_API"),
		DockerMinetestConfig:    os.Getenv("DOCKER_MINETEST_CONFIG"),
		DockerMinetestPort:      int(port),
		WASMMinetestHost:        os.Getenv("WASM_MINETEST_HOST"),
		DockerHostname:          os.Getenv("DOCKER_HOSTNAME"),
		DockerNetwork:           os.Getenv("DOCKER_NETWORK"),
		DockerWorlddir:          os.Getenv("DOCKER_WORLD_DIR"),
		DockerContainerPrefix:   os.Getenv("DOCKER_CONTAINER_PREFIX"),
		DockerAutoInstallEngine: os.Getenv("DOCKER_AUTOINSTALL_ENGINE") == "true",
	}
}
