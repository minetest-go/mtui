package app

import (
	"fmt"
	"mtui/dockerservice"
	"os"
	"path"
	"strings"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/mount"
	"github.com/docker/go-connections/nat"
)

func (app *App) SetupServices() {

	if app.Config.DockerContainerPrefix != "" {
		// docker management enabled, set up service utils
		no_proxy_env := []string{"HTTP_PROXY=", "HTTPS_PROXY=", "http_proxy=", "https_proxy="}

		// minetest
		portbinding := fmt.Sprintf("%d/udp", app.Config.DockerMinetestPort)
		os.MkdirAll(path.Join(app.Config.WorldDir, "textures"), 0777)

		app.ServiceEngine = dockerservice.New(&dockerservice.Config{
			ContainerName: fmt.Sprintf("%s_engine", app.Config.DockerContainerPrefix),
			Networks:      strings.Split(app.Config.DockerNetwork, ","),
			DefaultConfig: &container.Config{
				Cmd:  []string{"--world", "/world", "--config", "/minetest.conf"},
				Tty:  false,
				User: fmt.Sprintf("%d", os.Getuid()),
				Env:  no_proxy_env,
				ExposedPorts: nat.PortSet{
					nat.Port(fmt.Sprintf("%d/udp", app.Config.DockerMinetestPort)): {},
				},
			},
			DefaultHostConfig: &container.HostConfig{
				RestartPolicy: container.RestartPolicy{
					Name: "always",
				},
				Mounts: []mount.Mount{
					{
						Type:   mount.TypeBind,
						Source: app.Config.DockerWorlddir,
						Target: "/world",
					}, {
						Type:   mount.TypeBind,
						Source: app.Config.DockerMinetestConfig,
						Target: "/minetest.conf",
					}, {
						Type:   mount.TypeBind,
						Source: path.Join(app.Config.DockerWorlddir, "textures"),
						Target: "/root/.minetest/textures/server", //TODO: only works in uid=0 case
					},
				},
				PortBindings: nat.PortMap{
					nat.Port(portbinding): []nat.PortBinding{
						{
							HostIP:   "",
							HostPort: fmt.Sprintf("%d", app.Config.DockerMinetestPort),
						},
					},
				},
			},
		})

		// matterbridge
		app.ServiceMatterbridge = dockerservice.New(&dockerservice.Config{
			ContainerName: fmt.Sprintf("%s_matterbridge", app.Config.DockerContainerPrefix),
			Networks:      strings.Split(app.Config.DockerNetwork, ","),
			DefaultConfig: &container.Config{
				Env: no_proxy_env,
			},
			DefaultHostConfig: &container.HostConfig{
				RestartPolicy: container.RestartPolicy{
					Name: "always",
				},
				Mounts: []mount.Mount{
					{
						Type:   mount.TypeBind,
						Source: path.Join(app.Config.DockerWorlddir, "matterbridge.toml"),
						Target: "/etc/matterbridge/matterbridge.toml",
					},
				},
			},
		})

		// mapserver
		tfns := fmt.Sprintf("%s-mapserver", app.Config.DockerContainerPrefix)
		app.ServiceMapserver = dockerservice.New(&dockerservice.Config{
			ContainerName: fmt.Sprintf("%s_mapserver", app.Config.DockerContainerPrefix),
			Networks:      strings.Split(app.Config.DockerNetwork, ","),
			DefaultConfig: &container.Config{
				Env:        no_proxy_env,
				WorkingDir: "/world",
				Labels: map[string]string{
					"promtail":               "true",
					"traefik.enable":         "true",
					"traefik.docker.network": "terminator",
					"traefik.http.services." + tfns + ".loadbalancer.server.port":            "8080",
					"traefik.http.routers." + tfns + ".rule":                                 fmt.Sprintf("Host(`%s`) && PathPrefix(`/map`)", app.Config.CookieDomain),
					"traefik.http.routers." + tfns + ".entrypoints":                          "websecure",
					"traefik.http.routers." + tfns + ".tls.certresolver":                     "default",
					"traefik.http.routers." + tfns + ".middlewares":                          fmt.Sprintf("%s-stripprefix", tfns),
					"traefik.http.middlewares." + tfns + "-stripprefix.stripprefix.prefixes": "/map",
				},
			},
			DefaultHostConfig: &container.HostConfig{
				RestartPolicy: container.RestartPolicy{
					Name: "always",
				},
				Mounts: []mount.Mount{
					{
						Type:   mount.TypeBind,
						Source: app.Config.DockerWorlddir,
						Target: "/world",
					},
				},
			},
		})
	}

}
