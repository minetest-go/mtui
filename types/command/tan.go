package command

import "mtui/bridge"

const (
	COMMAND_TAN_SET    bridge.CommandType = "tan_set"
	COMMAND_TAN_REMOVE bridge.CommandType = "tan_remove"
)

// tan command from the engine
type TanCommand struct {
	Playername string `json:"playername"`
	TAN        string `json:"tan"`
}
