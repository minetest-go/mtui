package jobs

import "mtui/app"

func Start(a *app.App) {
	go logCleanup(a.Repos.LogRepository, a.MaintenanceMode)
	go metricCleanup(a.Repos.MetricRepository, a.MaintenanceMode)
	go mediaScan(a, a.MaintenanceMode)
}
