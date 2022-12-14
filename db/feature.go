package db

import (
	"database/sql"
	"mtui/types"

	"github.com/minetest-go/dbutil"
)

type FeatureRepository struct {
	DB dbutil.DBTx
}

func (r *FeatureRepository) Set(m *types.Feature) error {
	return dbutil.InsertOrReplace(r.DB, m)
}

func (r *FeatureRepository) GetAll() ([]*types.Feature, error) {
	return dbutil.SelectMulti(r.DB, func() *types.Feature { return &types.Feature{} }, "")
}

func (r *FeatureRepository) GetByName(name string) (*types.Feature, error) {
	f, err := dbutil.Select(r.DB, &types.Feature{}, "where name = $1", name)
	if err == sql.ErrNoRows {
		return nil, nil
	} else {
		return f, err
	}
}
