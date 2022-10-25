package jobs

import "mtui/app"

func Start(a *app.App) {
	go logcleanup(a.Repos.LogRepository)
}
