package jobs

import "mtui/app"

func Start(a *app.App) {
	go logCleanup(a)
	go chatlogCleanup(a)
	go metricCleanup(a)
	go mediaScan(a)
	go modAutoUpdate(a)
	go serviceLogs(a)

	if a.Config.LogStreamURL != "" {
		go logStream(a)
	}
}
