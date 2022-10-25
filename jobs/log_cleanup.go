package jobs

import (
	"fmt"
	"mtui/db"
	"time"
)

func logcleanup(r *db.LogRepository) {
	for {
		ts := time.Now().AddDate(0, 0, -30)
		err := r.RemoveBefore(ts.UnixMilli())
		if err != nil {
			fmt.Printf("Log cleanup error: %s\n", err.Error())
		}

		// re-schedule every minute
		time.Sleep(time.Minute)
	}
}
