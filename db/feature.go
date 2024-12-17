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
	return FindMulti[types.Feature](r.g)
}

func (r *FeatureRepository) GetByName(name string) (*types.Feature, error) {
	return FindSingle[types.Feature](r.g.Where(types.Feature{Name: name}))
}
