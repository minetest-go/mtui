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
		t, err := tail.TailFile(a.Config.TailEngineLogfile, tail.Config{Follow: true})
		if err != nil {
			logrus.WithError(err).Error("TailFile error")
		} else {
			// logfile opened
			if !processLines(a, t.Lines) {
				return
			}
			t.Cleanup()
		}

		// re-schedule
		time.Sleep(time.Second * 10)
	}
}

const errorFlushDelay = time.Millisecond * 500

func insertEngineLog(a *app.App, event, message string) error {
	return a.Repos.LogRepository.Insert(&types.Log{
		Category: types.CategoryService,
		Event:    event,
		Message:  message,
	})
}

// processLines consumes log lines, merging consecutive mod-error lines into a single multi-line log entry.
// returns false on a fatal insert error
func processLines(a *app.App, lines chan *tail.Line) bool {
	var pending []string

	flush := func() bool {
		if len(pending) == 0 {
			return true
		}
		err := insertEngineLog(a, "engine-mod-errors", strings.Join(pending, "\n"))
		pending = nil
		if err != nil {
			logrus.WithError(err).Error("could not insert log")
			return false
		}
		return true
	}

	for {
		var timeout <-chan time.Time
		if len(pending) > 0 {
			timeout = time.After(errorFlushDelay)
		}

		select {
		case <-timeout:
			if !flush() {
				return false
			}
		case line, ok := <-lines:
			if !ok {
				return flush()
			}
			if line.Err != nil {
				logrus.WithError(line.Err).Error("TailFile.lines error")
				flush()
				return true
			}

			if strings.Contains(line.Text, "ERROR[Main]") {
				pending = append(pending, line.Text)
				continue
			}

			// sequence ended
			if !flush() {
				return false
			}
			if err := insertEngineLog(a, "engine", line.Text); err != nil {
				logrus.WithError(err).Error("could not insert log")
				return false
			}
		}
	}
}
