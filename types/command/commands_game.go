package command

// commands from the game to the ui

import (
	"mtui/bridge"
)

const (
	COMMAND_CHAT_SEND_PLAYER bridge.CommandResponseType = "chat_send_player"
	COMMAND_CHAT_SEND_ALL    bridge.CommandResponseType = "chat_send_all"
	COMMAND_CHATCMD_RES      bridge.CommandResponseType = "execute_command"
)

type ChatSendPlayerNotification struct {
	Name string `json:"name"`
	Text string `json:"text"`
}

type ChatSendAllNotification struct {
	Text string `json:"text"`
}
