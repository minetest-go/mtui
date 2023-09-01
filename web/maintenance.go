package web

import (
	"mtui/types"
	"net/http"

	"github.com/sirupsen/logrus"
)

func (a *Api) GetMaintenanceMode(w http.ResponseWriter, r *http.Request, c *types.Claims) {
	maint_mode := a.app.MaintenanceMode.Load()
	Send(w, maint_mode, nil)
}

func (a *Api) EnableMaintenanceMode(w http.ResponseWriter, r *http.Request, c *types.Claims) {
	//TODO: this could lead to a race condition if called fast enough
	if a.app.MaintenanceMode.Load() {
		SendError(w, 500, "already in maintenance mode")
	}
	a.app.MaintenanceMode.Store(true)
	err := a.app.DetachDatabase()
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"err": err,
		}).Error("Detach database failed")
	}
	Send(w, true, err)
}

func (a *Api) DisableMaintenanceMode(w http.ResponseWriter, r *http.Request, c *types.Claims) {
	//TODO: this could lead to a race condition if called fast enough
	if !a.app.MaintenanceMode.Load() {
		SendError(w, 500, "maintenance mode already disabled")
	}
	a.app.MaintenanceMode.Store(false)
	err := a.app.AttachDatabase()
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"err": err,
		}).Error("Attach database failed")
	}
	Send(w, false, err)
}
