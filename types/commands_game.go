package types

// commands from the game to the ui

import (
	"encoding/json"
	"mtui/bridge"
)

const (
	COMMAND_STATS            bridge.CommandResponseType = "stats"
	COMMAND_TAN_SET          bridge.CommandResponseType = "tan_set"
	COMMAND_TAN_REMOVE       bridge.CommandResponseType = "tan_remove"
	COMMAND_CHAT_SEND_PLAYER bridge.CommandResponseType = "chat_send_player"
	COMMAND_CHAT_SEND_ALL    bridge.CommandResponseType = "chat_send_all"
	COMMAND_PING_RES         bridge.CommandResponseType = "ping"
	COMMAND_CHATCMD_RES      bridge.CommandResponseType = "execute_command"
)

// tan command from the engine
type TanCommand struct {
	Playername string `json:"playername"`
	TAN        string `json:"tan"`
}

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

type ChatSendPlayerNotification struct {
	Name string `json:"name"`
	Text string `json:"text"`
}

type ChatSendAllNotification struct {
	Text string `json:"text"`
}

// Parse an incoming command
func ParseCommand(cmd *bridge.CommandResponse) (interface{}, error) {
	var err error
	var payload interface{}

	switch cmd.Type {
	case COMMAND_PING_RES:
		payload = &PingCommand{}
	case COMMAND_STATS:
		payload = &StatsCommand{}
	case COMMAND_CHATCMD_RES:
		payload = &ExecuteChatCommandResponse{}
	case COMMAND_TAN_SET:
		payload = &TanCommand{}
	case COMMAND_TAN_REMOVE:
		payload = &TanCommand{}
	case COMMAND_CHAT_SEND_ALL:
		payload = &ChatSendAllNotification{}
	case COMMAND_CHAT_SEND_PLAYER:
		payload = &ChatSendPlayerNotification{}
	}

	if payload != nil {
		err = json.Unmarshal(cmd.Data, payload)
	}

	return payload, err
}
