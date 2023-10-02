package events

import (
	"encoding/json"
	"fmt"
	"mtui/app"
	"mtui/bridge"
	"mtui/types"
	"mtui/types/command"
)

func chatLoop(a *app.App, ch chan *bridge.CommandResponse) {
	for cmd := range ch {
		msg := &command.ChatMessage{}
		err := json.Unmarshal(cmd.Data, msg)
		if err != nil {
			fmt.Printf("Chat notification payload error: %s\n", err.Error())
			continue
		}

		err = a.Repos.ChatLogRepo.Insert(&types.ChatLog{
			Channel: msg.Channel,
			Name:    msg.Name,
			Message: msg.Message,
		})
		if err != nil {
			fmt.Printf("Chat insert error: %s\n", err.Error())
			continue
		}

	}
}
