package db

import (
	"mtui/types"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type LogRepository struct {
	g *gorm.DB
}

func (r *LogRepository) Insert(l *types.Log) error {

	if l.ID == "" {
		l.ID = uuid.NewString()
	}

	if l.Timestamp == 0 {
		l.Timestamp = time.Now().UnixMilli()
	}

	return r.g.Create(l).Error
}

func (r *LogRepository) Update(l *types.Log) error {
	return r.g.Model(l).Updates(l).Error
}

func (r *LogRepository) DeleteBefore(timestamp int64) error {
	return r.g.Where("timestamp < ?", timestamp).Delete(types.Log{}).Error
}

func (r *LogRepository) query(s *types.LogSearch) *gorm.DB {
	q := r.g.Model(types.Log{})

	if s.ID != nil {
		q = q.Where(types.Log{ID: *s.ID})
	}

	if s.Category != nil {
		q = q.Where(types.Log{Category: *s.Category})
	}

	if s.Event != nil {
		q = q.Where(types.Log{Event: *s.Event})
	}

	if s.Username != nil {
		q = q.Where(types.Log{Username: *s.Username})
	}

	if s.Nodename != nil {
		q = q.Where(types.Log{Nodename: s.Nodename})
	}

	if s.PosX != nil {
		q = q.Where(types.Log{PosX: types.JsonIntPtr(int64(*s.PosX))})
	}

	if s.PosY != nil {
		q = q.Where(types.Log{PosY: types.JsonIntPtr(int64(*s.PosY))})
	}

	if s.PosZ != nil {
		q = q.Where(types.Log{PosZ: types.JsonIntPtr(int64(*s.PosZ))})
	}

	if s.IPAddress != nil {
		q = q.Where(types.Log{IPAddress: s.IPAddress})
	}

	if s.GeoCountry != nil {
		q = q.Where(types.Log{GeoCountry: s.GeoCountry})
	}

	if s.GeoCity != nil {
		q = q.Where(types.Log{GeoCity: s.GeoCity})
	}

	if s.GeoASN != nil {
		q = q.Where(types.Log{GeoASN: s.GeoASN})
	}

	if s.FromTimestamp != nil {
		q = q.Where("timestamp > ?", *s.FromTimestamp)
	}

	if s.ToTimestamp != nil {
		q = q.Where("timestamp < ?", *s.ToTimestamp)
	}

	// limit result length to 1000 per default
	limit := 1000
	if s.Limit != nil && *s.Limit < limit {
		limit = *s.Limit
	}
	q = q.Limit(limit)

	if s.SortAscending {
		q = q.Order("timestamp ASC")
	} else {
		q = q.Order("timestamp DESC")
	}

	return q
}

func (r *LogRepository) Search(s *types.LogSearch) ([]*types.Log, error) {
	var list []*types.Log
	err := r.query(s).Find(&list).Error
	return list, err
}

func (r *LogRepository) Count(s *types.LogSearch) (int, error) {
	var c int64
	err := r.query(s).Count(&c).Error
	return int(c), err
}

func (r *LogRepository) GetEvents(c types.LogCategory) ([]string, error) {
	var list []string
	err := r.g.Raw("select event from log where category = ? group by event", c).Find(&list).Error
	return list, err
}
