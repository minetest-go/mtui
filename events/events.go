package events

import (
	"mtui/app"
	"mtui/types/command"
)

func Setup(app *app.App) error {
	go metricLoop(app.DB, app.Bridge.AddHandler(command.COMMAND_METRICS))
	go statsLoop(app.WSEvents, app.Bridge.AddHandler(command.COMMAND_STATS))
	go logLoop(app.Repos.LogRepository, app.GeoipResolver, app.Bridge.AddHandler(command.COMMAND_LOG))

	return nil
}
