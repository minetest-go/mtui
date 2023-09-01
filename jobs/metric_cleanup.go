package jobs

import (
	"fmt"
	"mtui/db"
	"sync/atomic"
	"time"
)

func metricCleanup(r *db.MetricRepository, maint_mode *atomic.Bool) {
	for {
		if !maint_mode.Load() {
			ts := time.Now().Add(time.Hour * -1)
			err := r.DeleteBefore(ts.UnixMilli())
			if err != nil {
				fmt.Printf("metric cleanup error: %s\n", err.Error())
			}
		}

		// re-schedule every minute
		time.Sleep(time.Minute)
	}
}
