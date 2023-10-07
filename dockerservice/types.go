package dockerservice

import "github.com/docker/docker/api/types/container"

type Config struct {
	ContainerName     string
	Networks          []string
	DefaultConfig     *container.Config
	DefaultHostConfig *container.HostConfig
}

type DockerService struct {
	cfg *Config
}

type Status struct {
	ID      string `json:"id"`
	Created bool   `json:"created"`
	Running bool   `json:"running"`
	Image   string `json:"image"`
	Version string `json:"version"`
}

type ServiceLog struct {
	Out string `json:"out"`
	Err string `json:"err"`
}
