package app

import (
	"fmt"
	"mtui/types/command"
	"time"

	"github.com/minetest-go/mtdb/player"
)

func (a *App) Healthcheck() error {
	// mtui db
	_, err := a.Repos.FeatureRepository.GetAll()
	if err != nil {
		return fmt.Errorf("mtui db error: %v", err)
	}

	// minetest db
	if a.DBContext.ModStorage != nil {
		_, err = a.DBContext.ModStorage.Count()
		if err != nil {
			return fmt.Errorf("modstorage db error: %v", err)
		}
	}
	_, err = a.DBContext.Player.Count(&player.PlayerSearch{})
	if err != nil {
		return fmt.Errorf("player db error: %v", err)
	}

	// bridge
	err = a.Bridge.ExecuteCommand(command.COMMAND_PING, nil, &command.PingCommand{}, time.Second*5)
	if err != nil {
		return fmt.Errorf("bridge unresponsive: %v", err)
	}

	return nil
}
