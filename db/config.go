package db

import (
	"mtui/types"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type ConfigRepository struct {
	g *gorm.DB
}

func (r *ConfigRepository) GetByKey(key types.ConfigKey) (*types.ConfigEntry, error) {
	return FindSingle[types.ConfigEntry](r.g.Where(types.ConfigEntry{Key: key}))
}

func (r *ConfigRepository) Set(c *types.ConfigEntry) error {
	return r.g.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "key"}},
		UpdateAll: true,
	}).Create(c).Error
}

func (r *ConfigRepository) Delete(key types.ConfigKey) error {
	return r.g.Delete(types.ConfigEntry{Key: key}).Error
}
