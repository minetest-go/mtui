package command

// commands from the game to the ui

import (
	"mtui/bridge"
)

const (
	COMMAND_CHAT_SEND_PLAYER bridge.CommandType = "chat_send_player"
	COMMAND_CHAT_SEND_ALL    bridge.CommandType = "chat_send_all"
)

type ChatSendPlayerNotification struct {
	Name string `json:"name"`
	Text string `json:"text"`
}

type ChatSendAllNotification struct {
	Text string `json:"text"`
}
