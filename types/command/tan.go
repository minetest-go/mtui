package command

import "mtui/bridge"

const (
	COMMAND_TAN_SET    bridge.CommandResponseType = "tan_set"
	COMMAND_TAN_REMOVE bridge.CommandResponseType = "tan_remove"
)

// tan command from the engine
type TanCommand struct {
	Playername string `json:"playername"`
	TAN        string `json:"tan"`
}
