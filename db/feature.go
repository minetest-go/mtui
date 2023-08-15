package db

import (
	"database/sql"
	"mtui/types"

	"github.com/minetest-go/dbutil"
)

type FeatureRepository struct {
	dbu *dbutil.DBUtil[*types.Feature]
}

func (r *FeatureRepository) Set(m *types.Feature) error {
	return r.dbu.InsertOrReplace(m)
}

func (r *FeatureRepository) GetAll() ([]*types.Feature, error) {
	return r.dbu.SelectMulti("")
}

func (r *FeatureRepository) GetByName(name string) (*types.Feature, error) {
	f, err := r.dbu.Select("where name = %s", name)
	if err == sql.ErrNoRows {
		return nil, nil
	} else {
		return f, err
	}
}
