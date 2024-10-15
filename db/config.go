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
	var list []*types.ConfigEntry
	err := r.g.Where(types.ConfigEntry{Key: key}).Limit(1).Find(&list).Error
	if len(list) == 0 {
		return nil, err
	}
	return list[0], err
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
