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
)

type PingCommand struct {
}

type TanCommand struct {
	Playername string `json:"playername"`
	TAN        string `json:"tan"`
}

type StatsCommand struct {
	Uptime      float64 `json:"uptime"`
	MaxLag      float64 `json:"max_lag"`
	TimeOfDay   float64 `json:"time_of_day"`
	PlayerCount float64 `json:"player_count"`
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
	}

	if payload != nil {
		err = json.Unmarshal(cmd.Data, payload)
	}

	return payload, err
}
