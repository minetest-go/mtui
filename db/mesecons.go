package db

import (
	"mtui/types"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type MeseconsRepository struct {
	g *gorm.DB
}

func (r *MeseconsRepository) Save(m *types.Mesecons) error {
	return r.g.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "poskey"}},
		UpdateAll: true,
	}).Create(m).Error
}

func (r *MeseconsRepository) GetByPlayerName(playername string) ([]*types.Mesecons, error) {
	return FindMulti[types.Mesecons](r.g.Where(types.Mesecons{PlayerName: playername}).Order("order_id ASC"))
}

func (r *MeseconsRepository) GetByPoskey(poskey string) (*types.Mesecons, error) {
	return FindSingle[types.Mesecons](r.g.Where(types.Mesecons{PosKey: poskey}))
}

func (r *MeseconsRepository) Remove(poskey string) error {
	return r.g.Delete(types.Mesecons{PosKey: poskey}).Error
}
