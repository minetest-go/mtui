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

func (r *ModRepository) GetAll(t types.ModType) ([]*types.Mod, error) {
	return SelectMulti(r.DB, func() *types.Mod { return &types.Mod{} }, "where mod_type = $1", t)
}

func (r *ModRepository) Update(m *types.Mod) error {
	return Update(r.DB, m, map[string]any{"id": m.ID})
}

func (r *ModRepository) Delete(id string) error {
	return Delete(r.DB, &types.Mod{}, "where id = $1", id)
}
