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

type MeseconsSet struct {
	Pos   *Pos  `json:"pos"`
	State State `json:"state"`
}

// game -> ui
const COMMAND_MESECONS_EVENT bridge.CommandType = "mesecons_event"

type MeseconsEvent struct {
	Pos   *Pos   `json:"pos"`
	State State  `json:"state"`
	Color string `json:"color"`
	Name  string `json:"name"`
}
