package db

import (
	"mtui/types"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type MetricTypeRepository struct {
	g *gorm.DB
}

func (r *MetricTypeRepository) Insert(mt *types.MetricType) error {
	return r.g.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "name"}},
		UpdateAll: true,
	}).Create(mt).Error
}

func (r *MetricTypeRepository) GetByName(name string) (*types.MetricType, error) {
	var list []*types.MetricType
	err := r.g.Where(types.MetricType{Name: name}).Find(&list).Error
	if len(list) == 0 {
		return nil, err
	}
	return list[0], err
}

func (r *MetricTypeRepository) GetAll() ([]*types.MetricType, error) {
	var list []*types.MetricType
	err := r.g.Find(&list).Error
	return list, err
}
