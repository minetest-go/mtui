package db

import (
	"database/sql"
	"mtui/types"

	"github.com/minetest-go/dbutil"
)

type ConfigRepository struct {
	DB dbutil.DBTx
}

func (r *ConfigRepository) GetByKey(key types.ConfigKey) (*types.ConfigEntry, error) {
	c, err := dbutil.Select(r.DB, &types.ConfigEntry{}, "where key = $1", key)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return c, err
}

func (r *ConfigRepository) Set(c *types.ConfigEntry) error {
	return dbutil.InsertOrReplace(r.DB, c)
}

func (r *ConfigRepository) Delete(key string) error {
	return dbutil.Delete(r.DB, &types.ConfigEntry{}, "where key = $1", key)
}
