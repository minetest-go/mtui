package jobs

import (
	"fmt"
	"mtui/app"
	"mtui/modmanager"
	"mtui/types"
	"mtui/types/command"
	"time"

	"github.com/sirupsen/logrus"
)

func checkAllMods(a *app.App) error {
	mods, err := a.Repos.ModRepo.GetAll()
	if err != nil {
		return err
	}

	changed_mods, err := modmanager.CheckUpdates(a.WorldDir, mods)
	if err != nil {
		return err
	}

	mods_changed := false

	for _, mod := range changed_mods {

		if mod.Version != mod.LatestVersion && mod.AutoUpdate {
			err = a.ModManager.Update(mod, mod.LatestVersion)
			if err != nil {
				logrus.WithError(err).WithFields(logrus.Fields{
					"name":            mod.Name,
					"version":         mod.Version,
					"lastest_version": mod.LatestVersion,
					"id":              mod.ID,
				}).Error("mod auto update failed")

				// log entry
				log := &types.Log{
					Category: types.CategoryUI,
					Event:    "mods",
					Message:  fmt.Sprintf("Auto-update failed for %s '%s' (%s) to version '%s'", mod.ModType, mod.Name, mod.SourceType, mod.LatestVersion),
				}
				err = a.Repos.LogRepository.Insert(log)
				if err != nil {
					return err
				}

				continue
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

			mods_changed = true
		}

		err = a.Repos.ModRepo.Update(mod)
		if err != nil {
			return err
		}
	}

	if mods_changed {
		err = a.Bridge.ExecuteCommand(command.COMMAND_NOTIFY_MODS_CHANGED, nil, nil, time.Second*5)
		if err != nil {
			// ignore error, just log
			logrus.WithError(err).Warn("mods updated notification failed")
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
