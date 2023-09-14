package jobs

import (
	"bytes"
	"encoding/json"
	"fmt"
	"mtui/app"
	"mtui/types"
	"net/http"
	"strconv"
	"time"

	"github.com/sirupsen/logrus"
)

func uploadLogs(a *app.App) error {
	entry, err := a.Repos.ConfigRepo.GetByKey(types.ConfigLogStreamTimestamp)
	if err != nil {
		return err
	}

	if entry == nil {
		// initialize
		entry = &types.ConfigEntry{
			Key:   types.ConfigLogStreamTimestamp,
			Value: fmt.Sprintf("%d", time.Now().UnixMilli()),
		}

		err = a.Repos.ConfigRepo.Set(entry)
		if err != nil {
			return err
		}
	}

	from_ts, err := strconv.ParseInt(entry.Value, 10, 64)
	if err != nil {
		return err
	}

	logs, err := a.Repos.LogRepository.Search(&types.LogSearch{
		SortAscending: true,
		FromTimestamp: &from_ts,
	})
	if err != nil {
		return err
	}
	if len(logs) == 0 {
		// nothing to do
		return nil
	}

	new_ts := from_ts
	for _, log := range logs {
		if log.Timestamp > new_ts {
			new_ts = log.Timestamp
		}
	}

	data, err := json.Marshal(logs)
	if err != nil {
		return err
	}

	req, err := http.NewRequest(http.MethodPost, a.Config.LogStreamURL, bytes.NewBuffer(data))
	if err != nil {
		return err
	}

	req.Header.Set("Authorization", a.Config.LogStreamAuthorization)
	req.Header.Set("Content-Type", "application/json")

	c := http.Client{}
	resp, err := c.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// set new value
	entry.Value = fmt.Sprintf("%d", new_ts)
	return a.Repos.ConfigRepo.Set(entry)
}

func logStream(a *app.App) {
	for {
		if !a.MaintenanceMode.Load() {
			err := uploadLogs(a)
			if err != nil {
				logrus.WithError(err).Warn("couldn't stream logs")
			}
		}

		// re-schedule
		time.Sleep(time.Second * 10)
	}
}
