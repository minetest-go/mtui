package db

import (
	"mtui/types"

	"github.com/minetest-go/dbutil"
)

type MetricRepository struct {
	DB dbutil.DBTx
}

func (r *MetricRepository) Create(m *types.Metric) error {
	return dbutil.Insert(r.DB, m)
}
