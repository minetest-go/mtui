package db

import (
	"database/sql"
	"mtui/types"

	"github.com/google/uuid"
	"github.com/minetest-go/dbutil"
)

type ModRepository struct {
	dbu *dbutil.DBUtil[*types.Mod]
}

func (r *ModRepository) Create(m *types.Mod) error {
	if m.ID == "" {
		m.ID = uuid.NewString()
	}
	return r.dbu.Insert(m)
}

func (r *ModRepository) GetAll() ([]*types.Mod, error) {
	return r.dbu.SelectMulti("")
}

func (r *ModRepository) GetByID(id string) (*types.Mod, error) {
	m, err := r.dbu.Select("where id = %s", id)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return m, err
}

func (r *ModRepository) Update(m *types.Mod) error {
	return r.dbu.Update(m, "where id = %s", m.ID)
}

func (r *ModRepository) Delete(id string) error {
	return r.dbu.Delete("where id = %s", id)
}

func (r *ModRepository) DeleteAll() error {
	return r.dbu.Delete("")
}
