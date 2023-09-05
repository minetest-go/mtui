package jobs

import (
	"fmt"
	"mtui/app"
	"time"
)

func logCleanup(a *app.App) {
	for {
		if !a.MaintenanceMode.Load() {
			ts := time.Now().AddDate(0, 0, -30)
			err := a.Repos.LogRepository.DeleteBefore(ts.UnixMilli())
			if err != nil {
				fmt.Printf("Log cleanup error: %s\n", err.Error())
			}
		}

		// re-schedule every minute
		time.Sleep(time.Second * 10)
	}
}
