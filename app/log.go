package app

import (
	"mtui/types"
	"net/http"
)

func (app *App) CreateUILogEntry(l *types.Log, r *http.Request) {
	if !app.MaintenanceMode.Load() {
		l.Category = types.CategoryUI
		app.GeoipResolver.ResolveLogGeoIP(l, r)
		app.Repos.LogRepository.Insert(l)
	}
}
