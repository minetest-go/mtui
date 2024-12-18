package app

import (
	"time"

	"github.com/sirupsen/logrus"
)

func (a *App) EnableMaintenanceMode() {
	if a.maintenanceMode.Load() {
		// already in maintenance mode
		return
	}
	a.maintenanceMode.Store(true)

	// grace time
	time.Sleep(time.Second * 2)

	err := a.DetachDatabase()
	if err != nil {
		logrus.WithError(err).Error("Enable maintenance mode")
	}
}

func (a *App) DisableMaintenanceMode() {
	if !a.maintenanceMode.Load() {
		// already disabled
		return
	}

	err := a.AttachDatabase()
	if err != nil {
		logrus.WithError(err).Error("Disable maintenance mode")
	}

	a.maintenanceMode.Store(false)
}

func (a *App) MaintenanceMode() bool {
	return a.maintenanceMode.Load()
}
