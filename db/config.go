package db

import (
	"database/sql"
	"mtui/types"
)

type ConfigRepository struct {
	DB *sql.DB
}

func (r *ConfigRepository) GetByKey(key string) (*types.ConfigEntry, error) {
	c, err := Select(r.DB, &types.ConfigEntry{}, "where key = $1", key)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return c, err
}

func (r *ConfigRepository) Set(c *types.ConfigEntry) error {
	return InsertOrReplace(r.DB, c)
}

func (r *ConfigRepository) Delete(key string) error {
	return Delete(r.DB, &types.ConfigEntry{}, "where key = $1", key)
}
