package db

import (
	"database/sql"
	"mtui/types"

	"github.com/google/uuid"
)

type ModRepository struct {
	db *sql.DB
}

func (r *ModRepository) columns() string {
	return "id,name,mod_type,source_type,url,version,auto_update"
}

func (r *ModRepository) scan(rows *sql.Rows) (*types.Mod, error) {
	entry := &types.Mod{}
	return entry, rows.Scan(
		&entry.ID,
		&entry.Name,
		&entry.ModType,
		&entry.SourceType,
		&entry.URL,
		&entry.Version,
		&entry.AutoUpdate,
	)
}

func (r *ModRepository) GetAll() ([]*types.Mod, error) {
	rows, err := r.db.Query("select " + r.columns() + " from mod")
	if err == sql.ErrNoRows {
		return nil, nil
	} else if err != nil {
		return nil, err
	}

	list := make([]*types.Mod, 0)
	for rows.Next() {
		m, err := r.scan(rows)
		if err != nil {
			return nil, err
		}
		list = append(list, m)
	}
	return list, nil
}

func (r *ModRepository) Delete(id string) error {
	_, err := r.db.Exec("delete from mod where id = $1", id)
	return err
}

func (r *ModRepository) Save(mod *types.Mod) error {
	if mod.ID == "" {
		mod.ID = uuid.NewString()
	}

	_, err := r.db.Exec(`
		insert or replace into mod(`+r.columns()+`)
		values($1,$2,$3,$4,$5,$6,$7);
	`,
		mod.ID,
		mod.Name,
		mod.ModType,
		mod.SourceType,
		mod.URL,
		mod.Version,
		mod.AutoUpdate,
	)

	return err
}