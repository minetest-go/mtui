package jobs

import (
	"fmt"
	"mtui/app"
	"mtui/dockerservice"
	"mtui/types"
	"strconv"
	"strings"
	"time"

	"github.com/sirupsen/logrus"
)

func collectServiceLogs(a *app.App, timestamp_key types.ConfigKey, event string, s *dockerservice.DockerService) {
	if s == nil {
		// nothing to collect
		return
	}
	e, err := a.Repos.ConfigRepo.GetByKey(timestamp_key)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"err":           err,
			"timestamp_key": timestamp_key,
		}).Error("could not get timestamp_key")
		return
	}

	var ts time.Time
	if e == nil {
		// start a minute ago by default
		ts = time.Now().Add(time.Minute * -1)
		e = &types.ConfigEntry{
			Key: timestamp_key,
		}
	} else {
		unixm, err := strconv.ParseInt(e.Value, 10, 64)
		if err != nil {
			logrus.WithFields(logrus.Fields{
				"err":   err,
				"value": e.Value,
			}).Error("could not parse timestamp millis")
			return
		}
		ts = time.UnixMilli(unixm)
	}
	now := time.Now()

	logs, err := s.GetLogs(ts, now)
	if err != nil {
		// silently ignore all docker related errors
		return
	}

	out_lines := strings.Split(logs.Out, "\n")
	err_lines := strings.Split(logs.Err, "\n")
	lines := append(out_lines, err_lines...)

	// log into log-db
	for _, line := range lines {
		if len(line) > 0 {
			l := &types.Log{
				Category: types.CategoryService,
				Event:    event,
				Message:  line,
			}
			err = a.Repos.LogRepository.Insert(l)
			if err != nil {
				logrus.WithError(err).Error("could not insert log")
				return
			}
		}
	}

	// save current timestamp for next query
	e.Value = fmt.Sprintf("%d", now.UnixMilli()+1)
	a.Repos.ConfigRepo.Set(e)
}

func serviceLogs(a *app.App) {
	for {
		if !a.MaintenanceMode.Load() {
			collectServiceLogs(a, "engine_log_timestamp", "engine", a.ServiceEngine)
			collectServiceLogs(a, "matterbridge_log_timestamp", "matterbridge", a.ServiceMatterbridge)
			collectServiceLogs(a, "mapserver_log_timestamp", "mapserver", a.ServiceMapserver)
			collectServiceLogs(a, "mtweb_log_timestamp", "mtweb", a.ServiceMTWeb)
		}

		// re-schedule
		time.Sleep(time.Second * 1)
	}
}
