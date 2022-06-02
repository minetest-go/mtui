package web

import (
	"mtui/app"
	"mtui/types"
)

type Api struct {
	app *app.App

	// outgoing commands to the engine
	tx_cmds chan *types.Command
	// incoming commands from the engine
	rx_cmds chan *types.Command
}

func NewApi(app *app.App) *Api {
	return &Api{
		app:     app,
		tx_cmds: make(chan *types.Command, 1000),
		rx_cmds: make(chan *types.Command, 1000),
	}
}
