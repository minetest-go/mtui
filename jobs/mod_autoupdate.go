package jobs

import (
	"fmt"
	"mtui/app"
	"mtui/types"
	"time"

	"github.com/sirupsen/logrus"
)

func checkAllMods(a *app.App) error {
	err := a.ModManager.CheckUpdates()
	if err != nil {
		return err
	}

	mods, err := a.Repos.ModRepo.GetAll()
	if err != nil {
		return err
	}

	for _, mod := range mods {
		if !mod.AutoUpdate {
			continue
		}

		if mod.Version != mod.LatestVersion {
			err = a.ModManager.Update(mod, mod.LatestVersion)
			if err != nil {
				return err
			}

			// create log entry
			log := &types.Log{
				Category: types.CategoryUI,
				Event:    "mods",
				Message:  fmt.Sprintf("Auto-updated the %s '%s' (%s) to version '%s'", mod.ModType, mod.Name, mod.SourceType, mod.LatestVersion),
			}
			err = a.Repos.LogRepository.Insert(log)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func modAutoUpdate(a *app.App) {
	for {
		if !a.MaintenanceMode() {
			err := checkAllMods(a)
			if err != nil {
				logrus.WithError(err).Warn("mod auto-update failed")
			}
		}
		time.Sleep(time.Minute * 30)
	}
}
