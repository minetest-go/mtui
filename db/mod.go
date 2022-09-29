package db

import (
	"database/sql"
	"mtui/types"
)

type ModRepository struct {
	DB *sql.DB
}

func (r *ModRepository) Create(m *types.Mod) error {
	return Insert(r.DB, m)
}

func (r *ModRepository) GetAll() ([]*types.Mod, error) {
	return SelectMulti(r.DB, func() *types.Mod { return &types.Mod{} }, "")
}

func (r *ModRepository) GetByID(id string) (*types.Mod, error) {
	m, err := Select(r.DB, &types.Mod{}, "where id = $1", id)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return m, err
}

func (r *ModRepository) Update(m *types.Mod) error {
	return Update(r.DB, m, map[string]any{"id": m.ID})
}

func (r *ModRepository) Delete(id string) error {
	return Delete(r.DB, &types.Mod{}, "where id = $1", id)
}

func (r *ModRepository) DeleteAll() error {
	return Delete(r.DB, &types.Mod{}, "")
}
