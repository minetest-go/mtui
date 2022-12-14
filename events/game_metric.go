package events

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"mtui/bridge"
	"mtui/db"
	"mtui/types"
	"mtui/types/command"
	"time"
)

func metricLoop(DB *sql.DB, ch chan *bridge.CommandResponse) {
	for cmd := range ch {
		metrics := make([]*command.GameMetric, 0)
		err := json.Unmarshal(cmd.Data, &metrics)
		if err != nil {
			fmt.Printf("Payload error: %s\n", err.Error())
			return
		}

		tx, err := DB.Begin()
		if err != nil {
			fmt.Printf("Tx begin error: %s\n", err.Error())
			return
		}
		defer tx.Rollback()

		repos := db.NewRepositories(tx)
		for _, metric := range metrics {
			err = repos.MetricTypeRepository.Insert(&types.MetricType{
				Name: metric.Name,
				Type: metric.Type,
				Help: metric.Help,
			})
			if err != nil {
				fmt.Printf("MetricType insert error: %s\n", err.Error())
				return
			}

			err = repos.MetricRepository.Insert(&types.Metric{
				Timestamp: time.Now().UnixMilli(),
				Name:      metric.Name,
				Value:     metric.Value,
			})
			if err != nil {
				fmt.Printf("Metric insert error: %s\n", err.Error())
				return
			}
		}

		err = tx.Commit()
		if err != nil {
			fmt.Printf("Commit error: %s\n", err.Error())
			return
		}
	}
}
