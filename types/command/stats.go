package command

import (
	"mtui/bridge"
	"mtui/types"
)

const (
	COMMAND_STATS bridge.CommandType = "stats"
)

// player stats from the engine
type PlayerStats struct {
	Name   string        `json:"name"`
	HP     types.JsonInt `json:"hp"`
	Breath types.JsonInt `json:"breath"`

	Pos *struct {
		X types.JsonInt `json:"x"`
		Y types.JsonInt `json:"y"`
		Z types.JsonInt `json:"z"`
	} `json:"pos"`

	Info *struct {
		Address              string        `json:"address"`
		IPVersion            float64       `json:"ip_version"`
		ConnectionUptime     float64       `json:"connection_uptime"`
		ProtocolVersion      types.JsonInt `json:"protocol_version"`
		FormspecVersion      types.JsonInt `json:"formspec_version"`
		LangCode             string        `json:"lang_code"`
		MinRTT               float64       `json:"min_rtt"`
		MaxRTT               float64       `json:"max_rtt"`
		AvgRTT               float64       `json:"avg_rtt"`
		SerializationVersion float64       `json:"ser_vers"`
		VersionString        string        `json:"vers_string"`
	} `json:"info"`
}

type GlobalStats struct {
	RegisteredNodes    types.JsonInt `json:"registered_nodes"`
	RegisteredItems    types.JsonInt `json:"registered_items"`
	RegisteredEntities types.JsonInt `json:"registered_entities"`
	RegisteredABMs     types.JsonInt `json:"registered_abms"`
}

// stats from the engine
type StatsCommand struct {
	Uptime      types.JsonInt  `json:"uptime"`
	MaxLag      float64        `json:"max_lag"`
	TimeOfDay   float64        `json:"time_of_day"`
	Mem         types.JsonInt  `json:"mem"`
	GlobalStats *GlobalStats   `json:"global_stats"`
	PlayerCount types.JsonInt  `json:"player_count"`
	Players     []*PlayerStats `json:"players"`
}
