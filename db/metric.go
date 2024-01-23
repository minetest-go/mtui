package db

import (
	"fmt"
	"mtui/types"

	"github.com/minetest-go/dbutil"
)

type MetricRepository struct {
	dbu *dbutil.DBUtil[*types.Metric]
}

func (r *MetricRepository) Insert(m *types.Metric) error {
	return r.dbu.Insert(m)
}

func (r *MetricRepository) buildWhereClause(s *types.MetricSearch, order bool) (string, []any) {
	q := "where true "
	args := make([]any, 0)

	if s.Name != nil {
		q += " and name = %s"
		args = append(args, *s.Name)
	}

	if s.FromTimestamp != nil {
		q += " and timestamp > %s"
		args = append(args, *s.FromTimestamp)
	}

	if s.ToTimestamp != nil {
		q += " and timestamp < %s"
		args = append(args, *s.ToTimestamp)
	}

	if order {
		// limit result length to 1000 per default
		limit := 1000
		if s.Limit != nil && *s.Limit < limit {
			limit = *s.Limit
		}
		q += fmt.Sprintf(" order by timestamp desc limit %d", limit)
	}

	return q, args
}

func (r *MetricRepository) Search(s *types.MetricSearch) ([]*types.Metric, error) {
	q, args := r.buildWhereClause(s, true)
	return r.dbu.SelectMulti(q, args...)
}

func (r *MetricRepository) Count(s *types.MetricSearch) (int, error) {
	q, args := r.buildWhereClause(s, false)
	return r.dbu.Count(q, args...)
}

func (r *MetricRepository) DeleteBefore(timestamp int64) error {
	return r.dbu.Delete("where timestamp < %s", timestamp)
}
