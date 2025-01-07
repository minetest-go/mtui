package events

import (
	"mtui/app"
	"mtui/types/command"
)

func Setup(app *app.App) error {
	go statsLoop(app.WSEvents, app.Bridge.AddHandler(command.COMMAND_STATS))
	go logLoop(app, app.Bridge.AddHandler(command.COMMAND_LOG))
	go chatLoop(app, app.Bridge.AddHandler(command.COMMAND_CHAT_NOTIFICATION))
	go meseconsRegister(app, app.Bridge.AddHandler(command.COMMAND_MESECONS_REGISTER))
	go meseconsEvent(app, app.Bridge.AddHandler(command.COMMAND_MESECONS_EVENT))

	return nil
}
