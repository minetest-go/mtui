package dockerservice

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strings"
	"time"

	dockertypes "github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/network"
	"github.com/docker/docker/pkg/stdcopy"
	"github.com/sirupsen/logrus"

	"github.com/docker/docker/api/types/filters"
)

func New(cfg *Config) *DockerService {
	return &DockerService{cfg: cfg}
}

func (s *DockerService) getContainer() (*dockertypes.Container, error) {
	cli, err := getDockerCli()
	if err != nil {
		return nil, err
	}
	defer cli.Close()

	ctx := context.Background()
	f := filters.NewArgs()
	f.Add("name", s.cfg.ContainerName)
	containers, err := cli.ContainerList(ctx, container.ListOptions{
		All:     true,
		Filters: f,
	})
	if err != nil {
		return nil, fmt.Errorf("container-list error: %v", err)
	}

	if len(containers) == 1 {
		// single container found
		return &containers[0], nil
	} else if len(containers) > 1 {
		return nil, fmt.Errorf("multiple containers found with name '%s'", s.cfg.ContainerName)
	}

	// no container found
	return nil, nil
}

func (s *DockerService) Status() (*Status, error) {
	status := &Status{}

	container, err := s.getContainer()
	if err != nil {
		return nil, fmt.Errorf("fetch container error: %v", err)
	}
	if container != nil {
		status.ID = container.ID
		status.Created = true
		status.Running = container.State == "running"
		parts := strings.Split(container.Image, ":")
		if len(parts) == 2 {
			status.Image = parts[0]
			status.Version = parts[1]
		} else if len(parts) == 1 {
			status.Image = parts[0]
			status.Version = "latest"
		} else {
			return nil, fmt.Errorf("could not parse image: '%s'", container.Image)
		}

	}

	return status, nil
}

func (s *DockerService) Stats() (*Stats, error) {
	container, err := s.getContainer()
	if err != nil {
		return nil, fmt.Errorf("fetch container error: %v", err)
	}

	st := &Stats{}

	if container == nil {
		return st, nil
	}

	cli, err := getDockerCli()
	if err != nil {
		return nil, err
	}
	defer cli.Close()

	statsres, err := cli.ContainerStats(context.Background(), container.ID, false)
	if err != nil {
		return nil, fmt.Errorf("could not get stats: %v", err)
	}
	defer statsres.Body.Close()

	stats := &dockertypes.StatsJSON{}
	err = json.NewDecoder(statsres.Body).Decode(stats)
	if err != nil {
		return nil, fmt.Errorf("could not decode stats: %v", err)
	}

	// cpu
	// https://github.com/moby/moby/blob/eb131c5383db8cac633919f82abad86c99bffbe5/cli/command/container/stats_helpers.go#L175-L188
	cpuDelta := float64(stats.CPUStats.CPUUsage.TotalUsage) - float64(stats.PreCPUStats.CPUUsage.TotalUsage)
	systemDelta := float64(stats.CPUStats.SystemUsage) - float64(stats.PreCPUStats.SystemUsage)
	if systemDelta > 0.0 && cpuDelta > 0.0 {
		st.CPUPercent = (cpuDelta / systemDelta) * 100.0
	}

	// memory
	st.MemoryUsage = stats.MemoryStats.Usage
	st.MemoryMax = stats.MemoryStats.Limit

	// net
	for _, e := range stats.Networks {
		st.NetworkRX += e.RxBytes
		st.NetworkTX += e.TxBytes
	}

	// disk io
	for _, e := range stats.BlkioStats.IoServiceBytesRecursive {
		switch e.Op {
		case "read":
			st.DiskRead += e.Value
		case "write":
			st.DiskWrite += e.Value
		}
	}

	return st, nil
}

func (s *DockerService) Create(image string) error {
	// https://docs.docker.com/engine/api/sdk/examples/

	// check if container already exists
	c, err := s.getContainer()
	if err != nil {
		return fmt.Errorf("fetch container error: %v", err)
	}
	if c != nil {
		return fmt.Errorf("container already created")
	}

	ctx := context.Background()
	cli, err := getDockerCli()
	if err != nil {
		return fmt.Errorf("docker client error: %v", err)
	}
	defer cli.Close()

	reader, err := cli.ImagePull(ctx, image, dockertypes.ImagePullOptions{})
	if err != nil {
		return fmt.Errorf("image pull error: %v", err)
	}
	defer reader.Close()
	_, err = io.Copy(os.Stdout, reader)
	if err != nil {
		return fmt.Errorf("io-copy error: %v", err)
	}

	logrus.WithFields(logrus.Fields{
		"image":          image,
		"networks":       s.cfg.Networks,
		"userid":         os.Getuid(),
		"container_name": s.cfg.ContainerName,
	}).Info("Creating minetest engine service")

	ccfg := s.cfg.DefaultConfig
	ccfg.Image = image

	resp, err := cli.ContainerCreate(ctx, s.cfg.DefaultConfig, s.cfg.DefaultHostConfig, nil, nil, s.cfg.ContainerName)
	if err != nil {
		return fmt.Errorf("could not create container: %v", err)
	}

	// normal network without alias
	for _, name := range s.cfg.Networks {
		err = cli.NetworkConnect(ctx, name, resp.ID, &network.EndpointSettings{
			NetworkID: name,
		})
		if err != nil {
			return fmt.Errorf("could not connect container %s to network %s: %v", resp.ID, name, err)
		}
	}

	// internal networks with alias
	for _, name := range s.cfg.InternalNetworks {
		aliases := []string{}
		if s.cfg.InternalName != "" {
			aliases = append(aliases, s.cfg.InternalName)
		}
		err = cli.NetworkConnect(ctx, name, resp.ID, &network.EndpointSettings{
			NetworkID: name,
			Aliases:   aliases,
		})
		if err != nil {
			return fmt.Errorf("could not connect container %s to internal network %s: %v", resp.ID, name, err)
		}
	}

	return nil
}

func (s *DockerService) Remove() error {
	c, err := s.getContainer()
	if err != nil {
		return fmt.Errorf("fetch container error: %v", err)
	}
	if c == nil {
		return fmt.Errorf("no container found")
	}

	cli, err := getDockerCli()
	if err != nil {
		return fmt.Errorf("docker client error: %v", err)
	}
	defer cli.Close()

	ctx := context.Background()
	return cli.ContainerRemove(ctx, c.ID, container.RemoveOptions{
		Force: true,
	})
}

func (s *DockerService) Start() error {
	c, err := s.getContainer()
	if err != nil {
		return fmt.Errorf("fetch container error: %v", err)
	}
	if c == nil {
		return fmt.Errorf("no container found")
	}

	cli, err := getDockerCli()
	if err != nil {
		return fmt.Errorf("docker client error: %v", err)
	}
	defer cli.Close()

	ctx := context.Background()
	return cli.ContainerStart(ctx, c.ID, container.StartOptions{})
}

func (s *DockerService) Stop() error {
	c, err := s.getContainer()
	if err != nil {
		return fmt.Errorf("fetch container error: %v", err)
	}
	if c == nil {
		return fmt.Errorf("no container found")
	}

	cli, err := getDockerCli()
	if err != nil {
		return fmt.Errorf("docker client error: %v", err)
	}
	defer cli.Close()

	ctx := context.Background()
	return cli.ContainerStop(ctx, c.ID, container.StopOptions{})
}

func (s *DockerService) GetLogs(since, until time.Time) (*ServiceLog, error) {
	c, err := s.getContainer()
	if err != nil {
		return nil, fmt.Errorf("fetch container error: %v", err)
	}
	if c == nil {
		return nil, fmt.Errorf("no container found")
	}

	cli, err := getDockerCli()
	if err != nil {
		return nil, fmt.Errorf("docker client error: %v", err)
	}
	defer cli.Close()

	ctx := context.Background()

	logs, err := cli.ContainerLogs(ctx, c.ID, container.LogsOptions{
		ShowStdout: true,
		ShowStderr: true,
		Since:      since.Format(time.RFC3339),
		Until:      until.Format(time.RFC3339),
	})
	if err != nil {
		return nil, fmt.Errorf("docker stdout log error: %v", err)
	}
	defer logs.Close()

	out_buf := bytes.NewBuffer([]byte{})
	err_buf := bytes.NewBuffer([]byte{})

	_, err = stdcopy.StdCopy(out_buf, err_buf, logs)
	if err != nil {
		return nil, fmt.Errorf("docker stdcopy error: %v", err)
	}

	slog := &ServiceLog{
		Out: out_buf.String(),
		Err: err_buf.String(),
	}
	return slog, nil
}
