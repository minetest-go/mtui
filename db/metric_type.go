package db

import (
	"database/sql"
	"mtui/types"

	"github.com/minetest-go/dbutil"
)

type MetricTypeRepository struct {
	dbu *dbutil.DBUtil[*types.MetricType]
}

func (r *MetricTypeRepository) Insert(mt *types.MetricType) error {
	return r.dbu.InsertOrReplace(mt)
}

func (r *MetricTypeRepository) GetByName(name string) (*types.MetricType, error) {
	mt, err := r.dbu.Select("where name = %s", name)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return mt, err
}

func (r *MetricTypeRepository) GetAll() ([]*types.MetricType, error) {
	return r.dbu.SelectMulti("")
}
