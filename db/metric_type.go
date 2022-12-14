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
	q := `insert or replace into metric_type(name, type, help) values($1,$2,$3)`
	_, err := r.DB.Exec(q, mt.Name, mt.Type, mt.Help)
	return err
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
