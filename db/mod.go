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
	return FindMulti[types.Mod](r.g)
}

func (r *ModRepository) GetByName(name string) (*types.Mod, error) {
	return FindSingle[types.Mod](r.g.Where(types.Mod{Name: name}))
}

func (r *ModRepository) GetByID(id string) (*types.Mod, error) {
	return FindSingle[types.Mod](r.g.Where(types.Mod{ID: id}))
}

func (r *ModRepository) GetByStatus(status types.ModStatus) ([]*types.Mod, error) {
	return FindMulti[types.Mod](r.g.Where(types.Mod{Status: status}))
}

func (r *ModRepository) Update(m *types.Mod) error {
	return r.g.Model(m).Select("*").Updates(m).Error
}

func (r *ModRepository) Delete(id string) error {
	return r.g.Delete(types.Mod{ID: id}).Error
}

func (r *ModRepository) DeleteAll() error {
	return r.g.Session(&gorm.Session{AllowGlobalUpdate: true}).Delete(types.Mod{}).Error
}
