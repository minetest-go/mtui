package db

import (
	"database/sql"
	"mtui/types"
)

type FeatureRepository struct {
	DB *sql.DB
}

func (r *FeatureRepository) Set(m *types.Feature) error {
	return InsertOrReplace(r.DB, m)
}

func (r *FeatureRepository) GetAll() ([]*types.Feature, error) {
	return SelectMulti(r.DB, func() *types.Feature { return &types.Feature{} }, "")
}

func (r *FeatureRepository) GetByName(name string) (*types.Feature, error) {
	f, err := Select(r.DB, &types.Feature{}, "where name = $1", name)
	if err == sql.ErrNoRows {
		return nil, nil
	} else {
		return f, err
	}
}
