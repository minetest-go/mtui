package events

import (
	"encoding/json"
	"fmt"
	"mtui/app"
	"mtui/bridge"
	"mtui/types"
)

func logLoop(a *app.App, ch chan *bridge.CommandResponse) {
	for cmd := range ch {
		log := &types.Log{}
		err := json.Unmarshal(cmd.Data, log)
		if err != nil {
			fmt.Printf("Payload error: %s\n", err.Error())
			continue
		}

		a.ResolveLogGeoIP(log, nil)
		err = a.Repos.LogRepository.Insert(log)
		if err != nil {
			fmt.Printf("DB error: %s\n", err.Error())
			continue
		}
	}
}
