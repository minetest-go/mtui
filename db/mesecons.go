package db

import (
	"mtui/types"

	"github.com/minetest-go/dbutil"
)

type MeseconsRepository struct {
	dbu *dbutil.DBUtil[*types.Mesecons]
}

func (r *MeseconsRepository) Set(m *types.Feature) error {
	return r.dbu.InsertOrReplace(m)
}

func (r *MeseconsRepository) GetByPlayerName(playername string) ([]*types.Mesecons, error) {
	return r.dbu.SelectMulti("where playername = %s", playername)
}
