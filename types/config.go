package types

import (
	"os"
	"strconv"
)

// env provided configuration flags
type Config struct {
	JWTKey                  string
	APIKey                  string
	CookieDomain            string
	CookieSecure            bool
	CookiePath              string
	Webdev                  bool
	Servername              string
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
		APIKey:                  os.Getenv("API_KEY"),
		Webdev:                  os.Getenv("WEBDEV") == "true",
		Servername:              os.Getenv("SERVER_NAME"),
		MinetestConfig:          os.Getenv("MINETEST_CONFIG"),
		DockerMinetestConfig:    os.Getenv("DOCKER_MINETEST_CONFIG"),
		DockerMinetestPort:      int(port),
		DockerHostname:          os.Getenv("DOCKER_HOSTNAME"),
		DockerNetwork:           os.Getenv("DOCKER_NETWORK"),
		DockerWorlddir:          os.Getenv("DOCKER_WORLD_DIR"),
		DockerMinetestContainer: os.Getenv("DOCKER_MINETEST_CONTAINER"),
	}
}
