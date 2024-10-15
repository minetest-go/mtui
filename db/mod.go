package db

import (
	"mtui/types"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ModRepository struct {
	g *gorm.DB
}

func (r *ModRepository) Create(m *types.Mod) error {
	if m.ID == "" {
		m.ID = uuid.NewString()
	}
	return r.g.Create(m).Error
}

func (r *ModRepository) GetAll() ([]*types.Mod, error) {
	var list []*types.Mod
	err := r.g.Find(&list).Error
	return list, err
}

func (r *ModRepository) GetByName(name string) (*types.Mod, error) {
	var list []*types.Mod
	err := r.g.Where(types.Mod{Name: name}).Limit(1).Find(&list).Error
	if len(list) == 0 {
		return nil, err
	}
	return list[0], err
}

func (r *ModRepository) GetByID(id string) (*types.Mod, error) {
	var list []*types.Mod
	err := r.g.Where(&types.Mod{ID: id}).Limit(1).Find(&list).Error
	if len(list) == 0 {
		return nil, err
	}
	return list[0], err
}

func (r *ModRepository) Update(m *types.Mod) error {
	return r.g.Model(m).Updates(m).Error
}

func (r *ModRepository) Delete(id string) error {
	return r.g.Delete(types.Mod{ID: id}).Error
}

func (r *ModRepository) DeleteAll() error {
	return r.g.Session(&gorm.Session{AllowGlobalUpdate: true}).Delete(types.Mod{}).Error
}
