package jobs

import "mtui/app"

func Start(a *app.App) {
	go logCleanup(a)
	go metricCleanup(a)
	go mediaScan(a)
	go modAutoUpdate(a)
}
