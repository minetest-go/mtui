package web

import (
	"mtui/app"
	"mtui/types"
	"mtui/types/command"
	"net/http"
)

type Api struct {
	app *app.App
}

func NewApi(app *app.App) *Api {
	return &Api{app: app}
}

func (api *Api) Setup() error {

	// start tan login listener
	go api.TanSetListener(api.app.Bridge.AddHandler(command.COMMAND_TAN_SET))
	go api.TanSetListener(api.app.Bridge.AddHandler(command.COMMAND_TAN_REMOVE))
	go api.StatsEventListener(api.app.Bridge.AddHandler(command.COMMAND_STATS))

	api.CreateUILogEntry(&types.Log{
		Event:   "system",
		Message: "mtui started",
	}, nil)

	return nil
}

func (api *Api) CreateUILogEntry(l *types.Log, r *http.Request) {
	if !api.app.MaintenanceMode.Load() {
		l.Category = types.CategoryUI
		api.app.GeoipResolver.ResolveLogGeoIP(l, r)
		api.app.Repos.LogRepository.Insert(l)
	}
}
