package db

import (
	"mtui/types"
	"time"

	"github.com/google/uuid"
	"github.com/minetest-go/dbutil"
)

type ChatLogRepository struct {
	dbu *dbutil.DBUtil[*types.ChatLog]
}

func (r *ChatLogRepository) Insert(l *types.ChatLog) error {
	if l.ID == "" {
		l.ID = uuid.NewString()
	}

	if l.Timestamp == 0 {
		l.Timestamp = time.Now().UnixMilli()
	}

	return r.dbu.Insert(l)
}

func (r *ChatLogRepository) Search(channel string, from, to int64) ([]*types.ChatLog, error) {
	return r.dbu.SelectMulti("where channel = %s and timestamp > %s and timestamp < %s order by timestamp asc limit 1000", channel, from, to)
}

func (r *ChatLogRepository) GetLatest(channel string, limit int) ([]*types.ChatLog, error) {
	return r.dbu.SelectMulti("where channel = %s order by timestamp asc limit %s", channel, limit)
}

func (r *ChatLogRepository) DeleteBefore(timestamp int64) error {
	return r.dbu.Delete("where timestamp < %s", timestamp)
}
