package db

import (
	"mtui/types"

	"gorm.io/gorm"
)

type MetricRepository struct {
	g *gorm.DB
}

func (r *MetricRepository) Insert(m *types.Metric) error {
	return r.g.Create(m).Error
}

func (r *MetricRepository) query(s *types.MetricSearch, order bool) *gorm.DB {
	q := r.g.Model(types.Metric{})

	if s.Name != nil {
		q = q.Where(types.Metric{Name: *s.Name})
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

	if order {
		q = q.Order("timestamp DESC")
	}

	return q
}

func (r *MetricRepository) Search(s *types.MetricSearch) ([]*types.Metric, error) {
	var list []*types.Metric
	err := r.query(s, true).Find(&list).Error
	return list, err
}

func (r *MetricRepository) Count(s *types.MetricSearch) (int, error) {
	var c int64
	err := r.query(s, false).Count(&c).Error
	return int(c), err
}

func (r *MetricRepository) DeleteBefore(timestamp int64) error {
	return r.g.Where("timestamp < ?", timestamp).Delete(types.Metric{}).Error
}
