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
	var list []*types.Mesecons
	err := r.g.Where(types.Mesecons{PlayerName: playername}).Order("order_id ASC").Find(&list).Error
	return list, err
}

func (r *MeseconsRepository) GetByPoskey(poskey string) (*types.Mesecons, error) {
	var list []*types.Mesecons
	err := r.g.Where(types.Mesecons{PosKey: poskey}).Find(&list).Error
	if len(list) == 0 {
		return nil, err
	}
	return list[0], err
}

func (r *MeseconsRepository) Remove(poskey string) error {
	return r.g.Delete(types.Mesecons{PosKey: poskey}).Error
}
