package db

import (
	"fmt"
	"mtui/types"

	"github.com/minetest-go/dbutil"
)

type MetricRepository struct {
	DB dbutil.DBTx
}

func (r *MetricRepository) Insert(m *types.Metric) error {
	return dbutil.Insert(r.DB, m)
}

func (r *MetricRepository) buildWhereClause(s *types.MetricSearch, order bool) (string, []interface{}) {
	q := "where true "
	args := make([]interface{}, 0)
	i := 1

	if s.Name != nil {
		q += fmt.Sprintf(" and name = $%d", i)
		args = append(args, *s.Name)
		i++
	}

	if s.FromTimestamp != nil {
		q += fmt.Sprintf(" and timestamp > $%d", i)
		args = append(args, *s.FromTimestamp)
		i++
	}

	if s.ToTimestamp != nil {
		q += fmt.Sprintf(" and timestamp < $%d", i)
		args = append(args, *s.ToTimestamp)
		i++
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
	return dbutil.SelectMulti(r.DB, func() *types.Metric { return &types.Metric{} }, q, args...)
}

func (r *MetricRepository) Count(s *types.MetricSearch) (int, error) {
	q, args := r.buildWhereClause(s, false)
	return dbutil.Count(r.DB, &types.Metric{}, q, args...)
}

func (r *MetricRepository) DeleteBefore(timestamp int64) error {
	return dbutil.Delete(r.DB, &types.Metric{}, "where timestamp < $1", timestamp)
}
