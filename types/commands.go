package types

import (
	"encoding/json"
	"mtui/bridge"
)

const (
	COMMAND_PING       bridge.CommandType = "ping"
	COMMAND_STATS      bridge.CommandType = "stats"
	COMMAND_CHATCMD    bridge.CommandType = "execute_command"
	COMMAND_TAN_SET    bridge.CommandType = "tan_set"
	COMMAND_TAN_REMOVE bridge.CommandType = "tan_remove"
	COMMAND_FEATURES   bridge.CommandType = "features"
)

type PingCommand struct {
}

// tan command from the engine
type TanCommand struct {
	Playername string `json:"playername"`
	TAN        string `json:"tan"`
}

// player stats from the engine
type PlayerStats struct {
	Name                 string  `json:"name"`
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
}

// stats from the engine
type StatsCommand struct {
	Uptime      float64 `json:"uptime"`
	MaxLag      float64 `json:"max_lag"`
	TimeOfDay   float64 `json:"time_of_day"`
	PlayerCount float64 `json:"player_count"`
	Players     []*PlayerStats
}

// ingame features
type FeaturesCommand struct {
	Mail bool `json:"mail"`
}

type ExecuteChatCommandRequest struct {
	Playername string `json:"playername"`
	Command    string `json:"command"`
}

type ExecuteChatCommandResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}

// Parse an incoming command
func ParseCommand(cmd *bridge.Command) (interface{}, error) {
	var err error
	var payload interface{}

	switch cmd.Type {
	case COMMAND_PING:
		payload = &PingCommand{}
	case COMMAND_STATS:
		payload = &StatsCommand{}
	case COMMAND_CHATCMD:
		payload = &ExecuteChatCommandResponse{}
	case COMMAND_TAN_SET:
		payload = &TanCommand{}
	case COMMAND_TAN_REMOVE:
		payload = &TanCommand{}
	case COMMAND_FEATURES:
		payload = &FeaturesCommand{}
	}

	if payload != nil {
		err = json.Unmarshal(cmd.Data, payload)
	}

	return payload, err
}
