package jobs

import "mtui/app"

func Start(a *app.App) {
	go logCleanup(a.Repos.LogRepository)
	go metricCleanup(a.Repos.MetricRepository)
	go mediaScan(a)
}
