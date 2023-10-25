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

type Stats struct {
	CPUPercent  float64 `json:"cpu_percent"`
	MemoryUsage uint64  `json:"memory_usage"`
	MemoryMax   uint64  `json:"memory_max"`
	// cumulative values
	NetworkRX uint64 `json:"network_rx"`
	NetworkTX uint64 `json:"network_tx"`
	DiskRead  uint64 `json:"disk_read"`
	DiskWrite uint64 `json:"disk_write"`
}

type ServiceLog struct {
	Out string `json:"out"`
	Err string `json:"err"`
}
