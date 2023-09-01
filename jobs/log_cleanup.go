package jobs

import (
	"fmt"
	"mtui/db"
	"sync/atomic"
	"time"
)

func logCleanup(r *db.LogRepository, maint_mode *atomic.Bool) {
	for {
		if !maint_mode.Load() {
			ts := time.Now().AddDate(0, 0, -30)
			err := r.DeleteBefore(ts.UnixMilli())
			if err != nil {
				fmt.Printf("Log cleanup error: %s\n", err.Error())
			}
		}

		// re-schedule every minute
		time.Sleep(time.Minute)
	}
}
