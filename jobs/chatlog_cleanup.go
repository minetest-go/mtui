package jobs

import (
	"fmt"
	"mtui/app"
	"time"
)

func chatlogCleanup(a *app.App) {
	for {
		if !a.MaintenanceMode() {
			ts := time.Now().AddDate(0, 0, -30)
			err := a.Repos.ChatLogRepo.DeleteBefore(ts.UnixMilli())
			if err != nil {
				fmt.Printf("ChatLog cleanup error: %s\n", err.Error())
			}
		}

		// re-schedule
		time.Sleep(time.Second * 10)
	}
}
