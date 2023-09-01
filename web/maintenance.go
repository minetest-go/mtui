package web

import (
	"mtui/types"
	"net/http"
)

func (a *Api) GetMaintenanceMode(w http.ResponseWriter, r *http.Request, c *types.Claims) {
	maint_mode := a.app.MaintenanceMode.Load()
	Send(w, maint_mode, nil)
}

func (a *Api) EnableMaintenanceMode(w http.ResponseWriter, r *http.Request, c *types.Claims) {
	a.app.MaintenanceMode.Store(true)
	//TODO: detach database
	Send(w, true, nil)
}

func (a *Api) DisableMaintenanceMode(w http.ResponseWriter, r *http.Request, c *types.Claims) {
	a.app.MaintenanceMode.Store(false)
	//TODO: attach database
	Send(w, false, nil)
}
