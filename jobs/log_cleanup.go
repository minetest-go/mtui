package jobs

import (
	"fmt"
	"mtui/app"
	"time"
)

func logCleanup(a *app.App) {
	log_retention := time.Hour * 24 * 7 // 7 days default log retention
	if a.Config.LogRetention != "" {
		var err error
		log_retention, err = time.ParseDuration(a.Config.LogRetention)
		if err != nil {
			fmt.Printf("Log retention parsing of '%s' failed: %v, defaulting to 7 days\n", a.Config.LogRetention, err)
			log_retention = time.Hour * 24 * 7
		}
	}
	for {
		if !a.MaintenanceMode.Load() {
			ts := time.Now().Add(log_retention * -1)
			err := a.Repos.LogRepository.DeleteBefore(ts.UnixMilli())
			if err != nil {
				fmt.Printf("Log cleanup error: %s\n", err.Error())
			}
		}

		// re-schedule
		time.Sleep(time.Second * 10)
	}
}
