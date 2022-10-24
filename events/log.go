package events

import (
	"encoding/json"
	"fmt"
	"mtui/bridge"
	"mtui/db"
	"mtui/types"
)

func logLoop(lr *db.LogRepository, ch chan *bridge.CommandResponse) {
	for cmd := range ch {
		log := &types.Log{}
		err := json.Unmarshal(cmd.Data, log)
		if err != nil {
			fmt.Printf("Payload error: %s\n", err.Error())
			return
		}

		err = lr.Insert(log)
		if err != nil {
			fmt.Printf("DB error: %s\n", err.Error())
			return
		}
	}
}
