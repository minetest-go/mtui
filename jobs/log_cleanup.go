package jobs

import (
	"fmt"
	"mtui/app"
	"os"
	"time"
)

func logCleanup(a *app.App) {
	log_retention_str := os.Getenv("LOG_RETENTION")
	log_retention := time.Hour * 24 * 7 // 7 days default log retention
	if log_retention_str != "" {
		var err error
		log_retention, err = time.ParseDuration(log_retention_str)
		if err != nil {
			fmt.Printf("Log retention parsing of '%s' failed: %v, defaulting to 7 days\n", log_retention_str, err)
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
