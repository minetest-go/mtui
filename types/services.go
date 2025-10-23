package types

var EngineServiceImages = map[string]string{
	"5.9.0":  "ghcr.io/minetest/minetest:5.9.0",
	"5.10.0": "ghcr.io/minetest/minetest:5.10.0",
	"5.11.0": "ghcr.io/luanti-org/luanti:5.11.0",
	"5.12.0": "ghcr.io/luanti-org/luanti:5.12.0",
	"5.13.0": "ghcr.io/luanti-org/luanti:5.13.0",
	"5.14.0": "ghcr.io/luanti-org/luanti:5.14.0",
}

// for auto install
var EngineServiceLatest = "5.14.0"

var MatterbridgeServiceImages = map[string]string{
	"1.26.0": "42wim/matterbridge:1.26.0",
}

var MapserverServiceImages = map[string]string{
	"4.8.0":  "minetestmapserver/mapserver:4.8.0",
	"4.9.1":  "ghcr.io/minetest-mapserver/mapserver:v4.9.1",
	"4.9.3":  "ghcr.io/minetest-mapserver/mapserver:v4.9.3",
	"4.9.4":  "ghcr.io/minetest-mapserver/mapserver:v4.9.4",
	"latest": "ghcr.io/minetest-mapserver/mapserver:latest",
}
