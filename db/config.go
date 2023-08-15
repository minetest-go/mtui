package db

import (
	"database/sql"
	"mtui/types"

	"github.com/minetest-go/dbutil"
)

type ConfigRepository struct {
	dbu *dbutil.DBUtil[*types.ConfigEntry]
}

func (r *ConfigRepository) GetByKey(key types.ConfigKey) (*types.ConfigEntry, error) {
	c, err := r.dbu.Select("where key = %s", key)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return c, err
}

func (r *ConfigRepository) Set(c *types.ConfigEntry) error {
	return r.dbu.InsertOrReplace(c)
}

func (r *ConfigRepository) Delete(key string) error {
	return r.dbu.Delete("where key = %s", key)
}
