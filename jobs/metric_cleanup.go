package jobs

import (
	"fmt"
	"mtui/app"
	"time"
)

func metricCleanup(a *app.App) {
	for {
		if !a.MaintenanceMode.Load() {
			ts := time.Now().Add(time.Hour * -1)
			err := a.Repos.MetricRepository.DeleteBefore(ts.UnixMilli())
			if err != nil {
				fmt.Printf("metric cleanup error: %s\n", err.Error())
			}
		}

		// re-schedule every minute
		time.Sleep(time.Second * 10)
	}
}
