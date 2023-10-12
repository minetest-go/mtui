package command

import (
	"mtui/bridge"
	"mtui/types"
)

type Pos struct {
	X types.JsonInt `json:"x"`
	Y types.JsonInt `json:"y"`
	Z types.JsonInt `json:"z"`
}

type State string

const (
	StateOn  State = "on"
	StateOff State = "off"
)

// ui -> game
const COMMAND_MESECONS_SET bridge.CommandType = "mesecons_set"

type MeseconsSetRequest struct {
	Pos      *Pos   `json:"pos"`
	State    State  `json:"state"`
	Nodename string `json:"nodename"`
}

type MeseconsSetRespone struct {
	Success          bool `json:"success"`
	NodenameMismatch bool `json:"nodename_mismatch"`
}

// game -> ui
const COMMAND_MESECONS_EVENT bridge.CommandType = "mesecons_event"

type MeseconsEvent struct {
	Pos      *Pos   `json:"pos"`
	State    State  `json:"state"`
	Color    string `json:"color"`
	Nodename string `json:"nodename"`
}
