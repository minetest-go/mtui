package command

// chat notifications and commands

import (
	"mtui/bridge"
)

const (
	// game -> ui
	COMMAND_CHAT_NOTIFICATION bridge.CommandType = "chat_notification"
	// ui -> game
	COMMAND_CHAT_SEND bridge.CommandType = "chat_send"
)

type ChatMessage struct {
	Channel string `json:"channel"`
	Name    string `json:"name"`
	Message string `json:"message"`
}
