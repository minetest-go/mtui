package command

import "mtui/bridge"

const (
	COMMAND_PING bridge.CommandType = "ping"
)

// ping command from the engine
type PingCommand struct{}
