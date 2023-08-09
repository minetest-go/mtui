package web

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"mtui/types"
	"net/http"
	"os"
	"strconv"
	"time"

	dockertypes "github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/api/types/mount"
	"github.com/docker/docker/api/types/network"
	"github.com/docker/docker/client"
	"github.com/docker/docker/pkg/stdcopy"
	"github.com/docker/go-connections/nat"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
)

// map versions with full image urls (in case the registry gets switched in the future)
var VersionImageMapping = map[string]string{
	"5.6.0": "registry.gitlab.com/minetest/minetest/server:5.6.0",
	"5.7.0": "registry.gitlab.com/minetest/minetest/server:5.7.0",
}

func (a *Api) GetEngineVersions(w http.ResponseWriter, r *http.Request, claims *types.Claims) {
	SendJson(w, VersionImageMapping)
}

type ServiceStatus struct {
	ID      string `json:"id"`
	Created bool   `json:"created"`
	Running bool   `json:"running"`
	Version string `json:"version"`
}

func getDockerCli() (*client.Client, error) {
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		return nil, fmt.Errorf("new client error: %v", err)
	}
	return cli, nil
}

func getMinetestContainer() (*dockertypes.Container, error) {
	cli, err := getDockerCli()
	if err != nil {
		return nil, err
	}

	ctx := context.Background()
	container_name := os.Getenv("DOCKER_MINETEST_CONTAINER")
	f := filters.NewArgs()
	f.Add("name", container_name)
	containers, err := cli.ContainerList(ctx, dockertypes.ContainerListOptions{
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
		return nil, fmt.Errorf("multiple containers found with name '%s'", container_name)
	}

	// no container found
	return nil, nil
}

func (a *Api) GetEngineStatus(w http.ResponseWriter, r *http.Request, claims *types.Claims) {
	status := &ServiceStatus{}

	container, err := getMinetestContainer()
	if err != nil {
		SendError(w, 500, fmt.Sprintf("fetch container error: %v", err))
		return
	}
	if container != nil {
		status.ID = container.ID
		status.Created = true
		status.Running = container.State == "running"
		for v, img := range VersionImageMapping {
			if img == container.Image {
				status.Version = v
				break
			}
		}
	}

	SendJson(w, status)
}

type CreateEngineRequest struct {
	Version string `json:"version"`
}

func (a *Api) CreateEngine(w http.ResponseWriter, r *http.Request, claims *types.Claims) {
	// https://docs.docker.com/engine/api/sdk/examples/

	// check if container already exists
	c, err := getMinetestContainer()
	if err != nil {
		SendError(w, 500, fmt.Sprintf("fetch container error: %v", err))
		return
	}
	if c != nil {
		SendError(w, 500, "container already created")
		return
	}

	cer := &CreateEngineRequest{}
	err = json.NewDecoder(r.Body).Decode(cer)
	if err != nil {
		SendError(w, 500, fmt.Sprintf("json error: %v", err))
		return
	}

	image := VersionImageMapping[cer.Version]
	if image == "" {
		SendError(w, 404, fmt.Sprintf("unknown version: %s", cer.Version))
		return
	}

	ctx := context.Background()
	cli, err := getDockerCli()
	if err != nil {
		SendError(w, 500, fmt.Sprintf("docker client error: %v", err))
		return
	}
	defer cli.Close()

	reader, err := cli.ImagePull(ctx, image, dockertypes.ImagePullOptions{})
	if err != nil {
		SendError(w, 500, fmt.Sprintf("image pull error: %v", err))
		return
	}
	defer reader.Close()
	_, err = io.Copy(os.Stdout, reader)
	if err != nil {
		SendError(w, 500, fmt.Sprintf("io-copy error: %v", err))
		return
	}

	world_dir := os.Getenv("DOCKER_WORLD_DIR")
	world_dir_container := "/world"
	minetest_config := os.Getenv("DOCKER_MINETEST_CONFIG")
	minetest_config_container := "/minetest.conf"
	if minetest_config == "" {
		SendError(w, 500, "minetest config not found")
		return
	}

	port_str := os.Getenv("DOCKER_MINETEST_PORT")
	port, _ := strconv.ParseInt(port_str, 10, 64)
	if port == 0 {
		SendError(w, 500, fmt.Sprintf("invalid port: '%s'", port_str))
		return
	}

	network_name := os.Getenv("DOCKER_NETWORK")
	container_name := os.Getenv("DOCKER_MINETEST_CONTAINER")

	logrus.WithFields(logrus.Fields{
		"world_dir":       world_dir,
		"minetest_config": minetest_config,
		"version":         cer.Version,
		"image":           image,
		"port":            port,
		"uid":             os.Getuid(),
		"network":         network_name,
		"container_name":  container_name,
	}).Info("Creating minetest engine service")

	// prefix world and config with /data inside container to prevent filename collisions
	resp, err := cli.ContainerCreate(ctx, &container.Config{
		Image: image,
		Cmd:   []string{"minetestserver", "--world", world_dir_container, "--config", minetest_config_container},
		Tty:   false,
		User:  fmt.Sprintf("%d", os.Getuid()),
	}, &container.HostConfig{
		RestartPolicy: container.RestartPolicy{
			Name: "always",
		},
		Mounts: []mount.Mount{
			{
				Type:   mount.TypeBind,
				Source: world_dir,
				Target: world_dir_container,
			}, {
				Type:   mount.TypeBind,
				Source: minetest_config,
				Target: minetest_config_container,
			},
		},
		PortBindings: nat.PortMap{
			"30000/udp": []nat.PortBinding{
				{
					HostIP:   "0.0.0.0",
					HostPort: fmt.Sprintf("%d", port),
				},
			},
		},
	}, &network.NetworkingConfig{
		EndpointsConfig: map[string]*network.EndpointSettings{
			network_name: {NetworkID: network_name},
		},
	}, nil, container_name)
	if err != nil {
		SendError(w, 500, fmt.Sprintf("container create error: %v", err))
		return
	}

	SendJson(w, &ServiceStatus{
		ID:      resp.ID,
		Created: true,
		Running: false,
		Version: cer.Version,
	})
}

func (a *Api) StartEngine(w http.ResponseWriter, r *http.Request, claims *types.Claims) {
	c, err := getMinetestContainer()
	if err != nil {
		SendError(w, 500, fmt.Sprintf("fetch container error: %v", err))
		return
	}
	if c == nil {
		SendError(w, 404, "no container found")
		return
	}

	cli, err := getDockerCli()
	if err != nil {
		SendError(w, 500, fmt.Sprintf("docker client error: %v", err))
		return
	}

	ctx := context.Background()
	err = cli.ContainerStart(ctx, c.ID, dockertypes.ContainerStartOptions{})
	Send(w, true, err)
}

func (a *Api) StopEngine(w http.ResponseWriter, r *http.Request, claims *types.Claims) {
	c, err := getMinetestContainer()
	if err != nil {
		SendError(w, 500, fmt.Sprintf("fetch container error: %v", err))
		return
	}
	if c == nil {
		SendError(w, 404, "no container found")
		return
	}

	cli, err := getDockerCli()
	if err != nil {
		SendError(w, 500, fmt.Sprintf("docker client error: %v", err))
		return
	}

	ctx := context.Background()
	err = cli.ContainerStop(ctx, c.ID, container.StopOptions{})
	Send(w, true, err)
}

func (a *Api) RemoveEngine(w http.ResponseWriter, r *http.Request, claims *types.Claims) {
	c, err := getMinetestContainer()
	if err != nil {
		SendError(w, 500, fmt.Sprintf("fetch container error: %v", err))
		return
	}
	if c == nil {
		SendError(w, 404, "no container found")
		return
	}

	cli, err := getDockerCli()
	if err != nil {
		SendError(w, 500, fmt.Sprintf("docker client error: %v", err))
		return
	}

	ctx := context.Background()
	err = cli.ContainerRemove(ctx, c.ID, dockertypes.ContainerRemoveOptions{})
	Send(w, true, err)
}

type ServiceLogResponse struct {
	Out string `json:"out"`
	Err string `json:"err"`
}

func (a *Api) GetEngineLogs(w http.ResponseWriter, r *http.Request, claims *types.Claims) {
	vars := mux.Vars(r)

	since_millis, err := strconv.ParseInt(vars["since"], 10, 64)
	if err != nil {
		SendError(w, 500, fmt.Sprintf("invalid since format: %s", vars["since"]))
		return
	}
	since := time.UnixMilli(since_millis).Format(time.RFC3339)

	until_millis, err := strconv.ParseInt(vars["until"], 10, 64)
	if err != nil {
		SendError(w, 500, fmt.Sprintf("invalid until format: %s", vars["until"]))
		return
	}
	until := time.UnixMilli(until_millis).Format(time.RFC3339)

	container, err := getMinetestContainer()
	if err != nil {
		SendError(w, 500, fmt.Sprintf("fetch container error: %v", err))
		return
	}
	if container == nil {
		SendError(w, 404, "container not found")
		return
	}

	cli, err := getDockerCli()
	if err != nil {
		SendError(w, 500, fmt.Sprintf("docker client error: %v", err))
		return
	}

	ctx := context.Background()
	logs, err := cli.ContainerLogs(ctx, container.ID, dockertypes.ContainerLogsOptions{
		ShowStdout: true,
		ShowStderr: true,
		Since:      since,
		Until:      until,
	})
	if err != nil {
		SendError(w, 500, fmt.Sprintf("docker stdout log error: %v", err))
		return
	}
	defer logs.Close()

	out_buf := bytes.NewBuffer([]byte{})
	err_buf := bytes.NewBuffer([]byte{})

	_, err = stdcopy.StdCopy(out_buf, err_buf, logs)
	if err != nil {
		SendError(w, 500, fmt.Sprintf("docker stdcopy error: %v", err))
		return
	}

	slogs := &ServiceLogResponse{
		Out: out_buf.String(),
		Err: err_buf.String(),
	}

	Send(w, slogs, nil)
}
