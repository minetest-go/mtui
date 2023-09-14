package types

import (
	"os"
	"strconv"
	"strings"
)

// env provided configuration flags
type Config struct {
	JWTKey                  string
	APIKey                  string
	CookieDomain            string
	CookieSecure            bool
	CookiePath              string
	DefaultTheme            string
	Webdev                  bool
	Servername              string
	EnabledFeatures         []string
	AdminUsername           string
	AdminPassword           string
	MinetestConfig          string
	DockerMinetestConfig    string
	DockerMinetestPort      int
	DockerHostname          string
	DockerNetwork           string
	DockerWorlddir          string
	DockerMinetestContainer string
}

func NewConfig() *Config {
	port, _ := strconv.ParseInt(os.Getenv("DOCKER_MINETEST_PORT"), 10, 64)

	return &Config{
		CookieDomain:            os.Getenv("COOKIE_DOMAIN"),
		CookieSecure:            os.Getenv("COOKIE_SECURE") == "true",
		CookiePath:              os.Getenv("COOKIE_PATH"),
		DefaultTheme:            os.Getenv("DEFAULT_THEME"),
		APIKey:                  os.Getenv("API_KEY"),
		Webdev:                  os.Getenv("WEBDEV") == "true",
		Servername:              os.Getenv("SERVER_NAME"),
		EnabledFeatures:         strings.Split(os.Getenv("ENABLE_FEATURES"), ","),
		AdminUsername:           os.Getenv("ADMIN_USERNAME"),
		AdminPassword:           os.Getenv("ADMIN_PASSWORD"),
		MinetestConfig:          os.Getenv("MINETEST_CONFIG"),
		DockerMinetestConfig:    os.Getenv("DOCKER_MINETEST_CONFIG"),
		DockerMinetestPort:      int(port),
		DockerHostname:          os.Getenv("DOCKER_HOSTNAME"),
		DockerNetwork:           os.Getenv("DOCKER_NETWORK"),
		DockerWorlddir:          os.Getenv("DOCKER_WORLD_DIR"),
		DockerMinetestContainer: os.Getenv("DOCKER_MINETEST_CONTAINER"),
	}
}
