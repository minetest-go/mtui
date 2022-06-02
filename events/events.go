package events

import (
	"mtui/app"
)

func Setup(app *app.App) error {
	c := app.Bridge.AddHandler()
	go statsLoop(app.WSEvents, c)

	return nil
}
