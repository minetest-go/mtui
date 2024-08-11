package types

var EngineServiceImages = map[string]string{
	"5.6.0": "registry.gitlab.com/minetest/minetest/server:5.6.0",
	"5.7.0": "registry.gitlab.com/minetest/minetest/server:5.7.0",
	"5.8.0": "ghcr.io/minetest-hosting/minetest-docker:5.8.0",
	"5.9.0": "ghcr.io/minetest/minetest:5.9.0",
}

// for auto install
var EngineServiceLatest = "5.9.0"

var MatterbridgeServiceImages = map[string]string{
	"1.26.0": "42wim/matterbridge:1.26.0",
}

var MapserverServiceImages = map[string]string{
	"4.7.0":  "minetestmapserver/mapserver:4.7.0",
	"4.8.0":  "minetestmapserver/mapserver:4.8.0",
	"4.9.1":  "ghcr.io/minetest-mapserver/mapserver:v4.9.1",
	"latest": "ghcr.io/minetest-mapserver/mapserver:latest",
}
