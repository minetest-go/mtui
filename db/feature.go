package db

import (
	"mtui/types"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type FeatureRepository struct {
	g *gorm.DB
}

func (r *FeatureRepository) Set(m *types.Feature) error {
	return r.g.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "name"}},
		UpdateAll: true,
	}).Create(m).Error
}

func (r *FeatureRepository) GetAll() ([]*types.Feature, error) {
	var list []*types.Feature
	err := r.g.Find(&list).Error
	return list, err
}

func (r *FeatureRepository) GetByName(name string) (*types.Feature, error) {
	var list []*types.Feature
	err := r.g.Where(types.Feature{Name: name}).Limit(1).Find(&list).Error
	if len(list) == 0 {
		return nil, err
	}
	return list[0], err
}
