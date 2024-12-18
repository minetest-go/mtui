package web

import (
	"fmt"
	"mtui/types"
	"net/http"
)

func (a *Api) GetMaintenanceMode(w http.ResponseWriter, r *http.Request, c *types.Claims) {
	maint_mode := a.app.MaintenanceMode()
	Send(w, maint_mode, nil)
}

func (a *Api) EnableMaintenanceMode(w http.ResponseWriter, r *http.Request, c *types.Claims) {
	// create log entry
	a.app.CreateUILogEntry(&types.Log{
		Username: c.Username,
		Event:    "maintenance",
		Message:  fmt.Sprintf("User '%s' enables the maintenance mode", c.Username),
	}, r)

	// clear current stats
	current_stats.Store(nil)

	a.app.EnableMaintenanceMode()

	Send(w, true, nil)
}

func (a *Api) DisableMaintenanceMode(w http.ResponseWriter, r *http.Request, c *types.Claims) {

	a.app.DisableMaintenanceMode()

	// create log entry
	a.app.CreateUILogEntry(&types.Log{
		Username: c.Username,
		Event:    "maintenance",
		Message:  fmt.Sprintf("User '%s' disabled the maintenance mode", c.Username),
	}, r)

	Send(w, false, nil)
}
