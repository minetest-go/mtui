package db

import (
	"mtui/types"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ChatLogRepository struct {
	g *gorm.DB
}

func (r *ChatLogRepository) Insert(l *types.ChatLog) error {
	if l.ID == "" {
		l.ID = uuid.NewString()
	}

	if l.Timestamp == 0 {
		l.Timestamp = time.Now().UnixMilli()
	}

	return r.g.Create(l).Error
}

func (r *ChatLogRepository) Search(channel string, from, to int64) ([]*types.ChatLog, error) {
	return FindMulti[types.ChatLog](r.g.Where("timestamp > ?", from).Where("timestamp < ?", to).Where(types.ChatLog{Channel: channel}))
}

func (r *ChatLogRepository) GetLatest(channel string, limit int) ([]*types.ChatLog, error) {
	return FindMulti[types.ChatLog](r.g.Where(types.ChatLog{Channel: channel}).Order("timestamp ASC"))
}

func (r *ChatLogRepository) DeleteBefore(timestamp int64) error {
	return r.g.Where("timestamp < ?", timestamp).Delete(types.ChatLog{}).Error
}
