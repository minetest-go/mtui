package command

import (
	"mtui/bridge"
	"mtui/types"
)

type State string

const (
	StateOn  State = "on"
	StateOff State = "off"
)

// ui -> game
const COMMAND_MESECONS_SET bridge.CommandType = "mesecons_set"

type MeseconsSetRequest struct {
	Pos      *types.Pos `json:"pos"`
	State    State      `json:"state"`
	Nodename string     `json:"nodename"`
}

type MeseconsSetResponse struct {
	Success bool `json:"success"`
}

// game -> ui
const COMMAND_MESECONS_EVENT bridge.CommandType = "mesecons_event"

type MeseconsEvent struct {
	Pos      *types.Pos `json:"pos"`
	State    State      `json:"state"`
	Nodename string     `json:"nodename"`
}

const COMMAND_MESECONS_REGISTER bridge.CommandType = "mesecons_register"

type MeseconsRegister struct {
	Pos        *types.Pos `json:"pos"`
	Playername string     `json:"playername"`
	Nodename   string     `json:"nodename"`
}
