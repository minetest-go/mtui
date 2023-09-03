package web

import (
	"mtui/app"
	"mtui/types/command"
)

type Api struct {
	app *app.App
}

func NewApi(app *app.App) *Api {
	return &Api{
		app: app,
	}
}

func (api *Api) Setup() error {

	// start tan login listener
	go api.TanSetListener(api.app.Bridge.AddHandler(command.COMMAND_TAN_SET))
	go api.TanSetListener(api.app.Bridge.AddHandler(command.COMMAND_TAN_REMOVE))
	go api.StatsEventListener(api.app.Bridge.AddHandler(command.COMMAND_STATS))

	return nil
}
