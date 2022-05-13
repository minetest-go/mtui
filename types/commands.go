package types

import "encoding/json"

type CommandType string

const (
	COMMAND_PING CommandType = "ping"
)

type Command struct {
	Type string          `json:"type"`
	Data json.RawMessage `json:"data"`
}

type PingCommand struct {
}
