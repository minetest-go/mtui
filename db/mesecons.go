package db

import (
	"database/sql"
	"mtui/types"

	"github.com/minetest-go/dbutil"
)

type MeseconsRepository struct {
	dbu *dbutil.DBUtil[*types.Mesecons]
}

func (r *MeseconsRepository) Save(m *types.Mesecons) error {
	return r.dbu.InsertOrReplace(m)
}

func (r *MeseconsRepository) GetByPlayerName(playername string) ([]*types.Mesecons, error) {
	return r.dbu.SelectMulti("where playername = %s order by order_id asc", playername)
}

func (r *MeseconsRepository) GetByPoskey(poskey string) (*types.Mesecons, error) {
	m, err := r.dbu.Select("where poskey = %s", poskey)
	if err == sql.ErrNoRows {
		return nil, nil
	} else {
		return m, err
	}
}

func (r *MeseconsRepository) Remove(poskey string) error {
	return r.dbu.Delete("where poskey = %s", poskey)
}
