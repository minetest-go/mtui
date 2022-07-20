package events

import (
	"mtui/app"
	"mtui/types/command"
)

func Setup(app *app.App) error {
	c := app.Bridge.AddHandler(command.COMMAND_STATS)
	go statsLoop(app.WSEvents, c)

	return nil
}
