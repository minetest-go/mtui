package types

import "encoding/json"

type CommandType string

const (
	COMMAND_PING  CommandType = "ping"
	COMMAND_STATS CommandType = "stats"
)

type Command struct {
	Type CommandType     `json:"type"`
	ID   *float64        `json:"id"`
	Data json.RawMessage `json:"data"`
}

type PingCommand struct {
}

type StatsCommand struct {
	Uptime    float64 `json:"uptime"`
	MaxLag    float64 `json:"max_lag"`
	TimeOfDay float64 `json:"time_of_day"`
}

// Parse an incoming command
func ParseCommand(cmd *Command) (interface{}, error) {
	var err error
	var payload interface{}

	switch cmd.Type {
	case COMMAND_PING:
		payload = &PingCommand{}
	case COMMAND_STATS:
		payload = &StatsCommand{}
	}

	if payload != nil {
		err = json.Unmarshal(cmd.Data, payload)
	}

	return payload, err
}
