package db

import (
	"database/sql"
	"mtui/types"

	"github.com/minetest-go/dbutil"
)

type MetricTypeRepository struct {
	DB dbutil.DBTx
}

func (r *MetricTypeRepository) CreateOrUpdate(mt *types.MetricType) error {
	return dbutil.InsertOrReplace(r.DB, mt)
}

func (r *MetricTypeRepository) GetByName(name string) (*types.MetricType, error) {
	mt, err := dbutil.Select(r.DB, &types.MetricType{}, "where name = $1", name)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return mt, err
}

func (r *MetricTypeRepository) GetAll() ([]*types.MetricType, error) {
	return dbutil.SelectMulti(r.DB, func() *types.MetricType { return &types.MetricType{} }, "")
}
