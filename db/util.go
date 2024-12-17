package db

import "gorm.io/gorm"

func FindSingle[T any](g *gorm.DB) (*T, error) {
	list := make([]*T, 0)
	err := g.Limit(1).Find(&list).Error
	if err != nil || len(list) == 0 {
		return nil, err
	} else {
		return list[0], err
	}
}

func FindMulti[T any](g *gorm.DB) ([]*T, error) {
	list := make([]*T, 0)
	err := g.Find(&list).Error
	return list, err
}
