package types

// commands from the ui to ingame

import (
	"mtui/bridge"
)

const (
	COMMAND_PING          bridge.CommandType = "ping"
	COMMAND_CHATCMD       bridge.CommandType = "execute_command"
	COMMAND_SEND_CHAT_MSG bridge.CommandType = "send_chat_message"
)

type PingCommand struct{}

type ExecuteChatCommandRequest struct {
	Playername string `json:"playername"`
	Command    string `json:"command"`
}

type ExecuteChatCommandResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}

type SendChatMessageRequest struct {
}
