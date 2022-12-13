package db

import (
	"database/sql"
	"mtui/types"
)

type MetricTypeRepository struct {
	DB *sql.DB
}

func (r *MetricTypeRepository) CreateOrUpdate(mt *types.MetricType) error {
	q := `insert or replace into metric_type(name, type, help) values($1,$2,$3)`
	_, err := r.DB.Exec(q, mt.Name, mt.Type, mt.Help)
	return err
}

func (r *MetricTypeRepository) GetByName(name string) (*types.MetricType, error) {
	q := `select name,type,help from metric_type where name = $1`
	row := r.DB.QueryRow(q, name)
	mt := &types.MetricType{}
	err := row.Scan(&mt.Name, &mt.Type, &mt.Help)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return mt, err
}
