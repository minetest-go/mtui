package command

// commands from the ui to ingame

import (
	"mtui/bridge"
)

const (
	COMMAND_CHATCMD bridge.CommandType = "execute_command"
)

type ExecuteChatCommandRequest struct {
	Playername string `json:"playername"`
	Command    string `json:"command"`
}

type ExecuteChatCommandResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}
