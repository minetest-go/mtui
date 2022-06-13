package types

// commands from the ui to ingame

import (
	"mtui/bridge"
)

const (
	COMMAND_PING_REQ      bridge.CommandRequestType = "ping"
	COMMAND_CHATCMD_REQ   bridge.CommandRequestType = "execute_command"
	COMMAND_SEND_CHAT_MSG bridge.CommandRequestType = "send_chat_message"
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
