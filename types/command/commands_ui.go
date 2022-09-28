package command

// commands from the ui to ingame

import (
	"encoding/json"
	"mtui/bridge"
)

const (
	COMMAND_CHATCMD bridge.CommandType = "execute_command"
	COMMAND_LUA     bridge.CommandType = "lua"
)

type ExecuteChatCommandRequest struct {
	Playername string `json:"playername"`
	Command    string `json:"command"`
}

type ExecuteChatCommandResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}

type LuaRequest struct {
	Code string `json:"code"`
}

type LuaResponse struct {
	Success bool            `json:"success"`
	Message string          `json:"message"`
	Result  json.RawMessage `json:"result"`
}
