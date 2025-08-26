package jobs

import (
	"mtui/app"
	"mtui/types"
	"strings"
	"time"

	"github.com/hpcloud/tail"
	"github.com/sirupsen/logrus"
)

func tailLogfile(a *app.App) {
	for {
		if !a.MaintenanceMode() {
			//TODO
			t, err := tail.TailFile(a.Config.TailEngineLogfile, tail.Config{Follow: true})
			if err != nil {
				logrus.WithError(err).Error("TailFile error")
			} else {
				// logfile opened
				for line := range t.Lines {
					if line.Err != nil {
						logrus.WithError(line.Err).Error("TailFile.lines error")
						break
					}
					if a.MaintenanceMode() {
						// switched to maintenance mode, skip log insertion
						break
					}

					eventname := "engine"
					if strings.Contains(line.Text, "ERROR[Main]") {
						eventname = "engine-mod-errors"
					}

					l := &types.Log{
						Category: types.CategoryService,
						Event:    eventname,
						Message:  line.Text,
					}
					err = a.Repos.LogRepository.Insert(l)
					if err != nil {
						logrus.WithError(err).Error("could not insert log")
						return
					}
				}
				t.Cleanup()
			}
		}

		// re-schedule
		time.Sleep(time.Second * 10)
	}
}
