package db

import (
	"database/sql"
	"mtui/types"

	"github.com/google/uuid"
	"github.com/minetest-go/dbutil"
)

type ModRepository struct {
	DB dbutil.DBTx
}

func (r *ModRepository) Create(m *types.Mod) error {
	if m.ID == "" {
		m.ID = uuid.NewString()
	}
	return dbutil.Insert(r.DB, m)
}

func (r *ModRepository) GetAll() ([]*types.Mod, error) {
	return dbutil.SelectMulti(r.DB, func() *types.Mod { return &types.Mod{} }, "")
}

func (r *ModRepository) GetByID(id string) (*types.Mod, error) {
	m, err := dbutil.Select(r.DB, &types.Mod{}, "where id = $1", id)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return m, err
}

func (r *ModRepository) Update(m *types.Mod) error {
	return dbutil.Update(r.DB, m, "where id = $1", m.ID)
}

func (r *ModRepository) Delete(id string) error {
	return dbutil.Delete(r.DB, &types.Mod{}, "where id = $1", id)
}

func (r *ModRepository) DeleteAll() error {
	return dbutil.Delete(r.DB, &types.Mod{}, "")
}
