package dockerservice

import (
	"fmt"

	"github.com/docker/docker/client"
)

func getDockerCli() (*client.Client, error) {
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		return nil, fmt.Errorf("new client error: %v", err)
	}
	return cli, nil
}
