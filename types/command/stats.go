package command

import "mtui/bridge"

const (
	COMMAND_STATS bridge.CommandResponseType = "stats"
)

// player stats from the engine
type PlayerStats struct {
	Name   string  `json:"name"`
	HP     float64 `json:"hp"`
	Breath float64 `json:"breath"`

	Pos struct {
		X float64 `json:"x"`
		Y float64 `json:"y"`
		Z float64 `json:"z"`
	} `json:"pos"`

	Info struct {
		Address              string  `json:"address"`
		IPVersion            float64 `json:"ip_version"`
		ConnectionUptime     float64 `json:"connection_uptime"`
		ProtocolVersion      float64 `json:"protocol_version"`
		FormspecVersion      float64 `json:"formspec_version"`
		LangCode             string  `json:"lang_code"`
		MinRTT               float64 `json:"min_rtt"`
		MaxRTT               float64 `json:"max_rtt"`
		AvgRTT               float64 `json:"avg_rtt"`
		SerializationVersion float64 `json:"ser_vers"`
		VersionString        string  `json:"vers_string"`
	} `json:"info"`
}

// stats from the engine
type StatsCommand struct {
	Uptime      float64 `json:"uptime"`
	MaxLag      float64 `json:"max_lag"`
	TimeOfDay   float64 `json:"time_of_day"`
	PlayerCount float64 `json:"player_count"`
	Players     []*PlayerStats
}
