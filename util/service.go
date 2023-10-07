package util

type DockerServiceConfig struct {
}

type DockerService struct {
	cfg *DockerServiceConfig
}

func NewDockerService(cfg *DockerServiceConfig) *DockerService {
	return &DockerService{cfg: cfg}
}
